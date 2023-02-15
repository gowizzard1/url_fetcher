package client

type IClient interface {
	Request(urls []string) []URLResponse
}
