package client

import  "github.com/zhouhui8915/go-socket.io-client"

type socketClient struct {
	client		*socketio_client.Client
	incoming 	<- chan chat
}

type chat struct {
	Name	string
	Message	string
}