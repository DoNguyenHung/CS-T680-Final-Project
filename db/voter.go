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
	RedisKeyPrefix       = "voter:"
)

type cache struct {
	cacheClient *redis.Client
	jsonHelper  *rejson.Handler
	context     context.Context
}

// ToDoItem is the struct that represents a single ToDo item
type voterPoll struct {
	PollID   uint      `json:"pollid"`
	VoteDate time.Time `json:"votedate"`
}

// VoterList is a type alias for a map of Voters.  The key
// will be the ToDoItem.Id and the value will be the ToDoItem
// type DbMap map[int]ToDoItem

type Voter struct {
	VoterId     uint        `json:"id"`
	FirstName   string      `json:"firstname"`
	LastName    string      `json:"lastname"`
	VoteHistory []voterPoll `json:"votehistory"`
}

type VoterList struct {
	//more things would be included in a real implementation

	//Redis cache connections
	cache
}

//------------------------------------------------------------
// THESE ARE THE PUBLIC FUNCTIONS THAT SUPPORT OUR TODO APP
//------------------------------------------------------------

func NewVoter(id uint, fn, ln string) *Voter {
	return &Voter{
		FirstName:   fn,
		LastName:    ln,
		VoteHistory: []voterPoll{},
	}
}

func (v *Voter) AddPoll(pollID uint) {
	v.VoteHistory = append(v.VoteHistory, voterPoll{PollID: pollID, VoteDate: time.Now()})
}

func NewVoterList() (*VoterList, error) {

	//We will use an override if the REDIS_URL is provided as an environment
	//variable, which is the preferred way to wire up a docker container
	redisUrl := os.Getenv("REDIS_URL")
	//This handles the default condition
	if redisUrl == "" {
		redisUrl = RedisDefaultLocation
	}
	return NewWithCacheInstance(redisUrl)

}

