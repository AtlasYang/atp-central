package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	lambda_wrapper "aigendrug.com/router-core/internal/shared/lambda-wrapper"
	"aigendrug.com/router-core/internal/tool/application/dto"
	"aigendrug.com/router-core/internal/tool/domain"
	"aigendrug.com/router-core/internal/tool/domain/entity"
	"aigendrug.com/router-core/internal/tool/domain/valueobject"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/google/uuid"
)

const DefaultFunctionSyncExecutionTimeout = 10 * time.Second

type FunctionExecutor interface {
	Sync(ctx context.Context, tool *entity.Tool, requestID int, executionRequest dto.ToolExecutionRequestDTO)
}

type functionExecutor struct {
	baseCtx      context.Context
	toolRepo     domain.ToolRepository
	lambdaClient lambda_wrapper.LambdaWrapperClient
	// inject other engine providers here (Azure, GCP, etc.)
}

func NewFunctionExecutor(
	baseCtx context.Context,
	toolRepo domain.ToolRepository,
	lambdaClient lambda_wrapper.LambdaWrapperClient,
) FunctionExecutor {
	return &functionExecutor{
		baseCtx:      baseCtx,
		toolRepo:     toolRepo,
		lambdaClient: lambdaClient,
	}
}

func (e *functionExecutor) InvokeLambdaFunction(
	ctx context.Context, functionName string, payload map[string]any, sync bool,
) (map[string]any, error) {
	invocationType := types.InvocationTypeRequestResponse
	if !sync {
		invocationType = types.InvocationTypeEvent
	}
	output, err := e.lambdaClient.Invoke(ctx, functionName, payload, invocationType, false)
	if err != nil {
		return nil, err
	}

	if output.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("lambda function %s returned status code %d", functionName, output.StatusCode)
	}

	if output.FunctionError != nil {
		return nil, fmt.Errorf("lambda function %s returned error: %s", functionName, *output.FunctionError)
	}

	var outputRes map[string]any
	err = json.Unmarshal(output.Payload, &outputRes)
	if err != nil {
		return nil, err
	}

	return outputRes, nil
}

func (e *functionExecutor) Sync(
	ctx context.Context, tool *entity.Tool, requestID int, executionRequest dto.ToolExecutionRequestDTO,
) {
	requestIdentifier, err := uuid.NewRandom()
	if err != nil {
		return
	}

	responseIdentifier, err := uuid.NewRandom()
	if err != nil {
		return
	}

	// Create independent context with timeout based on check status type
	var timeoutDuration time.Duration
	switch tool.EngineInterface.EngineInterfaceCheckStatusType {
	case valueobject.EngineInterfaceCheckStatusTypeNone:
		timeoutDuration = DefaultFunctionSyncExecutionTimeout
	case valueobject.EngineInterfaceCheckStatusTypeDelayed:
		delaySeconds := tool.EngineInterface.EngineImpl["delay_seconds"].(float64)
		timeoutDuration = time.Duration(delaySeconds) * time.Second
	default:
		// not implemented
		return
	}

	// Create timeout context for lambda execution only
	lambdaTimeoutCtx, cancel := context.WithTimeout(e.baseCtx, timeoutDuration)
	defer cancel()

	fmt.Printf("executing tool with timeout duration: %v\n", timeoutDuration)

	// Channel to receive execution result
	resultChan := make(chan map[string]any, 1)
	errorChan := make(chan error, 1)

	// Execute in goroutine with timeout
	go func() {
		var res map[string]any
		var err error

		switch tool.EngineInterface.EngineInterfaceType {
		case valueobject.EngineInterfaceAWSLambda:
			res, err = e.InvokeLambdaFunction(
				lambdaTimeoutCtx,
				tool.EngineInterface.EngineImpl["function_name"].(string),
				executionRequest.Payload,
				true,
			)
			if err != nil {
				errorChan <- err
				return
			}
		default:
			// not implemented
			errorChan <- fmt.Errorf("engine interface type not implemented")
			return
		}

		resultChan <- res
	}()

	// Wait for result or timeout
	var result map[string]any
	var status valueobject.ToolRequestStatus

	select {
	case result = <-resultChan:
		fmt.Printf("execution completed successfully: %v\n", result)
		status = valueobject.ToolRequestStatusSuccess
	case err := <-errorChan:
		fmt.Printf("execution error: %v\n", err)
		status = valueobject.ToolRequestStatusFailed
	case <-lambdaTimeoutCtx.Done():
		fmt.Printf("execution timeout after %v\n", timeoutDuration)
		status = valueobject.ToolRequestStatusFailed
	}

	// Use separate context for DB operations (with reasonable timeout)
	dbCtx, dbCancel := context.WithTimeout(e.baseCtx, 30*time.Second)
	defer dbCancel()

	// Update tool request with results
	toolRequest, err := e.toolRepo.FindToolRequestByID(dbCtx, requestID)
	if err != nil {
		fmt.Printf("failed to find tool request: %v\n", err)
		return
	}

	toolRequest.RequestData.RequestIdentifier = requestIdentifier.String()
	toolRequest.RequestData.Payload = executionRequest.Payload

	toolRequest.ResponseData.ResponseIdentifier = responseIdentifier.String()
	if result != nil {
		toolRequest.ResponseData.Payload = result
	}

	toolRequest.Status = status

	err = e.toolRepo.UpdateToolRequest(dbCtx, toolRequest)
	if err != nil {
		fmt.Printf("failed to update tool request: %v\n", err)
		return
	}
}
