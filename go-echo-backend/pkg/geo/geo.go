package geo

import (
	"context"
	"strings"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/location"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/rotisserie/eris"
	"github.com/thaitanloi365/go-utils/values"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/locationservice"
	"github.com/biter777/countries"
)

var instance *Client

type Client struct {
	config *config.Configuration
	logger *logger.Logger

	client *locationservice.LocationService
}

func New(config *config.Configuration) *Client {
	sess, err := session.NewSession(
		&aws.Config{
			Region:      aws.String(config.AWSS3Region),
			Credentials: credentials.NewStaticCredentials(config.AWSAccessKeyID, config.AWSSecretAccessKey, ""),
		},
	)
	if err != nil {
		panic("An error occured while init AWS map")
	}

	instance = &Client{
		config: config,
		logger: logger.New("utils/geo"),

		client: locationservice.New(sess),
	}
	return instance
}

func GetInstance() *Client {
	if instance == nil {
		panic("Must be call New() first")
	}
	return instance
}

type SearchPlaceIndexForTextParams struct {
	Address     string             `json:"address" param:"address" query:"address" form:"address"`
	CountryCode enums.CountryCode  `json:"country_code" param:"country_code" query:"country_code" form:"country_code"`
	Language    enums.LanguageCode `json:"language" param:"language" query:"language" form:"language"`
	BiasLat     *float64           `json:"bias_lat" param:"bias_lat" query:"bias_lat" form:"bias_lat"`
	BiasLng     *float64           `json:"bias_lng" param:"bias_lng" query:"bias_lng" form:"bias_lng"`
}

func (c *Client) SearchPlaceIndexForText(params SearchPlaceIndexForTextParams) ([]location.Coordinate, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	if params.Address == "" {
		return nil, eris.New("Address is required")
	}

	var country = params.CountryCode.DefaultIfInvalid()
	var countryAlpha3 = countries.ByName(string(country))
	var lan = params.Language.DefaultIfInvalid()
	c.logger.Debugf("Search address=%s country=%s country3=%s language=%s", params.Address, country, countryAlpha3.Alpha3(), lan)

	var request = &locationservice.SearchPlaceIndexForTextInput{
		Language:   aws.String(lan.ToLower()),
		IndexName:  aws.String(string(c.getPlaceIndex(country))),
		MaxResults: aws.Int64(10),
		FilterCountries: aws.StringSlice([]string{
			countryAlpha3.Alpha3(),
		}),
		Text: &params.Address,
	}

	if helper.IsLatLng(params.BiasLat, params.BiasLng) {
		request.BiasPosition = []*float64{
			params.BiasLng,
			params.BiasLat,
		}
	}

	resp, err := c.client.SearchPlaceIndexForTextWithContext(ctx, request)
	if err != nil {
		return nil, err
	}

	var records = c.mapToAdresses(resp, country, lan)
	if len(strings.TrimSpace(params.Address)) > 6 {
		var result = []location.Coordinate{
			{
				FormattedAddress: params.Address,
				CountryCode:      country,
			},
		}

		result = append(result, records...)
		return result, nil
	}

	return records, err

}

func (c *Client) mapToAdresses(resp *locationservice.SearchPlaceIndexForTextOutput, cc enums.CountryCode, lang enums.LanguageCode) []location.Coordinate {
	var records []location.Coordinate

	for _, result := range resp.Results {
		records = append(records, c.convertResultToAddress(result, cc, lang))
	}

	return records
}

func (c *Client) getPlaceIndex(cc enums.CountryCode) PlaceIndex {
	if enums.AsiaCountries.Contains(cc) {
		return PlaceIndexGrab
	}

	return PlaceIndexDefault
}

func (c *Client) convertResultToAddress(result *locationservice.SearchForTextResult, cc enums.CountryCode, lang enums.LanguageCode) (address location.Coordinate) {
	address.FormattedAddress = values.StringValue(result.Place.Label)
	address.AddressNumber = values.StringValue(result.Place.AddressNumber)
	address.Lng = result.Place.Geometry.Point[0]
	address.Lat = result.Place.Geometry.Point[1]

	address.PostalCode = values.StringValue(result.Place.PostalCode)
	address.Street = values.StringValue(result.Place.Street)
	address.CountryCode = cc

	if result.Place.TimeZone != nil {
		address.TimezoneName = enums.Timezone(*result.Place.TimeZone.Name)
		address.TimezoneOffset = int(*result.Place.TimeZone.Offset)
	}

	switch c.getPlaceIndex(cc) {
	case PlaceIndexDefault:
		address.Level1 = values.StringValue(result.Place.Municipality)
		address.Level2 = values.StringValue(result.Place.SubRegion)
		address.Level3 = values.StringValue(result.Place.Region)
	case PlaceIndexGrab:
		address.Level1 = values.StringValue(result.Place.SubRegion)
		address.Level2 = values.StringValue(result.Place.Municipality)
		address.Level3 = values.StringValue(result.Place.Neighborhood)
	}

	return
}
