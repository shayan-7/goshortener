package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shayan-7/goshortener/internal/auth"
	"github.com/shayan-7/goshortener/internal/db"
	"github.com/shayan-7/goshortener/internal/models"
	"gorm.io/gorm"
)

// Signup creates a member in db
func Signup(w http.ResponseWriter, r *http.Request) {
	member := &models.Member{}

	err := member.FromJSON(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(member)

	err = member.HashPassword(member.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = member.CreateRecord()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, _ := json.Marshal(map[string]string{"username": member.Username})
	w.Write(resp)
}

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse token response
type LoginResponse struct {
	Token string `json:"token"`
}

// Login logs users in
func Login(w http.ResponseWriter, r *http.Request) {
	payload := &LoginPayload{}
	var member models.Member

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := db.GlobalDB.Where("username = ?", payload.Username).First(&member)
	if result.Error == gorm.ErrRecordNotFound {
		http.Error(w, "Member not found", http.StatusUnauthorized)
		return
	}

	err = member.CheckPassword(payload.Password)
	if err != nil {
		http.Error(w, "invalid member credentials", http.StatusUnauthorized)
		return
	}

	jwtWrapper := auth.JwtWrapper{
		SecretKey:       "verysecretkey",
		Issuer:          "AuthService",
		ExpirationHours: 24,
	}

	signedToken, err := jwtWrapper.GenerateToken(member.Username)
	if err != nil {
		http.Error(w, "error signing token", http.StatusUnauthorized)
		return
	}

	tokenResponse := LoginResponse{
		Token: signedToken,
	}

	enc := json.NewEncoder(w)
	enc.Encode(tokenResponse)
}
