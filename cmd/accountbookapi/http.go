package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/shuufujita/account_book_api/infrastructure/persistance"
	"github.com/shuufujita/account_book_api/interfaces/handler"
	"github.com/shuufujita/account_book_api/usecases"

	"github.com/comail/colog"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// RunServer launch and run server.
func RunServer(port int64) error {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"*",
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAccessControlAllowCredentials,
			echo.HeaderCookie,
			echo.HeaderSetCookie,
		},
		AllowMethods: []string{
			echo.GET,
			echo.PUT,
			echo.POST,
			echo.DELETE,
		},
		AllowCredentials: true,
	}))

	err := persistance.InitializeAppDefault()
	if err != nil {
		return err
	}

	tokenRepository := persistance.NewTokenPersistance()
	tokenUsecase := usecases.NewTokenUsecase(tokenRepository)
	tokenHandler := handler.NewTokenHandler(tokenUsecase)

	g := e.Group("/v1", customLogger)
	g.POST("/token", tokenHandler.IssueToken)

	return e.Start(":" + strconv.FormatInt(port, 10))
}

func customLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		path, err := getLogFilePath()
		if err != nil {
			return err
		}

		file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
		if err != nil {
			panic(err)
		}
		colog.SetOutput(io.MultiWriter(file, os.Stdout))
		colog.SetFormatter(&colog.StdFormatter{
			Flag: log.Ldate | log.Ltime | log.Lshortfile,
		})
		colog.FixedValue("remoteAddr", c.Request().RemoteAddr)

		return next(c)
	}
}

func getLogFilePath() (string, error) {
	if _, err := os.Stat(os.Getenv("LOG_DIR_PATH")); os.IsNotExist(err) {
		err = os.Mkdir(os.Getenv("LOG_DIR_PATH"), 0777)
		if err != nil {
			return "", err
		}
		log.Println(fmt.Sprintf("%v: [%v] %v", "info", "http", "mkdir with "+os.Getenv("LOG_DIR_PATH")))
	}
	return os.Getenv("LOG_DIR_PATH") + "/" + time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format("20060102") + ".log", nil
}
