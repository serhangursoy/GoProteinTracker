package Gonome

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"sync/atomic"
)

// Database location goes here
var DATABASE = "./Data/ProteinDatabase.txt"
var IND int = 0

/*  Variables  */
var SEQUENCE_ARRAY []PartialSequence
var AA_SEQUENCE_ARRAY []string
var isFound uint32

type PartialSequence struct {
	HEADER   string
	SEQUENCE string
}

func StartSearch(file_name string, DBPARAM string, INDEX int) {
	DATABASE = DBPARAM
	IND = INDEX
	// First read the given source file...
	b, err := ioutil.ReadFile(file_name) // Fetch
	if err != nil {
		fmt.Print(err)
	} else {
		file_content := string(b) // convert content to a 'string'
		content_arr := strings.Split(file_content, "\n")
		var NA_SEQ string = ""
		var NA_HEA string = ""

		for i := 0; i < len(content_arr); i++ {
			if strings.Index(content_arr[i], ">") == -1 {
				NA_SEQ = NA_SEQ + content_arr[i]
			} else {
				if NA_HEA != "" {
					var TempSequence PartialSequence
					TempSequence.HEADER = NA_HEA
					TempSequence.SEQUENCE = NA_SEQ
					SEQUENCE_ARRAY = append(SEQUENCE_ARRAY, TempSequence)
					NA_HEA = content_arr[i]
					NA_SEQ = ""
				} else {
					// Means that its the first in the file
					NA_HEA = content_arr[i]
				}
			}
		}

		// Now we have all contents in the RNA Sequence.. We need to change their
		// Nucleic Acid codons to corresponding Amino Acid Representations.
		for i := 0; i < len(SEQUENCE_ARRAY); i++ {
			_AA := AATranslation(SEQUENCE_ARRAY[i].SEQUENCE)
			AA_SEQUENCE_ARRAY = append(AA_SEQUENCE_ARRAY, _AA)
		}

		var RESULT string = ""

		fmt.Println("\nFor (" + (SEQUENCE_ARRAY[IND].HEADER)[1:] + "): ")
		RESULT = RESULT + MatchWithDataBase(AA_SEQUENCE_ARRAY[IND], 2)
		if RESULT == "" {
			fmt.Println("No Match")
		} else {
			fmt.Println(RESULT)
		}
	}
}

func MatchWithDataBase(SEQ string, parallelCount int) string {

	// Create our channel
	c := make(chan string)

	b, err := ioutil.ReadFile(DATABASE) // Fetch DB
	if err != nil {
		fmt.Print(err)
	} else {
		file_content := string(b) // convert content to a 'string'
		content_arr := strings.Split(file_content, "\n")
		cont_length := len(content_arr) / parallelCount
		lastInd := 0

		for k := 0; k < parallelCount; k++ {
			go Matcher(SEQ, content_arr[lastInd:lastInd+cont_length], c)
			lastInd = lastInd + cont_length
		}
	}

	response := <-c
	close(c)
	return response
}

func Matcher(SEQ_DATA string, DATA_SLICE []string, c chan string) {
	for i := 0; i < len(DATA_SLICE)-1; i++ {
		if isFound == 0 {
			var header string = DATA_SLICE[i]
			if strings.Index(DATA_SLICE[i+1], SEQ_DATA) != -1 {
				var CONFIDENCE float32 = (float32)(len(SEQ_DATA)) / (float32)(len(DATA_SLICE[i+1])) * 100
				// Means that we have found it!
				atomic.AddUint32(&isFound, 1)
				c <- "Found--> " + header[1:] + " with confidence " + strconv.FormatFloat((float64)(CONFIDENCE), 'f', 2, 64) + "%\nSequence: " + DATA_SLICE[i+1]
				return
			}
		} else {
			break
		}
		i++
	}

	if isFound == 0 {
		c <- ""
	}
}

// Amino Acid Translation. Converts every codon with three Nucleic Acid into corresponding Amino Acid
func AATranslation(NA_SEQUENCE string) string {
	// var startIndex int = 0
	var AA_SEQ = ""
	var STARTED = true
	for i := 0; i < len(NA_SEQUENCE)-3; i += 3 {
		CODON := NA_SEQUENCE[i:(i + 3)]
		var AA = codonTranslate(CODON)
		if AA == "STOP" {
			STARTED = false
		} else if AA == "START" {
			STARTED = true
			// Start parsing.. But this wont come since DB already provide us the started sequence
		} else if STARTED {
			AA_SEQ = AA_SEQ + AA
		}
	}

	return AA_SEQ
}

// Simple Codon translation table.
func codonTranslate(CODON string) string {
	// This is tricky
	if CODON == "ATT" || CODON == "ATC" || CODON == "ATA" {
		return "I"
	} else if CODON == "CTT" || CODON == "CTC" || CODON == "CTA" || CODON == "CTG" || CODON == "TTA" || CODON == "TTG" {
		return "L"
	} else if CODON == "GTT" || CODON == "GTC" || CODON == "GTA" || CODON == "GTG" {
		return "V"
	} else if CODON == "TTT" || CODON == "TTC" {
		return "F"
	} else if CODON == "ATG" {
		return "M"
	} else if CODON == "TGT" || CODON == "TGC" {
		return "C"
	} else if CODON == "GCT" || CODON == "GCC" || CODON == "GCA" || CODON == "GCG" {
		return "A"
	} else if CODON == "GGT" || CODON == "GGC" || CODON == "GGA" || CODON == "GGG" {
		return "G"
	} else if CODON == "CCT" || CODON == "CCC" || CODON == "CCA" || CODON == "CCG" {
		return "P"
	} else if CODON == "ACT" || CODON == "ACC" || CODON == "ACA" || CODON == "ACG" {
		return "T"
	} else if CODON == "TCT" || CODON == "TCC" || CODON == "TCA" || CODON == "TCG" || CODON == "AGT" || CODON == "AGC" {
		return "S"
	} else if CODON == "TAT" || CODON == "TAC" {
		return "Y"
	} else if CODON == "TGG" {
		return "W"
	} else if CODON == "CAA" || CODON == "CAG" {
		return "Q"
	} else if CODON == "AAT" || CODON == "AAC" {
		return "N"
	} else if CODON == "CAT" || CODON == "CAC" {
		return "H"
	} else if CODON == "GAA" || CODON == "GAG" {
		return "E"
	} else if CODON == "GAT" || CODON == "GAC" {
		return "D"
	} else if CODON == "AAA" || CODON == "AAG" {
		return "K"
	} else if CODON == "CGT" || CODON == "CGC" || CODON == "CGA" || CODON == "CGG" || CODON == "AGA" || CODON == "AGG" {
		return "R"
	} else if CODON == "ATG" {
		return "START"
	} else if CODON == "TAA" || CODON == "TAG" || CODON == "TGA" {
		return "STOP"
	} else {
		return "NONE"
	}
}
