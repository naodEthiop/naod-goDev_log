package store

import (
	"sync"

	"github.com/naodEthiop/naod-goDev_log.git/internal/models"
)

var (
	UserCache  = make(map[int]models.User)
	CacheMutex sync.RWMutex
	NextID     = 1
)
