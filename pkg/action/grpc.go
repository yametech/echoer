package action

var _ Interface = &gRPC{}

type gRPC struct {
	Uri string `json:"uri"`
}

func (g gRPC) GrpcInterface() GrpcInterface {
	panic("implement me")
}

func (g gRPC) Call(params interface{}) Interface {
	panic("implement me")
}

func (g gRPC) Params(m map[string]interface{}) GrpcInterface {
	panic("implement me")
}

func (g gRPC) HttpInterface() HttpInterface {
	panic("implement me")
}
