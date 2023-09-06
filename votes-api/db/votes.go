// docker compose up
// load cache

// Docker repo: https://hub.docker.com/repository/docker/hungdo171/voter-container/general

package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"
)

const (
	RedisNilError        = "redis: nil"
	RedisDefaultLocation = "0.0.0.0:6379"
	RedisKeyPrefix       = "vote:"
)

type cache struct {
	cacheClient *redis.Client
	jsonHelper  *rejson.Handler
	context     context.Context
}

type pollOption struct {
	PollOptionID   uint   `json:"pollOptionId"`
	PollOptionText string `json:"pollOptionText"`
}

type Poll struct {
	PollID       uint         `json:"pollId"`
	PollTitle    string       `json:"pollTitle"`
	PollQuestion string       `json:"pollQuestion"`
	PollOptions  []pollOption `json:"pollOptions"`
}

type voterPoll struct {
	PollID   uint      `json:"pollid"`
	VoteDate time.Time `json:"votedate"`
}

type Voter struct {
	VoterId     uint        `json:"id"`
	FirstName   string      `json:"firstname"`
	LastName    string      `json:"lastname"`
	VoteHistory []voterPoll `json:"votehistory"`
}

type Vote struct {
	VoteID    uint `json:"voteId"`
	VoterID   uint `json:"voterId"`
	PollID    uint `json:"pollId"`
	VoteValue uint `json:"voteValue"`
}

type VoteList struct {
	//more things would be included in a real implementation

	//Redis cache connections
	cache
}

//------------------------------------------------------------
// THESE ARE THE PUBLIC FUNCTIONS THAT SUPPORT OUR TODO APP
//------------------------------------------------------------

// func NewSampleVote(voteId uint, pollId string, voterId string, voteValue uint) *Vote {
func NewSampleVote(voteId uint, pollId uint, voterId uint, voteValue uint) *Vote {
	return &Vote{
		VoteID:    voteId,
		PollID:    pollId,
		VoterID:   voterId,
		VoteValue: voteValue,
	}
}

func NewVoteList() (*VoteList, error) {

	//We will use an override if the REDIS_URL is provided as an environment
	//variable, which is the preferred way to wire up a docker container
	redisUrl := os.Getenv("CACHE_URL")
	//This handles the default condition
	if redisUrl == "" {
		redisUrl = RedisDefaultLocation
	}
	log.Println("using REDIS at: ", redisUrl)
	return NewWithCacheInstance(redisUrl)

}

// NewWithCacheInstance is a constructor function that returns a pointer to a new
// ToDo struct.  It accepts a string that represents the location of the redis
// cache.
func NewWithCacheInstance(location string) (*VoteList, error) {

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
	err := client.Ping(ctx).Err()
	if err != nil {
		log.Println("Error connecting to redis" + err.Error())
		return nil, err
	}

	//By default, redis manages keys and values, where the values
	//are either strings, sets, maps, etc.  Redis has an extension
	//module called ReJSON that allows us to store JSON objects
	//however, we need a companion library in order to work with it
	//Below we create an instance of the JSON helper and associate
	//it with our redis connnection
	jsonHelper := rejson.NewReJSONHandler()
	jsonHelper.SetGoRedisClientWithContext(ctx, client)

	//Return a pointer to a new ToDo struct
	return &VoteList{
		cache: cache{
			cacheClient: client,
			jsonHelper:  jsonHelper,
			context:     ctx,
		},
	}, nil
}

//------------------------------------------------------------
// REDIS HELPERS
//------------------------------------------------------------

// We will use this later, you can ignore for now
func isRedisNilError(err error) bool {
	return errors.Is(err, redis.Nil) || err.Error() == RedisNilError
}

// In redis, our keys will be strings, they will look like
// todo:<number>.  This function will take an integer and
// return a string that can be used as a key in redis
func redisKeyFromId(id int) string {
	return fmt.Sprintf("%s%d", RedisKeyPrefix, id)
}

// Helper to return a ToDoItem from redis provided a key
func (v *VoteList) getItemFromRedis(key string, item *Vote) error {

	//Lets query redis for the item, note we can return parts of the
	//json structure, the second parameter "." means return the entire
	//json structure
	itemObject, err := v.jsonHelper.JSONGet(key, ".")
	if err != nil {
		return err
	}

	//JSONGet returns an "any" object, or empty interface,
	//we need to convert it to a byte array, which is the
	//underlying type of the object, then we can unmarshal
	//it into our ToDoItem struct
	err = json.Unmarshal(itemObject.([]byte), item)
	if err != nil {
		return nil
	}

	return nil
}

//------------------------------------------------------------
// THESE ARE THE PUBLIC FUNCTIONS THAT SUPPORT OUR VOTER APP
//------------------------------------------------------------

func (lst *VoteList) AddVote(vote Vote) error {

	// lst.Voters[voter.VoterId] = voter
	// return nil

	//Before we add an item to the DB, lets make sure
	//it does not exist, if it does, return an error

	redisKey := redisKeyFromId(int(vote.VoteID))
	var existingVote Vote
	if err := lst.getItemFromRedis(redisKey, &existingVote); err == nil {
		return errors.New("voter already exists")
	}

	//Add item to database with JSON Set
	if _, err := lst.jsonHelper.JSONSet(redisKey, ".", vote); err != nil {
		return err
	}

	//If everything is ok, return nil for the error
	return nil
}

func (lst *VoteList) DeleteVote(id uint) error {

	pattern := redisKeyFromId(int(id))
	numDeleted, err := lst.cacheClient.Del(lst.context, pattern).Result()
	if err != nil {
		return err
	}
	if numDeleted == 0 {
		return errors.New("attempted to delete non-existent item")
	}

	return nil
}

func (lst *VoteList) DeleteAll() error {
	pattern := RedisKeyPrefix + "*"
	ks, _ := lst.cacheClient.Keys(lst.context, pattern).Result()
	//Note delete can take a collection of keys.  In go we can
	//expand a slice into individual arguments by using the ...
	//operator
	numDeleted, err := lst.cacheClient.Del(lst.context, ks...).Result()
	if err != nil {
		return err
	}

	if numDeleted != int64(len(ks)) {
		return errors.New("one or more items could not be deleted")
	}

	return nil
}

/*
Get a single voter resource with voterID=:id including their entire voting history.
POST version adds one to the "database"
*/
func (lst *VoteList) GetSingleVoterResource(id uint) (Vote, error) {

	// Check if item exists before trying to get it
	// this is a good practice, return an error if the
	// item does not exist
	var vote Vote
	pattern := redisKeyFromId(int(id))
	err := lst.getItemFromRedis(pattern, &vote)
	if err != nil {
		return Vote{}, err
	}

	return vote, nil
}

func (lst *VoteList) ChangeItemDoneStatus(id int, value bool) error {

	//update was successful
	return errors.New("not implemented")
}

/*
Get all voter resources including all voter history for each voter (note we will
discuss the concept of "paging" later, for now you can ignore)
*/
func (lst *VoteList) GetAllVotes() ([]Vote, error) {

	//Now that we have the DB loaded, lets crate a slice
	var voteList []Vote

	//Lets query redis for all of the items
	pattern := RedisKeyPrefix + "*"
	ks, _ := lst.cacheClient.Keys(lst.context, pattern).Result()
	for _, key := range ks {
		var vote Vote
		err := lst.getItemFromRedis(key, &vote)
		if err != nil {
			return nil, err
		}
		voteList = append(voteList, vote)
	}

	return voteList, nil
}
