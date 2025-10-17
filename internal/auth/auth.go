package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TokenType e un alias textual pentru a descrie "tipul" tokenului (issuer-ul).
// Poți avea pe viitor "chirpy-refresh", "chirpy-email", etc. dacă vei avea mai multe tipuri.
type TokenType string

const (
	// TokenTypeAccess identifică tokenurile de acces (folosite la autentificare standard).
	TokenTypeAccess TokenType = "chirpy-access"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hashedPassword, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return match, nil
}

// MakeJWT creează un JWT semnat HS256 ce conține:
// - Issuer: "chirpy-access" (tipul tokenului)
// - IssuedAt: timpul curent (UTC)
// - ExpiresAt: timpul curent + durata primită
// - Subject: ID-ul utilizatorului (stringificat)
func MakeJWT(
	userID uuid.UUID,
	tokenSecret string,
	expiresIn time.Duration,
) (string, error) {

	// Cheia simetrică folosită la HS256 e un []byte derivat din secretul tău.
	signingKey := []byte(tokenSecret)

	// Construim RegisteredClaims cu datele standard ale unui JWT.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    string(TokenTypeAccess),                             // "chirpy-access" (vezi const de mai sus)
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),                // iat
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)), // exp
		Subject:   userID.String(),                                     // sub: ID-ul userului ca string
	})

	// Semnăm tokenul cu HS256 + secretul și returnăm stringul compact JWS.
	return token.SignedString(signingKey)
}

// ValidateJWT validează un JWT (semnătura + claims) și întoarce userID-ul din Subject.
// Pași:
// 1) Parsează și verifică semnătura folosind secretul.
// 2) Extrage Subject (ID-ul userului) din claims.
// 3) Verifică Issuer == "chirpy-access" (tipul tokenului).
// 4) Parsează Subject ca UUID și îl returnează.
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	// Struct în care vrem să decodeze claims (RegisteredClaims standard).
	claimsStruct := jwt.RegisteredClaims{}

	// Parsează tokenul + validează semnătura folosind tokenSecret.
	// Notă: în v5, dacă 'exp' e prezent, expirarea e verificată automat (produce err dacă e expirat).
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct, // aici încarcă RegisteredClaims
		func(token *jwt.Token) (any, error) { return []byte(tokenSecret), nil }, // cheia HS256
	)
	if err != nil {
		return uuid.Nil, err // semnătură invalidă, token expirat, malformat, etc.
	}

	// Extragem subject (ID-ul userului) din interfața Claims a tokenului.
	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	// Verificăm Issuer-ul (practic tipul tokenului) — trebuie să fie "chirpy-access".
	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.Nil, err
	}
	if issuer != string(TokenTypeAccess) {
		return uuid.Nil, errors.New("invalid issuer")
	}

	// Convertim Subject (string) în UUID și îl returnăm.
	id, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	header := strings.TrimSpace(headers.Get("Authorization"))
	if header == "" {
		return "", errors.New("authorization header missing")
	}

	headerParts := strings.Fields(header)
	if len(headerParts) != 2 {
		return "", errors.New("invalid authorization header format")
	}

	token := strings.TrimSpace(headerParts[1])
	if token == "" {
		return "", errors.New("empty bearer token")
	}
	return token, nil
}

func MakeRefreshToken() (string, error) {
	byteToken := make([]byte, 32)
	data, err := rand.Read(byteToken)
	if err != nil {
		return "", err
	}

	if data != len(byteToken) {
		return "", fmt.Errorf("short read: got %d, want %d", data, len(byteToken))
	}

	refreshToken := hex.EncodeToString(byteToken)
	return refreshToken, nil
}

func GetAPIKey(headers http.Header) (string, error) {
	header := strings.TrimSpace(headers.Get("Authorization"))
	if header == "" {
		return "", errors.New("authorization header missing")
	}

	parts := strings.Fields(header)
	if len(parts) != 2 || parts[0] != "ApiKey" {
		return "", errors.New("invalid authorization header format")
	}

	key := strings.TrimSpace(parts[1])
	if key == "" {
		return "", errors.New("api key missing")
	}
	return key, nil
}
