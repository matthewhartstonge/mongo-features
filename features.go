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

	"github.com/Masterminds/semver"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// New expects a client with a valid connection to mongoDB.
func New(client *mongo.Client) *Features {
	f := &Features{}
	f.init(client)

	return f
}

// Features provides a selection of boolean switches for detection of mongodb
// features.
type Features struct {
	// HasSessions returns whether the mongo connected to supports the use of
	// server sessions via the mongo driver `mongo.NewSession`.
	//
	// Server sessions are only supported on a mongo replica/sharded set enabled
	// via the --replSet switch on `mongod`.
	//
	// Refer: https://docs.mongodb.com/manual/reference/server-sessions/
	HasSessions bool

	// HasTransactions returns whether the mongo connected to supports Distributed
	// Transactions/Multi-Document Transactions.
	//
	// Refer: https://docs.mongodb.com/manual/core/transactions/
	HasTransactions bool

	// MongoVersion returns the semver version of mongo connected to.
	MongoVersion *semver.Version
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
		f.MongoVersion = &semver.Version{}
		return
	}

	f.MongoVersion = semver.MustParse(result.Version)
}

type replInfo struct {
	// Ok returns whether the replSet is ready.
	Ok bool
}

// Sessions returns whether the mongo connected to supports the use of
// server sessions via the mongo driver `mongo.NewSession`.
//
// Server sessions are only supported on a mongo replica/sharded set enabled
// via the --replSet switch on `mongod`.
//
// Refer: https://docs.mongodb.com/manual/reference/server-sessions/
func (f *Features) canSession(ctx context.Context, adminDB *mongo.Database) {
	cmd := bson.D{
		{
			Key:   "replSetGetStatus",
			Value: 1,
		},
	}
	var result replInfo
	err := adminDB.RunCommand(ctx, cmd).Decode(&result)
	if err != nil {
		// assume we don't have session support on error..
		// error code 76 will be thrown if replSet is not enabled.
		return
	}

	f.HasSessions = result.Ok
}

// canTransact checks whether the mongo connected to supports Distributed
// Transactions/Multi-Document Transactions. This is done based on mongo
// version detection.
//
// Refer: https://docs.mongodb.com/manual/core/transactions/
func (f *Features) canTransact() {
	mongoV4 := semver.MustParse("4.0.0")

	if f.MongoVersion.GreaterThan(mongoV4) {
		f.HasTransactions = true
	}
}
