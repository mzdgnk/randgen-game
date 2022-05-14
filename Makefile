.PHONY: clean

local: docker-compose-up bin/randgen-game
	DATABASE_URL="user=user password=password host=localhost port=5432 dbname=db sslmode=disable" \
	heroku local web

FORCE:
bin/randgen-game: FORCE
	go build -o bin/randgen-game main.go

docker-compose-up:
	docker-compose up -d
	sleep 1

docker-compose-down:
	docker-compose down

clean: docker-compose-down
	rm -f bin/randgen-game

psql:
	PGPASSWORD=password psql -h localhost -p 5432 -U user -d db
