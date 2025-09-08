package resp

type Resp struct {
	Command RespCommand
	Value   any
}

func NewResp(c RespCommand, v any) *Resp {
	return &Resp{
		Command: c,
		Value:   v,
	}
}
