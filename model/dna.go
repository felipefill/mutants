package model

// DNA represents a DNA
type DNA struct {
	DNA []string `json:"Dna"`
}

const repetitionRequiredForSequence int = 4

// CheckSequenceToTheRight checks whether there's a repetition match to the right of given position
func (dna *DNA) CheckSequenceToTheRight(row, column int) bool {
	if len(dna.DNA[row])-column < repetitionRequiredForSequence {
		return false
	}

	requiredBase := dna.DNA[row][column]
	c := column + 1

	for loop := 0; loop < repetitionRequiredForSequence-1; loop++ {
		if dna.DNA[row][c] != requiredBase {
			return false
		}

		c++
	}

	return true
}

// CheckSequenceDown checks whether there's a repetition match looking down of given position
func (dna *DNA) CheckSequenceDown(row, column int) bool {
	if len(dna.DNA)-row < repetitionRequiredForSequence {
		return false
	}

	requiredBase := dna.DNA[row][column]
	r := row + 1

	for loop := 0; loop < repetitionRequiredForSequence-1; loop++ {
		if dna.DNA[r][column] != requiredBase {
			return false
		}

		r++
	}

	return true
}

// CheckSequenceDiagonalLeft checks whether there's a repetition match in the left diagonal of given position
func (dna *DNA) CheckSequenceDiagonalLeft(row, column int) bool {
	if row >= repetitionRequiredForSequence || column < repetitionRequiredForSequence-1 {
		return false
	}

	requiredBase := dna.DNA[row][column]
	r := row + 1
	c := column - 1

	for loop := 0; loop < repetitionRequiredForSequence-1; loop++ {
		if dna.DNA[r][c] != requiredBase {
			return false
		}

		r++
		c--
	}

	return true
}

// CheckSequenceDiagonalRight checks whether there's a repetition match in the right diagonal of given position
func (dna *DNA) CheckSequenceDiagonalRight(row, column int) bool {
	if row >= repetitionRequiredForSequence || column >= repetitionRequiredForSequence {
		return false
	}

	requiredBase := dna.DNA[row][column]
	r := row + 1
	c := column + 1

	for loop := 0; loop < repetitionRequiredForSequence-1; loop++ {
		if dna.DNA[r][c] != requiredBase {
			return false
		}

		r++
		c++
	}

	return true
}
