## Final Project

My name's Hung and this is my CS-T680 Final Project submission.

In order to run the project, you can go to the makefile and run every single command in order from top to bottom. This includes building the containers, running them using docker compose, loading the cache, and testing out the voter, poll methods. Finally, there are the vote methods to look for intra-API dependencies.

In other words, for example, run the make commands in the following order: 

```
      make build-poll-container
      make build-voter-container
      make build-votes-container
      make docker-compose
      make load-poll-cache
      make load-voter-cache
	  make load-votes-cache
      make get-all-votes
      make get-vote id=1
      make get-voter-by-vote id=1
      make get-poll-by-vote id=1

```

Please email me if you have any questions/problems running my code: hd386@cs.drexel.edu