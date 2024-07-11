package helper

import (
	"fmt"
	"strings"

	"github.com/thaitanloi365/go-utils/random"
)

func GeneratePurchaseOrderReferenceID() string {
	var numRan = random.String(5, random.Numerals)
	var strRan = strings.ToUpper(random.String(4, random.Alphabet))
	return fmt.Sprintf("PO-%s-%s", strRan, numRan)
}

func GenerateBulkPurchaseOrderReferenceID() string {
	var numRan = random.String(5, random.Numerals)
	var strRan = strings.ToUpper(random.String(4, random.Alphabet))
	return fmt.Sprintf("BPO-%s-%s", strRan, numRan)
}

func GenerateBulkPurchaseOrderGroupID() string {
	var numRan = random.String(5, random.Numerals)
	var strRan = strings.ToUpper(random.String(4, random.Alphabet))
	return fmt.Sprintf("GBPO-%s-%s", strRan, numRan)
}

func GenerateOrderGroupID() string {
	var numRan = random.String(5, random.Numerals)
	var strRan = strings.ToUpper(random.String(4, random.Alphabet))
	return fmt.Sprintf("COL-%s-%s", strRan, numRan)
}

func GenerateOrderReferenceID() string {
	var numRan = random.String(5, random.Numerals)
	var strRan = strings.ToUpper(random.String(4, random.Alphabet))
	return fmt.Sprintf("OR-%s-%s", strRan, numRan)
}

func GenerateInquiryReferenceID() string {
	var numRan = random.String(5, random.Numerals)
	var strRan = strings.ToUpper(random.String(4, random.Alphabet))
	return fmt.Sprintf("IQ-%s-%s", strRan, numRan)
}

func GenerateFabricCollectionReferenceID() string {
	var numRan = random.String(5, random.Numerals)
	var strRan = strings.ToUpper(random.String(4, random.Alphabet))
	return fmt.Sprintf("FC-%s-%s", strRan, numRan)
}

func GenerateFabricReferenceID() string {
	var numRan = random.String(5, random.Numerals)
	var strRan = strings.ToUpper(random.String(4, random.Alphabet))
	return fmt.Sprintf("FB-%s-%s", strRan, numRan)
}

func GeneratePaymentTransactionReferenceID() string {
	var numRan = random.String(5, random.Numerals)
	var strRan = strings.ToUpper(random.String(4, random.Alphabet))
	return fmt.Sprintf("PAY-%s-%s", strRan, numRan)
}

func GeneratePoRawMaterialReferenceID() string {
	var numRan = random.String(5, random.Numerals)
	var strRan = strings.ToUpper(random.String(4, random.Alphabet))
	return fmt.Sprintf("RAW-%s-%s", strRan, numRan)
}

func GenerateCheckoutSessionID() string {
	var numRan = random.String(5, random.Numerals)
	var strRan = strings.ToUpper(random.String(4, random.Alphabet))
	return fmt.Sprintf("CK-%s-%s", strRan, numRan)
}

func GenerateSampleRoundID() string {
	var numRan = random.String(5, random.Numerals)
	var strRan = strings.ToUpper(random.String(4, random.Alphabet))
	return fmt.Sprintf("SR-%s-%s", strRan, numRan)
}
