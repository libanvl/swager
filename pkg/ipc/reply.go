package ipc

type Version struct {
	Major                int    `json:"major"`
	Minor                int    `json:"minor"`
	Patch                int    `json:"patch"`
	HumanReadable        string `json:"human_readable"`
	LoadedConfigFileName string `json:"loaded_config_file_name"`
}

type Result struct {
	Success bool `json:"success"`
}

type Command struct {
	Result
	ParseError bool   `json:"parse_error"`
	Error      string `json:"error,omitempty"`
}

type Workspace struct {
	Num     int    `json:"num"`
	Name    string `json:"name"`
	Visible bool   `json:"visible"`
	Focused bool   `json:"focused"`
	Urgent  bool   `json:"urgent"`
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
	Dpms             bool    `json:"dpms"`
	Primary          bool    `json:"primary"`
	Scale            float64 `json:"scale"`
	SubpixelHinting  string  `json:"subpixel_hinting"`
	Transform        string  `json:"transform"`
	CurrentWorkspace string  `json:"current_workspace"`
	Modes            []Mode  `json:"modes"`
	CurrentMode      Mode    `json:"current_mode"`
	Rect             Rect    `json:"rect"`
}

type Rect struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type BindingState struct {
	Name string `json:"name"`
}
