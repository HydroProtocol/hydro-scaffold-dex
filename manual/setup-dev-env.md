# How to setup env for development

## Step0: prerequisite

Make sure you have:

- docker & docker-compose
- yarn
- go

## Step1: basic services

start redis, local-ethereum-node and PostgresDB

```shell
docker-compose up redis ethereum-node db
```

## Step2: hydro services

start other Hydro Backend Services

```shell
# change to sub-dir: backend
cd ./backend

# API
make api

# engine
make engine

# launcher
make launcher

# watcher
make watcher

# websocket
make ws

# go back to project-dir
cd ..
```

## Step3: web

```shell
# change to sub-dir: web
cd ./web

# install dependencies
yarn install

# start web
yarn start

# go back to project dir
cd ..
```

