package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
)

func GenerateOBFactoryProductTypeConstants() []*OnboardingConstant {
	var data []*OnboardingConstant
	var items = []enums.OBFactoryProductType{
		enums.OBFactoryProductTypeBlouse,
		enums.OBFactoryProductTypeBodysuit,
		enums.OBFactoryProductTypeCamisole,
		enums.OBFactoryProductTypeCrochetKnitted,
		enums.OBFactoryProductTypeDress,
		enums.OBFactoryProductTypeHoodie,
		enums.OBFactoryProductTypeJacketCoatBlaze,
		enums.OBFactoryProductTypeJumpsuit,
		enums.OBFactoryProductTypeLegging,
		enums.OBFactoryProductTypePant,
		enums.OBFactoryProductTypePoloShirt,
		enums.OBFactoryProductTypeSweatShirt,
		enums.OBFactoryProductTypeTShirt,
		enums.OBFactoryProductTypeSkirt,
		enums.OBFactoryProductTypeSkort,
		enums.OBFactoryProductTypeTankTop,
		enums.OBFactoryProductTypeUnderwearSwimwear,
		enums.OBFactoryProductTypeVestSuit,
		enums.OBFactoryProductTypeCrochetKnitting,
	}

	for _, item := range items {
		data = append(data, &OnboardingConstant{
			Name:  item.DisplayName(),
			Value: item.String(),
		})
	}

	data = append(data, &OnboardingConstant{
		Value: "other",
		Name:  "Other",
	})

	return data
}

func GenerateOBFabricTypeConstants() []*OnboardingConstant {
	var data []*OnboardingConstant
	var items = []enums.OBFabricType{
		enums.OBFabricTypeSequin,
		enums.OBFabricTypeVelvet,
		enums.OBFabricTypeSatin,
		enums.OBFabricTypeSilk,
		enums.OBFabricTypePuLeather,
		enums.OBFabricTypeFur,
		enums.OBFabricTypeKnit,
		enums.OBFabricTypeWoven,
		enums.OBFabricTypeDenim,
		enums.OBFabricTypeThickness,
	}

	for _, item := range items {
		data = append(data, &OnboardingConstant{
			Name:  item.DisplayName(),
			Value: item.String(),
		})
	}

	data = append(data, &OnboardingConstant{
		Value: "other",
		Name:  "Other",
	})

	return data
}

func GenerateOBMillFabricTypeConstants() []*OnboardingConstant {
	var data []*OnboardingConstant
	var items = []enums.OBMillFabricType{
		enums.OBMillFabricTypeTerryCloth,
		enums.OBMillFabricTypeSherpa,
		enums.OBMillFabricTypeJersey,
		enums.OBMillFabricTypePimaCotton,
		enums.OBMillFabricTypeSupimaCotton,
		enums.OBMillFabricTypePolyester,
		enums.OBMillFabricTypeCorduroy,
		enums.OBMillFabricTypeCottonFleece,
		enums.OBMillFabricTypeWool,
		enums.OBMillFabricTypeVelvet,
		enums.OBMillFabricTypeSatin,
		enums.OBMillFabricTypeTwill,
		enums.OBMillFabricTypeSlubbed,
		enums.OBMillFabricTypeSpandex,
		enums.OBMillFabricTypeTweed,
		enums.OBMillFabricTypeViscoseRayon,
		enums.OBMillFabricTypeSilk,
		enums.OBMillFabricTypeFauxFur,
		enums.OBMillFabricTypeDemin,
		enums.OBMillFabricTypeGingham,
		enums.OBMillFabricTypeFrenchTerry,
		enums.OBMillFabricTypeLeather,
		enums.OBMillFabricTypePique,
		enums.OBMillFabricTypeSuede,
		enums.OBMillFabricTypeChambray,
		enums.OBMillFabricTypeOrganicCotton,
		enums.OBMillFabricTypeNeoprene,
		enums.OBMillFabricTypeNylon,
		enums.OBMillFabricTypePoplin,
		enums.OBMillFabricTypeCalico,
		enums.OBMillFabricTypeChiffon,
		enums.OBMillFabricTypeFlannel,
		enums.OBMillFabricTypeOxford,
		enums.OBMillFabricTypeLinen,
		enums.OBMillFabricTypeFelt,
	}

	for _, item := range items {
		data = append(data, &OnboardingConstant{
			Name:     item.DisplayName(),
			Value:    item.String(),
			ImageUrl: item.ImageUrl(),
		})
	}

	data = append(data, &OnboardingConstant{
		Value: "other",
		Name:  "Other",
	})

	return data
}

