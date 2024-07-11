package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func ReadFrame(inputFile string) ([]byte, error) {
	var buf = bytes.NewBuffer(nil)

	err := ffmpeg.Input(inputFile).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", 1)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func MovToMP4(inputFile string) ([]byte, error) {
	var buf = bytes.NewBuffer(nil)
	var err = ffmpeg.Input(inputFile).
		Output("pipe:", ffmpeg.KwArgs{"q:v": "0"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

func ReadFrameAndUpload(inputFile string, fileKey string) (string, error) {
	var buf = bytes.NewBuffer(nil)

	err := ffmpeg.Input(inputFile).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", 1)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		return "", err
	}

	var ext = path.Ext(fileKey)
	var thumbnailFileKey = strings.ReplaceAll(fileKey, ext, ".jpeg")

	_, err = UploadS3(thumbnailFileKey, buf, http.DetectContentType(buf.Bytes()))
	if err != nil {
		return "", err
	}

	return thumbnailFileKey, nil
}
