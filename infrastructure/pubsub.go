package infrastructure

import (
	"context"
	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/caravan-inc/oshi-card-card-recommender/config"
)

type PubsubClient interface {
	CreateTopic(ctx context.Context, topicID string) (*pubsub.Topic, error)
	Topic(topicID string) *pubsub.Topic
	CreateSubscription(ctx context.Context, id string, cfg pubsub.SubscriptionConfig) (*pubsub.Subscription, error)
	Subscription(id string) *pubsub.Subscription
	Close() error
}

type pubsubClient struct {
	client *pubsub.Client
}

var (
	cli        *pubsubClient
	pubsubOnce sync.Once
)

func NewPubsubClient(ctx context.Context, cfg *config.GoogleCloudConfig) (PubsubClient, error) {
	var err error

	pubsubOnce.Do(func() {
		client, e := pubsub.NewClient(ctx, cfg.ProjectID)
		if e != nil {
			err = e
			return
		}

		cli = &pubsubClient{
			client: client,
		}
	})

	return cli, err
}

func (c *pubsubClient) CreateTopic(ctx context.Context, topicID string) (*pubsub.Topic, error) {
	return c.client.CreateTopic(ctx, topicID)
}

func (c *pubsubClient) Topic(topicID string) *pubsub.Topic {
	return c.client.Topic(topicID)
}

func (c *pubsubClient) CreateSubscription(ctx context.Context, id string, cfg pubsub.SubscriptionConfig) (*pubsub.Subscription, error) {
	return c.client.CreateSubscription(ctx, id, cfg)
}

func (c *pubsubClient) Subscription(id string) *pubsub.Subscription {
	return c.client.Subscription(id)
}

func (c *pubsubClient) Close() error {
	if c.client != nil {
		return c.client.Close()
	}

	return nil
}
