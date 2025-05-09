package db

import (
	"context"
	"errors"
	"time"

	"health-monitoring/log"
	"health-monitoring/types"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MDB *mongoDB = nil

type mongoDB struct {
	Mongo                  *mongo.Client
	deviceOnlineCollection *mongo.Collection
	deviceInfoCollection   *mongo.Collection
}

func InitMongo(ctx context.Context, uri, db string) error {
	opts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Log.Fatalf("Connect mongodb failed: %v", err)
	}
	if err = client.Ping(ctx, nil); err != nil {
		log.Log.Fatalf("Ping mongodb failed: %v", err)
	}
	MDB = &mongoDB{
		Mongo: client,
	}

	MDB.deviceOnlineCollection = client.Database(db).Collection("device_online")
	MDB.deviceInfoCollection = client.Database(db).Collection("device_info")
	return nil
}

func (db *mongoDB) Disconnect(ctx context.Context) {
	if err := db.Mongo.Disconnect(ctx); err != nil {
		panic(err)
	}
}

func (db *mongoDB) IsNodeOnline(ctx context.Context, nodeId string) bool {
	result := types.MDBDeviceOnline{}
	if err := db.deviceOnlineCollection.FindOne(ctx, bson.M{"device_id": nodeId}).Decode(&result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false
		}
		return false
	}
	return true
}

func (db *mongoDB) NodeOnline(ctx context.Context, nodeId string) error {
	res, err := db.deviceOnlineCollection.InsertOne(ctx, types.MDBDeviceOnline{
		DeviceId: nodeId,
		AddTime:  time.Now().UnixMilli(),
	})
	if err != nil {
		log.Log.WithFields(logrus.Fields{"node_id": nodeId}).Error("insert online failed:", err)
		return err
	}
	log.Log.WithFields(logrus.Fields{"node_id": nodeId}).Info("inserted online id", res.InsertedID)
	return nil
}

func (db *mongoDB) NodeOffline(ctx context.Context, nodeId string) error {
	result, err := db.deviceOnlineCollection.DeleteOne(ctx, bson.M{"device_id": nodeId})
	if err != nil {
		log.Log.WithFields(logrus.Fields{"node_id": nodeId}).Error("delete online failed:", err)
		return err
	}
	log.Log.WithFields(logrus.Fields{"node_id": nodeId}).Info("delete online count", result.DeletedCount)
	return nil
}

func (db *mongoDB) GetDeviceInfo(ctx context.Context, nodeId string) (*types.MDBDeviceInfo, error) {
	result := &types.MDBDeviceInfo{}
	if err := db.deviceInfoCollection.FindOne(ctx, bson.M{"device_id": nodeId}).Decode(result); err != nil {
		return nil, err
	}
	return result, nil
}

func (db *mongoDB) AddDeviceInfo(ctx context.Context, nodeId string, info types.WsMachineInfoRequest) error {
	result, err := db.deviceInfoCollection.InsertOne(
		ctx,
		types.MDBDeviceInfo{
			DeviceId:             nodeId,
			WsMachineInfoRequest: info,
			AddTime:              time.Now().UnixMilli(),
			UpdateTime:           time.Now().UnixMilli(),
		},
	)
	if err != nil {
		log.Log.WithFields(logrus.Fields{"node_id": nodeId}).Error("insert device info failed:", err)
		return err
	}
	log.Log.WithFields(logrus.Fields{"node_id": nodeId}).Info("inserted device info id", result.InsertedID)
	return nil
}

func (db *mongoDB) UpdateDeviceInfo(ctx context.Context, nodeId string, info types.WsMachineInfoRequest) error {
	result, err := db.deviceInfoCollection.UpdateOne(
		ctx,
		bson.M{"device_id": nodeId},
		bson.M{
			"$set": types.MDBDeviceInfo{
				WsMachineInfoRequest: info,
				UpdateTime:           time.Now().UnixMilli(),
			},
		},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		log.Log.WithFields(logrus.Fields{"node_id": nodeId}).Error("update device info failed:", err)
		return err
	}
	log.Log.WithFields(logrus.Fields{"node_id": nodeId}).Info("update device info count",
		result.MatchedCount, result.ModifiedCount, result.UpsertedCount)
	return nil
}
