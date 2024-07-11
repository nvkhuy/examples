package enums

import (
	"fmt"

	"github.com/engineeringinflow/inflow-backend/pkg/config"
)

type OBMillFabricType string

var (
	OBMillFabricTypeTerryCloth    OBMillFabricType = "terry_cloth"
	OBMillFabricTypeSherpa        OBMillFabricType = "sherpa"
	OBMillFabricTypeJersey        OBMillFabricType = "jersey"
	OBMillFabricTypePimaCotton    OBMillFabricType = "pima_cotton"
	OBMillFabricTypeSupimaCotton  OBMillFabricType = "supima_cotton"
	OBMillFabricTypePolyester     OBMillFabricType = "polyester"
	OBMillFabricTypeCorduroy      OBMillFabricType = "corduroy"
	OBMillFabricTypeCottonFleece  OBMillFabricType = "cotton_fleece"
	OBMillFabricTypeWool          OBMillFabricType = "wool"
	OBMillFabricTypeVelvet        OBMillFabricType = "velvet"
	OBMillFabricTypeSatin         OBMillFabricType = "satin"
	OBMillFabricTypeTwill         OBMillFabricType = "twill"
	OBMillFabricTypeSlubbed       OBMillFabricType = "slubbed"
	OBMillFabricTypeSpandex       OBMillFabricType = "spandex"
	OBMillFabricTypeTweed         OBMillFabricType = "tweed"
	OBMillFabricTypeViscoseRayon  OBMillFabricType = "viscose_rayon"
	OBMillFabricTypeSilk          OBMillFabricType = "silk"
	OBMillFabricTypeFauxFur       OBMillFabricType = "faux_fur"
	OBMillFabricTypeDemin         OBMillFabricType = "demin"
	OBMillFabricTypeGingham       OBMillFabricType = "gingham"
	OBMillFabricTypeFrenchTerry   OBMillFabricType = "french_terry"
	OBMillFabricTypeLeather       OBMillFabricType = "leather"
	OBMillFabricTypePique         OBMillFabricType = "pique"
	OBMillFabricTypeSuede         OBMillFabricType = "suede"
	OBMillFabricTypeChambray      OBMillFabricType = "chambray"
	OBMillFabricTypeOrganicCotton OBMillFabricType = "organic_cotton"
	OBMillFabricTypeNeoprene      OBMillFabricType = "neoprene"
	OBMillFabricTypeNylon         OBMillFabricType = "nylon"
	OBMillFabricTypePoplin        OBMillFabricType = "poplin"
	OBMillFabricTypeCalico        OBMillFabricType = "calico"
	OBMillFabricTypeChiffon       OBMillFabricType = "chiffon"
	OBMillFabricTypeFlannel       OBMillFabricType = "flannel"
	OBMillFabricTypeOxford        OBMillFabricType = "oxford"
	OBMillFabricTypeLinen         OBMillFabricType = "linen"
	OBMillFabricTypeFelt          OBMillFabricType = "felt"
)

func (p OBMillFabricType) String() string {
	return string(p)
}

func (p OBMillFabricType) DisplayName() string {
	var name = string(p)

	switch p {
	case OBMillFabricTypeTerryCloth:
		name = "Terry Cloth"
	case OBMillFabricTypeSherpa:
		name = "Sherpa"
	case OBMillFabricTypeJersey:
		name = "Jersey"
	case OBMillFabricTypePimaCotton:
		name = "Pima Cotton"
	case OBMillFabricTypeSupimaCotton:
		name = "SupimaCotton"
	case OBMillFabricTypePolyester:
		name = "Polyester"
	case OBMillFabricTypeCorduroy:
		name = "Corduroy"
	case OBMillFabricTypeCottonFleece:
		name = "Cotton Fleece"
	case OBMillFabricTypeWool:
		name = "Wool"
	case OBMillFabricTypeVelvet:
		name = "Velvet"
	case OBMillFabricTypeSatin:
		name = "Satin"
	case OBMillFabricTypeSatin:
		name = "Twill"
	case OBMillFabricTypeSlubbed:
		name = "Slubbed"
	case OBMillFabricTypeSpandex:
		name = "Spandex"
	case OBMillFabricTypeViscoseRayon:
		name = "Viscose/Rayon"
	case OBMillFabricTypeTweed:
		name = "Tweed"
	case OBMillFabricTypeSilk:
		name = "Silk"
	case OBMillFabricTypeFauxFur:
		name = "Faux Fur"
	case OBMillFabricTypeDemin:
		name = "Demin"
	case OBMillFabricTypeGingham:
		name = "Gingham"
	case OBMillFabricTypeFrenchTerry:
		name = "French Terry"
	case OBMillFabricTypeLeather:
		name = "Leather"
	case OBMillFabricTypePique:
		name = "Pique"
	case OBMillFabricTypeSuede:
		name = "Suede"
	case OBMillFabricTypeChambray:
		name = "Chambray"
	case OBMillFabricTypeOrganicCotton:
		name = "Organic Cotton"
	case OBMillFabricTypeNeoprene:
		name = "Neoprene"
	case OBMillFabricTypeNylon:
		name = "Nylon"
	case OBMillFabricTypePoplin:
		name = "Poplin"
	case OBMillFabricTypeCalico:
		name = "Calico"
	case OBMillFabricTypeChiffon:
		name = "Chiffon"
	case OBMillFabricTypeFlannel:
		name = "Flannel"
	case OBMillFabricTypeOxford:
		name = "Oxford"
	case OBMillFabricTypeLinen:
		name = "Linen"
	case OBMillFabricTypeFelt:
		name = "Felt"

	}

	return name
}

func (p OBMillFabricType) ImageUrl() string {
	var name = p.DisplayName()
	var cfg = config.GetInstance()

	return fmt.Sprintf("https://%s/onboarding/fabric-type-thumbnails/%s.jpg", cfg.CDNStaticURL, name)
}
