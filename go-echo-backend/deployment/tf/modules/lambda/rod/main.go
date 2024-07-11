package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
	"github.com/mitchellh/mapstructure"
)

func resp(body string, status int) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode:      status,
		Headers:         map[string]string{"Content-Type": "application/pdf"},
		Body:            body,
		IsBase64Encoded: true,
	}, nil
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var params GetPDFParams
	var err = mapstructure.Decode(request.QueryStringParameters, &params)
	if err != nil {
		log.Print(err)
		err := fmt.Sprintf("Failed to fetch page HTML. Error: %s", err.Error())
		return resp(err, 400)
	}

	// get the page's html
	data, err := getPDF(params)
	if err != nil {
		log.Print(err)
		err := fmt.Sprintf("Failed to fetch page HTML. Error: %s", err.Error())
		return resp(err, 500)
	}

	return resp(base64.StdEncoding.EncodeToString(data), 200)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if os.Getenv("IS_TEST_RFQ") == "1" {
		data, err := getPDF(GetPDFParams{
			URL:               "http://localhost:3000/invoices/print/603",
			Selector:          "#invoice-ready-to-print",
			JWTToken:          "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNnNWFucjJsbGttNmN0cHZxOGswIiwiYXVkIjoic3VwZXJfYWRtaW4iLCJpc3MiOiJjbXM4MHBmOGYzamg2YWI2bW5jMCIsInN1YiI6InN1cGVyX2FkbWluIn0.pPEXP5t__mekw65SANTjXu7jRZCJAq7rZJBWDJYvksU",
			Landscape:         "true",
			PrintBackground:   "true",
			PreferCSSPageSize: "true",
			PaperWidth:        "8.5",
			PaperHeight:       "11",
			MarginTop:         "0",
			MarginBottom:      "0",
			MarginLeft:        "0",
			MarginRight:       "0",
		})
		if err != nil {
			panic(err)
		}

		ioutil.WriteFile("test.pdf", data, 0664)
		return
	}

	if os.Getenv("IS_TEST_BULK") == "1" {
		data, err := getPDF(GetPDFParams{
			URL:               "http://localhost:3000/bulks/clmah0bb2hjafsli51t0/commercial-invoice",
			Selector:          "#commercial-invoice-ready-to-print",
			JWTToken:          "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ImNnNWFucjJsbGttNmN0cHZxOGswIiwiYXVkIjoic3VwZXJfYWRtaW4iLCJpc3MiOiJjbXM4MHBmOGYzamg2YWI2bW5jMCIsInN1YiI6InN1cGVyX2FkbWluIn0.pPEXP5t__mekw65SANTjXu7jRZCJAq7rZJBWDJYvksU",
			Landscape:         "true",
			PrintBackground:   "true",
			PreferCSSPageSize: "true",
			PaperWidth:        "8.5",
			PaperHeight:       "11",
			MarginTop:         "0",
			MarginBottom:      "0",
			MarginLeft:        "0",
			MarginRight:       "0",
		})
		if err != nil {
			panic(err)
		}

		ioutil.WriteFile("test.pdf", data, 0664)
		return
	}

	lambda.Start(handler)

}
