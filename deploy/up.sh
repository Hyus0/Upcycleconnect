#!/usr/bin/env sh
set -eu

docker compose -f deploy/docker-compose/docker-compose.yml up -d --build
