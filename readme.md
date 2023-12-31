## CS-T680 Final Project

My name's Hung and this is my CS-T680 Final Project submission.

In order to run the project, you can go to the makefile and run every single command in order from top to bottom. This includes building the containers, running them using docker compose, loading the cache, and testing out any of the voter and poll methods. Finally, there are the vote methods to look for intra-API dependencies. 

I used Git bash to run my make commands so please contact me if you aren't able to run via Unix. Additionally, I don't have any scripts that are used to point out specific errors, but, for example, if you want to test a duplicate voter, you can simply use the delete-voter-by-id command to get rid of any voter and then rerun the load-voter-cache method. You'll see that the voters that aren't deleted will have errors showing duplication. 

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

      make delete-all-resources
```

One final note is that I will be away on Friday for my friend's bachelor's party, but I will keep a close look at my email over the weekend. Thank you for your patience.

Please email me if you have any questions/problems running my code: hd386@cs.drexel.edu
