package httpapi

type ErrMsg struct {
	errcode string
	errmsg  string
}

func (e ErrMsg) Error() string {
	return e.errmsg
}
