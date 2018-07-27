Kalaxia Game API
===============

This repository is the Golang API for Kalaxia game.

It is used to develop, build and ship the game, but we recommend the [Docker Compose repository](https://github.com/Kalaxia/game-docker) for use purposes.

Requirements
------------

* Docker >= 18.03.0
* docker-compose >= 1.21.2


Installation
------------

### Docker installation

#### GNU/Linux

You can use your favorite packet manager to install the package `docker`

Notice that docker-compose is a package in your packet manager but it might be outdated. The authors advice following the [official documentation](https://docs.docker.com/compose/install/)

Then in order to use docker without being rooted create a new user group called docker using  
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
git clone https://github.com/Kalaxia/game-api.git
``` 
or
```Bash
git clone git@github.com:Kalaxia/game-api.git
```
Once the repository is cloned create a folder 'volumes' inside the 'game-api'  folder. You can use `mkdir volumes`.

Then clone the game-front inside 'volumes/app' using 
```Bash
git clone https://github.com/Kalaxia/game-front.git app
```
inside  'volumes'.

Navigate back to 'game-api'. Copy the files '.dist.env' to '.env' and 'kalaxia.dist.env' to 'kalaxia.env'. You can do that by using `cp .dist.env .env` and `cp kalaxia.dist.env kalaxia.env`. Optionally you can change both of the new files. But the authors recommend keeping as it is. And at this time the authors do not provide information in order to work properly if these files are changed.

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
where the `-d` flag means to detatch the container and run it in background. 

### Database setup

Now that the container is running you will need to create all the table in the database. In order create the database structure use
```Bash
docker exec -it kalaxia_api make migrate-latest
```
more informations are provided in the 'Database migrations' section below.

### Setup the game

In order to setup the game the API need to creat data in the database. In order to inialize these data you will need to do the folowing request
```JSON
url : "<address of your server>:<port>/api/servers",
methode : "POST",
Body : {
		"name": "<name>",
		"type" : "multiplayer",
		"signature":"<signature>",
		"factions": [
					 {"name":"<facion name>",
					 "description":"<description>",
					 "color":"<color>",
					 "banner":"<banner>"}
					 ]
		"map_size":100
},
Header : {
	"Authorization" : "Bearer <your Bearer token>"
}
```
 by default the port is 8004. You can also add multiple faction following the previous example.
 
 The proper authorization header for user kalaxia is required.
 TODO : authorization

Administration
------------

### Logs

You can use 
```Bash
docker logs -f kalaxia_api
```
to dispay the logs of the container.

### Container

 - Starting the container : `docker-compose start`
 - stopping the container : `docker-compose stop`
 - restating the container : `docker-compose restart`
 - recreating the container : `docker-compose up -d --force-recreate`
 - stopping and deleting the container :  `docker-compose down` **/!\  All data inside the container will be lost including the database**
 - (re)build and launch the container `docker-compose up -d --build` **/!\ All data inside the container will be lost including the database**

#### Updating the container 
TODO

### Database

You can connect to the database using 
```Bash
docker exec -it kalaxia_postgresql psql -U kalaxia kalaxia_game
```
in this mode you can directly type SQL commande. To quit type `\q` then press enter.

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

#### Advance database migrations

TODO

Troubleshooting
---------------------

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

If you give up on solving all the problem in the SQL command go to step 5.

Once this is finish type `\q` and enter to exit.

#### Step 3


Run
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