// NewWithCacheInstance is a constructor function that returns a pointer to a new
// ToDo struct.  It accepts a string that represents the location of the redis
// cache.
func NewWithCacheInstance(location string) (*VoterList, error) {

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
	return &VoterList{
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
func (v *VoterList) getItemFromRedis(key string, item *Voter) error {

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

func (lst *VoterList) AddVoter(voter Voter) error {

	// lst.Voters[voter.VoterId] = voter
	// return nil

	//Before we add an item to the DB, lets make sure
	//it does not exist, if it does, return an error

	redisKey := redisKeyFromId(int(voter.VoterId))
	var existingVoter Voter
	if err := lst.getItemFromRedis(redisKey, &existingVoter); err == nil {
		return errors.New("voter already exists")
	}

	//Add item to database with JSON Set
	if _, err := lst.jsonHelper.JSONSet(redisKey, ".", voter); err != nil {
		return err
	}

	//If everything is ok, return nil for the error
	return nil
}

func (lst *VoterList) DeleteVoter(id uint) error {

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

func (lst *VoterList) DeleteAll() error {
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

func (lst *VoterList) UpdateVoter(voter Voter) error {

	//Before we add an item to the DB, lets make sure
	//it does not exist, if it does, return an error
	redisKey := redisKeyFromId(int(voter.VoterId))
	var existingItem Voter
	if err := lst.getItemFromRedis(redisKey, &existingItem); err != nil {
		return errors.New("item does not exist")
	}

	//Add item to database with JSON Set.  Note there is no update
	//functionality, so we just overwrite the existing item
	if _, err := lst.jsonHelper.JSONSet(redisKey, ".", voter); err != nil {
		return err
	}

	//If everything is ok, return nil for the error
	return nil
}

/*
Get a single voter resource with voterID=:id including their entire voting history.
POST version adds one to the "database"
*/
func (lst *VoterList) GetSingleVoterResource(id uint) (Voter, error) {

	// Check if item exists before trying to get it
	// this is a good practice, return an error if the
	// item does not exist
	var voter Voter
	pattern := redisKeyFromId(int(id))
	err := lst.getItemFromRedis(pattern, &voter)
	if err != nil {
		return Voter{}, err
	}

	return voter, nil
}

// ChangeItemDoneStatus accepts an item id and a boolean status.
// It returns an error if the status could not be updated for any
// reason.  For example, the item itself does not exist, or an
// IO error trying to save the updated status.

// Preconditions:   (1) The database file must exist and be a valid
//
//					(2) The item must exist in the DB
//	    				because we use the item.Id as the key, this
//						function must check if the item already
//	    				exists in the DB, if not, return an error
//
// Postconditions:
//
//	 (1) The items status in the database will be updated
//		(2) If there is an error, it will be returned.
//		(3) This function MUST use existing functionality for most of its
//			work.  For example, it should call GetItem() to get the item
//			from the DB, then it should call UpdateItem() to update the
//			item in the DB (after the status is changed).
func (lst *VoterList) ChangeItemDoneStatus(id int, value bool) error {

	//update was successful
	return errors.New("not implemented")
}

// TODO
/*
Gets JUST the voter history for the voter with VoterID = :id
*/
func (lst *VoterList) GetVoterHistory(id uint) ([]voterPoll, error) {

	var voter Voter
	pattern := redisKeyFromId(int(id))
	err := lst.getItemFromRedis(pattern, &voter)
	if err != nil {
		return []voterPoll{}, err
	}

	return voter.VoteHistory, nil

}

/*
Get all voter resources including all voter history for each voter (note we will
discuss the concept of "paging" later, for now you can ignore)
*/
func (lst *VoterList) GetAllVoters() ([]Voter, error) {

	//Now that we have the DB loaded, lets crate a slice
	var voterList []Voter

	//Lets query redis for all of the items
	pattern := RedisKeyPrefix + "*"
	ks, _ := lst.cacheClient.Keys(lst.context, pattern).Result()
	for _, key := range ks {
		var voter Voter
		err := lst.getItemFromRedis(key, &voter)
		if err != nil {
			return nil, err
		}
		voterList = append(voterList, voter)
	}

	return voterList, nil
}

/*
Gets JUST the single voter poll data with PollID = :id and VoterID = :id.
*/
func (lst *VoterList) GetVoterPollData(voterId uint, pollId uint) (*voterPoll, error) {

	var currentVoter Voter
	pattern := redisKeyFromId(int(voterId))
	err := lst.getItemFromRedis(pattern, &currentVoter)
	if err != nil {
		return &voterPoll{}, err
	}

	for j := 0; j < len(currentVoter.VoteHistory); j++ {
		currentPoll := currentVoter.VoteHistory[j]
		if currentPoll.PollID == pollId {
			return &currentPoll, nil
		}
	}

	return nil, errors.New("item does not exist")

}

func (lst *VoterList) AddVoterPollData(voterId uint, pollId uint) error {

	var currentVoter Voter
	pattern := redisKeyFromId(int(voterId))
	err := lst.getItemFromRedis(pattern, &currentVoter)
	if err != nil {
		newVoter := Voter{
			VoterId: voterId,
		}
		// lst.Voters[voterId] = newVoter

		newPoll := voterPoll{
			PollID:   pollId,
			VoteDate: time.Now(),
		}
		newVoter.VoteHistory = append(newVoter.VoteHistory, newPoll)

		//Add item to database with JSON Set
		redisKey := redisKeyFromId(int(voterId))
		if _, err := lst.jsonHelper.JSONSet(redisKey, ".", newVoter); err != nil {
			return err
		}
		// lst.Voters[voterId] = newVoter

		return nil
	}

	newPoll := voterPoll{
		PollID:   pollId,
		VoteDate: time.Now(),
	}
	currentVoter.VoteHistory = append(currentVoter.VoteHistory, newPoll)

	// lst.Voters[voterId] = currentVoter
	//Add item to database with JSON Set
	redisKey := redisKeyFromId(int(voterId))
	if _, err := lst.jsonHelper.JSONSet(redisKey, ".", currentVoter); err != nil {
		return err
	}

	return nil
}

func (lst *VoterList) DeletePoll(voterId uint, pollId uint) error {

	var currentVoter Voter
	pattern := redisKeyFromId(int(voterId))
	err := lst.getItemFromRedis(pattern, &currentVoter)
	if err != nil {
		return err
	}

	index := 0
	for j := 0; j < len(currentVoter.VoteHistory); j++ {
		currentPoll := currentVoter.VoteHistory[j]
		if currentPoll.PollID == pollId {
			index = j
		}
	}

	currentVoter.VoteHistory = append(currentVoter.VoteHistory[:index], currentVoter.VoteHistory[index+1:]...)

	// lst.Voters[voterId] = currentVoter
	//Add item to database with JSON Set
	// TODO: How to override voter here?
	redisKey := redisKeyFromId(int(voterId))
	if _, err := lst.jsonHelper.JSONSet(redisKey, ".", currentVoter); err != nil {
		return err
	}

	return nil
}
