#!/bin/bash
docker run -d --rm --name voter-redis -p 6379:6379 -p 8001:8001 redis/redis-stack:latest