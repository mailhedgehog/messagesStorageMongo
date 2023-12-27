package messagesStorageMongo

import (
	"context"
	"fmt"
	"github.com/mailhedgehog/contracts"
	"github.com/mailhedgehog/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRoomsRepo struct {
	context *storageContext
}

func (repo *mongoRoomsRepo) List(offset, limit int) ([]contracts.Room, error) {
	cursor, err := repo.context.collection.Aggregate(context.TODO(), mongo.Pipeline{
		bson.D{
			{"$group", bson.M{"_id": "$room"}}},
		bson.D{{"$sort", bson.M{"_id": 1}}},
		bson.D{{"$skip", offset}},
		bson.D{{"$limit", limit}},
	},
	)
	logger.PanicIfError(err)

	var results []contracts.Room
	for cursor.Next(context.TODO()) {
		var result bson.M
		err := cursor.Decode(&result)
		logger.PanicIfError(err)
		results = append(results, result["_id"].(string))
	}

	logger.PanicIfError(cursor.Err())

	return results, nil
}

func (repo *mongoRoomsRepo) Count() int {
	groupStage := bson.A{
		bson.D{{"$group", bson.D{{"_id", "$room"}}}},
		bson.D{{"$group",
			bson.D{
				{"_id", 1},
				{"count", bson.D{{"$sum", 1}}},
			},
		}},
	}
	cursor, err := repo.context.collection.Aggregate(context.TODO(), groupStage)
	logger.PanicIfError(err)

	var results []bson.D
	err = cursor.All(context.TODO(), &results)
	logger.PanicIfError(err)

	if len(results) > 0 {
		return int(results[0][1].Value.(int32))
	}

	return 0
}

func (repo *mongoRoomsRepo) Delete(room contracts.Room) error {
	roomName := repo.context.roomName(room)
	result, err := repo.context.collection.DeleteMany(context.TODO(), bson.M{"room": roomName})

	logManager().Debug(fmt.Sprintf("Deleted room [%s] (%d items)", roomName, result.DeletedCount))

	return err
}
