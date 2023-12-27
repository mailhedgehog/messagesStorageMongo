package messagesStorageMongo

import (
	"fmt"
	"github.com/mailhedgehog/contracts"
	"github.com/mailhedgehog/gounit"
	"github.com/mailhedgehog/smtpMessage"
	"testing"
)

func TestRoomsCount(t *testing.T) {
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

		storage.MessagesRepo(room + fmt.Sprint(i)).Store(msg)
	}

	(*gounit.T)(t).AssertEqualsInt(15, storage.RoomsRepo().Count())
}

func TestRoomsDelete(t *testing.T) {
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

		storage.MessagesRepo(room + fmt.Sprint(i)).Store(msg)
	}

	(*gounit.T)(t).AssertEqualsInt(15, storage.RoomsRepo().Count())

	(*gounit.T)(t).AssertNotError(storage.RoomsRepo().Delete(room + "2"))
	(*gounit.T)(t).AssertNotError(storage.RoomsRepo().Delete(room + "3"))
	(*gounit.T)(t).AssertNotError(storage.RoomsRepo().Delete(room + "4"))

	(*gounit.T)(t).AssertEqualsInt(12, storage.RoomsRepo().Count())
}

func TestRoomsList(t *testing.T) {
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

		storage.MessagesRepo(room + fmt.Sprint(i)).Store(msg)
	}

	rooms, err := storage.RoomsRepo().List(2, 6)

	(*gounit.T)(t).AssertNotError(err)
	(*gounit.T)(t).AssertEqualsInt(6, len(rooms))
}
