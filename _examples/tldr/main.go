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
