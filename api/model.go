package api

import (
	"sync/atomic"
	"time"

	"github.com/MarianGheorghiu/chirpy/internal/database"
	"github.com/google/uuid"
)

type APIConfig struct {
	fileserverHits atomic.Int32
	Queries        *database.Queries
	Platform       string
	TokenSecret    string
	PolkaKey       string
}

type errorResp struct {
	Error string `json:"error"`
}

type userParameters struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	Password    string    `json:"-"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

type response struct {
	User
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

// input din request
type chirpCreateInput struct {
	Body string `json:"body"`
}

// DTO de răspuns
type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	Body      string    `json:"body"`
}

// payloadul primit de la Polka
type polkaWebhook struct {
	Event string           `json:"event"`
	Data  polkaWebhookData `json:"data"`
}

type polkaWebhookData struct {
	UserID string `json:"user_id"` // îl parsezi separat ca UUID
}
