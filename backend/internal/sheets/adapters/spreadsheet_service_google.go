package adapters

import (
	"github.com/doodocs/qaztrade/backend/internal/sheets/domain"
	"github.com/doodocs/qaztrade/backend/pkg/jwt"
	"github.com/doodocs/qaztrade/backend/pkg/qaztradeoauth2"
)

type SpreadsheetServiceGoogle struct {
	oauth2                *qaztradeoauth2.Client
	svcAccount            string
	reviewerAccount       string
	jwtcli                *jwt.Client
	originSpreadsheetID   string
	templateSpreadsheetID string
	destinationFolderID   string
}

var _ domain.SpreadsheetService = (*SpreadsheetServiceGoogle)(nil)

func NewSpreadsheetServiceGoogle(
	oauth2 *qaztradeoauth2.Client,
	svcAccount string,
	reviewerAccount string,
	jwtcli *jwt.Client,
	originSpreadsheetID string,
	templateSpreadsheetID string,
	destinationFolderID string,
) *SpreadsheetServiceGoogle {
	client := &SpreadsheetServiceGoogle{
		oauth2:                oauth2,
		svcAccount:            svcAccount,
		reviewerAccount:       reviewerAccount,
		jwtcli:                jwtcli,
		originSpreadsheetID:   originSpreadsheetID,
		templateSpreadsheetID: templateSpreadsheetID,
		destinationFolderID:   destinationFolderID,
	}

	return client
}
