package httpclient

import "net/http"

type HTTPClient struct {
	client *http.Client
}

func CreateClient() *HTTPClient {
	return &HTTPClient{client: &http.Client{}}
}

func (client *HTTPClient) CreateNewRequest(method, url string, header map[string]string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return req, err
	}
	for key, value := range header {
		req.Header.Set(key, value)
	}
	return req, nil
}

func (client *HTTPClient) DoRequest(req *http.Request) (*http.Response, error) {
	res, err := client.client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
