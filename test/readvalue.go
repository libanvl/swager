package test

import "github.com/libanvl/swager/test/union"

type ReadValue struct {
	union.Union[[]byte, error]
}

func (r ReadValue) Bytes() []byte {
	return *r.GetT()
}

func (r ReadValue) Error() error {
	return *r.GetU()
}

func (r ReadValue) WithBytes(p []byte) ReadValue {
	r.Union = r.Union.WithT(p)
	return r
}

func (r ReadValue) WithError(err error) ReadValue {
	r.Union = r.Union.WithU(err)
	return r
}
