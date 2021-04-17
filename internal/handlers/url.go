package handlers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/shayan-7/goshortener/internal/auth"
	"github.com/shayan-7/goshortener/internal/models"
	"github.com/shayan-7/goshortener/internal/util"
)

type URLHandler struct {
	l *log.Logger
	r *redis.Client
}

func NewURLHandler(l *log.Logger, r *redis.Client) *URLHandler {
	return &URLHandler{l, r}
}

func (u *URLHandler) GetURLs(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	item, err := models.FindItem(id, u.r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	err = item.ToJSON(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	isBrowser := isBrowserRequest(r)
	u.r.Publish("goshortener_stats", fmt.Sprintf("%s:%t", id, isBrowser))
}

func (u *URLHandler) AddURL(w http.ResponseWriter, r *http.Request) {
	status, err := authorize(r)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}
	item := &models.Item{}
	item.ID = util.GetHashID()
	err = item.FromJSON(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	models.AddItem(item, u.r)
	err = item.ToJSON(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type itemKey struct{}
type usernameKey struct{}

func authorize(r *http.Request) (int, error) {
	authorization, ok := r.Header["Authorization"]
	if !ok {
		code := http.StatusUnauthorized
		return code, errors.New(http.StatusText(code))
	}
	token := authorization[0]
	extractedToken := strings.Split(token, "Bearer ")
	if len(extractedToken) == 2 {
		token = strings.TrimSpace(extractedToken[1])
	} else {
		return 400, errors.New("Incorrect Format of Authorization Token")
	}

	// FIXME: The SecretKey and Issuer must be derived from config file
	jwtWrapper := auth.JwtWrapper{
		SecretKey: "verysecretkey",
		Issuer:    "AuthService",
	}
	claims, err := jwtWrapper.ValidateToken(token)
	if err != nil {
		return 401, err
	}
	ctx := context.WithValue(r.Context(), usernameKey{}, claims.Username)
	r = r.WithContext(ctx)
	return 200, nil
}

func isBrowserRequest(r *http.Request) bool {
	pattern := "\\bChrome\\b|\\bMozilla\\b|\\bAppleWebKit\\b"
	userAgent, ok := r.Header["User-Agent"]
	if !ok {
		return false
	}
	match, _ := regexp.MatchString(pattern, userAgent[0])
	return match
}
