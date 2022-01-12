package like

import (
	"fmt"
	"testing"
)

func TestLike(t *testing.T) {
	sStart,
		ratio,
		qStart,
		checkLen,
		countMatch,
		countFault,
		sMatch,
		qMatch := Like(
		//"TGGTTCAGTGCCACACATTGTAGATATTAAATATTTTATATTCAGTGACAGTCATAAACTTGTCCATTGTGTGTAAATAGTATTATGACTTTAACTCTGTGCACATTAGAATACAGTTCAGTTGGCGG",
		"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAGATAGGAAGAGCACACGTCTGAA",
		"GATCGGAAGAGCACACGTCTGAACTCCAGTCAC")
	fmt.Println(
		"sStart", sStart,
		"ratio", ratio,
		"checkLen", checkLen,
		"qStart", qStart,
		"countMatch", countMatch,
		"countFault", countFault,
		"sMatch", sMatch,
		"qMatch", qMatch)
}
