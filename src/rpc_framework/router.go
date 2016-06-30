package rpc_framework

type RouterInterface interface {
	Route(logid uint32, json []byte) ([]byte, error)
}

