package tests

import (
	"log"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/stretchr/testify/assert"
)

func TestProductTypesPriceRepo_FetchProductTypesPrice(t *testing.T) {
	var app = initApp("local")
	resp, err := repo.NewProductTypesPriceRepo(app.DB).
		WithSheetAPI(app.SheetAPI).FetchProductTypesPrice(&models.FetchProductTypesPriceParams{
		From: "A3",
		To:   "AE6",
	})
	if err != nil {
		return
	}
	helper.PrintJSON(resp)
}

func TestProductTypesPriceRepo_PatchSheetImageURL(t *testing.T) {
	var app = initApp("prod")
	_ = repo.NewProductTypesPriceRepo(app.DB).WithSheetAPI(app.SheetAPI).PatchSheetImageURL(&models.PatchSheetImageURLParams{
		Token:       "AC4w5VgosJUPvzUEWD6S0rYKPXPB_iXjqQ%3A1700126097867",
		Cookie:      "COMPASS=apps-spreadsheets=CmUACWuJV1WzaG0FlgKb1G2vLfelHALv1rBQS30zkudH3UpzBlLFyNSACBJAzxjBmOxxxg6OzRh7olq7pUEAD_MB9_lg7SexDF-0_GkMj19dGfniORtvyhgCqCnawZIeCGXqHezKbBDPvNeqBhqIAQAJa4lXOhIpbS6PMTkilRD_0ze29i1fKo2FUQgwTmGJbzAcXKaCetqy2eC2OoOkgG9Am5_ScQrEPlhsp9bQMVHv0ifhhHVpvpltrmW0q2fGbWQ6BEfSsz5MrWd8wKooYwYHrkRbzjMhlpAtyxfdL7-efLAKZyTk7DNtfsi2XM_Glwze7HZyHNg=; OSID=cgjGNodDGSfhaO-d6OHplEHH7fcvzN9gsS87YIudMX9zpuKgkXnLPd-piM8c87FnF_b1WA.; __Secure-OSID=cgjGNodDGSfhaO-d6OHplEHH7fcvzN9gsS87YIudMX9zpuKguRgy8LPzdDHHqpGxrKterA.; SID=cwjGNlCAIKosp-bG8FsS3v4UxvvpN56eABDwZPfB04QknOAJfgtO6Nb-XAu8__yRKFzcvg.; __Secure-1PSID=cwjGNlCAIKosp-bG8FsS3v4UxvvpN56eABDwZPfB04QknOAJH3XstwQLsdDX7hgCkonFpw.; __Secure-3PSID=cwjGNlCAIKosp-bG8FsS3v4UxvvpN56eABDwZPfB04QknOAJwKDQgOw7Es2q5LKREvbK3A.; HSID=ApsbI0ANNUpVNiGuW; SSID=AFwyem53z7_5CR7YY; APISID=cqrLiFrF2WN1iyUg/AGu-wCGTftlaO5Zwt; SAPISID=Tw6hTusCXD_p3gbl/Aw3SLtOfuaOTtxCfO; __Secure-1PAPISID=Tw6hTusCXD_p3gbl/Aw3SLtOfuaOTtxCfO; __Secure-3PAPISID=Tw6hTusCXD_p3gbl/Aw3SLtOfuaOTtxCfO; OTZ=7296691_28_28__28_; AEC=Ackid1Srvq7mfaacSMrIR2nMSwII9DXANjdJoICWPXi4A56EzupUT7UtIQ; NID=511=RgDao-aHqZ9ZYbri9Z76inlrFWbjDZuQQZLYM7Ww7WiFxJPBaTz1bcdDXQhZVJyQ_EaprPHGkPHEB4ghPKtUhp5F0DWO6VPdqVqGTTEOadwgsVo_reQZDpVkZ7DWhwRuSK_VrkSLovVEPZ1xtLJCq3X_IvCWlhzufZedMaqGTGFfsVRD5q5_UYpiYWSeTIGXKsaYAo-D3ExhCzufVtbjYuydSc8KqvZVhijTMTGO3XvxZWB47ub-FI359JLmrYT6D9UnwCqnEzr_ttzbyU7W06hUFpugyHRD8PYH3ZT8MalAVKKE8NjoNUS2x6ACu1b8-6mlzAySF8vi0iiG6UwEyf_j7OcszqG4dc9P6tFvX5gKb7aPAPRflyTNxfI8FF516osvgG-ASMQukThnwrVaek-DgoYK7S-0P3QvObawsNXlOcR4hBhZZR94QkOrCaxw9e78pjNMox_bxECr88NbeeCs-9100oW7vANqJ1ap1RUJm-XdqHKU8iGf; 1P_JAR=2023-11-16-09; GOOGLE_ABUSE_EXEMPTION=ID=d7e21f1ed2280056:TM=1700126085:C=r:IP=84.17.39.171-:S=jGdYx_fA2CHiPSKkTZNNTwc; SIDCC=ACA-OxMA1sMYA-BMNNeooT4GDNbVZui2Kup9wfHLlKxyhrQ5QYN3tivPQBACEVTLaPTzJJ48JA; __Secure-1PSIDCC=ACA-OxPEX0QatTc8_7ig91noh0uKhs0Hw9rtmOGlafxCVfCYzekrUOG8bsWD8lsxuBQIPnZ72w; __Secure-3PSIDCC=ACA-OxOLRwyAhnd9pmXRP0QHjM5zlMrIRFhEnoujgL0TAMSqych_29DXjNYarSkh7dchruB-xQ; NID=511=sssLrAnTFS6eVsEt8fMhN69TQrqnqigmmewcfmFek7hhPHHpwYzlgWr8Nl4cfwJyc4iZJmvWvQkAp2_lcBq7cqiNqScfzH-l2qCMpBoCcTl3TF0mq9jjRgN6ZDq1WVpYBUOY1s59EAGd3XAZu3OZ6W3-inS-IC_wuE53nUQ-F2I; SIDCC=ACA-OxO_zI7-Z9VXCPGII9vba950WdUY7fIcIJEGKxzXesdG898SqPbyfNYufen9chSS5La8FQ; __Secure-1PSIDCC=ACA-OxP8EFSnNOa3C9_BJMds8xDUwW4dMpDN7NQYlEbsOz9gnfOsNEcBdMD0GoLOlb-W5ZUToQ; __Secure-3PSIDCC=ACA-OxMgzm51BGE2ncyKDzpBKg7v9ZmX8utw1GHBOLRqLy7GbE7kO3cTGtdMYmnnh_7Csd_jsA",
		Concurrency: 10,
		From:        2,
		To:          662,
	})
}

