package enums

type DefaultImageType string

var (
	DefaultUserAvatar DefaultImageType = "default_user_avatar"
)

func (p DefaultImageType) String() string {
	return string(p)
}

func (p DefaultImageType) URL() string {
	var url = string(p)

	switch p {
	case DefaultUserAvatar:
		url = "https://dev-static.joininflow.io/common/default_category_icon.png"
	}

	return url
}
