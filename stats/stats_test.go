package main

import (
	"errors"
	"testing"

	"github.com/felipefill/mutants/utils"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestGetStatsSuccess(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	utils.InjectDatabase(db)

	mock.
		ExpectQuery("select count\\(id\\) count, type from dna group by type").
		WillReturnRows(
			sqlmock.NewRows([]string{"count", "type"}).
				AddRow(10, "mutant").
				AddRow(40, "ordinary"),
		)

	var expectedError error
	expectedStats := &Stats{
		HumanDNACount:  50,
		MutantDNACount: 10,
		Ratio:          0.2,
	}

	actualStats, actualError := GetStats()

	assert.EqualValues(t, expectedStats, actualStats, "Stats are not equal")
	assert.EqualValues(t, expectedError, actualError, "Error was not as expected")
}

func TestGetStatsWithEmptyDatabase(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	utils.InjectDatabase(db)

	mock.
		ExpectQuery("select count\\(id\\) count, type from dna group by type").
		WillReturnRows(
			sqlmock.NewRows([]string{"count", "type"}),
		)

	var expectedError error
	expectedStats := &Stats{
		HumanDNACount:  0,
		MutantDNACount: 0,
		Ratio:          0,
	}

	actualStats, actualError := GetStats()

	assert.EqualValues(t, expectedStats, actualStats, "Stats are not equal")
	assert.EqualValues(t, expectedError, actualError, "Error was not as expected")
}

func TestGetStatsFailsToQueryDatabase(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	utils.InjectDatabase(db)
	mock.
		ExpectQuery("select count\\(id\\) count, type from dna group by type").
		WillReturnError(sqlmock.ErrCancelled)

	var expectedStats *Stats
	expectedError := errors.New("Failed to query database")

	actualStats, actualError := GetStats()

	assert.EqualValues(t, expectedStats, actualStats, "Stats are not equal")
	assert.EqualValues(t, expectedError, actualError, "Error was not as expected")
}

func TestGetStatsFailsToParseResult(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	utils.InjectDatabase(db)

	mock.
		ExpectQuery("select count\\(id\\) count, type from dna group by type").
		WillReturnRows(
			sqlmock.NewRows([]string{"single_column"}).
				AddRow("just one column"),
		)

	var expectedStats *Stats
	expectedError := errors.New("Failed to retrieve status")

	actualStats, actualError := GetStats()

	assert.EqualValues(t, expectedStats, actualStats, "Stats are not equal")
	assert.EqualValues(t, expectedError, actualError, "Error was not as expected")
}

func TestStatsHandlerSuccess(t *testing.T) {
	request := events.APIGatewayProxyRequest{}
	db, mock, _ := sqlmock.New()
	defer db.Close()

	utils.InjectDatabase(db)

	mock.
		ExpectQuery("select count\\(id\\) count, type from dna group by type").
		WillReturnRows(
			sqlmock.NewRows([]string{"count", "type"}).
				AddRow(10, "mutant").
				AddRow(40, "ordinary"),
		)

	var expectedError error
	expectedResponse := events.APIGatewayProxyResponse{
		Body:       "{\"count_mutant_dna\":10,\"count_human_dna\":50,\"ratio\":0.2}",
		StatusCode: 200,
	}

	actualResponse, actualError := Handler(request)

	assert.EqualValues(t, expectedResponse, actualResponse, "Responses are not equal")
	assert.EqualValues(t, expectedError, actualError, "Error was not as expected")
}

func TestStatsHandlerFailsToRetrieveStats(t *testing.T) {
	request := events.APIGatewayProxyRequest{}
	db, mock, _ := sqlmock.New()
	defer db.Close()

	utils.InjectDatabase(db)
	mock.
		ExpectQuery("select count\\(id\\) count, type from dna group by type").
		WillReturnError(sqlmock.ErrCancelled)

	expectedResponse := events.APIGatewayProxyResponse{
		Body:       "Failed to retrieve stats",
		StatusCode: 500,
	}
	expectedError := errors.New("Failed to query database")

	actualResponse, actualError := Handler(request)

	assert.EqualValues(t, expectedResponse, actualResponse, "Responses are not equal")
	assert.EqualValues(t, expectedError, actualError, "Error was not as expected")
}
