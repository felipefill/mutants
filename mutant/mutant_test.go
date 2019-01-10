package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/aws/aws-lambda-go/events"
	"github.com/felipefill/mutants/utils"
	"github.com/stretchr/testify/assert"
)

func TestIsValidDNABase(t *testing.T) {
	assert.Equal(t, true, isValidDNABase('A'), "DNA base should be valid")
	assert.Equal(t, true, isValidDNABase('T'), "DNA base should be valid")
	assert.Equal(t, true, isValidDNABase('C'), "DNA base should be valid")
	assert.Equal(t, true, isValidDNABase('G'), "DNA base should be valid")

	assert.Equal(t, false, isValidDNABase('E'), "DNA base should not be valid")
	assert.Equal(t, false, isValidDNABase('X'), "DNA base should not be valid")
	assert.Equal(t, false, isValidDNABase('1'), "DNA base should not be valid")
	assert.Equal(t, false, isValidDNABase('9'), "DNA base should not be valid")
	assert.Equal(t, false, isValidDNABase('#'), "DNA base should not be valid")
}

func TestValidateDNAHasOnlyValidBases(t *testing.T) {
	validDNA := DNACheck{
		DNA: validDNASequence,
	}

	invalidDNA := DNACheck{
		DNA: invalidDNASequence,
	}

	assert.Equal(t, true, validDNA.validateDNAHasOnlyValidBases(), "DNA bases should be valid")
	assert.Equal(t, false, invalidDNA.validateDNAHasOnlyValidBases(), "DNA bases should not be valid")
}

func TestIsValidNxNTable(t *testing.T) {
	dnaWithNxN := DNACheck{
		DNA: tableNxN,
	}

	dnaWithMxN := DNACheck{
		DNA: tableMxN,
	}

	assert.Equal(t, true, dnaWithNxN.isValidNxNTable(), "Should be a valid NxN table")
	assert.Equal(t, false, dnaWithMxN.isValidNxNTable(), "Should not be a valid NxN table")
}

func TestValidate(t *testing.T) {
	validDNASequence := DNACheck{
		DNA: validDNASequence,
	}

	assert.Equal(t, nil, validDNASequence.validate(), "DNA sequence should be valid")
}

func TestValidateFailsWithInvalidTableSize(t *testing.T) {
	invalidDNAWithMxN := DNACheck{
		DNA: tableMxN,
	}

	expectedError := errors.New("DNA is not an NxN table")
	actualError := invalidDNAWithMxN.validate()

	assert.Equal(t, expectedError, actualError, "DNA sequence should not be valid")
}

func TestValidateFailsWithInvalidBases(t *testing.T) {
	invalidDNAWithWrongBases := DNACheck{
		DNA: invalidDNASequence,
	}

	expectedError := errors.New("DNA has invalid bases")
	actualError := invalidDNAWithWrongBases.validate()

	assert.Equal(t, expectedError, actualError, "DNA sequence should not be valid")
}

func TestHashFunction(t *testing.T) {
	check := DNACheck{
		DNA: validDNASequence,
	}

	expected := "f7bad0e12c11a6a23852bee23d64cc753bb51d83"
	actual := check.Hash()

	assert.Equal(t, expected, actual, "Hashes do not match")
}

func TestLookDNATypeInDatabaseFoundDNAType(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	utils.InjectDatabase(db)

	check := DNACheck{
		DNA: validDNASequence,
	}

	mock.
		ExpectQuery("select type from dna").
		WithArgs(check.Hash()).
		WillReturnRows(
			sqlmock.NewRows([]string{"type"}).
				AddRow("mutant"),
		)

	expected := "mutant"
	actual := check.lookDNATypeInDatabase()

	assert.Equal(t, expected, actual, "Should have found DNA type in DB")
}

func TestLookDNATypeInDatabaseNotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	utils.InjectDatabase(db)

	check := DNACheck{
		DNA: validDNASequence,
	}

	mock.
		ExpectQuery("select type from dna").
		WithArgs(check.Hash()).
		WillReturnError(sql.ErrNoRows)

	expected := "not found"
	actual := check.lookDNATypeInDatabase()

	assert.Equal(t, expected, actual, "Should not have found DNA type in DB")
}

func TestLookDNATypeInDatabasePanics(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	utils.InjectDatabase(db)

	check := DNACheck{
		DNA: validDNASequence,
	}

	mock.
		ExpectQuery("select type from dna").
		WithArgs(check.Hash()).
		WillReturnError(sql.ErrConnDone)

	assert.Panics(t, func() { check.lookDNATypeInDatabase() }, "Should have panicked")
}

