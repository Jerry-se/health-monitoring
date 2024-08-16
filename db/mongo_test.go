package db

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMongoDB(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var uri = "mongodb://localhost:27017"
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI))
	if err != nil {
		t.Fatalf("Connect mongodb failed: %v", err)
	}
	t.Log("Connect mongodb success")
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	collection := client.Database("test").Collection("testinfo")
	for i := 0; i < 10; i++ {
		res, err := collection.InsertOne(ctx, bson.M{
			"name":        strconv.Itoa(i),
			"score":       i * 10,
			"add_time":    time.Now().UnixMilli(),
			"update_time": time.Now().UnixMilli(),
		})
		if err != nil {
			t.Fatalf("Insert mongodb item failed: %v", err)
		}
		t.Logf("Inserted id %v", res.InsertedID)
	}

	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		t.Fatalf("Get iterater cursor failed: %v", err)
	}
	for cur.Next(ctx) {
		result := struct {
			Name  string
			Score int
		}{}
		if err := cur.Decode(&result); err != nil {
			t.Fatalf("decode cursor into struct failed: %v", err)
		}
		t.Logf("cursor %v -> %v", result, cur.Current.Lookup("_id"))
	}
	if err := cur.Err(); err != nil {
		t.Fatalf("cursor error: %v", err)
	}
	cur.Close(ctx)

	result := struct {
		Name       string
		Score      int
		AddTime    int64 `json:"add_time" bson:"add_time"`
		UpdateTime int64 `json:"update_time" bson:"update_time"`
	}{}
	if err := collection.FindOne(ctx, bson.M{"name": "6"}).Decode(&result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			t.Log("Find nothing in collection")
		} else {
			t.Fatalf("Find collection item failed: %v", err)
		}
	}
	t.Logf("Find collection item: %v", result)

	deleteRes, err := collection.DeleteOne(ctx, bson.M{"name": "6"})
	if err != nil {
		t.Fatalf("Delete collection item failed: %v", err)
	}
	t.Logf("Delete DeletedCount %v", deleteRes.DeletedCount)

	cur, err = collection.Find(ctx, bson.D{})
	if err != nil {
		t.Fatalf("Get iterater cursor failed: %v", err)
	}
	for cur.Next(ctx) {
		result := struct {
			Name       string
			Score      int
			Grade      int
			AddTime    int64 `json:"add_time" bson:"add_time"`
			UpdateTime int64 `json:"update_time" bson:"update_time"`
		}{}
		if err := cur.Decode(&result); err != nil {
			t.Fatalf("decode cursor into struct failed: %v", err)
		}
		t.Logf("cursor %v -> %v", result, cur.Current.Lookup("_id"))
	}
	if err := cur.Err(); err != nil {
		t.Fatalf("cursor error: %v", err)
	}
	cur.Close(ctx)

	updateRes, err := collection.UpdateOne(ctx, bson.M{"name": "6"}, bson.M{
		"$set": bson.M{
			"score":       6,
			"grade":       12,
			"update_time": time.Now().UnixMilli(),
		},
	}, options.Update().SetUpsert(true))
	if err != nil {
		t.Fatalf("Update collection item failed: %v", err)
	}
	t.Logf("Update MatchedCount %v ModifiedCount %v, UpsertedCount %v",
		updateRes.MatchedCount, updateRes.ModifiedCount, updateRes.UpsertedCount)

	cur, err = collection.Find(ctx, bson.D{})
	if err != nil {
		t.Fatalf("Get iterater cursor failed: %v", err)
	}
	for cur.Next(ctx) {
		result := struct {
			Name  string
			Score int
			Grade int
		}{}
		if err := cur.Decode(&result); err != nil {
			t.Fatalf("decode cursor into struct failed: %v", err)
		}
		t.Logf("cursor %v -> %v", result, cur.Current.String())
	}
	if err := cur.Err(); err != nil {
		t.Fatalf("cursor error: %v", err)
	}
	cur.Close(ctx)

	time.Sleep(2 * time.Second)

	if err := collection.Drop(ctx); err != nil {
		t.Fatalf("Drop test database failed: %v", err)
	}
	t.Log("Drop test database success")
}
