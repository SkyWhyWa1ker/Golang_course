package storage

import (
	"practice2/models"
	"sync"
)

var Tasks = make(map[int]models.Task)
var IDCounter = 1
var Mutex sync.Mutex
