package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/Naveenravi07/go-api/internal/store"
	"github.com/Naveenravi07/go-api/internal/utils"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	UserStore store.UserStore
	logger    *log.Logger
}

func NewUserHandler(us store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{UserStore: us, logger: logger}
}

type reqisterUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Bio      string `json:"bio"`
	Password string `json:"password"`
}

func validateUser(user *reqisterUserRequest) error {
	if user.Username == "" {
		return errors.New(" username is required")
	}
	if user.Email == "" {
		return errors.New(" email is required")
	}
	if user.Password == "" {
		return errors.New("Passowrd is required")
	}
	return nil
}

func (uh *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req reqisterUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		uh.logger.Printf("ERROR: decoding create user body %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid req body"})
		return
	}
	if err := validateUser(&req); err != nil {
		uh.logger.Printf("ERROR:  create user validation failed %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err.Error()})
		return
	}

	user := &store.User{
		Username: req.Username,
		Email:    req.Email,
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}
	err = user.PasswordHash.Set(req.Password)
	if err != nil {
		uh.logger.Printf("ERROR: hashing password failed %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	CreatedUser, err := uh.UserStore.CreateUser(user)
	if err != nil {
		uh.logger.Printf("ERROR: created user failed %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to create a new user"})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"data": CreatedUser})
}

func (uh *UserHandler) GetUserByUsernameHandler(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if username == "" {
		uh.logger.Printf("ERROR: readStringParam: username is empty")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid username"})
		return
	}

	user, err := uh.UserStore.GetUserByUsername(username)
	if err != nil {
		uh.logger.Printf("ERROR: GetUserByUsername: %v ", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "user did not exist"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"data": user})
}

func (uh *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user store.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		uh.logger.Printf("ERROR: invalid req.body %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": err})
		return
	}

	err = uh.UserStore.UpdateUser(&user)
	if err != nil {
		uh.logger.Printf("ERROR: failed to update user %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to update user"})
		return
	}
	utils.WriteJSON(w, http.StatusAccepted, utils.Envelope{"data": "user updated successfully"})
}
