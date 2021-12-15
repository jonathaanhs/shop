
test:
	go test ./...

run:
	docker-compose up --build -d

stop:
	docker-compose down