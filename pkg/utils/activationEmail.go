package utils

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"library/pkg/redis"
	"net/http"
	"strconv"
	"time"
)

func GenerateActivationLink(ctx context.Context, redisClient *redis.Client, userID int) (string, error) {
	log := GetLogger(ctx)

	token, err := generateRandomToken()
	if err != nil {
		log.Fatalf("Failed to generate random token: %v", err)
	}

	req, _ := http.NewRequest("GET", "https://localhost:5000/v1/users/activate", nil)
	q := req.URL.Query()

	q.Add("token", token)
	q.Add("userID", strconv.Itoa(userID))

	req.URL.RawQuery = q.Encode()
	activationLink := req.URL.String()

	err = redisClient.Client.Set(ctx, strconv.Itoa(userID), token, 24*time.Hour).Err()
	if err != nil {
		log.Errorf("failed to set activation link in Redis: %v", err)
		return "", err
	}

	return activationLink, nil
}

func generateRandomToken() (string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}

	token := base64.URLEncoding.EncodeToString(tokenBytes)

	return token, nil
}
