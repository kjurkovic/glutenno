#!/bin/bash

# build services
make -C ../backend/auth
make -C ../backend/recipes
make -C ../backend/comments
make -C ../backend/notification
make -C ../frontend

# run containers
docker compose up
