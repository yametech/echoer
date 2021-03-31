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

type HttpsInterface interface {
	Post(urls []string) HttpsInterface
	Params(map[string]interface{}) HttpsInterface
	Do() error
}

type Interface interface {
	HttpInterface() HttpInterface
	GrpcInterface() GrpcInterface
}

type HookClient struct {
	*http
	*gRPC
	*https
}

func NewHookClient() *HookClient {
	return &HookClient{
		http: newHttp(),
		gRPC: nil,
		https:newHttps(),
	}
}

func (hc *HookClient) GrpcInterface() GrpcInterface {
	return hc.gRPC
}

func (hc *HookClient) HttpInterface() HttpInterface {
	return hc.http
}


func (hc *HookClient) HttpsInterface() HttpsInterface {
	return hc.https
}