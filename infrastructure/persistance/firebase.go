package persistance

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
)

var appInstance *firebase.App

// InitializeAppDefault initialized app instance
func InitializeAppDefault() error {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
		return err
	}
	appInstance = app
	return nil
}

// GetAppInstance get firebase app instance
func GetAppInstance() *firebase.App {
	return appInstance
}
