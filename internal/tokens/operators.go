package tokens

func SetOp[T UnFunc | BiFunc](opfn T, tok *Token) {
	tok.literal = opfn
}

type UnFunc func(x interface{}) interface{}
type BiFunc func(x, y interface{}) interface{}

func (t *Token) UnFunc(x interface{}) interface{} {
	return t.literal.(UnFunc)(x)
}

func (t *Token) BiFunc(x, y interface{}) interface{} {
	return t.literal.(BiFunc)(x, y)
}
