rmdb:
	rm -f db/development.db
	rm -f db/test.db

initdb: rmdb
	psql -hlocalhost -dpostgres -upostgres < db/migrations/0001-init.up.sql

seed: initdb
	psql -hlocalhost -dpostgres -upostgres < db/seed.sql

test:
	go test ./... --count=1 --cover

api:
	go run ./cli/api/main.go

admin-api:
	go run ./cli/adminapi
ws:
	go run ./cli/websocket/main.go

watcher:
	go run ./cli/watcher/main.go

engine:
	go run ./cli/engine/main.go

launcher:
	go run ./cli/launcher/main.go

maker:
	go run ./cli/maker/main.go

clean:
	go clean

.PHONY: test api ws watcher engine launcher
