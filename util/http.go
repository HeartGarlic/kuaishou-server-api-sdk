package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// PostForm post form 数据请求
func PostForm(uri string, obj url.Values) ([]byte, error) {
	response, err := http.PostForm(uri, obj)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http get error : uri=%v , statusCode=%v", uri, response.StatusCode)
	}
	return ioutil.ReadAll(response.Body)
}

// PostJSON post json 数据请求
func PostJSON(uri string, obj interface{}) ([]byte, error) {
	marshal, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	response, err := http.Post(uri, "application/json;charset=utf-8", bytes.NewBuffer(marshal))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http get error : uri=%v , statusCode=%v", uri, response.StatusCode)
	}
	return ioutil.ReadAll(response.Body)
}

// JsonStructToMap ...
func JsonStructToMap(content interface{}) (map[string]interface{}, error) {
	var name map[string]interface{}
	if marshalContent, err := json.Marshal(content); err != nil {
		return name, err
	} else {
		d := json.NewDecoder(bytes.NewReader(marshalContent))
		d.UseNumber() // 设置将float64转为一个number
		if err := d.Decode(&name); err != nil {
			return name, err
		} else {
			for k, v := range name {
				name[k] = v
			}
		}
	}
	return name, nil
}
