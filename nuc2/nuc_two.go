// Package nuc2 Determine nucleotide value between Top strand and Bottom strand
package nuc2

var v [128][128]byte
var d [128][128]byte

// 변환 테이블 초기화
func init() {
	for i := 0; i < 128; i++ {
		for j := 0; j < 128; j++ {
			v[i][j] = '*'
			d[i][j] = '*'
		}
	}
	v['-']['-'] = '-'
	v['A']['A'] = 'A'
	v['A']['C'] = 'm'
	v['A']['G'] = 'r'
	v['A']['T'] = 'w'
	v['C']['A'] = 'm'
	v['C']['C'] = 'C'
	v['C']['G'] = 's'
	v['C']['T'] = 'y'
	v['G']['A'] = 'r'
	v['G']['C'] = 's'
	v['G']['G'] = 'G'
	v['G']['T'] = 'k'
	v['T']['A'] = 'w'
	v['T']['C'] = 'y'
	v['T']['G'] = 'k'
	v['T']['T'] = 'T'

	d['-']['-'] = '.'
	d['A']['A'] = '.'
	d['A']['C'] = 'm'
	d['A']['G'] = 'r'
	d['A']['T'] = 'w'
	d['C']['A'] = 'm'
	d['C']['C'] = '.'
	d['C']['G'] = 's'
	d['C']['T'] = 'y'
	d['G']['A'] = 'r'
	d['G']['C'] = 's'
	d['G']['G'] = '.'
	d['G']['T'] = 'k'
	d['T']['A'] = 'w'
	d['T']['C'] = 'y'
	d['T']['G'] = 'k'
	d['T']['T'] = '.'
}

func Nuc2(n1, n2 byte) byte {
	return v[n1][n2]
}

func Nuc2String(s1, s2 string) string {
	b1 := []byte(s1)
	b2 := []byte(s2)
	if len(b2) > len(b1) {
		b2 = b2[:len(b1)]
	} else if len(b1) > len(b2) {
		b1 = b1[:len(b2)]
	}
	var b3 []byte
	for ix := 0; ix < len(b1); ix++ {
		b3 = append(b3, v[b1[ix]][b2[ix]])
	}
	return string(b3)
}

func Nuc2D(n1, n2 byte) byte {
	return d[n1][n2]
}

func Nuc2DString(s1, s2 string) string {
	b1 := []byte(s1)
	b2 := []byte(s2)
	if len(b2) > len(b1) {
		b2 = b2[:len(b1)]
	} else if len(b1) > len(b2) {
		b1 = b1[:len(b2)]
	}
	var b3 []byte
	for ix := 0; ix < len(b1); ix++ {
		b3 = append(b3, d[b1[ix]][b2[ix]])
	}
	return string(b3)
}
