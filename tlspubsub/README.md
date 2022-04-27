Secure MQTT using SSL/TLS

1. Create a subdirectory called client-certs.
2. Locate all the clients keys and certs in client-certs directory.
3. Copy the mosquitto broker certs in the mosquitto directory such as followings: 

/Mosquitto/config/certs

update the  /mosquitto/config/mosquitto.conf for the followings

listener 8883

cafile /mosquitto/config/certs/ca.crt
certfile /mosquitto/config/certs/server.crt
keyfile /mosquitto/config/certs/server.key

require_certificate true
use_identity_as_username true

restart the mosquitto in docker

run the client program
go run mqttpubsub.go

You should be able to see the client subscribes to topic/security as well as publishing them.
