package mongofx

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
)

type MongoConfig struct {
	Dsn string `mapstructure:"dsn" validate:"required,uri"`
}

func NewMongoClient(config *MongoConfig) (*mongo.Client, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(config.Dsn))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func CleanupMongoClient(lifecycle fx.Lifecycle, client *mongo.Client) {
	lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return client.Disconnect(ctx)
		},
	})
}

var Module = fx.Options(
	fx.Provide(NewMongoClient),
	fx.Invoke(CleanupMongoClient),
)