func TestProductTypesPriceRepo_PaginateProductTypesPrice(t *testing.T) {
	var app = initApp("local")
	resp := repo.NewProductTypesPriceRepo(app.DB).WithSheetAPI(app.SheetAPI).PaginateProductTypesPrice(&models.PaginateProductTypesPriceParams{
		Product:           aws.String("CLOTHING"),
		Gender:            aws.String("MEN"),
		FabricType:        aws.String("KNIT"),
		Feature:           aws.String("TOP"),
		Category:          aws.String("T-shirt"),
		Description:       aws.String("V-neck"),
		KnitMaterial:      aws.String("Single Jersey"),
		KnitComposition:   aws.String("100% Cotton"),
		KnitWeight:        aws.Float64(250),
		FabricConsumption: aws.Float64(1.2),
	})
	helper.PrintJSON(resp)
}

func TestProductTypesPriceRepo_PaginateProductTypesPriceVine(t *testing.T) {
	var app = initApp("local")
	resp := repo.NewProductTypesPriceRepo(app.DB).
		PaginateProductTypesPriceVine(&models.PaginateProductTypesPriceParams{
			Product:         aws.String("CLOTHING"),
			Gender:          aws.String("MEN"),
			FabricType:      aws.String("KNIT"),
			Feature:         aws.String("TOP"),
			Category:        aws.String("T-shirt"),
			Description:     aws.String("V-neck"),
			KnitMaterial:    aws.String("Single Jersey"),
			KnitComposition: aws.String("100% Cotton"),
		})
	helper.PrintJSON(resp)
}

