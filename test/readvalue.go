package test

import "github.com/libanvl/swager/test/union"

type ReadValue struct {
	union.Union[[]byte, error]
}

func (r ReadValue) Bytes() []byte {
	t := r.GetT()
	if t == nil {
		return nil
	}

	return *t
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
