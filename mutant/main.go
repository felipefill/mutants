package main

import (
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
	if request.Body == "" {
		return events.APIGatewayProxyResponse{Body: "Empty body", StatusCode: 400}, nil
	}

	dnaCheck, err := NewDNACheckFromJSONString(request.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, nil
	}

	if !isMutant(dnaCheck.DNA) {
		return events.APIGatewayProxyResponse{Body: "", StatusCode: 403}, nil
	}

	return events.APIGatewayProxyResponse{Body: "", StatusCode: 200}, nil
}

func isMutant(data []string) bool {
	// This is being done this way so that the code complies with the requirements
	// Which is having a function with this signature
	// My apprach would be to have it related to DNACheck struct
	dnaCheck := DNACheck{DNA: data}

	return dnaCheck.IsMutant()
}

func main() {
	lambda.Start(Handler)
}
