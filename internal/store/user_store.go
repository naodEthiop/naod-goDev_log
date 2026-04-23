package store

import(
	"sync"
	"naod-goDev_log/internal/models"
)




var (
	userCache = make(map[int]models.User)
	cacheMutex  sync.RWMutex
	nextID = 1

)