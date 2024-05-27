package service

import (
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
)

var CustomHttpClient *http.Client

func NewCustomHttpClient() *http.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	client := retryClient.StandardClient()
	CustomHttpClient = client
	return CustomHttpClient
}
