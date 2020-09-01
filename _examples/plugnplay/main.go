/*
 * Copyright 2020 Matthew Hartstonge <matt@mykro.co.nz>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
