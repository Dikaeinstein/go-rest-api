#!/bin/bash
docker login --username=_ --password=$HEROKU_API_KEY registry.heroku.com
docker push registry.heroku.com/dikaeinstein-go-rest-api/web

echo $WEB_DOCKER_IMAGE_ID

# Release container via API
curl -n -X PATCH https://api.heroku.com/apps/dikaeinstein-go-rest-api/formation \
  -d '{
  "updates": [
    {
      "type": "web",
      "docker_image": "$WEB_DOCKER_IMAGE_ID"
    }
  ]
}' \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $HEROKU_API_KEY" \
  -H "Accept: application/vnd.heroku+json; version=3.docker-releases"
