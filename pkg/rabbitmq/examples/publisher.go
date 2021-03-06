package main

import (
	"fmt"
	rabbit "forklift/pkg/rabbitmq"
	"golang.org/x/net/context"
)

func main() {

	var rabbitClient = rabbit.NewRabbitMqClient([]string{"127.0.0.1"}, "guest", "guest", "", rabbit.RetryCount(2))

	rabbitClient.AddPublisher("PersonV1_Direct", rabbit.Direct, PersonV1{})

	var err = rabbitClient.Publish(context.TODO(), "123", PersonV1{
		Name:    "Emre",
		Surname: "Tiryaki",
		City:    City{},
		Count:   1,
	})

	fmt.Print(err)

}
