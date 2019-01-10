package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	stats, err := GetStats()
	if err != nil {
		return events.APIGatewayProxyResponse{Body: fmt.Sprintf("Failed to retrieve stats %s", err.Error()), StatusCode: 500}, err
	}

	json, err := json.Marshal(stats)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: fmt.Sprintf("Failed to retrieve stats %s", err.Error()), StatusCode: 500}, err
	}

	return events.APIGatewayProxyResponse{Body: string(json), StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
