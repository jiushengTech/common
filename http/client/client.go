package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type HttpClient struct {
	header      map[string]string
	body        map[string]any
	timeout     time.Duration
	redirectNum int           // HTTP重定向次数限制
	retryCount  int           // 请求重试次数
	retryDelay  time.Duration // 重试间隔时间
}

func (c *HttpClient) getClient() *http.Client {
	client := &http.Client{
		Transport: nil,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= c.redirectNum {
				return errors.New("stopped after redirects limit")
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
		header:      make(map[string]string),
		body:        make(map[string]any),
		timeout:     30 * time.Second, // 默认超时时间
		redirectNum: 3,                // 默认重定向次数限制
		retryCount:  0,                // 默认不重试
		retryDelay:  1 * time.Second,  // 默认重试间隔1秒
	}
	for _, o := range opts {
		o(srv)
	}
	return srv
}

// doRequestWithRetry 执行带重试的HTTP请求
func (c *HttpClient) doRequestWithRetry(req *http.Request) (*http.Response, error) {
	client := c.getClient()
	var lastErr error

	for attempt := 0; attempt <= c.retryCount; attempt++ {
		if attempt > 0 {
			// 等待重试间隔
			time.Sleep(c.retryDelay)
			log.Printf("重试第 %d 次请求: %s %s", attempt, req.Method, req.URL.String())
		}

		resp, err := client.Do(req)
		if err == nil {
			return resp, nil
		}

		lastErr = err
		// 如果这是最后一次尝试，不再重试
		if attempt == c.retryCount {
			break
		}
	}

	return nil, fmt.Errorf("请求失败，已重试 %d 次: %w", c.retryCount, lastErr)
}

// GetWithResp 发送GET请求并解析响应到指定结构体
func (c *HttpClient) GetWithResp(url string, t any) error {
	res, err := c.Get(url)
	if err != nil {
		return fmt.Errorf("HTTP GET 请求失败: %w", err)
	}

	err = json.Unmarshal(res, t)
	if err != nil {
		return fmt.Errorf("响应解析失败: %w\n响应体内容: %s", err, string(res))
	}

	return nil
}

// Get 发送GET请求并返回原始响应体
func (c *HttpClient) Get(url string) ([]byte, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建GET请求失败: %w", err)
	}

	// 添加自定义请求头
	for key, value := range c.header {
		request.Header.Set(key, value)
	}

	// 发送请求（带重试）
	resp, err := c.doRequestWithRetry(request)
	if err != nil {
		return nil, fmt.Errorf("发送GET请求失败: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("关闭响应体失败: %v", closeErr)
		}
	}()

	// 检查HTTP状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	return respBody, nil
}

// PostJSONWithResp 发送POST JSON请求并解析响应到指定结构体
func (c *HttpClient) PostJSONWithResp(url string, t any) error {
	res, err := c.PostJSON(url)
	if err != nil {
		return fmt.Errorf("HTTP POST JSON 请求失败: %w", err)
	}

	err = json.Unmarshal(res, t)
	if err != nil {
		return fmt.Errorf("响应解析失败: %w\n响应体内容: %s", err, string(res))
	}

	return nil
}

// PostJSON 发送POST JSON请求并返回原始响应体
func (c *HttpClient) PostJSON(url string) ([]byte, error) {
	// 构造请求体
	requestBody, err := json.Marshal(c.body)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %w", err)
	}

	// 创建HTTP请求
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("创建POST请求失败: %w", err)
	}

	// 设置Content-Type为JSON
	request.Header.Set("Content-Type", "application/json")

	// 添加自定义请求头
	for key, value := range c.header {
		request.Header.Set(key, value)
	}

	// 发送HTTP请求（带重试）
	response, err := c.doRequestWithRetry(request)
	if err != nil {
		return nil, fmt.Errorf("发送POST请求失败: %w", err)
	}
	defer func() {
		if closeErr := response.Body.Close(); closeErr != nil {
			log.Printf("关闭响应体失败: %v", closeErr)
		}
	}()

	// 检查HTTP状态码
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("请求失败，状态码: %d", response.StatusCode)
	}

	// 读取响应体
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	return responseBody, nil
}

// PostFormWithResp 发送POST表单请求并解析响应到指定结构体
func (c *HttpClient) PostFormWithResp(url string, formData map[string]string, t any) error {
	res, err := c.PostForm(url, formData)
	if err != nil {
		return fmt.Errorf("HTTP POST Form 请求失败: %w", err)
	}

	err = json.Unmarshal(res, t)
	if err != nil {
		return fmt.Errorf("响应解析失败: %w\n响应体内容: %s", err, string(res))
	}

	return nil
}

