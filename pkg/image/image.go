package image

type Image struct {
	Id   string `bson:"id"`
	Name string `bson:"name"`
	Data []byte `bson:"data"`
}
