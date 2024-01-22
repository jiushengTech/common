package client

type Option func(o *HttpClient)

func WithBody(k string, v string) Option {
	return func(s *HttpClient) {
		s.body[k] = v
	}
}

func WithHeader(k string, v string) Option {
	return func(s *HttpClient) {
		s.header[k] = v
	}
}

//func WithTimeout(timeout time.Duration) Option {
//	return func(s *HttpClient) {
//		s.timeout = timeout
//	}
//}
