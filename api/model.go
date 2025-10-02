package api

import "sync/atomic"

type APIConfig struct {
	fileserverHits atomic.Int32
}
