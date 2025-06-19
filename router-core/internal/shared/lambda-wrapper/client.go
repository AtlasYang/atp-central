package lambda_wrapper

import (
	"context"
	"encoding/json"

	"aigendrug.com/router-core/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

type LambdaWrapperClient struct {
	lambdaClient *lambda.Client
}

func NewLambdaWrapperClient(config *config.Config) LambdaWrapperClient {
	client := lambda.NewFromConfig(aws.Config{
		Region: config.AWS.Region,
		Credentials: aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     config.AWS.AccessKeyID,
				SecretAccessKey: config.AWS.SecretAccessKey,
			}, nil
		}),
	})

	return LambdaWrapperClient{
		lambdaClient: client,
	}
}

func (wrapper *LambdaWrapperClient) GetFunction(ctx context.Context, functionName string) (*types.State, error) {
	var state types.State
	funcOutput, err := wrapper.lambdaClient.GetFunction(ctx, &lambda.GetFunctionInput{
		FunctionName: aws.String(functionName),
	})
	if err != nil {
		return nil, err
	} else {
		state = funcOutput.Configuration.State
	}
	return &state, nil
}

func (wrapper *LambdaWrapperClient) Invoke(
	ctx context.Context, functionName string, parameters any, invocationType types.InvocationType, getLog bool,
) (*lambda.InvokeOutput, error) {
	logType := types.LogTypeNone
	if getLog {
		logType = types.LogTypeTail
	}
	payload, err := json.Marshal(parameters)
	if err != nil {
		return nil, err
	}
	invokeOutput, err := wrapper.lambdaClient.Invoke(ctx, &lambda.InvokeInput{
		FunctionName:   aws.String(functionName),
		LogType:        logType,
		Payload:        payload,
		InvocationType: invocationType,
	})
	if err != nil {
		return nil, err
	}
	return invokeOutput, nil
}
func (wrapper *LambdaWrapperClient) ListFunctions(ctx context.Context, maxItems int) ([]types.FunctionConfiguration, error) {
	var functions []types.FunctionConfiguration
	paginator := lambda.NewListFunctionsPaginator(wrapper.lambdaClient, &lambda.ListFunctionsInput{
		MaxItems: aws.Int32(int32(maxItems)),
	})
	for paginator.HasMorePages() && len(functions) < maxItems {
		pageOutput, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		functions = append(functions, pageOutput.Functions...)
	}
	return functions, nil
}
