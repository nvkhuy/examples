package customerio

type Identifiers struct {
	ID    string `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
	CIOID string `json:"cio_id"`
}
type Relationship struct {
	Identifiers Identifiers `json:"identifiers"`
}

func (client *Client) IdentifyFacility(id string, attributes map[string]interface{}, clientIDs []string) error {
	var body = map[string]interface{}{
		"type": "object",
		"identifiers": map[string]interface{}{
			"object_type_id": "1",
			"object_id":      id,
		},
		"action":     "identify",
		"attributes": attributes,
	}

	if len(clientIDs) > 0 {
		var relationships []*Relationship

		body["cio_relationships"] = []map[string]interface{}{}

		for _, id := range clientIDs {
			relationships = append(relationships, &Relationship{
				Identifiers: Identifiers{
					ID: id,
				},
			})
		}
		if relationships != nil {
			body["cio_relationships"] = relationships
		}
	}

	_, err := client.restyClient.R().
		SetBasicAuth(client.config.CustomerIOSiteID, client.config.CustomerIOApiTrackingKey).
		SetBody(&body).
		Post("https://track.customer.io/api/v2/entity")

	return err
}
