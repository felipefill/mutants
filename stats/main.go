package main

import (
	"encoding/json"

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
		return events.APIGatewayProxyResponse{Body: "Failed to retrieve stats", StatusCode: 500}, err
	}

	json, _ := json.Marshal(stats)

	return events.APIGatewayProxyResponse{Body: string(json), StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
