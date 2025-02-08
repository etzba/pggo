export ETZBA_DATABASE_USER=etzba
export ETZBA_DATABASE_PASSWORD=Pass1234
export ETZBA_DATABASE_DB=etzba
export ETZBA_DATABASE_PORT=5432
export ETZBA_DATABASE_HOST=localhost
export ETZBA_DATABASE_SSL=disable

docker-compose down
docker-compose up -d pg

sleep 6

go run main.go