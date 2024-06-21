package router

import (
	"net/http"
	"log"
	"github.com/althaafka/alk-proj-be.git/controllers/usercontroller"
	"github.com/althaafka/alk-proj-be.git/controllers/articlecontroller"
	"github.com/althaafka/alk-proj-be.git/controllers/commentcontroller"
	"github.com/althaafka/alk-proj-be.git/middlewares"
)


func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}


func SetupRouter() {
	http.Handle("/user/register", corsMiddleware(http.HandlerFunc(usercontroller.Register)))
	http.Handle("/user/login", corsMiddleware(http.HandlerFunc(usercontroller.Login)))

	http.Handle("/article/all", corsMiddleware(http.HandlerFunc(articlecontroller.GetArticles)))
	http.Handle("/article/my", corsMiddleware(middlewares.AuthMiddleware(http.HandlerFunc(articlecontroller.GetMyArticles))))
	http.Handle("/article", corsMiddleware(middlewares.AuthMiddleware(http.HandlerFunc(articlecontroller.GetArticle))))
	http.Handle("/article/create", corsMiddleware(middlewares.AuthMiddleware(http.HandlerFunc(articlecontroller.CreateArticle))))
    http.Handle("/article/edit", corsMiddleware(middlewares.AuthMiddleware(http.HandlerFunc(articlecontroller.EditArticle))))
	http.Handle("/article/delete", corsMiddleware(middlewares.AuthMiddleware(http.HandlerFunc(articlecontroller.DeleteArticle))))
	http.Handle("/article/like", corsMiddleware(middlewares.AuthMiddleware(http.HandlerFunc(articlecontroller.AddLikeToArticle))))

	http.Handle("/comment/create", corsMiddleware(middlewares.AuthMiddleware(http.HandlerFunc(commentcontroller.CreateComment))))
	http.Handle("/comment", corsMiddleware(http.HandlerFunc(commentcontroller.GetComments)))

	log.Println("Server started on localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}