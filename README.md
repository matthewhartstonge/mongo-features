# mongo-features

The official mongo driver is great - but it's a bit more low level than the mgo
community driver, so we need to check whether certain mongo features can be 
used. 

Enter `mongo (feat. features)`.

## tl;dr
Refer: [_examples/tldr](./_examples/tldr)
```go
package main

import (
    "context"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "github.com/matthewhartstonge/mongo-features"
)

func main() {
    ctx := context.Background()
    client, err := mongo.Connect(ctx)
    if err != nil {
        panic(err)
    }

    feat := features.New(client)
    fmt.Printf("I am running on mongo major version: %s\n", feat.Version.Major())
    fmt.Printf("I am running on mongo minor version: %s\n", feat.Version.Minor())
    fmt.Printf("I am running on mongo version: %s\n", feat.Version.String())
    fmt.Printf("I can perform server sessions: %t\n", feat.Sessions)
    fmt.Printf("I can perform multi-document acid transactions: %t\n", feat.Transactions)
}
```

## Made for Plug'n'Play
When creating your own mongo datastore API, you can plug this bad boy into your structs:

Refer: [_examples/plugnplay](./_examples/plugnplay)
```go
package main

import (
	"context"
	"fmt"

	feat "github.com/matthewhartstonge/mongo-features"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	db *mongo.Database

	mongo *feat.Features
}

type thing struct {
	Name string
}

func (s *Store) CreateThing(ctx context.Context, thing1 thing) {
	if s.mongo.Sessions {
		sess, err := s.db.Client().StartSession()
		if err != nil {
			panic(err)
		}
		ctx = mongo.NewSessionContext(ctx, sess)
		defer sess.EndSession(ctx)
	}

	// like, totally put your transact-able actions in here...
}


func New() *Store {
	ctx := context.Background()
	client, err := mongo.Connect(ctx)
	if err != nil {
		panic(err)
	}

	testDb := client.Database("test")
	return &Store{
		db:    testDb,
		mongo: feat.New(client),
	}
}

func main() {
	store := New()
	fmt.Printf("My datastore is running on mongo major version: %d\n", store.mongo.Version.Major())
	fmt.Printf("My datastore is running on mongo minor version: %d\n", store.mongo.Version.Minor())
	fmt.Printf("My datastore is running on mongo version: %s\n", store.mongo.Version.String())
	fmt.Printf("My datastore can perform server sessions: %t\n", store.mongo.Sessions)
	fmt.Printf("My datastore can perform multi-document acid transactions: %t\n", store.mongo.Transactions)
}
```
