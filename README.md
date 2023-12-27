# MailHedgehog package to save messages in MongoDB storage

All emails stored in MongoDB database. Useful if you have a lot of emails in application.

## Usage

```go
storage := CreateMongoDbStorage(createMongoTestCollection(), &contracts.MessagesStorageConfiguration{PerRoomLimit: 100})

msg, err := storage.MessagesRepo(room).Load("ID")
```

## Development

```shell
go mod tidy
go mod verify
go mod vendor
```

Test

```shell
docker-compose up -d
go test --cover
```

## Credits

- [![Think Studio](https://yaroslawww.github.io/images/sponsors/packages/logo-think-studio.png)](https://think.studio/)
