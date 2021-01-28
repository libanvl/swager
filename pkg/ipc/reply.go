package ipc

type Version struct {
	HumanReadable        string `json:"human_readable"`
	Variant              string
	Major                int
	Minor                int
	Patch                int
	LoadedConfigFileName string `json:"loaded_config_file_name"`
}

type Result struct {
	Success bool `json:"success"`
}

type Command struct {
	Result
	Error string `json:"error,omitempty"`
}

type Workspace struct {
	Num     int    `json:"num"`
	Name    string `json:"name"`
	Visible bool   `json:"visible"`
	Focused bool   `json:"focused"`
	Rect    Rect   `json:"rect"`
	Output  string `json:"output"`
}

type Mode struct {
	Width   int `json:"width"`
	Height  int `json:"height"`
	Refresh int `json:"refresh"`
}

type Output struct {
	Name             string  `json:"name"`
	Make             string  `json:"make"`
	Model            string  `json:"model"`
	Serial           string  `json:"serial"`
	Active           bool    `json:"active"`
	Primary          bool    `json:"primary"`
	Scale            float64 `json:"scale"`
	SubpixelHinting  string  `json:"subpixel_hinting"`
	Transform        string  `json:"transform"`
	CurrentWorkspace string  `json:"current_workspace"`
	Modes            []Mode  `json:"modes"`
	CurrentMode      Mode    `json:"current_mode"`
}

type Rect struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type WindowProperties struct {
	Class        string      `json:"class"`
	Instance     string      `json:"instance"`
	TransientFor interface{} `json:"transient_for"`
}

type NodeType string
const (
	RootNode        NodeType = "root"
	OutputNode      NodeType = "output"
	WorkspaceNode   NodeType = "workspace"
	ConNode         NodeType = "con"
	FloatingConNode NodeType = "floating_con"
)

type Node struct {
	ID                 int              `json:"id"`
	Name               string           `json:"name"`
	Rect               Rect             `json:"rect"`
	Focused            bool             `json:"focused"`
	Focus              []int            `json:"focus"`
	Border             string           `json:"border"`
	CurrentBorderWidth int              `json:"current_border_width"`
	Layout             string           `json:"layout"`
	Percent            float64          `json:"percent"`
	WindowRect         Rect             `json:"window_rect"`
	DecoRect           Rect             `json:"deco_rect"`
	Geometry           Rect             `json:"geometry"`
	Window             int              `json:"window"`
	Urgent             bool             `json:"urgent"`
	FloatingNodes      []Node           `json:"floating_nodes"`
	Type               NodeType         `json:"type"`
	Pid                int              `json:"pid"`
	AppID              string           `json:"app_id"`
	WindowProperties   WindowProperties `json:"window_properties"`
	Nodes              []Node           `json:"nodes"`
}
