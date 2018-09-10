package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

//Get get
func Get(url string, params map[string]interface{}) (*http.Response, error) {
	joinStr := JoinParams(params)
	if joinStr != "" {
		url = url + "?" + joinStr
	}
	return http.Get(url)
}

//PostJSON post json
func PostJSON(url string, data map[string]interface{}) (*http.Response, error) {

	bytesData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(bytesData)
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

//JoinParams params
func JoinParams(params map[string]interface{}) string {
	var results = []string{}
	for key, value := range params {
		var str = key + "=" + fmt.Sprintf("%s", value)
		results = append(results, str)
	}
	res := strings.Join(results, "&")
	return res
}

//ResponseToInterface response to model
func ResponseToInterface(r *http.Response, model *map[string]interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, model)
	if err != nil {
		return err
	}
	return nil
}
