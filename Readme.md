# GoLang setup
* Basic golang setup using keycloak and postgres for authentication
* reusable for quickly setting up a go backend service with keycloak Authentication

## How to run
### Development

Copy [dev.example.env](./dev.example.env) to .env

#### start all with docker

:warning: when running the app with docker all env variables are set from the [docker-compose.yml](./docker-compose.yml) file
```
docker compose up
```

#### start app locally
start db and keycloak using docker
```
docker compose up -d keycloak db
```

run app native for development
```
go run .
```


### Production
* deploy keycloak somewhere in production
* update the .env file/docker-compose with the keycloak url and clientID
* set env var to production
