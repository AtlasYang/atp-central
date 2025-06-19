package valueobject

type ToolClientPermissionLevel int

const (
	ToolClientPermissionLevelNone  ToolClientPermissionLevel = 0
	ToolClientPermissionLevelRead  ToolClientPermissionLevel = 1
	ToolClientPermissionLevelWrite ToolClientPermissionLevel = 2
)

func (t ToolClientPermissionLevel) Int() int {
	return int(t)
}