func TestCheckSequenceDiagonalRight(t *testing.T) {
	check := DNACheck{
		DNA: mutantWithAllCombinationsDNASequence,
	}

	assert.Equal(t, false, check.CheckSequenceDiagonalRight(4, 0))
	assert.Equal(t, false, check.CheckSequenceDiagonalRight(0, 4))
	assert.Equal(t, false, check.CheckSequenceDiagonalRight(0, 0))

	assert.Equal(t, true, check.CheckSequenceDiagonalRight(3, 0))
}

func TestCheckSequenceDiagonalLeft(t *testing.T) {
	check := DNACheck{
		DNA: mutantWithAllCombinationsDNASequence,
	}

	assert.Equal(t, false, check.CheckSequenceDiagonalLeft(4, 0))
	assert.Equal(t, false, check.CheckSequenceDiagonalLeft(0, 0))
	assert.Equal(t, false, check.CheckSequenceDiagonalLeft(0, 5))

	assert.Equal(t, true, check.CheckSequenceDiagonalLeft(0, 4))
}

func TestCheckSequenceDown(t *testing.T) {
	check := DNACheck{
		DNA: mutantWithAllCombinationsDNASequence,
	}

	assert.Equal(t, false, check.CheckSequenceDown(5, 0))
	assert.Equal(t, false, check.CheckSequenceDown(0, 5))

	assert.Equal(t, true, check.CheckSequenceDown(2, 6))
}

func TestCheckSequenceRight(t *testing.T) {
	check := DNACheck{
		DNA: mutantWithAllCombinationsDNASequence,
	}

	assert.Equal(t, false, check.CheckSequenceToTheRight(0, 5))
	assert.Equal(t, false, check.CheckSequenceToTheRight(0, 0))

	assert.Equal(t, true, check.CheckSequenceToTheRight(6, 3))
}

func TestIsMutantFindingMutantInDatabase(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	utils.InjectDatabase(db)

	check := DNACheck{
		DNA: validDNASequence,
	}

	mock.
		ExpectQuery("select type from dna").
		WithArgs(check.Hash()).
		WillReturnRows(
			sqlmock.NewRows([]string{"type"}).
				AddRow("mutant"),
		)

	assert.Equal(t, true, check.IsMutant())
}

func TestIsMutantFindingHumanInDatabase(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	utils.InjectDatabase(db)

	check := DNACheck{
		DNA: validDNASequence,
	}

	mock.
		ExpectQuery("select type from dna").
		WithArgs(check.Hash()).
		WillReturnRows(
			sqlmock.NewRows([]string{"type"}).
				AddRow("ordinary"),
		)

	assert.Equal(t, false, check.IsMutant())
}

func TestIsMutantFalse(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	utils.InjectDatabase(db)

	check := DNACheck{
		DNA: humanDNASequence,
	}

	sequenceAsJSON, _ := json.Marshal(&check.DNA)

	mock.
		ExpectQuery("select type from dna").
		WithArgs(check.Hash()).
		WillReturnError(sql.ErrNoRows)

	mock.
		ExpectExec("insert into dna").
		WithArgs(check.Hash(), "ordinary", sequenceAsJSON).
		WillReturnResult(sqlmock.NewResult(1, 1))

	assert.Equal(t, false, check.IsMutant())
}

func TestIsMutantTrue(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	utils.InjectDatabase(db)

	check := DNACheck{
		DNA: mutantDNASequence,
	}

	sequenceAsJSON, _ := json.Marshal(&check.DNA)

	mock.
		ExpectQuery("select type from dna").
		WithArgs(check.Hash()).
		WillReturnError(sql.ErrNoRows)

	mock.
		ExpectExec("insert into dna").
		WithArgs(check.Hash(), "mutant", sequenceAsJSON).
		WillReturnResult(sqlmock.NewResult(1, 1))

	assert.Equal(t, true, check.IsMutant())
}

func TestIsMutantPanics(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	utils.InjectDatabase(db)

	check := DNACheck{
		DNA: mutantDNASequence,
	}

	sequenceAsJSON, _ := json.Marshal(&check.DNA)

	mock.
		ExpectQuery("select type from dna").
		WithArgs(check.Hash()).
		WillReturnError(sql.ErrNoRows)

	mock.
		ExpectExec("insert into dna").
		WithArgs(check.Hash(), "mutant", sequenceAsJSON).
		WillReturnError(sql.ErrConnDone)

	assert.Panics(t, func() { check.IsMutant() })
}

