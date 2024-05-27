package http

import (
	"crypto/tls"
	"time"

	"github.com/go-resty/resty/v2"
)

func GetWithBasicAuth(url string, username, password string) (*resty.Response, error) {
	client := newClient()
	return client.R().
		SetHeader("Content-Type", "application/json").
		SetBasicAuth(username, password).
		Get(url)
}

func newClient() *resty.Client {
	client := resty.New()
	client.SetTimeout(time.Second * 3)
	client.SetTLSClientConfig(&tls.Config{
		InsecureSkipVerify: true,
	})
	return client
}

func PostJson(url string, headers map[string]string, body any) (*resty.Response, error) {
	client := newClient()
	request := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body)
	if len(headers) > 0 {
		request = request.SetHeaders(headers)
	}
	return request.Post(url)
}

func PostFormData(url string, headers map[string]string, formData map[string]string, files map[string]string) (*resty.Response, error) {
	client := newClient()
	request := client.R().
		SetHeader("Content-Type", "multipart/form-data")
	if len(headers) > 0 {
		request = request.SetHeaders(headers)
	}
	if len(formData) > 0 {
		request = request.SetFormData(formData)
	}
	if len(files) > 0 {
		request = request.SetFiles(files)
	}
	return request.Post(url)
}
