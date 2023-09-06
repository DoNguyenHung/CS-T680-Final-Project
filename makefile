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