package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"path"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"

	"github.com/rs/xid"
	uuid "github.com/satori/go.uuid"
	"github.com/thaitanloi365/go-utils/random"
)

func StringContains(list []string, value string) bool {
	for _, c := range list {
		if c == value {
			return true
		}
	}

	return false
}

func StringFoldContains(list []string, value string) bool {
	for _, c := range list {
		if strings.EqualFold(c, value) {
			return true
		}
	}

	return false
}

func SplitFirstAndLastName(name string) (firstName string, lastName string) {
	var names = strings.Split(name, " ")
	var lastNameParts []string
	for index, v := range names {
		if index == 0 {
			firstName = v
		} else {
			lastNameParts = append(lastNameParts, v)
		}
	}
	lastName = strings.Join(lastNameParts, " ")
	return
}

func GetColumnOfTable(item interface{}, tableName, prefix string, ignoreKey ...string) []string {
	var tags []string

	bytes, err := json.Marshal(item)
	if err != nil {
		return tags
	}

	var values map[string]interface{}
	err = json.Unmarshal(bytes, &values)
	if err != nil {
		return tags
	}

	for key := range values {
		var constains = StringContains(ignoreKey, key)
		if !constains {
			tags = append(tags, fmt.Sprintf("%s.%s AS %s%s", tableName, key, prefix, key))
		}
	}

	return tags
}

func GetRequestID() string {
	return random.String(32, "abcdefghijklmnopqrstuvwxyz"+random.Numerals)
}

func CompareBoolEqual(a *bool, b *bool) bool {
	return a != nil && b != nil && *a == *b
}

func GenerateUUID() string {
	id, err := uuid.NewV4()
	if err != nil {
		return ""
	}
	return id.String()
}

func GenerateXID() string {
	return xid.New().String()
}

func ParseToStruct(src interface{}, dest interface{}) error {
	bytes, err := json.Marshal(src)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, dest)
	return err

}

func GetTimeout(value, defaultValue time.Duration) time.Duration {
	if value > time.Second*30 {
		return value
	}

	return GetTimeout(defaultValue, time.Minute)
}

func ParseDate(dateTimeString string, loc *time.Location, fallback ...time.Time) time.Time {
	var fb = time.Now().In(loc)
	if len(fallback) > 0 {
		fb = fallback[0]
	}
	t, err := dateparse.ParseIn(dateTimeString, loc)
	if err != nil {
		return fb
	}

	return t
}

func GetTuningPoolSize(total int, maxSize ...int) int {
	if total == 0 {
		return 1
	}

	if total < 10 {
		return total
	}
	var v = total / 10

	var size = int(math.Round(float64(total/v))) + v/2

	if len(maxSize) > 0 {
		return int(math.Min(float64(size), float64(maxSize[0])))
	}

	return size
}

func GetContentType(seeker io.ReadSeeker) (string, error) {
	// At most the first 512 bytes of data are used:
	// https://golang.org/src/net/http/sniff.go?s=646:688#L11
	buff := make([]byte, 512)

	_, err := seeker.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	bytesRead, err := seeker.Read(buff)
	if err != nil && err != io.EOF {
		return "", err
	}

	// Slice to remove fill-up zero values which cause a wrong content type detection in the next step
	buff = buff[:bytesRead]

	return http.DetectContentType(buff), nil
}

func IsPDF(seeker io.ReadSeeker) bool {
	contentType, err := GetContentType(seeker)
	if err != nil {
		return false
	}

	return contentType == "application/pdf"
}

func PrintJSON(i interface{}) {
	data, _ := json.MarshalIndent(i, "", "   ")
	fmt.Println(string(data))
}

func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

func Distance(lat1, lon1, lat2, lon2 float64) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180

	r = 6378.1 // Earth radius in KM

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}

func StringContainsAny(list []string, value string) bool {
	for _, c := range list {
		if strings.ContainsAny(c, value) {
			return true
		}
	}

	return false
}

// Distance between 2 points
func DistanceInMeter(lat1, lon1, lat2, lon2 float64) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}

func ToJson(st interface{}) []byte {
	data, _ := json.Marshal(st)
	return data
}

func StringToInt(v string) int {
	intV, _ := strconv.Atoi(v)

	return intV
}

func IsLat(lat *float64) bool {
	return lat != nil && math.Abs(*lat) <= 90
}

func IsLng(lng *float64) bool {
	return lng != nil && math.Abs(*lng) <= 180
}

func IsLatLng(lat *float64, lng *float64) bool {
	return IsLat(lat) && IsLng(lng)
}

// func GetTaxPercentage(cc enums.Currency) float64 {
// 	switch cc {
// 	case enums.SGD:
// 		return 0

// 	case enums.VND:
// 		if time.Now().Year() >= 2024 {
// 			return 0.1
// 		}
// 		return 0.08
// 	}

// 	return 0
// }

func GetFuncName(frame int) string {
	pc, _, _, _ := runtime.Caller(frame)
	var parts = strings.Split(runtime.FuncForPC(pc).Name(), ".")

	if len(parts) > 0 {
		return parts[len(parts)-1]
	}

	return strings.Join(parts, ".")
}

func JoinNonEmptyStrings(separator string, s ...string) string {
	var list []string

	for _, v := range s {
		if v != "" {
			list = append(list, v)
		}
	}

	return strings.Join(list, separator)
}

func IsImageExt(fileKey string) bool {
	var ext = path.Ext(fileKey)
	var listExts = []string{
		".png",
		".jpg",
		".jpeg",
		".webp",
		".svg",
		".bmp",
		".avif",
	}

	for _, v := range listExts {
		if v == ext {
			return true
		}
	}

	return false
}

func IsVideoExt(fileKey string) bool {
	var ext = path.Ext(fileKey)
	var listExts = []string{
		".mp4",
		".mov",
	}

	for _, v := range listExts {
		if strings.Contains(ext, v) {
			return true
		}
	}

	return false
}

func PrintJSONBytes(src []byte) {
	var out = bytes.NewBuffer(nil)
	json.Indent(out, src, "", "   ")
	fmt.Println(out.String())
}

func DownloadImageFromURL(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch image, status code: %d", response.StatusCode)
	}

	imageData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return imageData, nil
}

func GetExtensionFromURL(imageURL string) string {
	return path.Ext(imageURL)

}

func StructToMap(item interface{}, tag ...string) map[string]interface{} {
	var defaultTag = "json"
	if len(tag) > 0 {
		defaultTag = tag[0]
	}

	res := map[string]interface{}{}
	if item == nil {
		return res
	}
	v := reflect.TypeOf(item)
	reflectValue := reflect.ValueOf(item)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		tag := v.Field(i).Tag.Get(defaultTag)

		// remove omitEmpty
		omitEmpty := false
		if strings.HasSuffix(tag, "omitempty") {
			omitEmpty = true
			idx := strings.Index(tag, ",")
			if idx > 0 {
				tag = tag[:idx]
			} else {
				tag = ""
			}
		}

		if !reflectValue.Field(i).CanInterface() {
			continue
		}

		field := reflectValue.Field(i).Interface()
		if tag != "" && tag != "-" {
			if v.Field(i).Type.Kind() == reflect.Struct {
				res[tag] = StructToMap(field)
			} else {
				if !(omitEmpty && reflectValue.Field(i).IsZero()) {
					res[tag] = field
				}
			}
		}
	}
	return res
}
