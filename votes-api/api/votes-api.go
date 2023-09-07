package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"drexel.edu/votes-api/db"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"

	"github.com/go-resty/resty/v2"
)

type cache struct {
	client  *redis.Client
	helper  *rejson.Handler
	context context.Context
}

type VoteAPI struct {
	cache
	voterAPIURL string
	pollAPIURL  string
	apiClient   *resty.Client
	db          *db.VoteList
}

func NewVoteAPI(location string, voterAPIURL string, pollAPIURL string) (*VoteAPI, error) {
	apiClient := resty.New()
	dbHandler, err := db.NewVoteList()

	//Connect to redis.  Other options can be provided, but the
	//defaults are OK
	client := redis.NewClient(&redis.Options{
		Addr: location,
	})

	//We use this context to coordinate betwen our go code and
	//the redis operaitons
	ctx := context.Background()

	//This is the reccomended way to ensure that our redis connection
	//is working
	err2 := client.Ping(ctx).Err()
	if err2 != nil {
		log.Println("Error connecting to redis" + err.Error())
		return nil, err2
	}

	jsonHelper := rejson.NewReJSONHandler()
	jsonHelper.SetGoRedisClientWithContext(ctx, client)

	//Return a pointer to a new ToDo struct
	return &VoteAPI{
		cache: cache{
			client:  client,
			helper:  jsonHelper,
			context: ctx,
		},
		voterAPIURL: voterAPIURL,
		pollAPIURL:  pollAPIURL,
		db:          dbHandler,
		apiClient:   apiClient,
	}, nil
}

func (v *VoteAPI) GetVote(c *gin.Context) {
	voteId := c.Param("id")
	if voteId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No vote ID provided"})
		return
	}

	cacheKey := "vote:" + voteId
	v1Bytes, err := v.helper.JSONGet(cacheKey, ".")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find vote in cache with id=" + cacheKey})
		return
	}

	var v1 db.Vote
	err = json.Unmarshal(v1Bytes.([]byte), &v1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cached data seems to be wrong type"})
		return
	}

	c.JSON(http.StatusOK, v1)
}

func (v *VoteAPI) GetVoterByVote(c *gin.Context) {
	voteId := c.Param("id")
	if voteId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No vote ID provided"})
		return
	}

	cacheKey := "vote:" + voteId
	var v1 db.Vote
	err := v.getItemFromRedis(cacheKey, &v1)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find vote in cache with id=" + cacheKey})
		return
	}

	voterLocation := v1.VoterID

	voterURL := v.voterAPIURL + "/voters/" + strconv.FormatUint(uint64(voterLocation), 10)
	log.Println(voterURL)
	var voter db.Voter

	_, err = v.apiClient.R().SetResult(&voter).Get(voterURL)
	if err != nil {
		emsg := "Could not get voter from API: (" + voterURL + ")" + err.Error()
		c.JSON(http.StatusNotFound, gin.H{"error": emsg})
		return
	}

	c.JSON(http.StatusOK, voter)

}

func (v *VoteAPI) GetPollByVote(c *gin.Context) {

	voteId := c.Param("id")
	if voteId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No vote ID provided"})
		return
	}

	cacheKey := "vote:" + voteId
	var v1 db.Vote
	err := v.getItemFromRedis(cacheKey, &v1)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Could not find vote in cache with id=" + cacheKey})
		return
	}

	pollLocation := v1.PollID

	// pollURL := v.pollAPIURL + "/polls/" + string(pollLocation)
	pollURL := v.pollAPIURL + "/polls/" + strconv.FormatUint(uint64(pollLocation), 10)
	log.Println(pollURL)
	var poll db.Poll

	_, err = v.apiClient.R().SetResult(&poll).Get(pollURL)
	if err != nil {
		emsg := "Could not get voter from API: (" + pollURL + ")" + err.Error()
		c.JSON(http.StatusNotFound, gin.H{"error": emsg})
		return
	}

	c.JSON(http.StatusOK, poll)
}

func (v *VoteAPI) GetAllVotes(c *gin.Context) {
	var voteList []db.Vote
	var voteItem db.Vote

	//Lets query redis for all of the items
	votePattern := "vote:*"
	voteKs, _ := v.client.Keys(v.context, votePattern).Result()
	for _, key := range voteKs {
		err := v.getItemFromRedis(key, &voteItem)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find reading list in cache with id=" + key})
			return
		}
		voteList = append(voteList, voteItem)
	}

	c.JSON(http.StatusOK, voteList)
}

// Helper to return a ToDoItem from redis provided a key
func (v *VoteAPI) getItemFromRedis(key string, rl *db.Vote) error {

	//Lets query redis for the item, note we can return parts of the
	//json structure, the second parameter "." means return the entire
	//json structure
	itemObject, err := v.helper.JSONGet(key, ".")
	if err != nil {
		return err
	}

	//JSONGet returns an "any" object, or empty interface,
	//we need to convert it to a byte array, which is the
	//underlying type of the object, then we can unmarshal
	//it into our ToDoItem struct
	err = json.Unmarshal(itemObject.([]byte), rl)
	if err != nil {
		return err
	}

	return nil
}

// implementation for DELETE /todo
// deletes all todos
func (v *VoteAPI) DeleteAllVotes(c *gin.Context) {

	if err := v.db.DeleteAll(); err != nil {
		log.Println("Error deleting all items: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

// implementation for POST /todo
// adds a new todo
func (v *VoteAPI) AddVote(c *gin.Context) {
	var vote db.Vote

	//With HTTP based APIs, a POST request will usually
	//have a body that contains the data to be added
	//to the database.  The body is usually JSON, so
	//we need to bind the JSON to a struct that we
	//can use in our code.
	//This framework exposes the raw body via c.Request.Body
	//but it also provides a helper function ShouldBindJSON()
	//that will extract the body, convert it to JSON and
	//bind it to a struct for us.  It will also report an error
	//if the body is not JSON or if the JSON does not match
	//the struct we are binding to.

	if err := c.ShouldBindJSON(&vote); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := v.db.AddVote(vote); err != nil {
		log.Println("Error adding item: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, vote)
}
