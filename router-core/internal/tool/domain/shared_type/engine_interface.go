package shared_type

import "aigendrug.com/router-core/internal/tool/domain/valueobject"

type EngineInterface struct {
	EngineInterfaceType            valueobject.EngineInterfaceType            `json:"engine_interface_type" validate:"required"`
	EngineInterfaceInvokeType      valueobject.EngineInterfaceInvokeType      `json:"engine_interface_invoke_type" validate:"required"`
	EngineInterfaceCheckStatusType valueobject.EngineInterfaceCheckStatusType `json:"engine_interface_check_status_type" validate:"required"`
	EngineImpl                     map[string]any                             `json:"engine_impl" validate:"required"`
}

// EngineImpl
//
// EngineImple stores engine-specific implementation as dynamic fields.
// Example fields:
// - "function_name": Name of the function to invoke. Provided when EngineInterfaceType is aws-lambda.
// - "url": URL of the HTTP server. Provided when EngineInterfaceType is http-server.
// - "aws_lambda_function_name": AWS Lambda function name. Provided when EngineInterfaceType is aws-lambda.
// - "aws_lambda_invoke_type": AWS Lambda invoke type. Provided when EngineInterfaceType is aws-lambda.
// - "aws_lambda_check_status_type": AWS Lambda check status type. Provided when EngineInterfaceType is aws-lambda.
// - "delay_seconds": Delay seconds. Provided when EngineInterfaceCheckStatusType is delayed.
// - "status_url": Status URL. Provided when EngineInterfaceCheckStatusType is poll-http.
// - "aws_s3_bucket": AWS S3 bucket name. Provided when EngineInterfaceCheckStatusType is aws-s3-trigger.
// - "aws_s3_key": AWS S3 key. Provided when EngineInterfaceCheckStatusType is aws-s3-trigger.
