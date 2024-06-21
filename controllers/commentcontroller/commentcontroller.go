package commentcontroller

import (
	"encoding/json"
	"net/http"

	"github.com/althaafka/alk-proj-be.git/database"
	"github.com/althaafka/alk-proj-be.git/helpers"
	"github.com/althaafka/alk-proj-be.git/models"
)

func CreateComment(w http.ResponseWriter, r *http.Request) {
    if !helpers.ValidateMethod(w, r, "POST") {
        return
    }

    userID, ok := r.Context().Value("userID").(uint)
    if !ok {
        helpers.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
        return
    }

    var comment models.Comment
    if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
        helpers.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
        return
    }
    defer r.Body.Close()

    comment.UserID = userID
    if err := database.DB.Create(&comment).Error; err != nil {
        helpers.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Server error"})
        return
    }

    helpers.RespondWithJSON(w, http.StatusCreated, comment)
}

func GetComments(w http.ResponseWriter, r *http.Request) {
    if !helpers.ValidateMethod(w, r, "GET") {
        return
    }

    articleID := r.URL.Query().Get("article_id")
    var comments []models.Comment
    if err := database.DB.Joins("User").Where("article_id = ?", articleID).Order("created_at ASC").Find(&comments).Error; err != nil {
        helpers.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Server error"})
        return
    }

	result := make([]map[string]interface{}, len(comments))
    for i, comment := range comments {
        result[i] = map[string]interface{}{
            "id":        comment.ID,
            "content":   comment.Content,
            "username":  comment.User.Username,
            "article_id": comment.ArticleID,
            "user_id":   comment.UserID,
			"created_at": comment.CreatedAt.Format("2006-01-02 15:04:05"),
        }
    }

    helpers.RespondWithJSON(w, http.StatusOK, result)
}