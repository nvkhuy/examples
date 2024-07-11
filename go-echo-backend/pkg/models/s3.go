package models

type ContentType string

var (
	ContentTypeImagePNG  ContentType = "image/png"
	ContentTypeImageJPG  ContentType = "image/jpg"
	ContentTypeImageJPEG ContentType = "image/jpeg"

	ContentTypePDF ContentType = "application/pdf"

	ContentTypeVideoMOV        ContentType = "video/mov"
	ContentTypeVideoMP4        ContentType = "video/mp4"
	ContentTypeVideoQuickTime  ContentType = "video/quicktime"
	ContentTypeVideoXQuickTime ContentType = "video/x-quicktime"

	ContentTypeAudioMPEG   ContentType = "audio/mpeg"
	ContentTypeAudioMPEG3  ContentType = "audio/mpeg3"
	ContentTypeAudioXMPEG3 ContentType = "audio/x-mpeg-3"
	ContentTypeAudioWav    ContentType = "audio/wav"
	ContentTypeAudioXWav   ContentType = "audio/x-wav"

	ContentTypeHTML ContentType = "text/html"
	ContentTypeXLSX ContentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	ContentTypeCSV  ContentType = "text/csv"
)

type S3SignatureForm struct {
	ContentType ContentType `json:"content_type" validate:"required"`
	Resource    string      `json:"resource"`
}

type S3SignatureForms struct {
	Records []*S3SignatureForm `jon:"records" validate:"dive,required"`
}

func (ct ContentType) IsMedia() bool {
	return true
}

func GetContentTypeFromExt(ext string) ContentType {
	switch ext {
	case ".png":
		return ContentTypeImagePNG

	case ".jpeg":
		return ContentTypeImageJPEG
	}

	return ContentTypeImageJPG
}
