#!/bin/bash
docker login --username=_ --password=$HEROKU_API_KEY registry.heroku.com
docker tag dikaeinstein-go-rest-api:$DOCKER_IMAGE_TAG registry.heroku.com/dikaeinstein-go-rest-api/web
docker push registry.heroku.com/dikaeinstein-go-rest-api/web

# Release container via heroku CLI
heroku container:release web -a dikaeinstein-go-rest-api
