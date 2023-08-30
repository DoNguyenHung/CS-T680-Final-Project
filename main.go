package main

import (
	"flag"
	"fmt"
	"os"

	"drexel.edu/todo/api"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Global variables to hold the command line flags to drive the todo CLI
// application
var (
	hostFlag string
	portFlag uint
)

// processCmdLineFlags parses the command line flags for our CLI
//
// TODO: This function uses the flag package to parse the command line
//		 flags.  The flag package is not very flexible and can lead to
//		 some confusing code.

//			 REQUIRED:     Study the code below, and make sure you understand
//						   how it works.  Go online and readup on how the
//						   flag package works.  Then, write a nice comment
//				  		   block to document this function that highights that
//						   you understand how it works.
//
//			 EXTRA CREDIT: The best CLI and command line processor for
//						   go is called Cobra.  Refactor this function to
//						   use it.  See github.com/spf13/cobra for information
//						   on how to use it.
//
//	 YOUR ANSWER: <GOES HERE>
func processCmdLineFlags() {

	//Note some networking lingo, some frameworks start the server on localhost
	//this is a local-only interface and is fine for testing but its not accessible
	//from other machines.  To make the server accessible from other machines, we
	//need to listen on an interface, that could be an IP address, but modern
	//cloud servers may have multiple network interfaces for scale.  With TCP/IP
	//the address 0.0.0.0 instructs the network stack to listen on all interfaces
	//We set this up as a flag so that we can overwrite it on the command line if
	//needed
	flag.StringVar(&hostFlag, "h", "0.0.0.0", "Listen on all interfaces")
	flag.UintVar(&portFlag, "p", 1080, "Default Port")

	flag.Parse()
}

// main is the entry point for our todo API application.  It processes
// the command line flags and then uses the db package to perform the
// requested operation
func main() {
	processCmdLineFlags()
	r := gin.Default()
	r.Use(cors.Default())

	apiHandler, err := api.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	r.GET("/voters", apiHandler.GetAllVoterResources)

	r.GET("/voters/:id", apiHandler.GetSingleVoterResource)
	// Create a voters resource with id = :id, initialize the polls slice to an
	// empty slice
	r.POST("/voters/:id", apiHandler.AddVoter)

	r.GET("/voters/:id/polls", apiHandler.GetVoterHistory)

	r.GET("/voters/:id/polls/:pollid", apiHandler.GetVoterPollData)
	// Look up the voter with id = :id, then add the poll with pollid = :pollid to
	// the internal poll slice
	// POST /voters/22/polls/3
	// Does voter 22 exist, if not return 404 error; if voter 22 exists, add
	// pollid 3 to the internal poll slice
	// ***** Does voter 22 exist, if not, create voter 22 (WHICH ASSUMES ALL THE
	// VOTER INFO IS IN THE PAYLOAD), then voter 22, then
	// add pollid 3 to the NEW voter 22 resource. If not, follow above
	r.POST("/voters/:id/polls/:pollid", apiHandler.AddVoterPollData)

	r.GET("/voters/health", apiHandler.HealthCheck)

	// Extra Credit
	r.DELETE("/voters", apiHandler.DeleteAllVoters)

	r.DELETE("/voters/:id", apiHandler.DeleteVoter)

	r.DELETE("/voters/:id/polls/:pollid", apiHandler.DeletePoll)

	r.PUT("/voters", apiHandler.UpdateVoter)

	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	r.Run(serverPath)
}
