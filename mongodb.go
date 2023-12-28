package messagesStorageMongo

import (
	"context"
	"fmt"
	"github.com/mailhedgehog/contracts"
	"github.com/mailhedgehog/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var configuredLogger *logger.Logger

func logManager() *logger.Logger {
	if configuredLogger == nil {
		configuredLogger = logger.CreateLogger("messagesStorageMongo")
	}
	return configuredLogger
}

type storageContext struct {
	storage      *Mongo
	perRoomLimit int
	collection   *mongo.Collection
}

func (c *storageContext) roomName(room contracts.Room) string {
	if len(room) <= 0 {
		room = "_default"
	}
	return string(room)
}

type Mongo struct {
	context *storageContext
}

func CreateMongoDbStorage(collection *mongo.Collection, storageConfig *contracts.MessagesStorageConfiguration) *Mongo {
	indexName, err := collection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{
			{"room", -1},
			{"id", 1},
		},
	})
	logger.PanicIfError(err)

	logManager().Debug(fmt.Sprintf("Index [%s] created", indexName))

	storage := &Mongo{
		context: &storageContext{
			collection:   collection,
			perRoomLimit: storageConfig.PerRoomLimit,
		}}

	storage.context.storage = storage

	return storage
}

func (repo *Mongo) RoomsRepo() contracts.RoomsRepo {
	return &mongoRoomsRepo{
		context: repo.context,
	}
}

func (repo *Mongo) MessagesRepo(room contracts.Room) contracts.MessagesRepo {
	return &mongoMessagesRepo{
		context: repo.context,
		room:    room,
	}
}
