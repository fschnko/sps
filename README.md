#The simple pub-sub data bus.

##Overview

The package sps provide a thread safe messaging data bus.
Currently all data stored in the memory and can't be used for a stateful applications. 

#Usage

###Publisher

```golang
func main(){
	bus := sps.New()
	defer bus.Close() // Release allocated resources and routines.

    // Publish the message.
	bus.Publish("topic name", []byte("some data"))
}
```

###Subscriber

```golang
func main() {
	bus := sps.New()
	defer bus.Close() // Release allocated resources and routines.

    // Subscribe to the particular topic.
    bus.Subscribe("topic name", "subscriber name")
    
    // Make sure to unsubscribe from the topic.
    defer bus.Unsubscribe("topic name", "subscriber name")

    // Each poll returns unread messages. It can be used for circular queries.
	data, err := bus.Poll("topic name", "subscriber name")
	if err != nil {
        log.Printf("poll: %v", err)
        return
	}

	for _, msg := range data {
		fmt.Println(string(msg))
	}	
}

```