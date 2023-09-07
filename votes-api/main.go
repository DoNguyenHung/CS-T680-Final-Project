package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"drexel.edu/votes-api/api"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Global variables to hold the command line flags to drive the todo CLI
// application
var (
	hostFlag    string
	portFlag    uint
	cacheURL    string
	voterAPIURL string
	pollAPIURL  string
)

func processCmdLineFlags() {
	flag.StringVar(&hostFlag, "h", "0.0.0.0", "Listen on all interfaces")
	flag.StringVar(&voterAPIURL, "voterapi", "http://localhost:1080", "Default endpoint for Voter API")
	flag.StringVar(&pollAPIURL, "pollapi", "http://localhost:1080", "Default endpoint for Poll API")
	flag.StringVar(&cacheURL, "c", "0.0.0.0:6379", "Default cache location")
	flag.UintVar(&portFlag, "p", 1080, "Default Port")

	flag.Parse()
}

func envVarOrDefault(envVar string, defaultVal string) string {
	envVal := os.Getenv(envVar)
	if envVal != "" {
		return envVal
	}
	return defaultVal
}

func setupParms() {
	//first process any command line flags
	processCmdLineFlags()

	//now process any environment variables
	cacheURL = envVarOrDefault("CACHE_URL", cacheURL)
	voterAPIURL = envVarOrDefault("VOTER_API_URL", voterAPIURL)
	pollAPIURL = envVarOrDefault("POLL_API_URL", pollAPIURL)
	hostFlag = envVarOrDefault("RLAPI_HOST", hostFlag)

	pfNew, err := strconv.Atoi(envVarOrDefault("RLAPI_PORT", fmt.Sprintf("%d", portFlag)))
	//only update the port if we were able to convert the env var to an int, else
	//we will use the default we got from the command line, or command line defaults
	if err == nil {
		portFlag = uint(pfNew)
	}

}

// main is the entry point for our todo API application.  It processes
// the command line flags and then uses the db package to perform the
// requested operation
func main() {
	//this will allow the user to override key parameters and also setup defaults
	setupParms()
	log.Println("Init/cacheURL: " + cacheURL)
	log.Println("Init/voterAPIURL: " + voterAPIURL)
	log.Println("Init/pollAPIURL: " + pollAPIURL)
	log.Println("Init/hostFlag: " + hostFlag)
	log.Printf("Init/portFlag: %d", portFlag)

	apiHandler, err := api.NewVoteAPI(cacheURL, voterAPIURL, pollAPIURL)

	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.Use(cors.Default())

	r.POST("/votes/:id", apiHandler.AddVote)
	r.GET("/votes", apiHandler.GetAllVotes)
	r.GET("/votes/:id", apiHandler.GetVote)
	r.GET("/votes/:id/voters/", apiHandler.GetVoterByVote)
	r.GET("/votes/:id/polls/", apiHandler.GetPollByVote)
	r.DELETE("/votes", apiHandler.DeleteAllVotes)

	//For now we will just support gets
	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	r.Run(serverPath)
}
