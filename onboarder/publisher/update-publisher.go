package publisher

import (
	"github.com/micro/go-micro"
	"fmt"
	"context"
	"time"
	"github.com/ZTP/onboarder/publisher/pubsub-proto"
	"github.com/pborman/uuid"
)

func PublishOnUpdate (topic string) {
	service := micro.NewService(
		micro.Name("PnPClient"),
	)

	service.Init()

	pub := micro.NewPublisher(topic, service.Client())

	ev := &publisher.Event{
		Id:        uuid.NewUUID().String(),
		Timestamp: time.Now().Unix(),
	}

	fmt.Printf("Publishing update event on topic %v", topic)
	if err := pub.Publish(context.Background(), ev); err != nil {
		fmt.Printf("Error publishing on topic %v", topic)
	}
}
