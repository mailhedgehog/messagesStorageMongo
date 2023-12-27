package messagesStorageMongo

import (
	"context"
	"fmt"
	"github.com/mailhedgehog/contracts"
	"github.com/mailhedgehog/gounit"
	"github.com/mailhedgehog/logger"
	"github.com/mailhedgehog/smtpMessage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

func TestStore(t *testing.T) {
	room := "foo_bar"

	storage := CreateMongoDbStorage(createMongoTestCollection(), &contracts.MessagesStorageConfiguration{PerRoomLimit: 100})

	for i := 0; i < 15; i++ {
		id := smtpMessage.MessageID(fmt.Sprint(i))
		msg := &smtpMessage.SmtpMessage{
			ID: id,
		}

		storedId, err := storage.MessagesRepo(room + fmt.Sprint(i)).Store(msg)
		(*gounit.T)(t).AssertEqualsString(string(id), string(storedId))
		(*gounit.T)(t).AssertNotError(err)

		_, _ = storage.MessagesRepo(room + fmt.Sprint(i)).Store(&smtpMessage.SmtpMessage{
			ID: smtpMessage.MessageID(fmt.Sprint(i + 1)),
		})
	}

	count, err := storage.context.collection.CountDocuments(context.TODO(), bson.D{}, options.Count().SetHint("_id_"))
	logger.PanicIfError(err)

	(*gounit.T)(t).AssertEqualsInt(30, int(count))

}

func TestStore_ClearIfOverLimit(t *testing.T) {
	room := "foo_bar"

	storage := CreateMongoDbStorage(createMongoTestCollection(), &contracts.MessagesStorageConfiguration{PerRoomLimit: 6})

	for i := 0; i < 15; i++ {
		id := smtpMessage.MessageID(fmt.Sprint(i))
		msg := &smtpMessage.SmtpMessage{
			ID: id,
		}

		storedId, err := storage.MessagesRepo(room).Store(msg)
		(*gounit.T)(t).AssertEqualsString(string(id), string(storedId))
		(*gounit.T)(t).AssertNotError(err)
	}

	count, err := storage.context.collection.CountDocuments(context.TODO(), bson.D{}, options.Count().SetHint("_id_"))
	logger.PanicIfError(err)

	(*gounit.T)(t).AssertEqualsInt(3, int(count))
}

func TestCount(t *testing.T) {
	room := "foo_bar"

	storage := CreateMongoDbStorage(createMongoTestCollection(), &contracts.MessagesStorageConfiguration{PerRoomLimit: 100})

	for i := 0; i < 15; i++ {
		id := smtpMessage.MessageID(fmt.Sprint(i))
		msg := &smtpMessage.SmtpMessage{
			ID: id,
		}

		storedId, err := storage.MessagesRepo(room).Store(msg)
		(*gounit.T)(t).AssertEqualsString(string(id), string(storedId))
		(*gounit.T)(t).AssertNotError(err)
	}

	for i := 0; i < 4; i++ {
		id := smtpMessage.MessageID(fmt.Sprint(i))
		msg := &smtpMessage.SmtpMessage{
			ID: id,
		}

		storedId, err := storage.MessagesRepo(room + "2").Store(msg)
		(*gounit.T)(t).AssertEqualsString(string(id), string(storedId))
		(*gounit.T)(t).AssertNotError(err)
	}

	(*gounit.T)(t).AssertEqualsInt(15, storage.MessagesRepo(room).Count())
	(*gounit.T)(t).AssertEqualsInt(4, storage.MessagesRepo(room+"2").Count())
}

func TestDelete(t *testing.T) {
	room := "foo_bar"

	storage := CreateMongoDbStorage(createMongoTestCollection(), &contracts.MessagesStorageConfiguration{PerRoomLimit: 100})

	for i := 0; i < 15; i++ {
		id := smtpMessage.MessageID(fmt.Sprint(i))
		msg := &smtpMessage.SmtpMessage{
			ID: id,
		}

		storedId, err := storage.MessagesRepo(room).Store(msg)
		(*gounit.T)(t).AssertEqualsString(string(id), string(storedId))
		(*gounit.T)(t).AssertNotError(err)
	}

	for i := 0; i < 4; i++ {
		id := smtpMessage.MessageID(fmt.Sprint(i))
		msg := &smtpMessage.SmtpMessage{
			ID: id,
		}

		storedId, err := storage.MessagesRepo(room + "2").Store(msg)
		(*gounit.T)(t).AssertEqualsString(string(id), string(storedId))
		(*gounit.T)(t).AssertNotError(err)
	}

	(*gounit.T)(t).AssertEqualsInt(15, storage.MessagesRepo(room).Count())
	(*gounit.T)(t).AssertEqualsInt(4, storage.MessagesRepo(room+"2").Count())

	(*gounit.T)(t).AssertNotError(storage.MessagesRepo(room).Delete("3"))
	(*gounit.T)(t).AssertNotError(storage.MessagesRepo(room + "2").Delete("1"))
	(*gounit.T)(t).AssertNotError(storage.MessagesRepo(room + "2").Delete("2"))

	(*gounit.T)(t).AssertEqualsInt(14, storage.MessagesRepo(room).Count())
	(*gounit.T)(t).AssertEqualsInt(2, storage.MessagesRepo(room+"2").Count())
}

func TestLoad(t *testing.T) {
	room := "foo_bar"

	storage := CreateMongoDbStorage(createMongoTestCollection(), &contracts.MessagesStorageConfiguration{PerRoomLimit: 100})

	for i := 0; i < 15; i++ {
		id := smtpMessage.MessageID(fmt.Sprint(i))
		msg := &smtpMessage.SmtpMessage{
			ID: id,
		}

		storedId, err := storage.MessagesRepo(room).Store(msg)
		(*gounit.T)(t).AssertEqualsString(string(id), string(storedId))
		(*gounit.T)(t).AssertNotError(err)
	}

	msg, err := storage.MessagesRepo(room).Load("3")

	(*gounit.T)(t).AssertNotError(err)

	(*gounit.T)(t).AssertEqualsString("3", string(msg.ID))
}

func TestList(t *testing.T) {
	room := "foo_bar"

	storage := CreateMongoDbStorage(createMongoTestCollection(), &contracts.MessagesStorageConfiguration{PerRoomLimit: 100})

	for i := 0; i < 15; i++ {
		id := smtpMessage.MessageID(fmt.Sprint(i))
		msg := &smtpMessage.SmtpMessage{
			ID: id,
		}
		msg.SetOrigin(fmt.Sprintf(`From: Rares <quix-%d@quib.com>
Date: Thu, 2 May 2019 11:25:35 +0300
Subject: Re: kern/54143 (virtualbox)
To: foo-%d@quib.com
Content-Type: multipart/mixed; boundary="0000000000007e2bb40587e36196"

--0000000000007e2bb40587e36196
Content-Type: text/html; charset="UTF-8"

<div dir="ltr"><div>html text part</div><div><br></div><div><br><br></div></div>

--0000000000007e2bb40587e36196--
`, i, i))

		storedId, err := storage.MessagesRepo(room).Store(msg)
		(*gounit.T)(t).AssertEqualsString(string(id), string(storedId))
		(*gounit.T)(t).AssertNotError(err)
	}

	msgs, count, err := storage.MessagesRepo(room).List(contracts.SearchQuery{
		contracts.SearchParamFrom: "quix-1",
	}, 0, 4)

	(*gounit.T)(t).AssertNotError(err)
	(*gounit.T)(t).AssertEqualsInt(6, count)
	(*gounit.T)(t).AssertEqualsInt(4, len(msgs))
	(*gounit.T)(t).AssertEqualsString("quix-11@quib.com", msgs[3].From.Address())
	(*gounit.T)(t).AssertEqualsString("foo-11@quib.com", msgs[3].To[0].Address())
}
