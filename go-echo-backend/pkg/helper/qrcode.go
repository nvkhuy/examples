package helper

import (
	"bytes"
	"io"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }

type GenerateQRCodeOptions struct {
	Content       string
	EncodeOptions []qrcode.EncodeOption
	ImageOptions  []standard.ImageOption
}

func GenerateQRCode(opts GenerateQRCodeOptions) (*bytes.Buffer, error) {
	var options = []qrcode.EncodeOption{
		qrcode.WithEncodingMode(qrcode.EncModeByte),
		qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionQuart),
	}
	options = append(options, opts.EncodeOptions...)
	qrc, err := qrcode.NewWith(opts.Content, options...)
	if err != nil {
		return nil, err
	}

	var buf = bytes.NewBuffer(nil)
	var wr = nopCloser{Writer: buf}
	var imageOptions = []standard.ImageOption{
		standard.WithQRWidth(10),
	}
	imageOptions = append(imageOptions, opts.ImageOptions...)
	w := standard.NewWithWriter(wr, imageOptions...)
	if err = qrc.Save(w); err != nil {
		return nil, err
	}

	return buf, err
}