func TestProductTypesPriceRepo_GetSheetImage(t *testing.T) {
	var app = initApp("local")
	token := "AC4w5VgosJUPvzUEWD6S0rYKPXPB_iXjqQ%3A1700126097867"
	cookie := "COMPASS=apps-spreadsheets=CmUACWuJV1WzaG0FlgKb1G2vLfelHALv1rBQS30zkudH3UpzBlLFyNSACBJAzxjBmOxxxg6OzRh7olq7pUEAD_MB9_lg7SexDF-0_GkMj19dGfniORtvyhgCqCnawZIeCGXqHezKbBDPvNeqBhqIAQAJa4lXOhIpbS6PMTkilRD_0ze29i1fKo2FUQgwTmGJbzAcXKaCetqy2eC2OoOkgG9Am5_ScQrEPlhsp9bQMVHv0ifhhHVpvpltrmW0q2fGbWQ6BEfSsz5MrWd8wKooYwYHrkRbzjMhlpAtyxfdL7-efLAKZyTk7DNtfsi2XM_Glwze7HZyHNg=; OSID=cgjGNodDGSfhaO-d6OHplEHH7fcvzN9gsS87YIudMX9zpuKgkXnLPd-piM8c87FnF_b1WA.; __Secure-OSID=cgjGNodDGSfhaO-d6OHplEHH7fcvzN9gsS87YIudMX9zpuKguRgy8LPzdDHHqpGxrKterA.; SID=cwjGNlCAIKosp-bG8FsS3v4UxvvpN56eABDwZPfB04QknOAJfgtO6Nb-XAu8__yRKFzcvg.; __Secure-1PSID=cwjGNlCAIKosp-bG8FsS3v4UxvvpN56eABDwZPfB04QknOAJH3XstwQLsdDX7hgCkonFpw.; __Secure-3PSID=cwjGNlCAIKosp-bG8FsS3v4UxvvpN56eABDwZPfB04QknOAJwKDQgOw7Es2q5LKREvbK3A.; HSID=ApsbI0ANNUpVNiGuW; SSID=AFwyem53z7_5CR7YY; APISID=cqrLiFrF2WN1iyUg/AGu-wCGTftlaO5Zwt; SAPISID=Tw6hTusCXD_p3gbl/Aw3SLtOfuaOTtxCfO; __Secure-1PAPISID=Tw6hTusCXD_p3gbl/Aw3SLtOfuaOTtxCfO; __Secure-3PAPISID=Tw6hTusCXD_p3gbl/Aw3SLtOfuaOTtxCfO; OTZ=7296691_28_28__28_; AEC=Ackid1Srvq7mfaacSMrIR2nMSwII9DXANjdJoICWPXi4A56EzupUT7UtIQ; NID=511=RgDao-aHqZ9ZYbri9Z76inlrFWbjDZuQQZLYM7Ww7WiFxJPBaTz1bcdDXQhZVJyQ_EaprPHGkPHEB4ghPKtUhp5F0DWO6VPdqVqGTTEOadwgsVo_reQZDpVkZ7DWhwRuSK_VrkSLovVEPZ1xtLJCq3X_IvCWlhzufZedMaqGTGFfsVRD5q5_UYpiYWSeTIGXKsaYAo-D3ExhCzufVtbjYuydSc8KqvZVhijTMTGO3XvxZWB47ub-FI359JLmrYT6D9UnwCqnEzr_ttzbyU7W06hUFpugyHRD8PYH3ZT8MalAVKKE8NjoNUS2x6ACu1b8-6mlzAySF8vi0iiG6UwEyf_j7OcszqG4dc9P6tFvX5gKb7aPAPRflyTNxfI8FF516osvgG-ASMQukThnwrVaek-DgoYK7S-0P3QvObawsNXlOcR4hBhZZR94QkOrCaxw9e78pjNMox_bxECr88NbeeCs-9100oW7vANqJ1ap1RUJm-XdqHKU8iGf; 1P_JAR=2023-11-16-09; GOOGLE_ABUSE_EXEMPTION=ID=d7e21f1ed2280056:TM=1700126085:C=r:IP=84.17.39.171-:S=jGdYx_fA2CHiPSKkTZNNTwc; SIDCC=ACA-OxMA1sMYA-BMNNeooT4GDNbVZui2Kup9wfHLlKxyhrQ5QYN3tivPQBACEVTLaPTzJJ48JA; __Secure-1PSIDCC=ACA-OxPEX0QatTc8_7ig91noh0uKhs0Hw9rtmOGlafxCVfCYzekrUOG8bsWD8lsxuBQIPnZ72w; __Secure-3PSIDCC=ACA-OxOLRwyAhnd9pmXRP0QHjM5zlMrIRFhEnoujgL0TAMSqych_29DXjNYarSkh7dchruB-xQ; NID=511=sssLrAnTFS6eVsEt8fMhN69TQrqnqigmmewcfmFek7hhPHHpwYzlgWr8Nl4cfwJyc4iZJmvWvQkAp2_lcBq7cqiNqScfzH-l2qCMpBoCcTl3TF0mq9jjRgN6ZDq1WVpYBUOY1s59EAGd3XAZu3OZ6W3-inS-IC_wuE53nUQ-F2I; SIDCC=ACA-OxO_zI7-Z9VXCPGII9vba950WdUY7fIcIJEGKxzXesdG898SqPbyfNYufen9chSS5La8FQ; __Secure-1PSIDCC=ACA-OxP8EFSnNOa3C9_BJMds8xDUwW4dMpDN7NQYlEbsOz9gnfOsNEcBdMD0GoLOlb-W5ZUToQ; __Secure-3PSIDCC=ACA-OxMgzm51BGE2ncyKDzpBKg7v9ZmX8utw1GHBOLRqLy7GbE7kO3cTGtdMYmnnh_7Csd_jsA"
	url, err := repo.NewProductTypesPriceRepo(app.DB).WithSheetAPI(app.SheetAPI).GetSheetImage(585, 2, token, cookie)
	if err != nil {
		return
	}
	log.Println(url)
}

