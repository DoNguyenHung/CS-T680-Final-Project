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

# Votes methods
.PHONY: get-all-votes
get-all-votes:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1082/votes

.PHONY: get-vote
get-vote:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1082/votes/$(id)

.PHONY: get-voter-by-vote
get-voter-by-vote:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1082/votes/$(id)/voters/

# make get-by-id id=2
.PHONY: get-poll-by-vote
get-poll-by-vote:
	curl -w "HTTP Status: %{http_code}\n" -H "Content-Type: application/json" -X GET http://localhost:1082/votes/$(id)/polls/