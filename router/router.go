package router

import (
	"net/http"
	"log"
	"github.com/althaafka/alk-proj-be.git/controllers/usercontroller"
	"github.com/althaafka/alk-proj-be.git/controllers/articlecontroller"
	"github.com/althaafka/alk-proj-be.git/middlewares"
)


func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}


func SetupRouter() {
	http.Handle("/users/register", corsMiddleware(http.HandlerFunc(usercontroller.Register)))
	http.Handle("/users/login", corsMiddleware(http.HandlerFunc(usercontroller.Login)))

	http.Handle("/articles", corsMiddleware(http.HandlerFunc(articlecontroller.GetArticles)))
	http.Handle("/articles/create", corsMiddleware(middlewares.AuthMiddleware(http.HandlerFunc(articlecontroller.CreateArticle))))
    http.Handle("/articles/edit", corsMiddleware(middlewares.AuthMiddleware(http.HandlerFunc(articlecontroller.EditArticle))))
    http.Handle("/articles/delete", corsMiddleware(middlewares.AuthMiddleware(http.HandlerFunc(articlecontroller.DeleteArticle))))

	log.Println("Server started on localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}