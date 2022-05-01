package email

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/sheets/v4"
	"net/http"
	"net/mail"
)

var router *gin.Engine

func init() {
	router := gin.Default()
	router.POST("/", SendEmail)
	//router.GET("/emails", getEmails)

	router.Run()

}

func EmailEntrypoint(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r) // ServeHTTP conforms to the http.Handler interface (https://godoc.org/github.com/gin-gonic/gin#Engine.ServeHTTP)
}

type EmailForm struct {
	Email string `form:"email" binding:"required"`
}

const (
	sheetsId = "1R-MhPdZLJ1W5JvuM0qkCOjlkCq2TeknOfV3ePhnPOPg"
)

func getEmails(context *gin.Context) {

}

func SendEmail(c *gin.Context) {
	var form EmailForm
	if err := c.Bind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	email, err := mail.ParseAddress(form.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email address."})
		return
	}
	saveToSheets(email.Address)
}

func saveToSheets(data string) {
	ctx := context.Background()
	sheetsService, err := sheets.NewService(ctx)

	if err != nil {
		panic(fmt.Errorf("unable to start spreadsheet service: %v", err))
	}

	resp, respErr := sheetsService.Spreadsheets.Get(sheetsId).Do()
	if respErr != nil {
		panic(fmt.Errorf("unable to retrieve data from sheet: %v", err))
	}

	// There is only one sheet page
	id := resp.Sheets[0].Properties.SheetId

	// Formatting centered text
	format := &sheets.CellFormat{
		HorizontalAlignment: "CENTER",
	}

	// Value equal to given email address (or any data as string if required)
	val := &sheets.CellData{
		UserEnteredValue: &sheets.ExtendedValue{
			StringValue: &data,
		},
		UserEnteredFormat: format,
	}

	req := []*sheets.Request{
		{
			AppendCells: &sheets.AppendCellsRequest{
				SheetId: id,
				Fields:  "*",
				Rows: []*sheets.RowData{
					{
						Values: []*sheets.CellData{val},
					},
				},
			},
		},
	}

	_, err = sheetsService.Spreadsheets.BatchUpdate(sheetsId, &sheets.BatchUpdateSpreadsheetRequest{
		IncludeSpreadsheetInResponse: false,
		Requests:                     req,
		ResponseIncludeGridData:      false,
	}).Do()

	if err != nil {
		panic(fmt.Errorf("unable to update: %v", err))
	}
}
