package httpmethods

type HttpRequestMethod string

const (
	Get    HttpRequestMethod = "GET"
	Post   HttpRequestMethod = "POST"
	Put    HttpRequestMethod = "PUT"
	Delete HttpRequestMethod = "DELETE"
)
