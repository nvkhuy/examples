package downloader

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"github.com/engineeringinflow/inflow-backend/pkg/retry"
	"github.com/rotisserie/eris"
)

func DownloadFile(URL string) ([]byte, error) {
	var data bytes.Buffer

	var err = retry.Exec(3, func() error {
		response, err := http.Get(URL)
		if err != nil {
			return err
		}
		defer response.Body.Close()
		if response.StatusCode != http.StatusOK {
			return errors.New(response.Status)
		}

		_, err = io.Copy(&data, response.Body)
		if err != nil {
			return err
		}

		if data.Len() < 1024 {
			return eris.Errorf("data size is too small, might be not a file")
		}
		return nil
	})
	return data.Bytes(), err
}

func DownloadMultipleFiles(urls []string) ([]io.ReadSeeker, error) {
	done := make(chan []byte, len(urls))
	errch := make(chan error, len(urls))
	for _, URL := range urls {
		go func(URL string) {
			b, err := DownloadFile(URL)
			if err != nil {
				errch <- err
				done <- nil
				return
			}
			done <- b
			errch <- nil
		}(URL)
	}
	bytesArray := make([]io.ReadSeeker, 0)
	var errStr string
	for i := 0; i < len(urls); i++ {

		if err := <-errch; err != nil {
			errStr = errStr + " " + err.Error()
		} else {
			bytesArray = append(bytesArray, bytes.NewReader(<-done))
		}
	}
	var err error
	if errStr != "" {
		err = errors.New(errStr)
	}
	return bytesArray, err
}
