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

package features

import (
	"context"

	semver "github.com/Masterminds/semver/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// New expects a client with a valid connection to mongoDB.
func New(client *mongo.Client) *Features {
	f := &Features{}
	f.init(client)

	return f
}

type Features struct {
	// Sessions returns whether the mongo connected to supports the use of
	// server sessions via the mongo driver `mongo.NewSession`.
	//
	// Server sessions are only supported on a mongo replica/sharded set enabled
	// via the --replSet switch on `mongod`.
	//
	// Refer: https://docs.mongodb.com/manual/reference/server-sessions/
	Sessions bool

	// Transactions returns whether the mongo connected to supports Distributed
	// Transactions/Multi-Document Transactions.
	//
	// Refer: https://docs.mongodb.com/manual/core/transactions/
	Transactions bool

	// Version returns the semver version of mongo connected to.
	Version *semver.Version
}

func (f *Features) init(c *mongo.Client) {
	ctx := context.Background()
	adminDB := c.Database("admin")

	f.getVersion(ctx, adminDB)
	f.canSession(ctx, adminDB)
	f.canTransact()
}

type buildInfo struct {
	Version string
}

// getVersion gets the connected mongo version and parses it using semver.
func (f *Features) getVersion(ctx context.Context, adminDB *mongo.Database) {
	cmd := bson.D{
		{
			Key:   "buildInfo",
			Value: 1,
		},
	}
	var result buildInfo
	err := adminDB.RunCommand(ctx, cmd).Decode(&result)
	if err != nil {
		f.Version = &semver.Version{}
		return
	}

	f.Version = semver.MustParse(result.Version)
}

// Sessions returns whether the mongo connected to supports the use of
// server sessions via the mongo driver `mongo.NewSession`.
//
// Server sessions are only supported on a mongo replica/sharded set enabled
// via the --replSet switch on `mongod`.
//
// Refer: https://docs.mongodb.com/manual/reference/server-sessions/
func (f *Features) canSession(ctx context.Context, adminDB *mongo.Database) {
	f.Sessions = true

	cmd := bson.D{
		{
			Key:   "replSetGetStatus",
			Value: 1,
		},
	}
	res := adminDB.RunCommand(ctx, cmd)
	if res.Err() != nil {
		if mErr, ok := res.Err().(mongo.CommandError); ok {
			if mErr.Code == 76 {
				f.Sessions = false
			}
		}
	}
}

// canTransact checks whether the mongo connected to supports Distributed
// Transactions/Multi-Document Transactions. This is done based on mongo
// version detection.
//
// Refer: https://docs.mongodb.com/manual/core/transactions/
func (f *Features) canTransact() {
	mongoV4 := semver.MustParse("4.0.0")

	if f.Version.GreaterThan(mongoV4) {
		f.Transactions = true
	}
}