func TestSaveDNA(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	utils.InjectDatabase(db)

	check := DNACheck{
		DNA: mutantDNASequence,
	}

	sequenceAsJSON, _ := json.Marshal(&check.DNA)

	mock.
		ExpectExec("insert into dna").
		WithArgs(check.Hash(), "mutant", sequenceAsJSON).
		WillReturnResult(sqlmock.NewResult(1, 1))

	check.Save("mutant")
}

func TestSaveDNAPanics(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	utils.InjectDatabase(db)

	check := DNACheck{
		DNA: mutantDNASequence,
	}

	sequenceAsJSON, _ := json.Marshal(&check.DNA)

	mock.
		ExpectExec("insert into dna").
		WithArgs(check.Hash(), "mutant", sequenceAsJSON).
		WillReturnError(sql.ErrConnDone)

	assert.Panics(t, func() { check.Save("mutant") })
}

func TestNewDNACheckFromJSONString(t *testing.T) {
	var expectedError error
	expectedCheck := DNACheck{
		DNA: validDNASequence,
	}

	actualCheck, actualError := NewDNACheckFromJSONString(validDNASequenceString)

	assert.Equal(t, expectedCheck, actualCheck)
	assert.Equal(t, expectedError, actualError)
}

func TestNewDNACheckFromJSONStringFailsInvalidFormat(t *testing.T) {
	expectedError := errors.New("Could not parse DNA check")
	expectedCheck := DNACheck{}

	actualCheck, actualError := NewDNACheckFromJSONString(invalidDNASequenceStringNotEvenAJSON)

	assert.Equal(t, expectedCheck, actualCheck)
	assert.Equal(t, expectedError, actualError)
}

func TestNewDNACheckFromJSONStringFailsInvalid(t *testing.T) {
	expectedError := errors.New("DNA has invalid bases")
	expectedCheck := DNACheck{}

	actualCheck, actualError := NewDNACheckFromJSONString(invalidDNASequenceStringWrongBases)

	assert.Equal(t, expectedCheck, actualCheck)
	assert.Equal(t, expectedError, actualError)
}

func TestHandlerMutantDNA(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	utils.InjectDatabase(db)

	check := DNACheck{
		DNA: mutantDNASequence,
	}

	mock.
		ExpectQuery("select type from dna").
		WithArgs(check.Hash()).
		WillReturnRows(
			sqlmock.NewRows([]string{"type"}).
				AddRow("mutant"),
		)

	request := events.APIGatewayProxyRequest{
		Body: mutantDNASequenceAsJSONString,
	}

	var expectedError error
	expectedResponse := events.APIGatewayProxyResponse{
		Body:       "",
		StatusCode: 200,
	}

	actualResponde, actualError := Handler(request)

	assert.Equal(t, expectedError, actualError)
	assert.Equal(t, expectedResponse, actualResponde)
}

func TestHandlerHumanDNA(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	utils.InjectDatabase(db)

	check := DNACheck{
		DNA: humanDNASequence,
	}

	mock.
		ExpectQuery("select type from dna").
		WithArgs(check.Hash()).
		WillReturnRows(
			sqlmock.NewRows([]string{"type"}).
				AddRow("ordinary"),
		)

	request := events.APIGatewayProxyRequest{
		Body: humanDNASequenceAsJSONString,
	}

	var expectedError error
	expectedResponse := events.APIGatewayProxyResponse{
		Body:       "",
		StatusCode: 403,
	}

	actualResponde, actualError := Handler(request)

	assert.Equal(t, expectedError, actualError)
	assert.Equal(t, expectedResponse, actualResponde)
}

func TestHandlerEmptyBody(t *testing.T) {
	request := events.APIGatewayProxyRequest{}

	var expectedError error
	expectedResponse := events.APIGatewayProxyResponse{
		Body:       "Empty body",
		StatusCode: 400,
	}

	actualResponde, actualError := Handler(request)

	assert.Equal(t, expectedError, actualError)
	assert.Equal(t, expectedResponse, actualResponde)
}

func TestHandlerInvalidDNA(t *testing.T) {
	request := events.APIGatewayProxyRequest{
		Body: invalidDNASequenceStringWrongBases,
	}

	var expectedError error
	expectedResponse := events.APIGatewayProxyResponse{
		Body:       "DNA has invalid bases",
		StatusCode: 400,
	}

	actualResponde, actualError := Handler(request)

	assert.Equal(t, expectedError, actualError)
	assert.Equal(t, expectedResponse, actualResponde)
}
