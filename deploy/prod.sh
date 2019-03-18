#!/bin/sh
echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
docker tag kalaxia/api:latest kalaxia/api:"$TRAVIS_TAG"
docker push kalaxia/api:latest
docker push kalaxia/api:"$TRAVIS_TAG"