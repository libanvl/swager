package ipc

type ClientRaw interface {
	Close() error
	WorkspacesRaw() (string, error)
	TreeRaw() (string, error)
	VersionRaw() (string, error)
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

func (c client) WorkspacesRaw() (string, error) {
	return c.ipccallraw(GetWorkspacesMessage, nil)
}

func (c client) TreeRaw() (string, error) {
	return c.ipccallraw(GetTreeMessage, nil)
}

func (c client) VersionRaw() (string, error) {
	return c.ipccallraw(GetVersionMessage, nil)
}