func TestProductTypesPriceRepo_SaveImageToS3(t *testing.T) {
	var app = initApp("local")
	url := "https://lh7-us.googleusercontent.com/u/0/sheets/ACTFsxS7kE6M9sW2zcuQImBh3qz2-Nfa35xW2TsiYKOObhuLgmXITKZ4bVVvYZENyNIbdURjTMYPry8U-q2vnOtYefmv-0_FdXz2_8AykqDVsWI2iuGWd7zvJkOxw4xNqsbne8iujo2-yBgVg1tOYwxAadNF6XFYcopIG8dCuCtT2JeOXcGVQN93FljfA65_iYObTtxU-I77ceHOdqz9URdxzXaA_YirxauoZnmd8jzDcp3gRll1vepvN5PCvsoSCRWio3hk4PYvefh2Ge_j9NHxWTZ2pO5U8t6t91JmewzjujprTw-gUpED23yK8h8rS4RK3lV-Eyn9w-A4MhO4r8g_xv0GEic0mVndHrrwidlwpbKAwVgVnQTPxemB-0twuCx7GBBCx_Kk7oqm74YoH1CITFwv4blu8T-UczgLcu-o-7YW1yZzM3NiY6vM0zlQQunHgoNRKLq9xgv37cvvy0tIJVS7NSQ5Eu7kWcugBDejmi8MaH75-aySB0_S0xC50MC-RUVTGerMG3rfOO3LfmpqY2R8xqBnpcoAkx_fsbO4CdmKEic37OG0nvfksIKEkyiCaw7ZFD5CWZ-tVHgCPJhgdjMs6ZppbN_yThdCACgQ9HvQCOP-cofBzm8p4qc1VECfdURA7hyaRs8aF4in7lPl7w1u_Qc-IQ1Eo9au4B-8rEJvIBW5iHBNTEdl4p9RB4ttnHV6A4Atg1LdzgVkn-NUflUuhMjZEeUastpPWS82l_m71oRKTQi-Uyy1-OV9ofyL-qwMR_Rf8T9GAtISUIsfhdJz6f-g4Y_FPEUS825gCcLCnF3T464xluyfleYEWX0Xf10IcVNk1x3x3DvxR-QsqX12lQDCWQGp0T3Z7CwPPROLk8CDZF3Sx5a_mK7DxHjJ0YY5jQe4ysVBGh7zf8o0uVyDNYeQ854rF3HcyWs5CZZK07AGSJ5_j6cYnqTxE7IWw6YF"
	imgURL, err := repo.NewProductTypesPriceRepo(app.DB).SaveImageToS3(585, 2, url)
	assert.NoError(t, err)
	helper.PrintJSON(imgURL)
}

func TestProductTypesPriceRepo_PaginateProductTypesPriceQuote(t *testing.T) {
	var app = initApp("dev")
	result := repo.NewProductTypesPriceRepo(app.DB).PaginateProductTypesPriceQuote(&models.PaginateProductTypesPriceQuoteParams{
		Product:     aws.String("CLOTHING"),
		Gender:      aws.String("WOMEN"),
		FabricType:  aws.String("KNIT"),
		Feature:     aws.String("ACTIVEWEAR"),
		Category:    aws.String("Activewear"),
		Item:        aws.String("Bikini"),
		Form:        aws.String(""),
		Description: aws.String("Ruffle"),
		Material:    aws.String("Knit Polo Interlock filagen FL4885-305081"),
		Composition: aws.String("POLY/ POLY BLEND"),
		Weight:      aws.Float64(140),
		CutWidth:    aws.Float64(170),
	})
	helper.PrintJSON(result)
}
