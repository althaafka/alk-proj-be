package router

import (
	"net/http"
	"log"
	"github.com/althaafka/alk-proj-be.git/controllers/usercontroller"
	"github.com/althaafka/alk-proj-be.git/controllers/articlecontroller"
)

func SetupRouter() {
	http.HandleFunc("/users/register", usercontroller.Register)
	http.HandleFunc("/users/login", usercontroller.Login)

	http.HandleFunc("/articles/create", articlecontroller.CreateArticle)
	http.HandleFunc("/articles", articlecontroller.GetArticles)
	http.HandleFunc("/articles/edit", articlecontroller.EditArticle)
	http.HandleFunc("/articles/delete", articlecontroller.DeleteArticle)

	log.Println("Server started at localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}