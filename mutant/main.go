package main

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/felipefill/mutants/model"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

const repetitionRequiredForSequence int = 4

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.Body == "" {
		return events.APIGatewayProxyResponse{Body: "", StatusCode: 400}, nil
	}

	dna, err := validateRequest(request.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 400}, nil
	}

	if !isMutant(dna.DNA) {
		return events.APIGatewayProxyResponse{Body: "", StatusCode: 403}, nil
	}

	return events.APIGatewayProxyResponse{Body: "", StatusCode: 200}, nil
}

func isMutant(dna []string) bool {
	count := 0

	for row := 0; row < len(dna); row++ {
		for column := 0; column < len(dna[row]); column++ {
			if checkSequenceToTheRight(dna[row], row) {
				count++
			}

			if checkSequenceDown(dna, row, column) {
				count++
			}

			if checkSequenceDiagonalLeft(dna, row, column) {
				count++
			}

			if checkSequenceDiagonalRight(dna, row, column) {
				count++
			}

			//TODO: break when count > 1
		}
	}

	return count > 1
}

func checkSequenceToTheRight(dnaRow string, column int) bool {
	if len(dnaRow)-column < repetitionRequiredForSequence {
		return false
	}

	requiredBase := dnaRow[column]
	c := column + 1

	for loop := 0; loop < repetitionRequiredForSequence-1; loop++ {
		if dnaRow[c] != requiredBase {
			return false
		}

		c++
	}

	return true
}

func checkSequenceDown(dna []string, row, column int) bool {
	if len(dna)-row < repetitionRequiredForSequence {
		return false
	}

	requiredBase := dna[row][column]
	r := row + 1

	for loop := 0; loop < repetitionRequiredForSequence-1; loop++ {
		if dna[r][column] != requiredBase {
			return false
		}

		r++
	}

	return true
}

func checkSequenceDiagonalLeft(dna []string, row, column int) bool {
	if row >= repetitionRequiredForSequence || column < repetitionRequiredForSequence-1 {
		return false
	}

	requiredBase := dna[row][column]
	r := row + 1
	c := column - 1

	for loop := 0; loop < repetitionRequiredForSequence-1; loop++ {
		if dna[r][c] != requiredBase {
			return false
		}

		r++
		c--
	}

	return true
}

func checkSequenceDiagonalRight(dna []string, row, column int) bool {
	if row >= repetitionRequiredForSequence || column >= repetitionRequiredForSequence {
		return false
	}

	requiredBase := dna[row][column]
	r := row + 1
	c := column + 1

	for loop := 0; loop < repetitionRequiredForSequence-1; loop++ {
		if dna[r][c] != requiredBase {
			return false
		}

		r++
		c++
	}

	return true
}

func validateRequest(body string) (model.DNA, error) {
	data, err := parseDNAFromString(body)
	if err != nil {
		return model.DNA{}, errors.New("Could not parse request :" + err.Error())
	}

	if !isValidNxNTable(data.DNA) {
		return data, errors.New("DNA is not an NxN table")
	}

	if !validateDNAHasOnlyValidBases(data.DNA) {
		return data, errors.New("DNA has invalid bases")
	}

	return data, nil
}

func isValidNxNTable(dna []string) bool {
	// This will give the number of rows
	hypothesis := len(dna)

	for row := 0; row < hypothesis; row++ {
		// If we have a number of columns that's different from our hypothesis
		if len(dna[row]) != hypothesis {
			return false
		}
	}

	return true
}

func parseDNAFromString(str string) (model.DNA, error) {
	data := model.DNA{}
	err := json.Unmarshal([]byte(str), &data)

	return data, err
}

func validateDNAHasOnlyValidBases(dna []string) bool {
	//TODO: This could be done while checking DNA
	for row := 0; row < len(dna); row++ {
		for column := 0; column < len(dna[row]); column++ {
			currentChar := dna[row][column]
			if !isValidDNABase(currentChar) {
				return false
			}
		}
	}

	return true
}

func isValidDNABase(c byte) bool {
	switch c {
	case 'A', 'T', 'C', 'G':
		return true
	default:
		return false
	}
}

func main() {
	lambda.Start(Handler)
}
