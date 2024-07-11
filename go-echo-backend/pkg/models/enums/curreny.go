package enums

type Currency string

var (
	USD Currency = "USD"
	SGD Currency = "SGD"
	VND Currency = "VND"
)

func (c Currency) DefaultIfInvalid() Currency {
	var result = c
	if result == "" {
		result = USD
	}

	return result
}

func (c Currency) GetCountryCode() CountryCode {
	switch c {

	default:
		return CountryCodeUS
	}
}

/*
af, ar, az, be, bg, bn, bs, ca, cs, cy, da, de, de-AT, de-CH, de-DE, el, el-CY, en, en-AU, en-CA, en-GB, en-IE, en-IN, en-NZ, en-US, en-ZA, en-CY, en-TT,
eo, es, es-419, es-AR, es-CL, es-CO, es-CR, es-EC, es-ES, es-MX, es-NI, es-PA, es-PE, es-US, es-VE, et, eu, fa, fi, fr, fr-CA, fr-CH, fr-FR, gl, he,
hi, hi-IN, hr, hu, id, is, it, it-CH, ja, ka, km, kn, ko, lb, lo, lt, lv, mk, ml, mn, mr-IN, ms, nb, ne, nl, nn, oc, or, pa, pl, pt, pt-BR, rm, ro, ru,
sk, sl, sq, sr, st, sw, ta, te, th, tl, tr, tt, ug, ur, uz, vi, wo, zh-CN, zh-HK, zh-TW, zh-YUE
*/
func (c Currency) GetCustomerIOCode() string {
	switch c {
	case VND:
		return "vi"
	default:
		return "en-US"
	}
}
