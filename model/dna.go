package model

import (
	"encoding/json"
	"errors"
)

const repetitionRequiredForSequence int = 4

// DNACheck represents a DNA check
type DNACheck struct {
	DNA []string `json:"Dna"`
}

// NewDNACheckFromJSONString creates a DNA check from a json string
func NewDNACheckFromJSONString(data string) (DNACheck, error) {
	dnaCheck := DNACheck{}
	err := json.Unmarshal([]byte(data), &dnaCheck)

	if err != nil {
		return DNACheck{}, errors.New("Could not parse DNA check :" + err.Error())
	}

	err = dnaCheck.validate()
	if err != nil {
		return DNACheck{}, err
	}

	return dnaCheck, nil
}

// IsMutant checks whether this is a DNA sequence from a mutant
func (dnaCheck *DNACheck) IsMutant() bool {
	count := 0

	for row := 0; row < len(dnaCheck.DNA) && count < 2; row++ {
		for column := 0; column < len(dnaCheck.DNA[row]) && count < 2; column++ {
			if dnaCheck.CheckSequenceToTheRight(row, column) {
				count++
			}

			if dnaCheck.CheckSequenceDown(row, column) {
				count++
			}

			if dnaCheck.CheckSequenceDiagonalLeft(row, column) {
				count++
			}

			if dnaCheck.CheckSequenceDiagonalRight(row, column) {
				count++
			}
		}
	}

	return count > 1
}

// CheckSequenceToTheRight checks whether there's a repetition match to the right of given position
func (dnaCheck *DNACheck) CheckSequenceToTheRight(row, column int) bool {
	if len(dnaCheck.DNA[row])-column < repetitionRequiredForSequence {
		return false
	}

	requiredBase := dnaCheck.DNA[row][column]
	c := column + 1

	for loop := 0; loop < repetitionRequiredForSequence-1; loop++ {
		if dnaCheck.DNA[row][c] != requiredBase {
			return false
		}

		c++
	}

	return true
}

// CheckSequenceDown checks whether there's a repetition match looking down of given position
func (dnaCheck *DNACheck) CheckSequenceDown(row, column int) bool {
	if len(dnaCheck.DNA)-row < repetitionRequiredForSequence {
		return false
	}

	requiredBase := dnaCheck.DNA[row][column]
	r := row + 1

	for loop := 0; loop < repetitionRequiredForSequence-1; loop++ {
		if dnaCheck.DNA[r][column] != requiredBase {
			return false
		}

		r++
	}

	return true
}

// CheckSequenceDiagonalLeft checks whether there's a repetition match in the left diagonal of given position
func (dnaCheck *DNACheck) CheckSequenceDiagonalLeft(row, column int) bool {
	if row >= repetitionRequiredForSequence || column < repetitionRequiredForSequence-1 {
		return false
	}

	requiredBase := dnaCheck.DNA[row][column]
	r := row + 1
	c := column - 1

	for loop := 0; loop < repetitionRequiredForSequence-1; loop++ {
		if dnaCheck.DNA[r][c] != requiredBase {
			return false
		}

		r++
		c--
	}

	return true
}

// CheckSequenceDiagonalRight checks whether there's a repetition match in the right diagonal of given position
func (dnaCheck *DNACheck) CheckSequenceDiagonalRight(row, column int) bool {
	if row >= repetitionRequiredForSequence || column >= repetitionRequiredForSequence {
		return false
	}

	requiredBase := dnaCheck.DNA[row][column]
	r := row + 1
	c := column + 1

	for loop := 0; loop < repetitionRequiredForSequence-1; loop++ {
		if dnaCheck.DNA[r][c] != requiredBase {
			return false
		}

		r++
		c++
	}

	return true
}

func (dnaCheck *DNACheck) validate() error {
	if !dnaCheck.isValidNxNTable() {
		return errors.New("DNA is not an NxN table")
	}

	if !dnaCheck.validateDNAHasOnlyValidBases() {
		return errors.New("DNA has invalid bases")
	}

	return nil
}

func (dnaCheck *DNACheck) isValidNxNTable() bool {
	// This will give the number of rows
	hypothesis := len(dnaCheck.DNA)

	for row := 0; row < hypothesis; row++ {
		// If we have a number of columns that's different from our hypothesis
		if len(dnaCheck.DNA[row]) != hypothesis {
			return false
		}
	}

	return true
}

func (dnaCheck *DNACheck) validateDNAHasOnlyValidBases() bool {
	//TODO: This could be done while checking DNA
	for row := 0; row < len(dnaCheck.DNA); row++ {
		for column := 0; column < len(dnaCheck.DNA[row]); column++ {
			currentChar := dnaCheck.DNA[row][column]
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
