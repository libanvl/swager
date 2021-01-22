package ipc

type ClientRaw interface {
	Close() error
	VersionRaw() (string, error)
	WorkspacesRaw() (string, error)
}

func init() {
	var _ ClientRaw = client{}
}

func (c client) ipccallraw(pt PayloadType, payload []byte) (string, error) {
	res, err := c.ipccall(pt, payload)
	if err != nil {
		return "", nil
	}
	return string(res), nil
}

func (c client) VersionRaw() (string, error) {
	return c.ipccallraw(GetVersionMessage, nil)
}

func (c client) WorkspacesRaw() (string, error) {
	return c.ipccallraw(GetWorkspacesMessage, nil)
}
