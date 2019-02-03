#!/bin/bash
docker login --username=_ --password=$HEROKU_API_KEY registry.heroku.com
docker push registry.heroku.com/dikaeinstein-go-rest-api/web

# Release container via heroku CLI
heroku container:release web -a dikaeinstein-go-rest-api
