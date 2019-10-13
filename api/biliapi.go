package api

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type API interface {
	_get(string, time.Duration) (string, error)
	_post(string) (string, error)
}

// 发送GET请求
// url:请求地址
// response:请求返回的内容
func _get(url string, timeout time.Duration) (string, error) {
	client := http.Client{}
	client.Timeout = time.Second * timeout
	resp, err := client.Get(url)
	if err != nil {
		// panic(err)
		return "", err
	}

	defer resp.Body.Close()

	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}

	return result.String(), nil
}

// 发送POST请求
// url:请求地址，data:POST请求提交的数据,contentType:请求体格式，如：application/json
// content:请求放回的内容
func Post(url string, data interface{}, contentType string) string {
	jsonStr, _ := json.Marshal(data)
	req, err := http.NewRequest(`POST`, url, bytes.NewBuffer(jsonStr))
	req.Header.Add(`content-type`, contentType)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	return string(result)
}
