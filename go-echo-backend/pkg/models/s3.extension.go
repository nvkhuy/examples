package models

func (ct ContentType) GetExtension() string {
	var extension = ""
	switch ct {
	case ContentTypeImageJPG:
		extension = ".jpg"
	case ContentTypeImagePNG:
		extension = ".png"
	case ContentTypeImageJPEG:
		extension = ".jpeg"
	case ContentTypePDF:
		extension = ".pdf"
	case ContentTypeVideoMP4:
		extension = ".mp4"
	case ContentTypeVideoMOV:
		extension = ".mov"
	case ContentTypeVideoQuickTime, ContentTypeVideoXQuickTime:
		extension = ".mov"
	case ContentTypeAudioMPEG3, ContentTypeAudioXMPEG3:
		extension = ".mp3"
	case ContentTypeAudioWav, ContentTypeAudioXWav:
		extension = ".wav"
	case ContentTypeHTML:
		extension = ".html"
	case ContentTypeXLSX:
		extension = ".xlsx"
	case ContentTypeCSV:
		extension = ".csv"
	}

	return extension
}