// PostForm 发送POST表单请求并返回原始响应体
func (c *HttpClient) PostForm(url string, formData map[string]string) ([]byte, error) {
	// 构造表单数据
	form := make(map[string][]string)
	for k, v := range formData {
		form[k] = []string{v}
	}

	// 创建HTTP请求
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(encodeForm(formData)))
	if err != nil {
		return nil, fmt.Errorf("创建POST Form请求失败: %w", err)
	}

	// 设置Content-Type为表单
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 添加自定义请求头
	for key, value := range c.header {
		request.Header.Set(key, value)
	}

	// 发送HTTP请求（带重试）
	response, err := c.doRequestWithRetry(request)
	if err != nil {
		return nil, fmt.Errorf("发送POST Form请求失败: %w", err)
	}
	defer func() {
		if closeErr := response.Body.Close(); closeErr != nil {
			log.Printf("关闭响应体失败: %v", closeErr)
		}
	}()

	// 检查HTTP状态码
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("请求失败，状态码: %d", response.StatusCode)
	}

	// 读取响应体
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	return responseBody, nil
}

// PutJSONWithResp 发送PUT JSON请求并解析响应到指定结构体
func (c *HttpClient) PutJSONWithResp(url string, t any) error {
	res, err := c.PutJSON(url)
	if err != nil {
		return fmt.Errorf("HTTP PUT JSON 请求失败: %w", err)
	}

	err = json.Unmarshal(res, t)
	if err != nil {
		return fmt.Errorf("响应解析失败: %w\n响应体内容: %s", err, string(res))
	}

	return nil
}

// PutJSON 发送PUT JSON请求并返回原始响应体
func (c *HttpClient) PutJSON(url string) ([]byte, error) {
	// 构造请求体
	requestBody, err := json.Marshal(c.body)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %w", err)
	}

	// 创建HTTP请求
	request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("创建PUT请求失败: %w", err)
	}

	// 设置Content-Type为JSON
	request.Header.Set("Content-Type", "application/json")

	// 添加自定义请求头
	for key, value := range c.header {
		request.Header.Set(key, value)
	}

	// 发送HTTP请求（带重试）
	response, err := c.doRequestWithRetry(request)
	if err != nil {
		return nil, fmt.Errorf("发送PUT请求失败: %w", err)
	}
	defer func() {
		if closeErr := response.Body.Close(); closeErr != nil {
			log.Printf("关闭响应体失败: %v", closeErr)
		}
	}()

	// 检查HTTP状态码
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("请求失败，状态码: %d", response.StatusCode)
	}

	// 读取响应体
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	return responseBody, nil
}

// DeleteWithResp 发送DELETE请求并解析响应到指定结构体
func (c *HttpClient) DeleteWithResp(url string, t any) error {
	res, err := c.Delete(url)
	if err != nil {
		return fmt.Errorf("HTTP DELETE 请求失败: %w", err)
	}

	err = json.Unmarshal(res, t)
	if err != nil {
		return fmt.Errorf("响应解析失败: %w\n响应体内容: %s", err, string(res))
	}

	return nil
}

// Delete 发送DELETE请求并返回原始响应体
func (c *HttpClient) Delete(url string) ([]byte, error) {
	request, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建DELETE请求失败: %w", err)
	}

	// 添加自定义请求头
	for key, value := range c.header {
		request.Header.Set(key, value)
	}

	// 发送请求（带重试）
	resp, err := c.doRequestWithRetry(request)
	if err != nil {
		return nil, fmt.Errorf("发送DELETE请求失败: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("关闭响应体失败: %v", closeErr)
		}
	}()

	// 检查HTTP状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	return respBody, nil
}

// SetBody 设置请求体数据
func (c *HttpClient) SetBody(key string, value any) *HttpClient {
	c.body[key] = value
	return c
}

// SetHeader 设置请求头
func (c *HttpClient) SetHeader(key, value string) *HttpClient {
	c.header[key] = value
	return c
}

// ClearBody 清空请求体
func (c *HttpClient) ClearBody() *HttpClient {
	c.body = make(map[string]any)
	return c
}

// ClearHeaders 清空请求头
func (c *HttpClient) ClearHeaders() *HttpClient {
	c.header = make(map[string]string)
	return c
}

// encodeForm 编码表单数据
func encodeForm(data map[string]string) string {
	if len(data) == 0 {
		return ""
	}

	values := url.Values{}
	for key, value := range data {
		values.Set(key, value)
	}

	return values.Encode()
}
