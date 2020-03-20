//INCOMPLETE PROGRAM

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Albums struct {
	Name string `json:"name"`
	Image []Image `json:"image"`
}

type Image struct {
	Name string `json:"name"`
}

var albums []Albums
var Session *gocql.Session

func init() {
	var err error
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "albumdb"
	Session, err = cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	fmt.Println("cassandra init done")
}

//OK
//Show all albums
func showAlbum(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Displaying album names:\n")

	iter := Session.Query(`SELECT albumname FROM albumtable`).Iter()
	var albumname string
	for iter.Scan(&albumname) {
		fmt.Println("Album Name:", albumname,)
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}
}

//OK
//Create a new album
func addAlbum(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")
	param := mux.Vars(r)
	fmt.Fprintf(w, "Creating Album",param["album"],"\n")

	if err := Session.Query("INSERT INTO albumtable(albumname) VALUES(?)", param["album"]).Exec(); err != nil {
		fmt.Println("Error while inserting")
		fmt.Println(err)
	}
}

//OK
//Delete an existing album
func deleteAlbum(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	param := mux.Vars(r)

	if err := Session.Query("DELETE FROM albumtable WHERE albumname=?",param["album"]).Exec(); err != nil {
		fmt.Println("Error while deleting")
		fmt.Println(err)
	}
}

//OK
//Show all images in an album
func showImagesInAlbum(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	param := mux.Vars(r)
	str := param["album"]

	for i:=1; i<=5; i++ {
		img := "image"
		query := fmt.Sprintf(`SELECT %s%d FROM albumtable WHERE albumname='%s';`,img,i,str)
		iter := Session.Query(query).Iter()
		var imagename string
		for iter.Scan(&imagename) {
			fmt.Println("Image Name:", imagename)
		}
		if err := iter.Close(); err != nil {
			log.Fatal(err)
		}
	}
}

//TODO
//Show a particular image inside an album
func showImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	param := mux.Vars(r)
	for idx:=0; idx < len(albums);idx++ {
		if albums[idx].Name == param["album"] {
			for i:=0;i<len(albums[idx].Image);i++ {
				if albums[idx].Image[i].Name == param["image"] {
					fmt.Fprintf(w, "Displaying",param["image"],"in album", param["album"],"\n")
					json.NewEncoder(w).Encode(albums[idx].Image)
					return
				}
			}
		}
	}
	fmt.Fprintf(w, "ERROR:",param["image"],"image does not exist in album",param["album"],"\n")
}

//TODO
//Create an image in an album
func addImage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Add Image")
	w.Header().Set("Content-Type","application/json")
	param := mux.Vars(r)

	str := param["album"]
	for i:=1; i<=5; i++ {
		img := "image"
		query := fmt.Sprintf(`SELECT %s%d FROM albumtable WHERE albumname='%s';`,img,i,str)
		iter := Session.Query(query).Iter()
		var imagename string
		for iter.Scan(&imagename) {
			if imagename == "null" {
				if err := Session.Query("INSERT INTO albumtable(albumname,imagename) VALUES(?,?)",param["album"],param["image"]).Exec(); err != nil {
					fmt.Println("Error while inserting")
					fmt.Println(err)
				}
			}
			//fmt.Println("Image Name:", imagename)
		}
		if err := iter.Close(); err != nil {
			log.Fatal(err)
		}
	}



	image := Image{Name: param["image"]}
	for idx,item := range albums {
		if item.Name == param["album"] {
			albums[idx].Image = append(albums[idx].Image, image)
			json.NewEncoder(w).Encode(albums)
			return
		}
	}
	fmt.Fprintf(w, "ERROR:",param["album"],"album does not exist. Hence, image",param["image"],"cannot be added.")
}

//TODO
//Delete an image in an album
func deleteImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	param := mux.Vars(r)
	var alb []Albums
	for idx, item := range albums {
		if item.Name == param["album"] {
			for i:=0;i<len(item.Image);i++ {
				if item.Image[i].Name == param["image"] {
					fmt.Fprintf(w, "Deleting",param["image"],"in album", param["album"],"\n")
					item.Image = append(item.Image[:i],item.Image[i+1:]...)
					alb = append(albums[:idx], Albums{Name: param["album"], Image: item.Image})
					albums = append(alb, albums[idx+1:]...)
					break
				}
			}
		}
	}
	json.NewEncoder(w).Encode(albums)
}

func main() {
	//Initialize Router
	myRouter := mux.NewRouter().StrictSlash(true)

	//Sample Data
	/*
	alb1 := Albums{Name: "car", Image:[]Image{{Name: "amaze"},{Name: "ciaz"}}}
	alb2 := Albums{Name: "car", Image:[]Image{{Name: "amaze"},{Name: "ciaz"}}}
	alb3 := Albums{Name: "car", Image:[]Image{{Name: "amaze"},{Name: "ciaz"}}}
	alb4 := Albums{Name: "car", Image:[]Image{{Name: "amaze"},{Name: "ciaz"}}}
	*/
	/*
	albums = append(albums, Albums{Name: "car", Image: []Image{{Name: "amaze"},{Name: "ciaz"}}})
	albums = append(albums, Albums{Name: "bike", Image: []Image{{Name: "apache"}}})
	albums = append(albums, Albums{Name: "mountain", Image: []Image{{Name: "everest"}}})
	albums = append(albums, Albums{Name: "ocean", Image: []Image{{Name: "pacific"}}})
	*/

	//Show all albums
	myRouter.HandleFunc("/",showAlbum).Methods(http.MethodGet)
	//Create a new album
	myRouter.HandleFunc("/{album}",addAlbum).Methods(http.MethodPost)
	//Delete an existing album
	myRouter.HandleFunc("/{album}",deleteAlbum).Methods(http.MethodDelete)

	//Show all images in an album
	myRouter.HandleFunc("/{album}",showImagesInAlbum).Methods(http.MethodGet)
	//Show a particular image inside an album
	myRouter.HandleFunc("/{album}/{image}",showImage).Methods(http.MethodGet)
	//Create an image in an album
	myRouter.HandleFunc("/{album}/{image}",addImage).Methods(http.MethodPost)
	//Delete an image in an album
	myRouter.HandleFunc("/{album}/{image}",deleteImage).Methods(http.MethodDelete)
	log.Fatal(http.ListenAndServe(":8085",myRouter))
}
