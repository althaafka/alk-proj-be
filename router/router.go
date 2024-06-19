package router

import (
	"net/http"
	"log"
	"github.com/althaafka/alk-proj-be.git/controllers/usercontroller"
	"github.com/althaafka/alk-proj-be.git/controllers/articlecontroller"
	"github.com/althaafka/alk-proj-be.git/middlewares"
)

func SetupRouter() {
	http.HandleFunc("/users/register", usercontroller.Register)
	http.HandleFunc("/users/login", usercontroller.Login)

	http.HandleFunc("/articles", articlecontroller.GetArticles)
	http.Handle("/articles/create", middlewares.AuthMiddleware(http.HandlerFunc(articlecontroller.CreateArticle)))
    http.Handle("/articles/edit", middlewares.AuthMiddleware(http.HandlerFunc(articlecontroller.EditArticle)))
    http.Handle("/articles/delete", middlewares.AuthMiddleware(http.HandlerFunc(articlecontroller.DeleteArticle)))

	log.Println("Server started at localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}