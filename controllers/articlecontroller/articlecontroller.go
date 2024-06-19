package articlecontroller

import (
	"encoding/json"
	"net/http"

	"github.com/althaafka/alk-proj-be.git/database"
	"github.com/althaafka/alk-proj-be.git/helpers"
	"github.com/althaafka/alk-proj-be.git/models"
)

func CreateArticle(w http.ResponseWriter, r *http.Request) {
	if !helpers.ValidateMethod(w, r, "POST") {
		return
	}

	var article models.Article
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		helpers.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}

	defer r.Body.Close()

	if err := database.DB.Create(&article).Error; err != nil {
		helpers.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Server error"})
		return
	}

	helpers.RespondWithJSON(w, http.StatusCreated, article)
}

func GetArticles(w http.ResponseWriter, r *http.Request) {
	if !helpers.ValidateMethod(w, r, "GET") {
		return
	}

	var articles []models.Article
	if err := database.DB.Find(&articles).Error; err != nil {
        helpers.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Server error"})
        return
    }

	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		helpers.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Server error"})
		return
	}

	userMap := make(map[uint]models.User)
    for _, user := range users {
        userMap[user.ID] = user
    }

	result := make([]map[string]interface{}, len(articles))
	for i, article := range articles {
		user := userMap[article.UserID]
		result[i] = map[string]interface{}{
			"id": article.ID,
			"title": article.Title,
			"content": article.Content,
			"user_id": article.UserID,
			"username": user.Username,
			"email": user.Email,
		}
	}

	helpers.RespondWithJSON(w, http.StatusOK, result)
}

// func GetArticle(w http.ResponseWriter, r *http.Request) {
// 	if !helpers.ValidateMethod(w, r, "GET") {
// 		return
// 	}

// 	id := r.URL.Query().Get("id")
// 	var article models.Article
// 	if err := database.DB.Where("id = ?", id).First(&article).Error; err != nil {
// 		helpers.RespondWithJSON(w, http.StatusNotFound, map[string]string{"error": "Article not found"})
// 		return
// 	}

// 	helpers.RespondWithJSON(w, http.StatusOK, article)
// }

func EditArticle(w http.ResponseWriter, r *http.Request) {
	if !helpers.ValidateMethod(w, r, "PUT") {
		return
	}

	var article models.Article
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		helpers.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}

	defer r.Body.Close()

	if err := database.DB.Save(&article).Error; err != nil {
		helpers.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Server error"})
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, article)
}

func DeleteArticle(w http.ResponseWriter, r *http.Request) {
	if !helpers.ValidateMethod(w, r, "DELETE") {
		return
	}

	id := r.URL.Query().Get("id")
	var article models.Article
	if err := database.DB.Where("id = ?", id).First(&article).Error; err != nil {
		helpers.RespondWithJSON(w, http.StatusNotFound, map[string]string{"error": "Article not found"})
		return
	}

	if err := database.DB.Delete(&article).Error; err != nil {
		helpers.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Server error"})
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Article deleted"})
}