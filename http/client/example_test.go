package client

import (
	"fmt"
	"testing"
	"time"
)

// ExampleResponse 示例响应结构体
type ExampleResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// TestHttpClientUsage 展示HttpClient的各种用法
func TestHttpClientUsage(t *testing.T) {
	// 创建HTTP客户端
	client := NewHttpClient(
		WithTimeout(10*time.Second),
		WithRedirectNum(5),            // 设置重定向次数限制
		WithRetryCount(3),             // 设置重试次数
		WithRetryDelay(2*time.Second), // 设置重试间隔
		WithHeader("User-Agent", "MyApp/1.0"),
		WithHeader("Accept", "application/json"),
	)

	// 示例1: GET请求 - 返回原始响应
	fmt.Println("=== GET请求示例（带重试） ===")
	data, err := client.Get("https://httpbin.org/get")
	if err != nil {
		t.Logf("GET请求失败: %v", err)
	} else {
		t.Logf("GET响应长度: %d", len(data))
	}

	// 示例2: GET请求 - 解析到结构体
	var getResp ExampleResponse
	err = client.GetWithResp("https://httpbin.org/get", &getResp)
	if err != nil {
		t.Logf("GET结构体解析可能失败(这是正常的，因为httpbin返回格式不同): %v", err)
	}

	// 示例3: POST JSON请求
	fmt.Println("=== POST JSON请求示例（带重试） ===")
	client.SetBody("name", "张三").
		SetBody("age", 25).
		SetBody("email", "zhangsan@example.com")

	postData, err := client.PostJSON("https://httpbin.org/post")
	if err != nil {
		t.Logf("POST JSON请求失败: %v", err)
	} else {
		t.Logf("POST JSON响应长度: %d", len(postData))
	}

	// 示例4: POST表单请求
	fmt.Println("=== POST表单请求示例（带重试） ===")
	formData := map[string]string{
		"username": "testuser",
		"password": "testpass",
	}

	formResp, err := client.PostForm("https://httpbin.org/post", formData)
	if err != nil {
		t.Logf("POST表单请求失败: %v", err)
	} else {
		t.Logf("POST表单响应长度: %d", len(formResp))
	}

	// 示例5: 测试重试功能 - 使用一个可能失败的URL
	fmt.Println("=== 重试功能测试 ===")
	retryClient := NewHttpClient(
		WithTimeout(3*time.Second),
		WithRetryCount(2),             // 重试2次
		WithRetryDelay(1*time.Second), // 每次重试间隔1秒
	)

	// 使用一个可能超时或失败的URL来测试重试
	_, err = retryClient.Get("https://httpbin.org/delay/5") // 这个会超时，触发重试
	if err != nil {
		t.Logf("重试测试完成，最终失败（这是预期的）: %v", err)
	}

	// 示例6: 不同的重试配置
	fmt.Println("=== 快速重试配置示例 ===")
	fastRetryClient := NewHttpClient(
		WithTimeout(5*time.Second),
		WithRetryCount(1),                    // 只重试1次
		WithRetryDelay(500*time.Millisecond), // 重试间隔500毫秒
	)

	data, err = fastRetryClient.Get("https://httpbin.org/get")
	if err != nil {
		t.Logf("快速重试配置请求失败: %v", err)
	} else {
		t.Logf("快速重试配置请求成功，响应长度: %d", len(data))
	}

	// 示例7: 链式调用配置
	fmt.Println("=== 链式调用示例 ===")
	client.ClearHeaders().
		ClearBody().
		SetHeader("Authorization", "Bearer token123").
		SetHeader("Content-Type", "application/json").
		SetBody("action", "chain_test").
		SetBody("data", map[string]interface{}{
			"key1": "value1",
			"key2": 123,
		})

	t.Log("链式调用设置完成")
}

// TestHttpClientRetryDemo 专门演示重试功能
func TestHttpClientRetryDemo(t *testing.T) {
	t.Log("=== HTTP客户端重试功能演示 ===")

	// 配置说明
	t.Log("配置说明:")
	t.Log("- redirectNum: HTTP重定向次数限制（默认3次）")
	t.Log("- retryCount: 请求失败时的重试次数（默认0次，即不重试）")
	t.Log("- retryDelay: 每次重试之间的间隔时间（默认1秒）")

	// 创建有重试功能的客户端
	client := NewHttpClient(
		WithTimeout(2*time.Second),    // 较短的超时时间，容易触发重试
		WithRetryCount(3),             // 重试3次
		WithRetryDelay(1*time.Second), // 每次重试间隔1秒
		WithRedirectNum(5),            // 允许5次重定向
	)

	t.Log("开始测试重试功能...")
	start := time.Now()

	// 这个请求应该会成功
	_, err := client.Get("https://httpbin.org/get")
	duration := time.Since(start)

	if err != nil {
		t.Logf("请求失败: %v (耗时: %v)", err, duration)
	} else {
		t.Logf("请求成功 (耗时: %v)", duration)
	}
}

// BenchmarkHttpClient 性能测试
func BenchmarkHttpClient(b *testing.B) {
	client := NewHttpClient(
		WithTimeout(5*time.Second),
		WithRetryCount(0), // 性能测试时不使用重试
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.Get("https://httpbin.org/get")
		if err != nil {
			b.Errorf("请求失败: %v", err)
		}
	}
}
