package persistance

import (
	"github.com/shuufujita/account_book_api/domain/repository"
)

// TokenPersistance token persistance
type tokenPersistance struct{}

// NewTokenPersistance token persistance instance
func NewTokenPersistance() repository.TokenRepository {
	return &tokenPersistance{}
}

func (tp tokenPersistance) GetRefreshToken(userID string) (string, error) {
	refreshToken, err := RedisGet("refresh:" + userID)
	if err != nil {
		return "", err
	}
	return refreshToken, nil
}

func (tp tokenPersistance) SaveRefreshToken(userID, token string, ttl int) error {
	return RedisSet("refresh:"+userID, token, ttl)
}
