package handler

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/shuufujita/account_book_api/infrastructure/persistance"
	"github.com/shuufujita/account_book_api/interfaces/response"
	"github.com/shuufujita/account_book_api/usecases"
)

// TokenHandler token handler
type TokenHandler interface {
	IssueToken(c echo.Context) error
}

type tokenHandler struct {
	tokenUsecase usecases.TokenUsecase
}

// NewTokenHandler return token handler instance
func NewTokenHandler(tu usecases.TokenUsecase) TokenHandler {
	return &tokenHandler{
		tokenUsecase: tu,
	}
}

// IssueTokenRequest issue token request
type IssueTokenRequest struct {
	IDToken string `json:"id_token"`
}

func (th tokenHandler) IssueToken(c echo.Context) error {
	request := &IssueTokenRequest{}
	if err := c.Bind(request); err != nil {
		return response.ErrorResponse(c, "INVALID_PARAMETER", err.Error())
	}

	// FirebaseのIDトークンを検証する
	app := persistance.GetAppInstance()
	client, err := app.Auth(c.Request().Context())
	if err != nil {
		return response.ErrorResponse(c, "INTERNAL_SERVER_ERROR", err.Error())
	}
	verifyToken, err := client.VerifyIDToken(c.Request().Context(), request.IDToken)
	if err != nil {
		return response.ErrorResponse(c, "AUTHORIZATION_ERROR", err.Error())
	}

	// トークン情報を生成する
	token, err := th.tokenUsecase.Generate(verifyToken.UID)
	if err != nil {
		return response.ErrorResponse(c, "TOKEN_GENERATE_ERROR", err.Error())
	}

	return c.JSON(http.StatusOK, token)
}
