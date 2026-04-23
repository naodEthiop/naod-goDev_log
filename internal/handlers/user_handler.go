package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/naodEthiop/naod-goDev_log.git/internal/models"
    "github.com/naodEthiop/naod-goDev_log.git/internal/store"
)

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	store.CacheMutex.RLock()

	users := make([]models.User, 0, len(store.UserCache))
	for _, user := range store.UserCache {
		users = append(users, user)
	}

	store.CacheMutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	store.CacheMutex.RLock()
	user, ok := store.UserCache[id]
	store.CacheMutex.RUnlock()

	if !ok {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if user.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	store.CacheMutex.Lock()
	user.ID = store.NextID
	store.UserCache[user.ID] = user
	store.NextID++
	store.CacheMutex.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	store.CacheMutex.Lock()
	defer store.CacheMutex.Unlock()

	if _, ok := store.UserCache[id]; !ok {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	delete(store.UserCache, id)
	w.WriteHeader(http.StatusNoContent)
}