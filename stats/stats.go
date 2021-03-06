package main

import (
	"errors"

	"github.com/felipefill/mutants/utils"
)

// Stats struct that holds DNA status information
type Stats struct {
	MutantDNACount int     `json:"count_mutant_dna"`
	HumanDNACount  int     `json:"count_human_dna"`
	Ratio          float64 `json:"ratio"`
}

// GetStats retrieve status regarding the number of mutant and ordinary human DNAs
func GetStats() (*Stats, error) {
	db := utils.GetDB()

	rows, err := db.Query("select count(id) count, type from dna group by type")
	if err != nil {
		return nil, errors.New("Failed to query database")
	}

	mutantCount := 0
	humanCount := 0
	ratio := float64(0)

	for rows.Next() {
		var dnaType string
		var count int

		err = rows.Scan(&count, &dnaType)
		if err != nil {
			return nil, errors.New("Failed to retrieve status")
		}

		if dnaType == "mutant" {
			mutantCount = count
			humanCount = humanCount + count
		} else {
			humanCount = humanCount + count
		}
	}

	if humanCount > 0 {
		ratio = float64(mutantCount) / float64(humanCount)
	}

	stats := &Stats{
		HumanDNACount:  humanCount,
		MutantDNACount: mutantCount,
		Ratio:          ratio,
	}

	return stats, nil
}
