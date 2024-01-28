package datastore

import (
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"sync"
	"time"

	"github.com/s-li1/remarkable-screen-share/internal"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Page struct {
  ImageUri string `bson:"image_uri"`
}

// func download(imageBufPool *sync.Pool) error {
//   imgBuffer := imageBufPool.Get().(image.Gray)
// 
//   timestamp := time.Now().Unix()
//   rand, err := internal.GenerateRandomHex(10)
//   if err != nil {
//     return err
//   }
// 
//   file, err := os.Create(fmt.Sprintf("%d-%s.jpg", timestamp, rand))
//   if err != nil {
//     return err
//   }
//   defer file.Close()
// 
//   mongoDbUri := os.Getenv("MONGODB_URI")
// 
//   mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoDbUri))
//   if err != nil {
//     log.Panicf("Something went wrong with connecting to DB: %v", err)
// 	} 
// 
//   defer func() {
//     if err = mongoClient.Disconnect(context.TODO()); err != nil {
//       log.Panicf("Something went wrong with disconnecting to DB: %v", err)
//     }
//   }()
// 
//   newPage := Page{ImageUri: fmt.Sprintf("file:///Users/steven.li/remarkable-screen-share/%s", file.Name())}
//   collection := mongoClient.Database("remarkable").Collection("pages")
//   if _, err := collection.InsertOne(context.TODO(), newPage); err != nil { 
//     log.Panicf("Something went wrong with inserting into DB: %v", err)
//   }
// 
//   return jpeg.Encode(file, &imgBuffer, nil)
// }
