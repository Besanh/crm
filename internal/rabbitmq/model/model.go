package rabbitmq

type Payload struct {
	HttpMethod string `json:"http_method"`
	Uri        string `json:"uri"`
	Body       any    `json:"body"`
}
