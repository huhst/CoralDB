#!/bin/bash
docker rmi -f yan259128/algdb
docker build --no-cache  -t yan259128/algdb .
