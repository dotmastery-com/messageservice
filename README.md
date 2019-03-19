# MessageSocketSevice

The MessageSocketService microservice provides a single websocket endpoint:

- localhost:12345 

A new *Client* object is created to track every browser connection and these are registered with the *ClientManager* object.

Clients will listen for incoming messages and pass them to the ClientManager, which broadcasts them out to all resgistered clients for display.

Messages are also put onto a Kafka queue for dispersal to other application instances.
We listen for incoming messages on the same kafka queue and give them to the ClientManager to broadcast to our registered clients.

Messages are tagged with a unique location code, so we can discard messages that originate locally that we've already broadcast.


##Prerequisites

Install Golang as per the install instructions here: https://golang.org/doc/install

##Running the code

Clone the project and make sure its part of the $GOPATH
From the command line run go run main.go
You should see a confirmation message that the service is up and running
