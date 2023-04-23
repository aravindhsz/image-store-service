package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image-store-service/pkg/album"
	"image-store-service/pkg/image"
	"io"
	"strings"

	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

//import "time"

type Server struct {
	server   *http.Server
	dbclient *mongo.Client
}

func newServer(address string, client *mongo.Client) *Server {
	return &Server{
		server: &http.Server{
			Addr: address,
		},
		dbclient: client,
	}
}

func Start(client *mongo.Client) error {
	s := newServer(":8080", client)
	fmt.Println("starting server")
	router := mux.NewRouter()
	router.HandleFunc("/api/album", s.createAlbum).Methods(http.MethodPost)
	router.HandleFunc("/api/album", s.getAlbum).Methods(http.MethodGet)
	router.HandleFunc("/api/album/{id}", s.deleteAlbum).Methods(http.MethodDelete)
	router.HandleFunc("/api/album/{id}/image", s.getAlbumDetails).Methods(http.MethodGet)
	router.HandleFunc("/api/album/{id}/image", s.createImage).Methods(http.MethodPost)
	router.HandleFunc("/api/album/{id}/image/{imageid}", s.getImage).Methods(http.MethodGet)
	router.HandleFunc("/api/album/{id}/image/{imageid}", s.deleteImage).Methods(http.MethodDelete)
	s.server.Handler = router
	err := s.server.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) createAlbum(w http.ResponseWriter, r *http.Request) {
	var alb album.Album

	err := json.NewDecoder(r.Body).Decode(&alb)
	if err != nil {
		io.WriteString(w, "unable to create album")
	}

	alb.Id = uuid.New().String()
	alb.Image = []image.Image{}
	err = album.CreateAlbum(s.dbclient, &alb)
	type response struct {
		Id  string `json:"id"`
		Err string `json:"err"`
	}
	var resp response
	if err != nil {
		resp.Id = ""
		resp.Err = err.Error()

	} else {
		resp.Id = alb.Id
		resp.Err = ""
	}

	res, err := json.Marshal(resp)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(res)
}

func (s *Server) getAlbum(w http.ResponseWriter, r *http.Request) {

	result, err := album.GetAlbums(s.dbclient)
	type response struct {
		Albums []album.Album `json:"albums"`
		Err    string        `json:"err"`
	}
	var resp response
	if err != nil {
		resp.Err = err.Error()
	} else {
		resp.Albums = result
		resp.Err = ""
	}
	res, err := json.Marshal(resp)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(res)
}

func (s *Server) getAlbumDetails(w http.ResponseWriter, r *http.Request) {
	arr := strings.Split(r.URL.Path, "/")
	id := arr[3]
	result, err := album.GetAlbumById(s.dbclient, id)
	type response struct {
		Album album.Album `json:"album"`
		Err   string      `json:"err"`
	}
	var resp response
	if err != nil {
		resp.Err = err.Error()
	} else {
		resp.Album = *result
		resp.Err = ""
	}
	res, err := json.Marshal(resp)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(res)

}

func (s *Server) deleteAlbum(w http.ResponseWriter, r *http.Request) {

	id := strings.TrimPrefix(r.URL.Path, "/api/album/")
	err := album.DeleteAlbum(s.dbclient, id)
	type response struct {
		Id  string `json:"id"`
		Err string `json:"err"`
	}
	var resp response
	if err != nil {
		resp.Id = ""
		resp.Err = err.Error()

	} else {
		resp.Id = id
		resp.Err = ""
	}

	res, err := json.Marshal(resp)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(res)
}

func (s *Server) createImage(w http.ResponseWriter, r *http.Request) {
	var img image.Image

	r.ParseForm()
	file, _, err := r.FormFile("data")

	defer file.Close()
	if err != nil {
		fmt.Println("err is", err)
	}

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		fmt.Println("err is", err)
	}

	img.Name = r.FormValue("name")
	img.Data = buf.Bytes()

	arr := strings.Split(r.URL.Path, "/")
	albumid := arr[3]
	imgId, err := album.CreateImage(s.dbclient, albumid, &img)
	type response struct {
		Id  string `json:"id"`
		Err string `json:"err"`
	}
	var resp response
	if err != nil {
		resp.Id = ""
		resp.Err = err.Error()

	} else {
		resp.Id = imgId
		resp.Err = ""
	}

	res, err := json.Marshal(resp)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(res)
}

func (s *Server) deleteImage(w http.ResponseWriter, r *http.Request) {
	arr := strings.Split(r.URL.Path, "/")
	albumid := arr[3]
	imageid := arr[5]

	err := album.DeleteImage(s.dbclient, albumid, imageid)
	type response struct {
		Id  string `json:"id"`
		Err string `json:"err"`
	}
	var resp response
	if err != nil {
		resp.Id = ""
		resp.Err = err.Error()

	} else {
		resp.Id = imageid
		resp.Err = ""
	}
	res, err := json.Marshal(resp)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(res)
}

func (s *Server) getImage(w http.ResponseWriter, r *http.Request) {

	arr := strings.Split(r.URL.Path, "/")
	albumid := arr[3]
	imageid := arr[5]
	result, err := album.GetImage(s.dbclient, albumid, imageid)
	type response struct {
		Image image.Image `json:"image"`
		Err   string      `json:"err"`
	}
	var resp response
	if err != nil {
		resp.Err = err.Error()
	} else {
		resp.Image = *result
		resp.Err = ""
	}
	res, err := json.Marshal(resp)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(res)
}
