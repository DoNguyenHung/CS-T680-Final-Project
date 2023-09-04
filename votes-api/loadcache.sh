#!/bin/bash
curl -d '{ "voteId": 1, "pollId": 1, "voterId": 1, "voteValue": 1 }' -H "Content-Type: application/json" -X POST http://localhost:1082/votes/1
curl -d '{ "voteId": 2, "pollId": 2, "voterId": 2, "voteValue": 2 }' -H "Content-Type: application/json" -X POST http://localhost:1082/votes/2
curl -d '{ "voteId": 3, "pollId": 3, "voterId": 3, "voteValue": 3 }' -H "Content-Type: application/json" -X POST http://localhost:1082/votes/3