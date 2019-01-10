package utils

import (
	"os"
	"testing"

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
