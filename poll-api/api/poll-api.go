package api

import (
	"log"
	"net/http"
	"strconv"

	"drexel.edu/poll-api/db"
	"github.com/gin-gonic/gin"
)

// The api package creates and maintains a reference to the data handler
// this is a good design practice
type PollAPI struct {
	db *db.PollList
}

func New() (*PollAPI, error) {
	dbHandler, err := db.NewPollList()
	if err != nil {
		return nil, err
	}

	return &PollAPI{db: dbHandler}, nil
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
func (p *PollAPI) GetAllPollResources(c *gin.Context) {

	pollList, err := p.db.GetAllPolls()
	if err != nil {
		log.Println("Error Getting All Voters: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	//Note that the database returns a nil slice if there are no items
	//in the database.  We need to convert this to an empty slice
	//so that the JSON marshalling works correctly.  We want to return
	//an empty slice, not a nil slice. This will result in the json being []
	if pollList == nil {
		pollList = make([]db.Poll, 0)
	}

	c.JSON(http.StatusOK, pollList)
}

// implementation for GET /todo/:id
// returns a single todo
func (p *PollAPI) GetSinglePollResource(c *gin.Context) {

	//Note go is minimalistic, so we have to get the
	//id parameter using the Param() function, and then
	//convert it to an int64 using the strconv package
	idS := c.Param("id")
	if idS == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No poll ID provided"})
		return
	}
	id64, err := strconv.ParseInt(idS, 10, 32)
	if err != nil {
		log.Println("Error converting id to int64: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	//Note that ParseInt always returns an int64, so we have to
	//convert it to an int before we can use it.
	poll, err := p.db.GetSinglePollResource(uint(id64))
	if err != nil {
		log.Println("Item not found: ", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	//Git will automatically convert the struct to JSON
	//and set the content-type header to application/json
	c.JSON(http.StatusOK, poll)
}

// implementation for POST /todo
// adds a new todo
func (p *PollAPI) AddPoll(c *gin.Context) {
	var poll db.Poll

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

	if err := c.ShouldBindJSON(&poll); err != nil {
		log.Println("Error binding JSON: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := p.db.AddPoll(poll); err != nil {
		log.Println("Error adding item: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, poll)
}

// implementation for DELETE /todo/:id
// deletes a todo
func (p *PollAPI) DeletePoll(c *gin.Context) {
	idS := c.Param("id")
	if idS == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No poll ID provided"})
		return
	}
	id64, _ := strconv.ParseInt(idS, 10, 32)

	if err := p.db.DeletePoll(uint(id64)); err != nil {
		log.Println("Error deleting item: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

// implementation for DELETE /todo
// deletes all todos
func (p *PollAPI) DeleteAllPolls(c *gin.Context) {

	if err := p.db.DeleteAll(); err != nil {
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
func (p *PollAPI) CrashSim(c *gin.Context) {
	//panic() is go's version of throwing an exception
	panic("Simulating an unexpected crash")
}

// implementation of GET /health. It is a good practice to build in a
// health check for your API.  Below the results are just hard coded
// but in a real API you can provide detailed information about the
// health of your API with a Health Check
func (p *PollAPI) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK,
		gin.H{
			"status":             "ok",
			"version":            "1.0.0",
			"uptime":             100,
			"users_processed":    1000,
			"errors_encountered": 10,
		})
}
