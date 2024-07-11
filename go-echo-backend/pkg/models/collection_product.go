package models

// Faq Faq's model
type CollectionProduct struct {
	Model

	ProductID    string `json:"product_id,omitempty"`
	CollectionID string `json:"collection_id,omitempty"`
}

type CollectionProductUpdateForm struct {
	ProductID    string `json:"product_id,omitempty"`
	CollectionID string `json:"collection_id,omitempty"`
}

type CollectionProductIDsForm struct {
	Products []string `json:"products,omitempty"`
}
