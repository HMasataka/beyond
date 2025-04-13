package main

import (
	"context"
	"fmt"

	"github.com/caravan-inc/oshi-card-card-recommender/config"
	"github.com/caravan-inc/oshi-card-card-recommender/infrastructure"
)

func main() {
	ctx := context.Background()

	cfg, err := config.NewGoogleCloudConfig()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Google Cloud Config: %+v\n", cfg)

	client, err := infrastructure.NewPubsubClient(ctx, cfg)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	topic, err := client.CreateTopic(ctx, cfg.PubsubConfig.TopicID)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Created topic: %s\n", topic.ID())
}
