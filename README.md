# mongo-features

[![PkgGoDev](https://pkg.go.dev/badge/github.com/matthewhartstonge/mongo-features)](https://pkg.go.dev/github.com/matthewhartstonge/mongo-features) [![Go Report Card](https://goreportcard.com/badge/github.com/matthewhartstonge/mongo-features)](https://goreportcard.com/report/github.com/matthewhartstonge/mongo-features)

The official mongo driver is great - but it's a bit more low level than the mgo
community driver, so we need to check whether certain mongo features can be 
used. 

Enter `mongo (feat. features)`.

> Session detection requires the consumer to have `clusterMonitor` mongodb permissions.
> Sessions should work on a single node, but there is [an in progress fix for this](https://github.com/mongodb/mongo-go-driver/pull/497).

## tl;dr
Refer: [_examples/tldr](./_examples/tldr)
```go
package main

import (
	"context"
	"fmt"

	"github.com/matthewhartstonge/mongo-features"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx)
	if err != nil {
		panic(err)
	}

	featureSet := features.New(client)
	fmt.Printf("I am running on mongo major version: %d\n", featureSet.MongoVersion.Major())
	fmt.Printf("I am running on mongo minor version: %d\n", featureSet.MongoVersion.Minor())
	fmt.Printf("I am running on mongo version: %s\n", featureSet.MongoVersion.String())
	fmt.Printf("I can perform server sessions: %t\n", featureSet.HasSessions)
	fmt.Printf("I can perform multi-document acid transactions: %t\n", featureSet.HasTransactions)
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

	*feat.Features
}

type thing struct {
	Name string
}

func (s *Store) CreateThing(ctx context.Context, thing1 thing) {
	if s.HasSessions {
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
		db:       testDb,
		Features: feat.New(client),
	}
}

func main() {
	store := New()
	fmt.Printf("My datastore is running on mongo major version: %d\n", store.MongoVersion.Major())
	fmt.Printf("My datastore is running on mongo minor version: %d\n", store.MongoVersion.Minor())
	fmt.Printf("My datastore is running on mongo version: %s\n", store.MongoVersion.String())
	fmt.Printf("My datastore can perform server sessions: %t\n", store.HasSessions)
	fmt.Printf("My datastore can perform multi-document acid transactions: %t\n", store.HasTransactions)
}
```
