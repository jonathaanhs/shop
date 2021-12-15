# Shop

Checkout services for shopping with promos.

## Requirements
- Go 1.17
- [golang-migrate](https://github.com/golang-migrate/migrate)
- Docker
- Docker Compose

## Quick Start
```bash
make test # to run unit test
make run # run the app
make migrate-up  # run data migration
```

## Call API
```bash
curl --header "Content-Type: application/json" \
  --request POST \
  --data '[{"product_id": 2,"qty": 1}]' \
  http://localhost:8080/checkout
```