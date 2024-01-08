package mongodb

import (
	"bot/internal/config"
	"bot/pkg/logging"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type MongoDB struct {
	Client    *mongo.Client
	log       *logging.Logger
	cfg       *config.MongoDBConfig
	connected bool
}

func New(log *logging.Logger, cfg *config.MongoDBConfig) *MongoDB {
	return &MongoDB{
		Client:    nil,
		log:       log,
		cfg:       cfg,
		connected: false,
	}
}

func (m *MongoDB) Connect() error {
	var err error
	uri := m.cfg.ConnString()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	m.Client, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		m.log.Error("Failed to connect to MongoDB", zap.Error(err))
		return err
	}

	var result bson.M
	if err = m.Client.Database("admin").RunCommand(context.TODO(), bson.M{"ping": 1}).Decode(&result); err != nil {
		m.log.Error("Failed to ping MongoDB", zap.Error(err))
		return err
	}

	m.connected = true
	return nil
}

func (m *MongoDB) IsConnected() bool {
	return m.connected
}

func (m *MongoDB) Disconnect() {
	if m.Client != nil {
		err := m.Client.Disconnect(context.TODO())
		if err != nil {
			m.log.Error("Failed to disconnect from MongoDB", zap.Error(err))
		}
	}
}

func (m *MongoDB) DB() *mongo.Database {
	return m.Client.Database(m.cfg.DbName)
}

func (m *MongoDB) Collection(name string) *mongo.Collection {
	return m.DB().Collection(name)
}
