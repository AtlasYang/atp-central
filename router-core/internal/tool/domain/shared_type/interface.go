package shared_type

type ToolRequestData struct {
	RequestIdentifier string         `json:"request_identifier"`
	Payload           map[string]any `json:"payload"`
}

// ResponseData
//
// ResponseIdentifier: Unique identifier of the response. (etc: UUID)
//
// CheckStatusImpl: Implementation of the status checking strategy. Arranged with EngineInterfaceCheckStatusType.
// - Use "type" field to determine the implementation. (same with EngineInterfaceCheckStatusType of correspoding tool)
// - If EngineInterfaceCheckStatusType is none, this field is not required.
// - If EngineInterfaceCheckStatusType is delayed, use "delay_seconds" field to determine the delay time.
// - If EngineInterfaceCheckStatusType is poll-http, retrieve response payload from "url" field using GET method.
// - If EngineInterfaceCheckStatusType is aws-s3-trigger, retrieve response payload from "aws-s3-bucket" and "aws-s3-key" fields.
//
// Payload: Payload of the response.
type ToolRequestResponseData struct {
	ResponseIdentifier string         `json:"response_identifier"`
	CheckStatusImpl    map[string]any `json:"check_status_impl"`
	Payload            map[string]any `json:"payload"`
}
