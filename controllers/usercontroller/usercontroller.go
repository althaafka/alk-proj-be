package usercontroller

import (
	"encoding/json"
	"net/http"

	"github.com/althaafka/alk-proj-be.git/database"
	"github.com/althaafka/alk-proj-be.git/helpers"
	"github.com/althaafka/alk-proj-be.git/models"
	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request){
	if !helpers.ValidateMethod(w, r, "POST") {
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		helpers.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}
	defer r.Body.Close()

	var existingUser models.User
	if err := database.DB.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		helpers.RespondWithJSON(w, http.StatusConflict, map[string]string{"error": "User already exists"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	if err := database.DB.Create(&user).Error; err != nil {
		helpers.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Server error"})
		return
	}

	helpers.RespondWithJSON(w, http.StatusCreated, user)
}

func Login(w http.ResponseWriter, r *http.Request){
	if !helpers.ValidateMethod(w, r, "POST") {
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		helpers.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}
	defer r.Body.Close()

	var existingUser models.User
	if err := database.DB.Where("username = ?", user.Username).First(&existingUser).Error; err != nil {
		helpers.RespondWithJSON(w, http.StatusNotFound, map[string]string{"error": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password)); err != nil {
		helpers.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
		return
	}

	token, err := helpers.GenerateToken(existingUser.ID)
	if err != nil {
		helpers.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Server error"})
		return
	}

	helpers.RespondWithJSON(w, http.StatusOK, map[string]string{"username": existingUser.Username, "token": token})
}
