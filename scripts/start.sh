export DATABASE_USER=etzba
export DATABASE_PASSWORD=Pass1234
export DATABASE_DB=etzba
export DATABASE_PORT=5432
export DATABASE_HOST=localhost
export DATABASE_SSL=disable

docker-compose down
docker-compose up -d pg

sleep 6

go run main.go