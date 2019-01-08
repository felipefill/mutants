package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/felipefill/mutants/model"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.Body == "" {
		return events.APIGatewayProxyResponse{Body: "", StatusCode: 400}, nil
	}

	data, err := parseDNAFromString(request.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "", StatusCode: 400}, nil
	}

	//TODO: validate it's a NxN table

	if !validateDNAHasOnlyValidBases(data.DNA) {
		return events.APIGatewayProxyResponse{Body: "", StatusCode: 400}, nil
	}

	if !isMutant(data.DNA) {
		return events.APIGatewayProxyResponse{Body: "", StatusCode: 403}, nil
	}

	return events.APIGatewayProxyResponse{Body: "", StatusCode: 200}, nil
}

func isMutant(dna []string) bool {
	return false
}

func parseDNAFromString(str string) (model.DNA, error) {
	data := model.DNA{}
	err := json.Unmarshal([]byte(str), &data)

	return data, err
}

func validateDNAHasOnlyValidBases(dna []string) bool {
	for row := 0; row < len(dna); row++ {
		for column := 0; column < len(dna[row]); column++ {
			currentChar := dna[row][column]
			//TODO: Maybe we can use some kind of switch here
			if currentChar != 'A' && currentChar != 'T' && currentChar != 'C' && currentChar != 'G' {
				return false
			}
		}
	}

	return true
}

func main() {
	lambda.Start(Handler)
}
