package album

import (
	"context"
	"errors"
	"fmt"
	"image-store-service/pkg/image"
	"log"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Album struct {
	Id          string        `bson:"id"`
	Name        string        `bson:"name"`
	Description string        `bson:"description"`
	Image       []image.Image `bson:"image"`
}

type ReadAlbum struct {
	AlbumName string `bson:"albumName"`
	ImageName string `bson:"imageName"`
}

func CreateAlbum(client *mongo.Client, album *Album) error {

	albumsCollection := client.Database("imageStore").Collection("albums")

	filter := bson.D{
		{"name", album.Name},
	}

	// retrieve all the documents that match the filter
	res := albumsCollection.FindOne(context.TODO(), filter)

	if res.Err() == nil {
		return errors.New("Album name already exists")
	}
	img := []image.Image{}
	albumrecord := bson.D{{Key: "id", Value: album.Id}, {Key: "name", Value: album.Name}, {Key: "description", Value: album.Description}, {Key: "image", Value: img}}

	_, err := albumsCollection.InsertOne(context.TODO(), albumrecord)
	if err != nil {
		return err
	}
	return nil

}

func GetAlbums(client *mongo.Client) ([]Album, error) {

	albumsCollection := client.Database("imageStore").Collection("albums")

	filter := bson.D{}

	// retrieve all the documents that match the filter
	cur, err := albumsCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var results []Album
	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var record Album
		err := cur.Decode(&record)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, record)

	}

	return results, nil

}

func GetAlbumById(client *mongo.Client, id string) (*Album, error) {

	albumsCollection := client.Database("imageStore").Collection("albums")

	filter := bson.D{
		{"id", id},
	}

	// retrieve all the documents that match the filter
	record := albumsCollection.FindOne(context.TODO(), filter)
	if record.Err() != nil {
		return nil, errors.New("Album name doesn't exists")
	}
	var alb Album
	err := record.Decode(&alb)
	if err != nil {
		log.Fatal(err)
	}

	return &alb, nil

}

func DeleteAlbum(client *mongo.Client, id string) error {

	albumsCollection := client.Database("imageStore").Collection("albums")

	filter := bson.D{
		{"id", id},
	}

	// delete the first document that match the filter
	result, err := albumsCollection.DeleteOne(context.TODO(), filter)
	// check for errors in the deleting
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("Record not found")
	}
	return nil
}

func CreateImage(client *mongo.Client, albumId string, img *image.Image) (string, error) {

	albumsCollection := client.Database("imageStore").Collection("albums")

	filter := bson.D{
		{Key: "id", Value: albumId},
	}

	// retrieve all the documents that match the filter
	alb_res := albumsCollection.FindOne(context.TODO(), filter)

	if alb_res.Err() != nil {
		return "", errors.New("Album doesn't exists")
	}
	var doc Album
	err := alb_res.Decode(&doc)
	emptyArray := []image.Image{}
	if doc.Image != nil && len(doc.Image) != 0 {
		emptyArray = doc.Image
	}

	for _, v := range emptyArray {
		if v.Name == img.Name {
			return "", errors.New("Image name already exists in given album")
		}
	}
	img.Id = uuid.New().String()
	emptyArray = append(emptyArray, *img)
	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "id", Value: albumId},
				{Key: "image", Value: emptyArray},
			},
		},
	}
	_, err = albumsCollection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return "", err
	}
	return img.Id, nil
}

func DeleteImage(client *mongo.Client, albumId string, imageId string) error {

	albumsCollection := client.Database("imageStore").Collection("albums")

	filter := bson.D{
		{Key: "id", Value: albumId},
	}

	// retrieve all the documents that match the filter
	alb_res := albumsCollection.FindOne(context.TODO(), filter)

	if alb_res.Err() != nil {
		return errors.New("Album doesn't exists")
	}
	var doc Album
	err := alb_res.Decode(&doc)
	if err != nil {
		fmt.Println(err)
	}

	emptyArray := []image.Image{}
	if doc.Image != nil && len(doc.Image) != 0 {
		emptyArray = doc.Image
	}
	flag := false
	for i, v := range emptyArray {
		if v.Id == imageId {
			emptyArray = append(emptyArray[:i], emptyArray[i+1:]...)
			flag = true
			break
		}
	}
	if !flag {
		return errors.New("Image name doesn't exists in given album")
	}

	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "id", Value: albumId},
				{Key: "image", Value: emptyArray},
			},
		},
	}
	_, err = albumsCollection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return err
	}
	return nil

}

func GetImage(client *mongo.Client, albumId string, imageId string) (*image.Image, error) {

	albumsCollection := client.Database("imageStore").Collection("albums")
	filter := bson.D{
		{Key: "id", Value: albumId},
	}
	// retrieve all the documents that match the filter
	alb_res := albumsCollection.FindOne(context.TODO(), filter)

	if alb_res.Err() != nil {
		return nil, errors.New("Album doesn't exists")
	}
	var doc Album
	err := alb_res.Decode(&doc)
	if err != nil {
		fmt.Println(err)
	}

	emptyArray := []image.Image{}
	if doc.Image != nil && len(doc.Image) != 0 {
		emptyArray = doc.Image
	}
	flag := false
	var response image.Image
	for _, v := range emptyArray {
		if v.Id == imageId {
			response = v
			flag = true
			break
		}

	}
	if !flag {
		return nil, errors.New("Image name doesn't exists in given album")
	}
	return &response, nil

}

func GetAllImages(client *mongo.Client, albumName string) error {

	albumsCollection := client.Database("imageStore").Collection("albums")

	filter := bson.D{
		{Key: "name", Value: albumName},
	}

	// retrieve all the documents that match the filter
	alb_res := albumsCollection.FindOne(context.TODO(), filter)

	if alb_res.Err() != nil {
		return errors.New("Album doesn't exists")
	}
	var doc Album
	err := alb_res.Decode(&doc)
	if err != nil {
		fmt.Println(err)
	}

	emptyArray := []image.Image{}
	if doc.Image != nil && len(doc.Image) != 0 {
		emptyArray = doc.Image
	}

	if len(emptyArray) == 0 {
		return errors.New("Album doesn't have any images")
	}
	return nil

}
