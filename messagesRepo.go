package messagesStorageMongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/mailhedgehog/contracts"
	"github.com/mailhedgehog/email"
	"github.com/mailhedgehog/logger"
	"github.com/mailhedgehog/smtpMessage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/mail"
	"time"
)

type Message struct {
	ID      smtpMessage.MessageID    `bson:"id"`
	Room    string                   `bson:"room"`
	From    []*mail.Address          `bson:"from"`
	To      []*mail.Address          `bson:"to"`
	Subject string                   `bson:"subject"`
	Date    time.Time                `bson:"date"`
	Message *smtpMessage.SmtpMessage `bson:"message"`
}

type mongoMessagesRepo struct {
	context *storageContext
	room    contracts.Room
}

func (repo *mongoMessagesRepo) SetRoom(room contracts.Room) {
	repo.room = room
}

func (repo *mongoMessagesRepo) Store(message *smtpMessage.SmtpMessage) (smtpMessage.MessageID, error) {
	if repo.context.perRoomLimit > 0 && repo.context.perRoomLimit <= repo.Count() {
		repo.context.storage.RoomsRepo().Delete(repo.room)
	}

	emailMessage := message.GetEmail()
	if emailMessage == nil {
		emailMessage = &email.Email{}
	}

	insertResult, err := repo.context.collection.InsertOne(context.TODO(), Message{
		message.ID,
		repo.context.roomName(repo.room),
		emailMessage.From,
		emailMessage.To,
		emailMessage.Subject,
		emailMessage.Date,
		message,
	})

	logManager().Debug(fmt.Sprintf("New message saved, mongo _id='%s'", insertResult.InsertedID))

	return message.ID, err
}

func (repo *mongoMessagesRepo) List(query contracts.SearchQuery, offset, limit int) ([]smtpMessage.SmtpMessage, int, error) {
	opts := options.Find().SetSort(bson.M{"date": -1}).SetSkip(int64(offset)).SetLimit(int64(limit))

	textsMatch := bson.A{}
	for criteria, queryValue := range query {
		switch criteria {
		case contracts.SearchParamTo:
			textsMatch = append(
				textsMatch,
				bson.M{"to.name": primitive.Regex{Pattern: queryValue, Options: ""}},
				bson.M{"to.address": primitive.Regex{Pattern: queryValue, Options: ""}},
			)
		case contracts.SearchParamFrom:
			textsMatch = append(
				textsMatch,
				bson.M{"from.name": primitive.Regex{Pattern: queryValue, Options: ""}},
				bson.M{"from.address": primitive.Regex{Pattern: queryValue, Options: ""}},
			)
		case contracts.SearchParamContent:
			textsMatch = append(
				textsMatch,
				bson.M{"subject": primitive.Regex{Pattern: queryValue, Options: ""}},
			)
		}
	}

	filterQuery := bson.A{
		bson.M{"room": repo.context.roomName(repo.room)},
	}
	if len(textsMatch) > 0 {
		filterQuery = append(filterQuery, bson.M{"$or": textsMatch})
	}
	filter := bson.M{"$and": filterQuery}

	totalCount, err := repo.context.collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return nil, 0, err
	}

	cursor, err := repo.context.collection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, 0, err
	}

	var emailsList []smtpMessage.SmtpMessage
	var results []Message
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, 0, err
	}
	for _, result := range results {
		emailsList = append(emailsList, *result.Message)
	}

	return emailsList, int(totalCount), nil
}

func (repo *mongoMessagesRepo) Count() int {
	count, err := repo.context.collection.CountDocuments(context.TODO(), bson.D{{
		"room",
		repo.context.roomName(repo.room),
	}})
	logger.PanicIfError(err)

	return int(count)
}

func (repo *mongoMessagesRepo) Delete(messageId smtpMessage.MessageID) error {
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"room", repo.context.roomName(repo.room)}},
				bson.D{{"id", messageId}},
			}},
	}
	result, err := repo.context.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount != 1 {
		return errors.New(fmt.Sprintf("Unexpected count of deleted items, extected 1, got %d", result.DeletedCount))
	}

	return nil
}

func (repo *mongoMessagesRepo) Load(messageId smtpMessage.MessageID) (*smtpMessage.SmtpMessage, error) {
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"room", repo.context.roomName(repo.room)}},
				bson.D{{"id", messageId}},
			}},
	}
	var result Message
	err := repo.context.collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.Message, nil
}
