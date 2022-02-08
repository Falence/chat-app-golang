package main

type room struct {
	// forward is a channel that holds incoming messages
	// that should be forwarded to the other clients.
	forward chan []byte
	// join is a channel for clients wishing to join the room
	join chan *client
	// leave is a channerl for clients wishing to leave the room
	leave chan *client
	// clients hold all current clients in this room
	clients map[*client]bool
}