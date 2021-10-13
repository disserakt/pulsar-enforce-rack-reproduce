#!/bin/bash

set -uxe

docker-compose --compatibility up -d --build
sleep 180
docker-compose kill broker2 bookie7 bookie8 bookie9
docker-compose logs -f client