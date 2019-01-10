package utils

import (
	"os"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestMustGetEnvVar(t *testing.T) {
	varName := "MUST_GET_ENV_VAR_TEST"
	expected := "yay"

	os.Setenv(varName, expected)

	actual := MustGetEnvVar(varName)
	os.Clearenv()

	if expected != actual {
		t.Errorf("Expected %s and got %s", expected, actual)
	}
}

func TestMustGetEnvVarPanicsWhenFails(t *testing.T) {
	os.Clearenv()
	assert.Panics(t, func() { MustGetEnvVar("ANY_VAR") }, "Code should have panicked")
}

func TestInjectDatabase(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer db.Close()

	InjectDatabase(db)

	assert.Equal(t, db, GetDB())
}

func TestGETDBPanics(t *testing.T) {
	os.Clearenv()
	_db = nil
	assert.Panics(t, func() { GetDB() })
}

func TestGetDatabaseInfo(t *testing.T) {
	expectedHost := "localhost"
	expectedName := "mutants"
	expectedUser := "felipefill"
	expectedPswd := "secret"

	os.Setenv("DB_HOST", expectedHost)
	os.Setenv("DB_NAME", expectedName)
	os.Setenv("DB_USER", expectedUser)
	os.Setenv("DB_PSWD", expectedPswd)

	actualHost, actualName, actualUser, actualPswd := getDatabaseInfo()
	os.Clearenv()

	assert.Equal(t, expectedHost, actualHost)
	assert.Equal(t, expectedName, actualName)
	assert.Equal(t, expectedUser, actualUser)
	assert.Equal(t, expectedPswd, actualPswd)
}