func GenerateOBSewingAccessoryTypeConstants() []*OnboardingConstant {
	var data []*OnboardingConstant
	var items = []enums.OBSewingAccessoryType{
		enums.OBSewingAccessoryTypeButton,
		enums.OBSewingAccessoryTypeZipper,
		enums.OBSewingAccessoryTypeLining,
		enums.OBSewingAccessoryTypeInterlining,
		enums.OBSewingAccessoryTypeSnapButton,
		enums.OBSewingAccessoryTypeHaspsAndSlider,
		enums.OBSewingAccessoryTypeEmbroidery,
		enums.OBSewingAccessoryTypeApplique,
		enums.OBSewingAccessoryTypeBeads,
		enums.OBSewingAccessoryTypeGlitter,
		enums.OBSewingAccessoryTypeRhinestones,
		enums.OBSewingAccessoryTypeSequins,
		enums.OBSewingAccessoryTypeDrawstring,
		enums.OBSewingAccessoryTypeWaistTies,
		enums.OBSewingAccessoryTypeBows,
		enums.OBSewingAccessoryTypeFringe,
		enums.OBSewingAccessoryTypePomPom,
		enums.OBSewingAccessoryTypeTassel,
		enums.OBSewingAccessoryTypeLabel,
		enums.OBSewingAccessoryTypeMainLabel,
		enums.OBSewingAccessoryTypePULabel,
		enums.OBSewingAccessoryTypePatch,
		enums.OBSewingAccessoryTypeHookAndLoop,
		enums.OBSewingAccessoryTypeEyeletOrGrommet,
		enums.OBSewingAccessoryTypeHookAndEye,
		enums.OBSewingAccessoryTypePadding,
		enums.OBSewingAccessoryTypeElastic,
		enums.OBSewingAccessoryTypeLaceFabric,
		enums.OBSewingAccessoryTypeTwillTape,
		enums.OBSewingAccessoryTypeRib,
		enums.OBSewingAccessoryTypeBelt,
		enums.OBSewingAccessoryTypeStrapping,
	}

	for _, item := range items {
		data = append(data, &OnboardingConstant{
			Name:  item.DisplayName(),
			Value: item.String(),
		})
	}

	return data
}

func GenerateOBPackingAccessoryTypeConstants() []*OnboardingConstant {
	var data []*OnboardingConstant
	var items = []enums.OBPackingAccessoryType{
		enums.OBPackingAccessoryTypeHanger,
		enums.OBPackingAccessoryTypeSafetyPin,
		enums.OBPackingAccessoryTypeScotchTape,
		enums.OBPackingAccessoryTypePolybag,
		enums.OBPackingAccessoryTypeCarton,
		enums.OBPackingAccessoryTypeTags,
		enums.OBPackingAccessoryTypeTissuePaper,
		enums.OBPackingAccessoryTypeButterPaper,
		enums.OBPackingAccessoryTypePlasticClip,
		enums.OBPackingAccessoryTypePaperBoard,
		enums.OBPackingAccessoryTypeButterfly,
		enums.OBPackingAccessoryTypeShirtCollarSupport,
		enums.OBPackingAccessoryTypeShirtBackSupport,
		enums.OBPackingAccessoryTypeTagPin,
		enums.OBPackingAccessoryTypePriceTag,
		enums.OBPackingAccessoryTypeBallHeadPin,
		enums.OBPackingAccessoryTypeInnerBox,
		enums.OBPackingAccessoryTypeFoam,
		enums.OBPackingAccessoryTypeTagGun,
		enums.OBPackingAccessoryTypeClip,
		enums.OBPackingAccessoryTypePlasticAdjuster,
		enums.OBPackingAccessoryTypeShirtBox,
	}

	for _, item := range items {
		data = append(data, &OnboardingConstant{
			Name:  item.DisplayName(),
			Value: item.String(),
		})
	}

	return data
}
