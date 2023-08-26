package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jcesardev/go-gorm-restapi/db"
	"github.com/jcesardev/go-gorm-restapi/models"
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	var users = []models.User{}
	db.DB.Find(&users)
	json.NewEncoder(w).Encode(&users)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := params["id"]
	var user = models.User{}
	db.DB.First(&user, userId)
	if user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("User with ID %s not found", userId)))
		return
	}
	db.DB.Model(&user).Association("Tasks").Find(&user.Tasks)
	json.NewEncoder(w).Encode(&user)
}

func PostUserHandler(w http.ResponseWriter, r *http.Request) {
	var user = models.User{}
	json.NewDecoder(r.Body).Decode(&user)
	createdUser := db.DB.Create(&user)
	if createdUser.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(createdUser.Error.Error()))
	}
	json.NewEncoder(w).Encode(&user)
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	var params = mux.Vars(r)
	var user = models.User{}
	db.DB.First(&user, params["id"])
	if user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("User with ID %s not found", params["id"])))
		return
	}
	db.DB.Delete(&user) // Soft delete
	// db.DB.Unscoped().Delete(&user) // Delete user permanently
	w.WriteHeader(http.StatusOK)
}
