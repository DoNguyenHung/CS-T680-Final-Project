package api

import (
	"log"
	"net/http"
	"strconv"

	"drexel.edu/todo/db"
	"github.com/gin-gonic/gin"
)

// The api package creates and maintains a reference to the data handler
// this is a good design practice
type VoterAPI struct {
	db *db.VoterList
}

func New() (*VoterAPI, error) {
	dbHandler, err := db.NewVoterList()
	if err != nil {
		return nil, err
	}

	return &VoterAPI{db: dbHandler}, nil
}

//Below we implement the API functions.  Some of the framework
//things you will see include:
//   1) How to extract a parameter from the URL, for example
//	  the id parameter in /todo/:id
//   2) How to extract the body of a POST request
//   3) How to return JSON and a correctly formed HTTP status code
//	  for example, 200 for OK, 404 for not found, etc.  This is done
//	  using the c.JSON() function
//   4) How to return an error code and abort the request.  This is
//	  done using the c.AbortWithStatus() function

// implementation for GET /todo
// returns all todos
func (v *VoterAPI) GetAllVoterResources(c *gin.Context) {

	voterList, err := v.db.GetAllVoters()
	if err != nil {
		log.Println("Error Getting All Voters: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	//Note that the database returns a nil slice if there are no items
	//in the database.  We need to convert this to an empty slice
	//so that the JSON marshalling works correctly.  We want to return
	//an empty slice, not a nil slice. This will result in the json being []
	if voterList == nil {
		voterList = make([]db.Voter, 0)
	}

	c.JSON(http.StatusOK, voterList)
}

// implementation for GET /todo/:id
// returns a single todo
func (v *VoterAPI) GetSingleVoterResource(c *gin.Context) {

	//Note go is minimalistic, so we have to get the
	//id parameter using the Param() function, and then
	//convert it to an int64 using the strconv package
	idS := c.Param("id")
	id64, err := strconv.ParseInt(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	//Note that ParseInt always returns an int64, so we have to
	//convert it to an int before we can use it.
	voter, err := v.db.GetSingleVoterResource(uint(id64))
	if err != nil {
		log.Println("Item not found: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	//Git will automatically convert the struct to JSON
	//and set the content-type header to application/json
	c.JSON(http.StatusOK, voter)
}

func (v *VoterAPI) GetVoterHistory(c *gin.Context) {

	//Note go is minimalistic, so we have to get the
	//id parameter using the Param() function, and then
	//convert it to an int64 using the strconv package
	idS := c.Param("id")
	id64, err := strconv.ParseInt(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	//Note that ParseInt always returns an int64, so we have to
	//convert it to an int before we can use it.
	voter, err := v.db.GetVoterHistory(uint(id64))
	if err != nil {
		log.Println("Item not found: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	//Git will automatically convert the struct to JSON
	//and set the content-type header to application/json
	c.JSON(http.StatusOK, voter)
}

func (v *VoterAPI) GetVoterPollData(c *gin.Context) {

	//Note go is minimalistic, so we have to get the
	//id parameter using the Param() function, and then
	//convert it to an int64 using the strconv package
	idS := c.Param("id")
	idP := c.Param("pollid")
	id64_1, err_1 := strconv.ParseInt(idS, 10, 32)
	id64_2, err_2 := strconv.ParseInt(idP, 10, 32)
	if err_1 != nil {
		log.Println("Error converting voterid to int64: ", err_1)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err_2 != nil {
		log.Println("Error converting pollid to int64: ", err_2)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	//Note that ParseInt always returns an int64, so we have to
	//convert it to an int before we can use it.
	voter, err := v.db.GetVoterPollData(uint(id64_1), uint(id64_2))
	if err != nil {
		log.Println("Item not found: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	//Git will automatically convert the struct to JSON
	//and set the content-type header to application/json
	c.JSON(http.StatusOK, voter)
}

func (v *VoterAPI) AddVoterPollData(c *gin.Context) {

	//Note go is minimalistic, so we have to get the
	//id parameter using the Param() function, and then
	//convert it to an int64 using the strconv package
	idS := c.Param("id")
	idP := c.Param("pollid")
	id64_1, err_1 := strconv.ParseInt(idS, 10, 32)
	id64_2, err_2 := strconv.ParseInt(idP, 10, 32)
	if err_1 != nil {
		log.Println("Error converting voterid to int64: ", err_1)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err_2 != nil {
		log.Println("Error converting voterid to int64: ", err_2)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	//Note that ParseInt always returns an int64, so we have to
	//convert it to an int before we can use it.
	err2 := v.db.AddVoterPollData(uint(id64_1), uint(id64_2))
	if err2 != nil {
		log.Println("Item not found: ", err_1)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
}

func (v *VoterAPI) DeletePoll(c *gin.Context) {

	//Note go is minimalistic, so we have to get the
	//id parameter using the Param() function, and then
	//convert it to an int64 using the strconv package
	idS := c.Param("id")
	idP := c.Param("pollid")
	id64_1, err_1 := strconv.ParseInt(idS, 10, 32)
	id64_2, err_2 := strconv.ParseInt(idP, 10, 32)
	if err_1 != nil {
		log.Println("Error converting voterid to int64: ", err_1)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err_2 != nil {
		log.Println("Error converting voterid to int64: ", err_2)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	//Note that ParseInt always returns an int64, so we have to
	//convert it to an int before we can use it.
	err2 := v.db.DeletePoll(uint(id64_1), uint(id64_2))
	if err2 != nil {
		log.Println("Item not found: ", err_1)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
}

// implementation for POST /todo
// adds a new todo
func (v *VoterAPI) AddVoter(c *gin.Context) {
	var voter db.Voter

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

	if err := c.ShouldBindJSON(&voter); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := v.db.AddVoter(voter); err != nil {
		log.Println("Error adding item: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, voter)
}

// implementation for PUT /todo
// Web api standards use PUT for Updates
func (v *VoterAPI) UpdateVoter(c *gin.Context) {
	var voter db.Voter
	if err := c.ShouldBindJSON(&voter); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := v.db.UpdateVoter(voter); err != nil {
		log.Println("Error updating item: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, voter)
}

// implementation for DELETE /todo/:id
// deletes a todo
func (v *VoterAPI) DeleteVoter(c *gin.Context) {
	idS := c.Param("id")
	id64, _ := strconv.ParseInt(idS, 10, 32)

	if err := v.db.DeleteVoter(uint(id64)); err != nil {
		log.Println("Error deleting item: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

// implementation for DELETE /todo
// deletes all todos
func (v *VoterAPI) DeleteAllVoters(c *gin.Context) {

	if err := v.db.DeleteAll(); err != nil {
		log.Println("Error deleting all items: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

/*   SPECIAL HANDLERS FOR DEMONSTRATION - CRASH SIMULATION AND HEALTH CHECK */

// implementation for GET /crash
// This simulates a crash to show some of the benefits of the
// gin framework
func (v *VoterAPI) CrashSim(c *gin.Context) {
	//panic() is go's version of throwing an exception
	panic("Simulating an unexpected crash")
}

// implementation of GET /health. It is a good practice to build in a
// health check for your API.  Below the results are just hard coded
// but in a real API you can provide detailed information about the
// health of your API with a Health Check
func (v *VoterAPI) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK,
		gin.H{
			"status":             "ok",
			"version":            "1.0.0",
			"uptime":             100,
			"users_processed":    1000,
			"errors_encountered": 10,
		})
}
