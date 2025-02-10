#curl http://localhost:8080/locations
#curl http://localhost:8080/locations/2
#curl -X POST http://localhost:8080/location -d '{ "name": "Etz", "address": "Davrewasdf",  "longtitude": 34.43212123, "latitude": 12.321312234 }'
#curl http://localhost:8080/locations
#curl http://localhost:8080/locations/1
#curl -X POST http://localhost:8080/location -d '{ "name": "Etz", "address": "Davwer",  "longtitude": 76.655645321, "latitude": 14.11312234 }'
#curl http://localhost:8080/locations
#curl http://localhost:8080/locations/2
etz api --exec=scripts/executions.yaml -w 10 -r 20 -d 3s