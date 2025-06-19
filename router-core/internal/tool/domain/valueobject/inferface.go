package valueobject

type EngineInterfaceType string
type EngineInterfaceInvokeType string
type EngineInterfaceCheckStatusType string

// EngineInterfaceType defines source of each tool
const (
	// provided by AWS
	EngineInterfaceAWSLambda EngineInterfaceType = "aws-lambda"

	// provided by http server endpoint
	EngineInterfaceHTTPServer EngineInterfaceType = "http-server"
)

// EngineInterfaceInvokeType defines how to invoke the tool
const (
	// invoke the tool and wait for the response
	EngineInterfaceInvokeTypeSyncWait EngineInterfaceInvokeType = "sync-wait"

	// invoke the tool and return the response immediately
	EngineInterfaceInvokeTypeAsyncEvent EngineInterfaceInvokeType = "async-event"
)

// EngineInterfaceCheckStatusType defines how to check the status of the tool
const (
	// (sync-wait) no status check provided
	EngineInterfaceCheckStatusTypeNone EngineInterfaceCheckStatusType = "none"

	// (sync-wait) no status check provided, wait for response artifact after certain time
	EngineInterfaceCheckStatusTypeDelayed EngineInterfaceCheckStatusType = "delayed"

	// (async-event) poll the status of the tool using HTTP GET method
	EngineInterfaceCheckStatusTypePollHTTP EngineInterfaceCheckStatusType = "poll-http"

	// (async-event) poll the status of the tool using AWS S3 trigger
	EngineInterfaceCheckStatusTypeAWSS3Trigger EngineInterfaceCheckStatusType = "aws-s3-trigger"
)
