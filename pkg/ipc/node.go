package ipc

type Node struct {
	ID                 int                 `json:"id"`
	Name               string              `json:"name"`
	Type               NodeType            `json:"type"`
	Border             BorderType          `json:"border"`
	CurrentBorderWidth int                 `json:"current_border_width"`
	Layout             LayoutType          `json:"layout"`
	Orientation        OrientationType     `json:"orientation"`
	Percent            *float64            `json:"percent"`
	Rect               Rect                `json:"rect"`
	WindowRect         Rect                `json:"window_rect"`
	DecoRect           Rect                `json:"deco_rect"`
	Geometry           Rect                `json:"geometry"`
	Urgent             *bool               `json:"urgent"`
	Sticky             bool                `json:"sticky"`
	Marks              []string            `json:"marks"`
	Focused            bool                `json:"focused"`
	Focus              []int               `json:"focus"`
	Nodes              []*Node             `json:"nodes"`
	FloatingNodes      []*Node             `json:"floating_nodes"`
	Representation     *string             `json:"representation"`
	FullscreenMode     *FullscreenModeType `json:"fullscreen_mode"`
	AppID              *string             `json:"app_id"`
	Pid                *int                `json:"pid"`
	Visible            *bool               `json:"visible"`
	Shell              *string             `json:"shell"`
	Window             *int                `json:"window"`
	WindowProperties   *WindowProperties   `json:"window_properties"`
}

type NodeType string

const (
	RootNode        NodeType = "root"
	OutputNode      NodeType = "output"
	WorkspaceNode   NodeType = "workspace"
	ConNode         NodeType = "con"
	FloatingConNode NodeType = "floating_con"
)

type BorderType string

const (
	NormalBorder BorderType = "normal"
	NoneBorder   BorderType = "none"
	PixelBorder  BorderType = "pixel"
	CSDBorder    BorderType = "csd"
)

type LayoutType string

const (
	SplitHLayout  LayoutType = "splith"
	SplitVLayout  LayoutType = "splitv"
	StackedLayout LayoutType = "stacked"
	TabbedLayout  LayoutType = "tabbed"
	OutputLayout  LayoutType = "output"
)

type OrientationType string

const (
	VerticalOrientation   OrientationType = "vertical"
	HorizontalOrientation OrientationType = "horizontal"
	NoneOrientation       OrientationType = "none"
)

//go:generate go run golang.org/x/tools/cmd/stringer -type=FullscreenModeType
type FullscreenModeType uint8

const (
	NoneFullscreenMode      FullscreenModeType = 0
	WorkspaceFullscreenMode FullscreenModeType = 1
	GlobalFullscreenMode    FullscreenModeType = 2
)

type WindowProperties struct {
	Class        string      `json:"class"`
	Instance     string      `json:"instance"`
	TransientFor interface{} `json:"transient_for"`
}
