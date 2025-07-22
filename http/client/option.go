package client

import "time"

type Option func(o *HttpClient)

func WithBody(k string, v any) Option {
	return func(s *HttpClient) {
		s.body[k] = v
	}
}

func WithHeader(k string, v string) Option {
	return func(s *HttpClient) {
		s.header[k] = v
	}
}

// WithRedirectNum 设置HTTP重定向次数限制（默认3次）
func WithRedirectNum(redirectNum int) Option {
	return func(s *HttpClient) {
		s.redirectNum = redirectNum
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(s *HttpClient) {
		s.timeout = timeout
	}
}

// WithRetryCount 设置请求重试次数（默认0次，即不重试）
func WithRetryCount(retryCount int) Option {
	return func(s *HttpClient) {
		s.retryCount = retryCount
	}
}

// WithRetryDelay 设置重试间隔时间（默认1秒）
func WithRetryDelay(retryDelay time.Duration) Option {
	return func(s *HttpClient) {
		s.retryDelay = retryDelay
	}
}
