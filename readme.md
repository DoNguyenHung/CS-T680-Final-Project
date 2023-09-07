## CS-T680 Final Project

My name's Hung and this is my CS-T680 Final Project submission.

In order to run the project, you can go to the makefile and run every single command in order from top to bottom. This includes building the containers, running them using docker compose, loading the cache, and testing out any of the voter and poll methods. Finally, there are the vote methods to look for intra-API dependencies.

You can view cache as you run by using this link: http://localhost:8001/redis-stack/browser

In other words, for example, run the make commands in the following order: 

```
      make build-poll-container
      make build-voter-container
      make build-votes-container
      make docker-compose
      make load-poll-cache
      make load-voter-cache
      make load-votes-cache

      make get-poll-by-id id=1
      make get-all-polls
      make delete-all-polls
      make delete-poll-by-id id=1

      make get-voter-by-id id=1
      make get-all-voters
      make get-voter-history id=1
      make get-voter-poll id=1 pollid=59231
      make get-health
      make add-voter-poll id=1 pollid=45678
      make delete-all-voters
      make delete-voter-by-id id=1

      make get-all-votes
      make get-vote id=1
      make get-voter-by-vote id=1
      make get-poll-by-vote id=1
```

Please email me if you have any questions/problems running my code: hd386@cs.drexel.edu