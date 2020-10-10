package action

type HttpInterface interface {
	Post(urls []string) HttpInterface
	Params(map[string]interface{}) HttpInterface
	Do() error
}

type GrpcInterface interface {
	Call(params interface{}) Interface
	Params(map[string]interface{}) GrpcInterface
}

type Interface interface {
	HttpInterface() HttpInterface
	GrpcInterface() GrpcInterface
}

type HookClient struct {
	*http
	*gRPC
}

func NewHookClient() *HookClient {
	return &HookClient{
		http: newHttp(),
		gRPC: nil,
	}
}

func (hc *HookClient) GrpcInterface() GrpcInterface {
	return hc.gRPC
}

func (hc *HookClient) HttpInterface() HttpInterface {
	return hc.http
}
