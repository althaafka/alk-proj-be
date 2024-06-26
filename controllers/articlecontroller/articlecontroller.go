package articlecontroller

import (
	"encoding/json"
	"net/http"
	"io"
	"os"
	"path/filepath"

	"github.com/althaafka/alk-proj-be.git/database"
	"github.com/althaafka/alk-proj-be.git/helpers"
	"github.com/althaafka/alk-proj-be.git/models"
)

func CreateArticle(w http.ResponseWriter, r *http.Request) {
	if !helpers.ValidateMethod(w, r, "POST") {
		return
	}

	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		helpers.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		return
	}

	var article models.Article
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		helpers.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}
	defer r.Body.Close()

	article.UserID = userID

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
	if err := database.DB.Order("ID desc").Find(&articles).Error; err != nil {
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
			"likes": article.Likes,
			"photo_url": article.PhotoURL,
		}
	}

	helpers.RespondWithJSON(w, http.StatusOK, result)
}

func GetMyArticles(w http.ResponseWriter, r *http.Request) {
	if !helpers.ValidateMethod(w, r, "GET") {
		return
	}

	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		helpers.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		return
	}

	var articles []models.Article
	if err := database.DB.Where("user_id = ?", userID).Order("ID desc").Find(&articles).Error; err != nil {
		helpers.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Server error"})
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, articles)
}

func GetArticle(w http.ResponseWriter, r *http.Request) {
	if !helpers.ValidateMethod(w, r, "GET") {
		return
	}

	id := r.URL.Query().Get("id")
	var article models.Article
	if err := database.DB.Where("id = ?", id).First(&article).Error; err != nil {
		helpers.RespondWithJSON(w, http.StatusNotFound, map[string]string{"error": "Article not found"})
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, article)
}

func EditArticle(w http.ResponseWriter, r *http.Request) {
	if !helpers.ValidateMethod(w, r, "PUT") {
		return
	}

	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
        helpers.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
        return
    }

	var article models.Article
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		helpers.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}

	defer r.Body.Close()

	var existingArticle models.Article
    if err := database.DB.Where("id = ?", article.ID).First(&existingArticle).Error; err != nil {
        helpers.RespondWithJSON(w, http.StatusNotFound, map[string]string{"error": "Article not found"})
        return
    }

    if existingArticle.UserID != userID {
        helpers.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
        return
    }

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

	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		helpers.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
		return
	}

	id := r.URL.Query().Get("id")
	var article models.Article
	if err := database.DB.Where("id = ?", id).First(&article).Error; err != nil {
		helpers.RespondWithJSON(w, http.StatusNotFound, map[string]string{"error": "Article not found"})
		return
	}

	if article.UserID != userID {
		helpers.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}
	
	if err := database.DB.Delete(&article).Error; err != nil {
		helpers.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Server error"})
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Article deleted"})
}

func AddLikeToArticle(w http.ResponseWriter, r *http.Request) {
	if !helpers.ValidateMethod(w, r, "POST") {
		return
	}

	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		helpers.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}
	defer r.Body.Close()

	articleID, ok := data["article_id"].(float64)
	if !ok {
		helpers.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid article ID"})
		return
	}

	var article models.Article
	if err := database.DB.First(&article, uint(articleID)).Error; err != nil {
		helpers.RespondWithJSON(w, http.StatusNotFound, map[string]string{"error": "Article not found"})
		return
	}

	article.Likes++

	if err := database.DB.Save(&article).Error; err != nil {
		helpers.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Server error"})
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, article)
}

func UploadArticlePhoto(w http.ResponseWriter, r*http.Request) {
	if !helpers.ValidateMethod(w,r, "POST"){
		return
	}

	err := r.ParseMultipartForm(10 << 20) //Max 10MB
	if err != nil {
		helpers.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Failes to parse multipart form"})
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		helpers.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Failed to get file"})
		return
	}
	defer file.Close()

	os.MkdirAll("./assets/images", os.ModePerm)

	filePath := filepath.Join("./assets/images/", handler.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		helpers.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create file"})
	}
	defer dst.Close()

	if _,err := io.Copy(dst, file); err != nil {
		helpers.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to save file"})
		return
	}

	photoUrl := "/assets/images/" + handler.Filename
	helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"photo_url": photoUrl})
} 