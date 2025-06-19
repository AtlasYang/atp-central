package valueobject

type ToolRequestStatus string
type ToolExecutionStatus string

const (
	ToolRequestStatusPending ToolRequestStatus = "pending"
	ToolRequestStatusSuccess ToolRequestStatus = "success"
	ToolRequestStatusFailed  ToolRequestStatus = "failed"
)

const (
	ToolExecutionStatusSuccess      ToolExecutionStatus = "success"
	ToolExecutionStatusUnauthorized ToolExecutionStatus = "unauthorized"
	ToolExecutionStatusFailed       ToolExecutionStatus = "failed"
)

func (t ToolRequestStatus) String() string {
	return string(t)
}
