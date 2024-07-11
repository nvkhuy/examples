package googlesheet

import (
	"context"
	"log"
	"os"

	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const Scope = "https://spreadsheets.google.com/feeds"

func New(config *config.Configuration) (srv *sheets.Service) {
	ctx := context.Background()
	b, err := os.ReadFile(config.GoogleClientSecretURL)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
		return
	}

	googleConfig, err := google.JWTConfigFromJSON(b, Scope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
		return
	}
	client := googleConfig.Client(context.TODO())

	srv, err = sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
		return
	}
	return
}
