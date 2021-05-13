package fmsspack

var nucDigits = []byte{'A', 'C', 'G', 'T'}
var nucDigitsBase = uint32(len(nucDigits))

//var shiftString = []string{"", "T", "GT", "AGT", "CAGT"}
var shiftStringDecode = []string{"----", "T---", "GT--", "AGT-", "CAGT"}

func nuc2Quad(nucleotides string) uint32 {
	quadNumber := uint32(0)
	for _, b := range []byte(nucleotides) {
		quadNumber *= nucDigitsBase
		switch b {
		case 'A', 'a':
			quadNumber += 0
		case 'C', 'c':
			quadNumber += 1
		case 'G', 'g':
			quadNumber += 2
		case 'T', 't':
			quadNumber += 3
		}
	}
	return quadNumber
}

func quad2Nuc(quadValue uint32) string {
	var nucBytes []byte
	for quadValue > 0 {
		ix := quadValue % nucDigitsBase
		quadValue /= nucDigitsBase
		nucBytes = append(nucBytes, nucDigits[ix])
	}
	reverseBytes(nucBytes)
	return string(nucBytes)
}

func reverseBytes(numbers []byte) {
	for i, j := 0, len(numbers)-1; i < j; i, j = i+1, j-1 {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
}

func Encode(fileNumber int, molecularIndex string, shift string, strand string) uint32 {
	result := uint32(fileNumber) << 18
	// Molecular Index - 7 ; 14bit
	result |= nuc2Quad(molecularIndex) << 4
	// none|T|GT|AGT|CAGT - 0~4; 3bit
	result |= uint32(len(shift)) << 1
	// strand
	if strand == "C" {
		result |= 1
	}
	return result
}

func Decode(val uint32) (fileNumber int, molecularIndex string, shift string, strand string) {
	// strand
	if (val & 1) == 0 {
		strand = "A"
	} else {
		strand = "C"
	}
	// shift
	shift1 := (val >> 1) & 7
	shift = shiftStringDecode[shift1]
	// molecular index
	molecularIndexDecimal := (val >> 4) & 0x3fff
	molecularIndex = quad2Nuc(molecularIndexDecimal)
	molecularIndex = "AAAAAAA"[0:7-len(molecularIndex)] + molecularIndex // fill 7 bps
	// File sequence
	fileNumber = int(val >> 18) // int(val >> 18 & 0x1fff)

	return fileNumber, molecularIndex, shift, strand
}
