package client

import (
	"bytes"
	"encoding/json"
	"github.com/jiushengTech/common/log"
	"io"
	"net/http"
)

type HttpClient struct {
	header map[string]string
	body   map[string]string
}

func NewHttpClient(opts ...Option) *HttpClient {
	srv := &HttpClient{
		header: make(map[string]string),
		body:   make(map[string]string),
	}
	for _, o := range opts {
		o(srv)
	}
	return srv
}
func (c *HttpClient) Get(url string) (data []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Error("http get error:", err)
		return data, err
	}
	respBody, err := io.ReadAll(resp.Body)
	return respBody, err
}

// url第三方接口地址, t用于接收响应体的结构体
func (c *HttpClient) PostJSONWithResp(url string, t any) error {
	res, err := c.PostJSON(url)
	if err != nil {
		log.Errorf("Error sending request: %+v", err)
		return err
	}
	// 将响应体映射到形参 t
	err = json.Unmarshal(res, t)
	if err != nil {
		log.Errorf("Error unmarshalling response body: %+v", err)
		return err
	}
	return nil
}

func (c *HttpClient) PostJSON(url string) (data []byte, err error) {
	// 构造请求体
	requestBody, err := json.Marshal(c.body)
	if err != nil {
		log.Errorf("Error marshalling request body: %+v", err)
		return data, err
	}
	// 创建 HTTP 请求
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Errorf("Error creating request: %+v", err)
		return data, err
	}
	// 添加自定义请求头
	for key, value := range c.header {
		request.Header.Set(key, value)
	}
	// 发送 HTTP 请求
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Errorf("Error sending request: %+v", err)
		return data, err
	}
	defer response.Body.Close()
	// 读取响应体
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Errorf("Error reading response body: %+v", err)
		return data, err
	}
	return responseBody, err
}
