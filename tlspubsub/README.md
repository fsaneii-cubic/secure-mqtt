Secure MQTT using SSL/TLS

1. All the clients keys are placed in client-certs
2. All the broker keys are located in mosquitto-certs. Copy the mosquitto certs in the mosquitto directory some location similar to the following: 

/Mosquitto/config/certs

update the  /mosquitto/config/mosquitto.conf for the followings

listener 8883

cafile /mosquitto/config/certs/ca.crt
certfile /mosquitto/config/certs/server.crt
keyfile /mosquitto/config/certs/server.key

require_certificate true
use_identity_as_username true

run the mosquitto in docker

run the client program
go run mqttpubsub.go

You should able to see the client subscribes to topic/security as well as publishing them.
