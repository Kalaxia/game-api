
# Table of content
<!-- TOC depthFrom:1 depthTo:6 withLinks:1 updateOnSave:1 orderedList:0 -->

- [Table of content](#table-of-content)
- [Kalaxia Game API](#kalaxia-game-api)
	- [Introduction](#introduction)
	- [Requirements](#requirements)
	- [Installation](#installation)
		- [Docker installation](#docker-installation)
			- [GNU/Linux](#gnulinux)
			- [Windows](#windows)
				- [Windows server 2016](#windows-server-2016)
			- [Mac](#mac)
		- [Repository setup](#repository-setup)
		- [Container setup](#container-setup)
			- [Building the container](#building-the-container)
			- [Launching the container](#launching-the-container)
		- [Database setup](#database-setup)
		- [Setup the game](#setup-the-game)
	- [Administration](#administration)
		- [Logs](#logs)
		- [Container](#container)
			- [Updating the container](#updating-the-container)
				- [Method 1](#method-1)
				- [Method 2](#method-2)
		- [Database](#database)
			- [Database migrations](#database-migrations)
			- [Advanced database migrations](#advanced-database-migrations)
	- [Troubleshooting](#troubleshooting)
		- [Error given by `docker-compose up (-d)`](#error-given-by-docker-compose-up-d)
		- [Web page](#web-page)
		- [Database migration](#database-migration)
			- [Step 0 (Optional)](#step-0-optional)
			- [Step 1](#step-1)
				- [Error cause by an up migration](#error-cause-by-an-up-migration)
				- [Error cause by a down migration](#error-cause-by-a-down-migration)
			- [Step 2](#step-2)
			- [Step 3](#step-3)
			- [Step 4](#step-4)
			- [Step 5](#step-5)
	- [The code](#the-code)
		- [Containers](#containers)
		- [API code](#api-code)
			- [Routes](#routes)
			- [Controller (or shipController)](#controller-or-shipcontroller)
			- [Manager](#manager)
			- [Model](#model)
			- [Resources](#resources)
		- [Library used](#library-used)

<!-- /TOC -->
---------


# Kalaxia Game API

This repository is the Golang API for Kalaxia game.

It is used to develop, build and ship the game, but we recommend the [Docker Compose repository](https://github.com/Kalaxia/game-docker) for use purposes.

## Introduction

What is kalaxia ?
Kalaxia is an old school multilayer browser strategy game on real time with semi persistent session.
This is an open source game. You can find more information on the project on our website [kalaxia.org](https://kalaxia.org), or on our [discord](https://discordapp.com/invite/bSQ3WV).



## Requirements

* Docker >= 18.03.0
* docker-compose >= 1.21.2



## Installation

### Docker installation

#### GNU/Linux

You can use your favorite packet manager to install the package `docker`

Notice that docker-compose is a package in your packet manager but it might be outdated. The authors advice following the [official documentation](https://docs.docker.com/compose/install/)

Normally your user should have permission to use docker. If it is not the case enter
```Bash
sudo groupadd docker
```
and then add the user to this group using
```Bash
sudo usermod -aG docker <user>
```
 You need to log off and log back in before you will be able to use the `docker` command. If you have any issue please refer to the [official documentation](https://docs.docker.com/install/linux/linux-postinstall/)

#### Windows
 Please refer to the [official documentation](https://docs.docker.com/docker-for-windows/install/). Note that docker-compose is already included in this build **excepte** for Windows server 2016.

##### Windows server 2016
If you are using Windows server 2016 please read the [official documentation for the installation of docker-compose](https://docs.docker.com/compose/install/#install-compose).

####  Mac

Please refer to the [official documentation](https://docs.docker.com/docker-for-mac/install/). Note that docker-compose is already included in this build.


### Repository setup
The first step is to clone the repository using
```Bash
git clone git@github.com:Kalaxia/game-api.git
```
Once the repository is cloned you will need to clone the front in 'volumes/app' using
```Bash
git clone git@github.com:Kalaxia/game-front.git volumes/app
```
inside  'volumes'.

Navigate back to 'game-api'. Copy the files '.dist.env' to '.env' and 'kalaxia.dist.env' to 'kalaxia.env'. You can do that by using `cp .dist.env .env` and `cp kalaxia.dist.env kalaxia.env`.
Now open the file `.env` using a text editor ( like vim, emacs, nano, gedit, ...) and change the  `NGINX_PORT` and `NGINX_HTTPS_PORT` to your liking. `NGINX_PORT`is the port nginx will listen to.
Optionally you can change the other values both of the new files.  But the authors recommend keeping as it is except for `NGINX_PORT` et `NGINX_HTTPS_PORT`. And at this time the authors do not provide information in order to work properly if these files are changed.

### Container setup

**Note that the following step could lead to data lost inside container with the same name as defined in `docker-compose.yml`**

#### Building the container
To build the Docker image with your new code compiled, use the following command:

```Bash
docker-compose build
```

Compilation errors will be displayed during the build.

To use the created image with the Docker Compose repository, you must tag it:

```Bash
docker tag kalaxiagameapi_api kalaxia/api
```

#### Launching the container

To launch the container use
```Bash
docker-compose up -d
```
where the `-d` flag means to detach the container and run it in background.

### Database setup

Now that the container is running you will need to create all the table in the database. In order create the database structure use
```Bash
docker exec -it kalaxia_api make migrate-latest
```
more informations are provided in the 'Database migrations' section below.

### Setup the game

To setup the game you will need to setup the [portal](https://github.com/Kalaxia/portal).

## Administration

### Logs

You can use
```Bash
docker logs -f kalaxia_api
```
to display the logs of the container.

### Container

 - Starting the container : `docker-compose start`
 - stopping the container : `docker-compose stop`
 - restating the container : `docker-compose restart`
 - recreating the container : `docker-compose up -d --force-recreate`
 - stopping and deleting the container :  `docker-compose down` **/!\  All data inside the container will be lost including the database**
 - (re)build and launch the container `docker-compose up -d --build` **/!\ All data inside the container will be lost including the database**

#### Updating the container

##### Method 1

Update your files locally (by instance using `git pull`).
Then you type the command
```Bash
docker-compose up -d --build api
```
to rebuild the container and launch it. This will only rebuild the api so your database should be safe.   
Note that the previous container will be still on your machine but will be stopped. The delete the old image please refer to the docker documentation.

##### Method 2

TODO


### Database

You can connect to the database using
```Bash
docker exec -it kalaxia_postgresql psql -U kalaxia kalaxia_game
```
in this mode you can directly type SQL queries. To quit type `\q` then press enter.

#### Database migrations

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

```Bash
docker exec -it kalaxia_api make migrate-latest
```

You can adapt the full-command in this file to do your stuff, for example rollback or use a specific version.

#### Advanced database migrations

TODO

## Troubleshooting


### Error given by `docker-compose up (-d)`

> ```
> WARNING: The NGINX_PORT variable is not set. Defaulting to a blank string.
> ```

Did you create both `.env` and `kalaxia.env` as mentioned in the Installation process? Check that these files defines the same variable as `.dist.env` and `kalaxia.dist.env`.

----

>`
>ERROR: The Compose file './docker-compose.yml' is invalid because: networks.game value Additional properties are not allowed ('name' was unexpected) services.nginx.ports contains an invalid type, it should be a number, or an object
>`

You might have an oudated version of docker-compose. Check the version using
```Bash
docker-compose -v
```
The version must be greater or equal to 1.21.2. If it where outdated uninstall docker-compose and reinstall a newer version. ( See section 'Docker installation' in this document).



### Web page

> When I connect to the main web page I only see the nginx page
>
Check your `.env` and `kalaxia.env` files. By default you need to use the user kalaxia as provided in the `.dist.env` and `kalaxia.dist.env` files.


### Database migration

>`error: Dirty database version <v>. Fix and force version.`

Your database did not migrate correctly and it is in a dirty state the following steps present how to resolve and go back to a clean state.

Let `<v>` the version given by the error message.  And `<v+1>` the version number plus one.

#### Step 0 (Optional)
Debug your migration files.
Once they are debug you still need to folow the next step to go back to a clean version.


#### Step 1
##### Error cause by an up migration

Read the file in 'build/migrations' which is name `<v>_*.up.sql`  and copy the command somewhere easily accessible for the next step.


##### Error cause by a down migration

Read the file in 'build/migrations' which is name `<v+1>_*.down.sql`  and copy the command somewhere easily accesible for the next setp.
#### Step 2
Run
```Bash
docker exec -it kalaxia_postgresql psql -U kalaxia kalaxia_game
```
and there type one by one the SQL request you copy before, correcting where there is problem. Error like this table does not exist or this column does not exist can be safely ignored.

If you give up on solving all the problem in the SQL command go to ste5.



#### Step 3

Now update the table  `schema_migrations` in the postgresql environment using
```SQL
UPDATE schema_migrations SET dirty=f WHERE dirty=t;
```

Or you can run
```Bash
docker exec -it kalaxia_api migrate -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:5432/${POSTGRES_DB}?sslmode=disable -source file://build/migrations force <v>
```
This will tell the program that it is in the version `<v>` but this not run the SQL command. This is use to remove the dirty version state and let you continue your migration.

#### Step 4

After that you can resume your migration.
If you the error persist continue on step 5. Otherwise you can stop here.

#### Step 5

This step is only to apply if the database is very corrupted.
**The following step will erase all data in your database.**
You will need to erase the database by running
```Bash
docker exec -it kalaxia_api migrate -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:5432/${POSTGRES_DB}?sslmode=disable -source file://build/migrations drop
```
Your database is now empty. You can resume your migration.

## API documentation

The API documentation is written in RAML. You can use [tools](https://raml.org/projects) converter to read the documentation.

 **What is up with these `#-nolink.raml` ?**
 
 Due to how we retrieve data in the database we send partially linked object. To reflect this we have different data type files with the different level of depth of the links.
 As a rule of thumb, a file with `nolink` does not link to any other object in the DB. A `lowlink` link to the `nolink` object. And a `highlink` links to the `lowlink` object.

## The code

This section present basically the architecture of the code to simplify your potential contribution.

### Containers

First if you look at the `./docker-compose.yml` you see that we use three docker container, one for the this code, one for the postgreSQL database and one for nginx.

If your contribution include some data that need to be stored in the database read the section [Database](#database).

### API code

#### Routes

In order for the application to call some function associated to an http request you need to âdd a new route inside `./route/route.go`.   
By instance the route to get the planets a player control
```Go
Route{
    "Get Player Planets", // name of the route
    "GET", // Method (GET,PUT,POST,DELETE,..)
    "/api/players/{id}/planets", // URL pattern
    controller.GetPlayerPlanets, // the function called when the request is received
    true, // if the user need to be authenticated
}
```
Note that in the URL pattern you can pass variable like in this case using `{id}`. This is only meant to pass few and simple argument. For more complex argument prefer using json in the request.

The function called should be in the package 'controller' (or 'shipController') located under the folder `./controller/`

When a new route is added it will be greatly appreciated if you could complete the documentation under `./doc/`.

#### Controller (or shipController)

This is where the check if the player is allowed to see the information and where the request is parsed.   
Function of this package should not do request to the database directly but use functions in the package 'manager' (or 'shipManager')

#### Manager

This package is located under `./manager/`. This is the place where the request to the data base is performed and where the data are modified.

Moreover there may be some function called `init`. These function does not take any argument and are performed when the application start. Note that previous schedule are not retained when the application stops, thus this is the place to reschedule them.

#### Model

The package model is located under `./model/`. It regroup all the structure definition of the object used.

As an example look at the following structure.
```Go
SystemOrbit struct {
	TableName struct{} `json:"-" sql:"map__system_orbits"`

	Id uint16 `json:"id"`
	Radius uint16 `json:"radius"`
	SystemId uint16 `json:"-"`
	System *System `json:"system"`
}
```
Notice the json attribute and SQL attribute. If your object need to be stored in the database do not forget to add the name of the table in the object.

#### Resources

Under `./resources/` you can add some json file containing variable changing different aspect of the game, like the price of a building.

### Library used

 - [pg](https://github.com/go-pg/pg): PostgreSQL ORM for Golang
 - TODO
