package repository

// TokenRepository token repository
type TokenRepository interface {
	GetRefreshToken(userID string) (string, error)
	SaveRefreshToken(userID, token string, ttl int) error
}
