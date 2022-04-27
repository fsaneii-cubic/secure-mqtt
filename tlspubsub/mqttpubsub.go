package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	clientID   = "ssl_mqtt_client"
	brokerPort = 8883
	brokerAddr = "172.17.0.3"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Message %s received on topic %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	opts := client.OptionsReader()
	fmt.Printf("client %s is Connected.\n", opts.ClientID())
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection lost: %s\n", err.Error())
}

func newTLSConfig() *tls.Config {
	certPool := x509.NewCertPool()
	certPath := "client-certs"
	ca, err := ioutil.ReadFile(fmt.Sprintf("%s/ca.crt", certPath))
	if err != nil {
		panic(err.Error())
	}

	certPool.AppendCertsFromPEM(ca)

	//import client certificate/key pair
	//LoadX509KeyPair(certFile, keyFile)
	cert, err := tls.LoadX509KeyPair(fmt.Sprintf("%s/client.crt", certPath), fmt.Sprintf("%s/client.key", certPath))

	if err != nil {
		panic(err)
	}

	//just print out the client certificate..
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		panic(err)
	}
	fmt.Println(cert.Leaf)
	fmt.Println()

	//Create tls config with tls properties
	return &tls.Config{
		RootCAs:            certPool,
		ClientAuth:         tls.NoClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert},
	}
}

func main() {
	var broker = fmt.Sprintf("ssl://%s:%d", brokerAddr, brokerPort)
	options := mqtt.NewClientOptions()

	//broker IP and port
	options.AddBroker(broker)
	options.SetClientID("ssl_mqtt_client")
	options.SetDefaultPublishHandler(messagePubHandler)
	options.OnConnect = connectHandler
	options.OnConnectionLost = connectLostHandler
	options.SetTLSConfig(newTLSConfig())

	client := mqtt.NewClient(options)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	//subscription
	topic := "topic/security"
	token = client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic %s\n\n", topic)

	//publish
	num := 20
	for i := 0; i < num; i++ {
		text := fmt.Sprintf("%d", i)
		token := client.Publish(topic, 0, false, text)
		token.Wait()
		time.Sleep(time.Second)
	}

	client.Disconnect(150)
}
