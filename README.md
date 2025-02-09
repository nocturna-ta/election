# Election API

## Contributing
Please refer to each project's style and contribution guidelines for submitting patches and additions. In general, we follow the "clone-and-pull" Git workflow.
1. **Clone** the project to your own machine
2. **Commit** changes to your own branch
3. **Push** your work back up to your branch
4. Submit a **Merge Request** so that anyone can review your changes

NOTE: Be sure to merge the latest from "upstream" before making a merge request!

## Getting started

### Requirements
- go version >= 1.20
- Makefile
- mockery v2.32.0
- swag


## Usage

### Config
Clone config file `config.yaml.example` from directory `/config/files`, put it on the same directory and rename it to `config.yaml`

You can also define config from env.

### API Server
running http server using
```
make run-api
```
or go command
```
go run main.go serve-http
```

### Mockery
#### Install the Mockery
https://vektra.github.io/mockery/latest/installation/#go-install
```
$ go install github.com/vektra/mockery/v2@v2.32.0
```

Check version
```
$ mockery --version
16 Jul 23 19:16 WIB INF couldn't read any config file version=v2.32.0
16 Jul 23 19:16 WIB INF Starting mockery dry-run=false version=v2.32.0
16 Jul 23 19:16 WIB INF Using config:  dry-run=false version=v2.32.0
v2.32.0
```

#### Remock Interface
You will need mock interface for unit test. You can run command to make mock for all required directories:
```
make remock
```


### Swagger
#### Install the CLI
```
$ go get -u github.com/swaggo/swag/cmd/swag@v1.8.12
```

Register on bash or zsh
```
$ sudo nano ~/.zshrc

// ...
export PATH=$(go env GOPATH)/bin:$PATH

// save 
$ source ~/.zshrc
```

#### Access the swagger docs
```
http://{HOST}/docs/index.html
```

####
update swagger docs file manually by running command
```
make swag-init
```

accessing swagger docs using
```
http://localhost:8900/docs/index.html
```

enable/disable swagger docs by configuration properties, default is enabled


### DB Migration
Reference:
- [go-migrate](https://github.com/golang-migrate/migrate)
- [go-migrate for PostgreSQL](https://github.com/golang-migrate/migrate/tree/master/database/postgres)

#### Install Go Migrate
```
$ go install -tags "postgres" github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```
#### Create Migration File
```
$ migrate create -ext sql -dir db/migrations [name_of_migration_file]
```
#### Migrate up
```
$ migrate -database "postgres://[user]:[password]@[host]:[port]/[dbname]?query" -path db/migrations up
```
#### Migrate down
```
$ migrate -database "postgres://[user]:[password]@[host]:[port]/[dbname]?query" -path db/migrations down
```