package usecases

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/shuufujita/account_book_api/domain/model"
	"github.com/shuufujita/account_book_api/domain/repository"
)

// TokenUsecase usecase of token
type TokenUsecase interface {
	Generate(userID string) (model.Token, error)
	Parse(tokenString string) (*jwt.Token, error)
}

type tokenUsecase struct {
	repository repository.TokenRepository
}

// NewTokenUsecase return TokenUsecase instance
func NewTokenUsecase(tr repository.TokenRepository) TokenUsecase {
	return &tokenUsecase{
		repository: tr,
	}
}

func (tu tokenUsecase) Generate(userID string) (model.Token, error) {
	now := time.Now().Unix()

	// アクセストークンを生成する
	accessTokenExpire, err := getAccessTokenExpire(now)
	if err != nil {
		return model.Token{}, err
	}
	accessTokenString, err := generateTokenString(userID, now, accessTokenExpire)
	if err != nil {
		log.Println(fmt.Sprintf("%v: [%v] %v", "error", userID, err.Error()))
		return model.Token{}, err
	}

	// リフレッシュトークンを生成する
	refreshTokenExpire, err := getRefreshTokenExpire(now)
	if err != nil {
		return model.Token{}, err
	}
	refreshTokenString, err := generateTokenString(userID, now, refreshTokenExpire)
	if err != nil {
		log.Println(fmt.Sprintf("%v: [%v] %v", "error", userID, err.Error()))
		return model.Token{}, err
	}

	// リフレッシュトークンをキャッシュに保存する
	err = tu.repository.SaveRefreshToken(userID, refreshTokenString, int(refreshTokenExpire))
	if err != nil {
		log.Println(fmt.Sprintf("%v: [%v] %v", "error", userID, err.Error()))
		return model.Token{}, err
	}

	tokens := model.Token{
		AccessToken:         accessTokenString,
		AccessTokenExpires:  time.Unix(accessTokenExpire, 0).Format(time.RFC3339),
		RefreshToken:        refreshTokenString,
		RefreshTokenExpires: time.Unix(refreshTokenExpire, 0).Format(time.RFC3339),
	}
	return tokens, nil
}

func getAccessTokenExpire(now int64) (int64, error) {
	tokenExpireMinutes, err := strconv.ParseInt(os.Getenv("ACCESS_TOKEN_EXPIRATION_MINUTES"), 10, 64)
	if err != nil {
		return 0, err
	}
	return time.Unix(now, 0).Add(time.Minute * time.Duration(tokenExpireMinutes)).Unix(), nil
}

func getRefreshTokenExpire(now int64) (int64, error) {
	tokenExpireMinutes, err := strconv.ParseInt(os.Getenv("REFRESH_TOKEN_EXPIRATION_MINUTES"), 10, 64)
	if err != nil {
		return 0, err
	}
	return time.Unix(now, 0).Add(time.Minute * time.Duration(tokenExpireMinutes)).Unix(), nil
}

func generateTokenString(userID string, issuedAt int64, expire int64) (string, error) {
	claims := &jwt.StandardClaims{
		Issuer:    "AccountBookAPI",
		Subject:   userID,
		Audience:  "AccountBookAPI",
		IssuedAt:  issuedAt,
		NotBefore: issuedAt,
		ExpiresAt: expire,
	}

	keyData, err := ioutil.ReadFile(os.Getenv("PRYVATE_KEY_PATH"))
	if err != nil {
		log.Println(fmt.Sprintf("%v: [%v] %v", "error", "privateKey", err.Error()))
		return "", nil
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		log.Println(fmt.Sprintf("%v: [%v] %v", "error", "privateKey", err.Error()))
		return "", nil
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return jwtToken.SignedString(privateKey)
}

func (tu tokenUsecase) Parse(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		keyData, err := ioutil.ReadFile(os.Getenv("PUBLIC_KEY_PATH"))
		if err != nil {
			log.Println(fmt.Sprintf("%v: [%v] %v", "error", "publicKey", err.Error()))
			return nil, err
		}
		key, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
		if err != nil {
			log.Println(fmt.Sprintf("%v: [%v] %v", "error", "publicKey", err.Error()))
			return nil, err
		}
		return key, nil
	})
}
