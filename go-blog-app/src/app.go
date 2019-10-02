package main

import (
	"context"
	"fmt"
	"html/template"

	//"encoding/==json"

	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Blog struct {
	title string `bson:"title" json:"title"`

	image_url string `bson:image_url json:"image_url"`

	description string `bson : description json:"description"`

	contact_detail string `bson: contact_detail json:"contact_detail"`
}

func blogHandler(res http.ResponseWriter, r *http.Request) {
	p := Blog{title: "check out your favourite blogs"}
	t, _ := template.ParseFiles("../views/index.htm")
	t.Execute(res, p)

}
func create(res http.ResponseWriter, req *http.Request) {
	//	p := Blog{title: "check out your favourite blogs"}
	t, _ := template.ParseFiles("../views/create.htm")
	t.Execute(res, nil)
}
func delete(res http.ResponseWriter, req *http.Request) {
	//	p := Blog{title: "check out your favourite blogs"}
	t, _ := template.ParseFiles("../views/delete.htm")
	t.Execute(res, nil)
}
func createdBlog(res http.ResponseWriter, req *http.Request) {

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)

	collection := client.Database("blog").Collection("blogs")
	req.ParseForm()

	insertResult, err := collection.InsertOne(context.TODO(), bson.M{
		"title":          req.Form["$title"],
		"image_url":      req.Form["$image_url"],
		"description":    req.Form["$description"],
		"contact_detail": req.Form["$contact_detail"],
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult)
	t, _ := template.ParseFiles("../views/created.htm")

	t.Execute(res, nil)

}
func updateBlog(res http.ResponseWriter, req *http.Request) {

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)

	collection := client.Database("blog").Collection("blogs")

	filter := bson.M{{ "title" , req.Form["$title"] }}

	req.ParseForm()

	update := bson.M{
		"$set": bson.M{
			"image_url":      req.Form["$image_url"],
			"description":    req.Form["$description"],
			"contact_detail": req.Form["$contact_detail"],
		},
	}
	
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	t, _ := template.ParseFiles("../views/updated.htm")
	t.Execute(res, nil)
}
func deleteBlog(res http.ResponseWriter, req *http.Request) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database("blog").Collection("blogs")
	//	var titl string
	req.ParseForm()
	titl := req.Form["title"]
	filter := bson.M{"title": titl}

	result, err1 := collection.DeleteMany(context.TODO(), filter)
	if err1 != nil {
		log.Fatal(err1)
	}
	fmt.Println("succesfully deleted", result)
	p := Blog{title: "your blog with title is deleted"}
	t, _ := template.ParseFiles("../views/deleted.htm")
	t.Execute(res, p)

}
func findBlogs(res http.ResponseWriter, req *http.Request) {

	findOptions := options.Find()
	//findOptions.SetLimit(2)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)

	collection := client.Database("blog").Collection("blogs")

	var results []*Blog
	filter := bson.M{}
	cur, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {

		var elem Blog
		fmt.Println(cur)
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		//	fmt.Println("title is", elem.title, "description is ", elem.description)
		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.TODO())

	var c1=struct{
		result[] Blog 
	}{
		result:results
	}	
	t, _ := template.ParseFiles("../views/find.htm")

	t1 := t.Lookup("find.htm")

	t1.Execute(res, c1)
}

func UpdateYourBlog(res http.ResponseWriter, req *http.Request) {
	//	p := Blog{title: "check out your favourite blogs"}
	t, _ := template.ParseFiles("../views/update.htm")
	t.Execute(res, nil)
}
func main() {

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	http.HandleFunc("/blogs", findBlogs)

	http.HandleFunc("/", blogHandler)
	http.HandleFunc("/create", create)
	http.HandleFunc("/created", createdBlog)

	http.HandleFunc("/delete", delete)

	http.HandleFunc("/deleted", deleteBlog)

	http.HandleFunc("/update", UpdateYourBlog)

	http.HandleFunc("/updated", updateBlog)

	log.Fatal(http.ListenAndServe(":3000", nil))
}
