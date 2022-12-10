package mongodb

import (
	"context"
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/nam-truong-le/lambda-utils-go/pkg/aws/ssm"
	"github.com/nam-truong-le/lambda-utils-go/pkg/logger"
)

var (
	initClient sync.Once
	client     *mongo.Client
	database   string
)

// newClient returns mongo client. Stage must be in context.
func newClient(ctx context.Context) (*mongo.Client, error) {
	log := logger.FromContext(ctx)

	initClient.Do(func() {
		log.Infof("init mongodb client")
		mongoHost, err := ssm.GetParameter(ctx, "/mongo/host", false)
		if err != nil {
			return
		}
		mongoUsername, err := ssm.GetParameter(ctx, "/mongo/username", false)
		if err != nil {
			return
		}
		mongoPassword, err := ssm.GetParameter(ctx, "/mongo/password", true)
		if err != nil {
			return
		}
		database, err = ssm.GetParameter(ctx, "/mongo/db", false)
		if err != nil {
			return
		}
		log.Infof("host = %s, user = %s, db = %s", mongoHost, mongoUsername, database)
		mongoFullUrl := fmt.Sprintf("mongodb+srv://%s/?retryWrites=true&w=majority", mongoHost)

		c, err := mongo.NewClient(
			options.Client().ApplyURI(mongoFullUrl).SetAuth(options.Credential{
				Username: mongoUsername,
				Password: mongoPassword,
			}),
		)
		if err != nil {
			log.Errorf("failed to create mongodb client: %s", err)
			return
		}
		err = c.Connect(ctx)
		if err != nil {
			log.Errorf("failed to connect to mongodb: %s", err)
			return
		}

		client = c
	})

	if client == nil {
		return nil, fmt.Errorf("failed to init mongodb client")
	}
	return client, nil
}

// Collection returns collection. Stage must be in context.
func Collection(ctx context.Context, name string) (*mongo.Collection, error) {
	log := logger.FromContext(ctx)
	log.Infof("get collection [%s]", name)
	c, err := newClient(ctx)
	if err != nil {
		return nil, err
	}
	return c.Database(database).Collection(name), nil
}

// Disconnect disconnects from db
func Disconnect(ctx context.Context) {
	log := logger.FromContext(ctx)
	if client != nil {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Errorf("failed to disconnect mongodb client: %s", err)
		} else {
			log.Infof("mongodb client disconnected")
			client = nil
			initClient = sync.Once{}
		}
	}
}
