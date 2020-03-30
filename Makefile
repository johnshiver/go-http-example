build: stop
	docker-compose build server db

gen: clean
	go generate ./...

clean: stop
	git clean -dfx

run: clean build
	docker-compose up db server

stop:
	docker-compose down

test: clean gen
	docker-compose build db_test
	docker-compose up -d db_test
    # wait for db to come up
	sleep 5
	go test -v ./...
	docker-compose down

test-race: clean gen
	docker-compose build db_test
	docker-compose up -d db_test
    # wait for db to come up
	sleep 5
	go test -race -v ./...
	docker-compose down

new-migration:
	touch sql/$(shell date +%s)_$(name).sql
