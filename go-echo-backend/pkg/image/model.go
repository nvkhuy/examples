package image

// import (
// 	"strconv"
// 	"strings"
// )

// type ThumbnailSize string

// type Dimension struct {
// 	Width  int64 `json:"width"`
// 	Height int64 `json:"height"`
// }

// func (s ThumbnailSize) GetDimension() *Dimension {
// 	var dimension = new(Dimension)
// 	if strings.Contains(string(s), "w") {
// 		if w, err := strconv.ParseInt(strings.ReplaceAll(string(s), "w", ""), 10, 64); err == nil {
// 			dimension.Width = w
// 		}
// 	}

// 	if strings.Contains(string(s), "h") {
// 		if h, err := strconv.ParseInt(strings.ReplaceAll(string(s), "h", ""), 10, 64); err == nil {
// 			dimension.Height = h
// 		}
// 	}

// 	if strings.Contains(string(s), "x") {
// 		var parts = strings.Split(string(s), "x")
// 		if len(parts) == 2 {
// 			if w, err := strconv.ParseInt(parts[0], 10, 64); err == nil {
// 				dimension.Width = w
// 			}

// 			if h, err := strconv.ParseInt(parts[1], 10, 64); err == nil {
// 				dimension.Height = h
// 			}
// 		}
// 	}

// 	return dimension
// }

// func (d Dimension) IsValid() bool {
// 	return d.Width > 0 || d.Height > 0
// }
