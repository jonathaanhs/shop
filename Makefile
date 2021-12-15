
test:
	go test ./...

run:
	docker-compose up --build -d

stop:
	docker-compose down

migrate-up:
	migrate -database "postgresql://shop:shop@:5432/shop?sslmode=disable" -path ./sql up

migrate-down:
	migrate -database "postgresql://shop:shop@:5432/shop?sslmode=disable" -path ./sql down