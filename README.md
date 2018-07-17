Kalaxia Game API
===============

This repository is the Golang API for Kalaxia game.

It is used to develop, build and ship the game, but we recommend the [Docker Compose repository](https://github.com/Kalaxia/game-docker) for use purposes.

Requirements
------------

* Docker

Usage
------------

To build the Docker image with your new code compiled, use the following command:

```
docker-compose build
```

Compilation errors will be displayed during the build.

To use the created image with the Docker Compose repository, you must tag it:

```
docker tag kalaxiagameapi_api kalaxia/api
```

Build Docker Image with Travis
-----------------------------

As Docker doesn't work with some members of the team, we made a test server to run the environment.

To build and push the API image, it is necessary to fork this repository and then configure Travis to build it.

You will need to install Ruby and Gem on your computer.

Then, you need to copy the Travis configuration file:

```
cp .travis.dist.yml .travis.yml
```

Then you must set the different environment variables needed for the build:

```
travis env set DOCKER_IMAGE your-docker-hub-repository
travis encrypt DOCKER_USERNAME=your-docker-hub-username --add env.matrix
travis encrypt DOCKER_PASSWORD=your-docker-hub-password --add env.matrix
```

Travis will log in Docker Hub and push the built image on your repository.

Then you will be able to pull the image on the test server :) !

Database migrations
-------------------

To update the database schema, the game uses an [external migration package](https://github.com/mattes/migrate).

If you have new changes to apply, create a new file following the naming convention described here:

```
${version}_${model_type}.up.sql
${version}_${model_type}.down.sql
```

These two files are mandatory: the ```*.up.sql``` contains the SQL statements which apply your changes.

``*.down.sql`` on the other hand is meant to rollback these changes if anything goes wrong.

The way to validate that your migration works properly is to be able to do several up and rollbacks on the same file.

``${version}`` is a simple number next to the latest migration file.

``${model_type}`` is the model structure table you are working on.

For example, if I want to add a field for ``Player`` structure, I must create ``3_player.up.sql`` and ``3_player.down.sql`` and set the SQL statements inside.

A shortcut command has been implemented in the repository's ``Makefile`` to quickly update your database schema with the latest version. To use it, type the following command:

```
make migrate-latest
```

You can adapt the full-command in this file to do your stuff, for example rollback or use a specific version.
