package client

import (
	"bytes"
	"encoding/json"
	"errors"
	log "github.com/jiushengTech/common/log/zap/logger"
	"io"
	"net/http"
	"time"
)

type HttpClient struct {
	header      map[string]string
	body        map[string]any
	timeout     time.Duration
	redirectNum int
}

func (c *HttpClient) getClient() *http.Client {
	client := &http.Client{
		Transport: nil,
		CheckRedirect: func(req *http.Request, via []*http.Request) (err error) {
			if len(via) >= c.redirectNum {
				return errors.New("stopped after 3 redirects")
			}
			return nil
		},
		Jar:     nil,
		Timeout: c.timeout,
	}
	return client
}

func NewHttpClient(opts ...Option) *HttpClient {
	srv := &HttpClient{
		header: make(map[string]string),
		body:   make(map[string]any),
	}
	for _, o := range opts {
		o(srv)
	}
	return srv
}

//func (c *HttpClient) Get(url string) (data []byte, err error) {
//	resp, err := http.Get(url)
//	if err != nil {
//		log.Error("http get error:", err)
//		return data, err
//	}
//	respBody, err := io.ReadAll(resp.Body)
//	return respBody, err
//}

func (c *HttpClient) GetWithResp(url string, t any) (err error) {
	res, err := c.Get(url)
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
	return err
}

func (c *HttpClient) Get(url string) (data []byte, err error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("Error creating request: %+v", err)
		return data, err
	}
	// 添加自定义请求头
	for key, value := range c.header {
		request.Header.Set(key, value)
	}
	// 获取客户端
	client := c.getClient()
	resp, err := client.Do(request)
	if err != nil {
		log.Errorf("Error sending request: %+v", err)
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
	// 获取客户端
	client := c.getClient()
	response, err := client.Do(request)
	if err != nil {
		log.Errorf("Error sending request: %+v", err)
		return data, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error("http client close error:", err)
		}
	}(response.Body)
	// 读取响应体
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Errorf("Error reading response body: %+v", err)
		return data, err
	}
	return responseBody, err
}
