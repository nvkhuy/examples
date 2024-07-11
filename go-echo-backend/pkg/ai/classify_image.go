package ai

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type ImageClassifyResponse struct {
	Image struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"image"`
	Predictions []struct {
		X          float64 `json:"x"`
		Y          float64 `json:"y"`
		Width      int     `json:"width"`
		Height     int     `json:"height"`
		Confidence float64 `json:"confidence"`
		Class      string  `json:"class"`
		ClassID    int     `json:"class_id"`
	} `json:"predictions"`
}

type ImageClassifyParams struct {
	Image      string  `json:"image"`
	Size       int64   `json:"size"`
	Confidence float64 `json:"confidence"`
	Overlap    float64 `json:"overlap"`
}

func ClassifyImage(params ImageClassifyParams) (resp ImageClassifyResponse, err error) {
	url := "https://dev-classify.joininflow.io/classify"
	method := "POST"

	// Convert params to JSON
	jsonParams, err := json.Marshal(params)
	if err != nil {
		return
	}
	payload := bytes.NewBuffer(jsonParams)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &resp)
	return
}

// Current Maximum Concurrency is 6
// AWS Configure Concurrency
// https://ap-southeast-1.console.aws.amazon.com/lambda/home?region=ap-southeast-1#/functions/inflow-dev-classify?tab=configure
func ClassifyMultiImage(params []ImageClassifyParams) (resp []ImageClassifyResponse, err error) {
	type data struct {
		Resp ImageClassifyResponse
		Err  error
	}
	f := func(params ImageClassifyParams, c chan data) {
		_resp, _err := ClassifyImage(params)
		c <- data{
			Resp: _resp,
			Err:  _err,
		}
	}

	num := len(params)
	classifyChan := make(chan data)
	for i := 0; i < num; i++ {
		go f(params[i], classifyChan)
	}
	for i := 0; i < num; i++ {
		v := <-classifyChan
		if v.Err == nil {
			resp = append(resp, v.Resp)
		}
	}
	return
}
