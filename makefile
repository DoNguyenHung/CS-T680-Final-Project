SHELL := /bin/bash

.PHONY: help
help:
	@echo "Usage make <TARGET>"
	@echo ""
	@echo "  Targets:"
	@echo "	   build				Build the voters executable"
	@echo "	   run					Run the voters program from code"
	@echo "	   run-bin				Run the voters executable"
	@echo "	   load-db				Add sample data via curl"
	@echo "	   get-by-id			Get a voters by id pass id=<id> on command line"
	@echo "	   get-all				Get all voterss"
	@echo "	   update-2				Update record 2, pass a new title in using title=<title> on command line"
	@echo "	   delete-all			Delete all voterss"
	@echo "	   delete-by-id			Delete a voters by id pass id=<id> on command line"
	@echo "	   get-v2				Get all voterss by done status pass done=<true|false> on command line"
	@echo "	   get-v2-all			Get all voterss using version 2"
	@echo "	   build-amd64-linux	Build amd64/Linux executable"
	@echo "	   build-arm64-linux	Build arm64/Linux executable"

.PHONY: start-redis
start-redis:
	./poll-api/start-redis.sh

# Poll-API
.PHONY: build-poll-container
build-poll-container:
	cd poll-api/ && ./build-basic-docker.sh

.PHONY: run-poll-container
run-poll-container:
	./poll-api/run-basic-docker.sh

.PHONY: load-poll-cache
load-poll-cache:
	cd poll-api/ && ./loadcache.sh

# Voter-API
.PHONY: build-voter-container
build-voter-container:
	cd voter-api/ && ./build-basic-docker.sh

.PHONY: run-voter-container
run-voter-container:
	./voter-api/run-basic-docker.sh

.PHONY: load-voter-cache
load-voter-cache:
	./voter-api/loadcache.sh

# Votes-API 
.PHONY: build-votes-container
build-votes-container:
	cd votes-api && ./build-basic-docker.sh

.PHONY: run-votes-container
run-votes-container:
	./votes-api/run-basic-docker.sh

.PHONY: load-votes-cache
load-votes-cache:
	./votes-api/loadcache.sh

.PHONY: docker-compose
docker-compose:
	docker compose up

#
.PHONY: get-all-votes
get-all-votes:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1080/votes

# make get-by-id id=2
.PHONY: get-vote-by-voter
get-vote-by-voter:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1080/voters/$(id)

# # make get-by-id id=2
# .PHONY: get-by-id
# get-by-id:
# 	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1080/voters/$(id)

# .PHONY: get-all
# get-all:
# 	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1080/voters 

# # make get-by-id id=2
# .PHONY: get-voter-history
# get-voter-history:
# 	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1080/voters/$(id)/polls

# .PHONY: get-voter-poll
# get-voter-poll:
# 	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1080/voters/$(id)/polls/$(pollid)

# .PHONY: get-health
# get-health:
# 	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1080/voters/health

# .PHONY: add-voter-poll
# add-voter-poll:
# 	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X POST http://localhost:1080/voters/$(id)/polls/$(pollid)

# # Extra credit
# .PHONY: delete-all
# delete-all:
# 	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X DELETE http://localhost:1080/voters 

# .PHONY: delete-by-id
# delete-by-id:
# 	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X DELETE http://localhost:1080/voters/$(id) 

# .PHONY: delete-by-pollid
# delete-by-pollid:
# 	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X DELETE http://localhost:1080/voters/$(id)/polls/$(pollid)

# .PHONY: update-1
# update-1:
# 	curl -d '{ "VoterId": 1, "FirstName": "$(fn)", "LastName": "$(ln)", "VoteHistory": [{"PollID": 59231, "VoteDate": "2021-08-15T14:30:45.00Z"}] }' -H "Content-Type: application/json" -X PUT http://localhost:1080/voters

# .PHONY: update-2
# update-2:
# 	curl -d '{ "VoterId": 2, "FirstName": "$(fn)", "LastName": "$(ln)", "VoteHistory": [{"PollID": 12345, "VoteDate": "2021-08-16T14:30:45.00Z"}] }' -H "Content-Type: application/json" -X PUT http://localhost:1080/voters

# .PHONY: update-3
# update-3:
# 	curl -d '{ "VoterId": 3, "FirstName": "$(fn)", "LastName": "$(ln)", "VoteHistory": [{"PollID": 54321, "VoteDate": "2021-08-17T14:30:45.00Z"}] }' -H "Content-Type: application/json" -X PUT http://localhost:1080/voters