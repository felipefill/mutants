package main

var validDNASequence = []string{"ATCGAAA", "TTGATGA", "GTACCCG", "AAATAAG", "AATTGGG", "AAACCCG", "GTTACCC"}
var invalidDNASequence = []string{"ATCGAAA", "TTGATGG", "GTCCCCA", "ATAT$AT", "AATTG2C", "AAACCCG", "GTTACCX"}

var mutantWithAllCombinationsDNASequence = []string{"ATCGAAA", "TTGATGA", "GTACCCG", "AAATAAG", "AATTGGG", "AAACCCG", "GTTAAAA"}

var mutantDNASequence = []string{"ATCGAAA", "TTGATGA", "GTACCCG", "AAATAAG", "AATTGGG", "AAACCCG", "GTTAAAA"}
var mutantDNASequenceAsJSONString = "{\"Dna\": [\"ATCGAAA\", \"TTGATGA\", \"GTACCCG\", \"AAATAAG\", \"AATTGGG\", \"AAACCCG\", \"GTTAAAA\"]}"
var humanDNASequence = []string{"ATGCGA", "CAGTGC", "TTATTT", "AGACGG", "GCGTCA", "TCACTG"}
var humanDNASequenceAsJSONString = "{\"Dna\": [\"ATGCGA\", \"CAGTGC\", \"TTATTT\", \"AGACGG\", \"GCGTCA\", \"TCACTG\"]}"

var validDNASequenceString = "{\"Dna\": [\"ATCGAAA\", \"TTGATGA\", \"GTACCCG\", \"AAATAAG\", \"AATTGGG\", \"AAACCCG\", \"GTTACCC\"]}"
var invalidDNASequenceStringNotEvenAJSON = "This is clearly not a DNA sequence"
var invalidDNASequenceStringWrongBases = "{\"Dna\": [\"ATGCXA\", \"CAGTGC\", \"TTATTT\", \"AGACGG\", \"GCGTCA\", \"TCACTG\"]}"

var tableNxN = []string{"123", "321", "213"}
var tableMxN = []string{"XXXX", "YYY", "ZZ"}
