package command

type ErrorReply struct {
	Message interface{}
}

func (e *ErrorReply) Value() interface{} {
	return e.Message
}

type OkReply struct {
	Message []byte
}

func (o *OkReply) Value() interface{} {
	return o.Message
}

type RawReply struct {
	Message []byte
}

func (r *RawReply) Value() interface{} {
	return r.Message
}
