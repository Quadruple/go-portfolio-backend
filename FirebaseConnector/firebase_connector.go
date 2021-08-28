package firebase_connector

import (
	"context"
	"io/ioutil"
	"log"
	"time"

	gstorage "cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/storage"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type image struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

func createFirebaseApp() (*firebase.App, error) {
	config := &firebase.Config{
		StorageBucket: "<your-bucket-name>",
	}
	opt := option.WithCredentialsFile("<path-to-firebase-key-json>")
	app, err := firebase.NewApp(context.Background(), config, opt)
	return app, err
}

func getStorage() (*storage.Client, error) {
	app, err := createFirebaseApp()
	if err != nil {
		log.Fatalln(err)
	}
	client, err := app.Storage(context.Background())
	return client, err
}

func getBucket() (*gstorage.BucketHandle, error) {
	client, err := getStorage()
	if err != nil {
		log.Fatalln(err)
	}
	bucket, err := client.DefaultBucket()
	return bucket, err
}

func getBucketAndAttributes() (*gstorage.BucketHandle, *gstorage.BucketAttrs, error) {
	bucket, err := getBucket()
	if err != nil {
		log.Fatalln(err)
	}
	attrs, err := bucket.Attrs(context.Background())
	return bucket, attrs, err
}

func obtainConfigs() (*jwt.Config, error) {
	key, err := ioutil.ReadFile("<path-to-firebase-key-json>")
	if err != nil {
		log.Fatalln(err)
	}
	cfg, err := google.JWTConfigFromJSON(key)
	if err != nil {
		log.Fatalln(err)
	}
	return cfg, err
}

func GetImages() []image {
	images := make([]image, 0)
	configs, err := obtainConfigs()
	if err != nil {
		log.Fatalln(err)
	}
	bucket, bucketAttrs, err := getBucketAndAttributes()
	if err != nil {
		log.Fatalln(err)
	}
	method := "GET"
	expires := time.Now().Add(time.Second * 60)
	it := bucket.Objects(context.Background(), nil)
	for {
		objectAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		url, err := gstorage.SignedURL(bucketAttrs.Name, objectAttrs.Name, &gstorage.SignedURLOptions{
			GoogleAccessID: configs.Email,
			PrivateKey:     configs.PrivateKey,
			Method:         method,
			Expires:        expires,
		})
		if err != nil {
			log.Fatalln(err)
		}
		images = append(images, image{
			Name: objectAttrs.Name,
			Link: url,
		})
	}
	return images
}
