SHELL := /bin/bash

.PHONY: help
help:
	@echo "Usage make <TARGET>"
	@echo ""
	@echo "  Targets:"
	@echo "	   build-poll-container			Build Poll-API"
	@echo "	   build-voter-container		Build Voter-API"
	@echo "	   build-votes-container		Build Votes-API"
	@echo "	   docker-compose			Run containers"
	@echo "	   load-poll-cache			Load cache for poll"
	@echo "	   load-voter-cache			Load cache for voters"
	@echo "	   load-votes-cache			Load cache for votes"
	@echo "	   get-all-votes			Retrieve all votes"
	@echo "	   get-vote				Get vote based on vote ID"
	@echo "	   get-voter-by-vote			Get voter based on vote ID"
	@echo "	   get-poll-by-vote			Get poll based on vote ID"

# Build Poll-API
.PHONY: build-poll-container
build-poll-container:
	cd poll-api/ && ./build-basic-docker.sh

# Build Voter-API
.PHONY: build-voter-container
build-voter-container:
	cd voter-api/ && ./build-basic-docker.sh

# Build Votes-API 
.PHONY: build-votes-container
build-votes-container:
	cd votes-api && ./build-basic-docker.sh

# Run containers
.PHONY: docker-compose
docker-compose:
	docker compose up

# Load cache for poll, voter, and votes
.PHONY: load-poll-cache
load-poll-cache:
	cd poll-api/ && ./loadcache.sh

.PHONY: load-voter-cache
load-voter-cache:
	./voter-api/loadcache.sh

.PHONY: load-votes-cache
load-votes-cache:
	./votes-api/loadcache.sh

# Poll methods
.PHONY: get-poll-by-id
get-poll-by-id:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1080/polls/$(id)

.PHONY: get-all-polls
get-all-polls:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1080/polls

.PHONY: delete-all-polls
delete-all-polls:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X DELETE http://localhost:1080/polls

.PHONY: delete-poll-by-id
delete-poll-by-id:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X DELETE http://localhost:1080/polls/$(id) 

# Voter methods
.PHONY: get-voter-by-id
get-voter-by-id:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1081/voters/$(id)

.PHONY: get-all-voters
get-all-voters:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1081/voters 

.PHONY: get-voter-history
get-voter-history:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1081/voters/$(id)/polls

.PHONY: get-voter-poll
get-voter-poll:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1081/voters/$(id)/polls/$(pollid)

.PHONY: get-health
get-health:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1081/voters/health

.PHONY: add-voter-poll
add-voter-poll:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X POST http://localhost:1081/voters/$(id)/polls/$(pollid)

.PHONY: delete-all-voters
delete-all-voters:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X DELETE http://localhost:1081/voters

.PHONY: delete-voter-by-id
delete-voter-by-id:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X DELETE http://localhost:1081/voters/$(id) 

# Votes methods
.PHONY: get-all-votes
get-all-votes:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1082/votes

.PHONY: get-vote-by-id
get-vote-by-id:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1082/votes/$(id)

.PHONY: get-voter-by-vote
get-voter-by-vote:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1082/votes/$(id)/voters/

.PHONY: get-poll-by-vote
get-poll-by-vote:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1082/votes/$(id)/polls/

.PHONY: delete-all-votes
delete-all-votes:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X DELETE http://localhost:1082/votes

# Delete all
.PHONY: delete-all-resources
delete-all-resources:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X DELETE http://localhost:1080/polls
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X DELETE http://localhost:1081/voters
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X DELETE http://localhost:1082/votes
