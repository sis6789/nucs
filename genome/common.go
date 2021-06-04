package genome

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/sis6789/nucs/nuc2"
)

type Genome struct {
	RecKey          string      `bson:"_id"`
	GChr            int         `bson:"gchr"`               //염색체번호
	GFrom           int         `bson:"gfrom"`              //서열 시작위치
	GTo             int         `bson:"gto"`                //서열 종료위치
	Len             int         `bson:"len"`                // 서열 길이
	Human           string      `bson:"human"`              // Poll result, 인간 게놈 자료
	Poll            string      `bson:"poll"`               //종합 판결 서열
	PollQuality     string      `bson:"pollquality"`        //종합 판결 정확도(0-9) 제일 많은 염기 비율; "." 전부 일치
	CountNucs       int         `bson:"countnucs"`          //판결에 사용한 전체 염기수, 평평화된 것의 합
	CountRead       int         `bson:"countread"`          //사용된 library 수량
	ASequence       []string    `bson:"asequence"`          //each Fms information, 각 FMS 시작과 끝에 맞추어 조정한 서열 (전부 동일한 길이)
	AQuality        []string    `bson:"aquality"`           //각 FMS 시작과 끝에 맞추어 조정한 품질 (전부 동일한 길이)
	PollDifference  []string    `bson:"polldifference"`     //종합 판결 서열과 다른 값들 표시
	Ratio           []float32   `bson:"ratio"`              //전체 염기수 대비 대표서열과 상이한 비율
	FmsList         []FlatFMS   `bson:"fmslist"`            //판결에 사용된 모든 FMS 정보
	CatNames        []string    `bson:"catnames"`           //File group information
	CatCountRead    []int       `bson:"catcountread"`       //
	CatCountNuc     []int       `bson:"catcountnuc"`        //
	CatCount        int         `bson:"catcount"`           //
	CountFms        int         `bson:"count_fms"`          //평편화된 FMS 관련 수치
	CountFmsNucs    int         `bson:"count_fms_nucs"`     //
	CatCountFmsRead []int       `bson:"cat_count_fms_read"` //
	CatCountFmsNucs []int       `bson:"cat_count_fms_nucs"` //
	UniqueFms       []DiffMerge `bson:"unique_fms"`         // Unique sequence set
}

type DiffMerge struct {
	FmsLink         []int  `bson:"fms_link"`  // related fms index number
	CatCount        int    `bson:"cat_count"` // count of active file category
	Difference      string `bson:"difference"`
	Modify          string `bson:"modify"`
	CountRead       int    `bson:"count_read"`
	CountNuc        int    `bson:"count_nuc"`
	CatRead         []int  `bson:"cat_read"`
	CatNuc          []int  `bson:"cat_nuc"`
	CountFms        int    `bson:"count_fms"`
	CountFmsNucs    int    `bson:"count_fms_nucs"`
	CatCountFmsRead []int  `bson:"cat_count_fms_read"`
	CatCountFmsNucs []int  `bson:"cat_count_fms_nucs"`
	IsMut           bool   `bson:"is_mut"`
}

func (x *Genome) MakeUniqueFms() {
	m := make(map[string]DiffMerge)
	for ix := 0; ix < len(x.ASequence); ix++ {
		d, ok := m[x.PollDifference[ix]]
		if !ok {
			d.CatRead = make([]int, len(x.CatNames))
			d.CatNuc = make([]int, len(x.CatNames))
			d.CatCountFmsRead = make([]int, len(x.CatNames))
			d.CatCountFmsNucs = make([]int, len(x.CatNames))
			d.Difference = x.PollDifference[ix]
			d.IsMut = false
		}
		var catNum int
		for catNum = 0; x.CatCountFmsRead[catNum] == 0; catNum++ {
		}
		d.FmsLink = append(d.FmsLink, ix)
		d.CatRead[catNum] += x.FmsList[ix].CountRead
		d.CatNuc[catNum] += x.FmsList[ix].CountNucs
		d.CountRead += x.FmsList[ix].CountRead
		d.CountNuc += x.FmsList[ix].CountNucs
		d.Modify += x.FmsList[ix].ModifySequence
		d.CountFms++
		d.CountFmsNucs += x.FmsList[ix].FlatNucCount
		d.CatCountFmsRead[catNum]++
		d.CatCountFmsNucs[catNum] += x.FmsList[ix].FlatNucCount
		m[x.PollDifference[ix]] = d
	}
	for _, v := range m {
		x.UniqueFms = append(x.UniqueFms, v)
	}
}

type FlatFMS struct {
	RecKey         string    `bson:"_id"`
	GChr           int       `bson:"gchr"`           //염색체번호
	GFrom          int       `bson:"gfrom"`          //서열 시작위치
	GTo            int       `bson:"gto"`            //서열 종료위치
	Fms            uint32    `bson:"fms"`            //해당 FMS 값
	Sequence       string    `bson:"sequence"`       //판결 서열
	Quality        string    `bson:"quality"`        //판결 정확도(0-9) 제일 많은 염기 비율; "." 전부 일치
	QualityTB      [2]string `bson:"qualitytb"`      //판결 정확도(0-9) 제일 많은 염기 비율; "." 전부 일치
	ModifySequence string    `bson:"modifysequence"` //BTOP 변형 정보
	CountNucs      int       `bson:"countnucs"`      //사용된 모든 염기 갯수
	CountNucsTB    [2]int    `bson:"countnucstb"`    //사용된 모든 상위와 하위 염기 갯수
	CountRead      int       `bson:"countread"`      //사용된 library 갯수
	CountReadTB    [2]int    `bson:"countreadtb"`    //사용된 library 상위와 하위 갯수
	QueryID        []int     `bson:"queryid"`        //사용된 시료 번호 리스트
	QueryCount     []int     `bson:"querycount"`     //사용된 시료 번호 갯수
	FlatNucCount   int       `bson:"flat_nuc_count"` //편평화된 fms 구성 염기 수 (정상만 계수)
}

func (x *FlatFMS) Set(v NucSeq) {
	x.GChr = v.Chr
	x.GFrom = v.PosStart
	if v.PosLast[0] >= v.PosLast[1] {
		x.GTo = v.PosLast[0]
	} else {
		x.GTo = v.PosLast[1]
	}
	x.Fms = uint32(v.Fms)
	x.Sequence = v.NucAll
	x.Quality = v.NucAllQuality
	x.QualityTB[0] = v.NucQualityTB[0]
	x.QualityTB[1] = v.NucQualityTB[1]
	modString := ""
	// sort
	if len(v.Modify) > 0 {
		var kSlice []string
		for k := range v.Modify {
			kSlice = append(kSlice, k)
		}
		sort.Slice(kSlice, func(i, j int) bool {
			return kSlice[i][1:] < kSlice[j][1:]
		})
		for _, k := range kSlice {
			modString += fmt.Sprintf("%s%02d", k, v.Modify[k])
		}
	}
	x.ModifySequence = modString
	x.CountNucs = v.NucCount[0] + v.NucCount[1]
	x.CountNucsTB[0] = v.NucCount[0]
	x.CountNucsTB[1] = v.NucCount[1]
	rCnt := v.AllReadCount()
	x.CountRead = rCnt[0] + rCnt[1]
	x.CountReadTB[0] = rCnt[0]
	x.CountReadTB[1] = rCnt[1]
	for k, v := range v.Qid[0] {
		x.QueryID = append(x.QueryID, k)
		x.QueryCount = append(x.QueryCount, v)
	}
	for k, v := range v.Qid[1] {
		x.QueryID = append(x.QueryID, k)
		x.QueryCount = append(x.QueryCount, v)
	}
}

type NucLine struct {
	Gchr      int    //염색체번호
	Gpos      int    //염기위치
	Fms       int    //시료 FMS
	Side      int    //시료 Top/Bottom
	CountRead int    //염기위치 총 수량
	Qid       int    //시료 번호
	Qstart    int    //시료 시작위치
	Qpos      int    //해당염기의 시작위치부터 위치
	Seq       string //염기 값, BTOP 변이 값은 2개 이상의 염기로 표현한다.
}

func (r *NucLine) ParseLine(line string) *NucLine {
	_, err := fmt.Sscanf(line, "%d %d %d %d %d %d %d %d %s",
		&r.Gchr, &r.Gpos, &r.Fms, &r.Side, &r.CountRead, &r.Qid, &r.Qstart, &r.Qpos, &r.Seq)
	if err != nil {
		log.Fatalln(err)
	}
	return r
}

type NucSeq struct {
	Fms           int
	Side          int
	Chr           int
	PosStart      int
	PosLast       [2]int
	ReadCount     [2]int
	NucCount      [2]int
	NucAll        string
	NucAllQuality string
	NucTB         [2]string
	NucQualityTB  [2]string
	Modify        map[string]int
	Qid           [2]map[int]int
	IsStarted     bool
}

func (x *NucSeq) SameGroup(y NucSeq) bool {
	cmp := x.Fms == y.Fms && x.Chr == y.Chr && (y.PosLast[y.Side] == x.PosLast[y.Side] || y.PosLast[y.Side] == x.PosLast[y.Side]+1)
	return cmp
}

func (x *NucSeq) SameNuc(y NucSeq) bool {
	return x.Fms == y.Fms && x.Chr == y.Chr && x.PosLast[y.Side] == y.PosStart
}

func (x *NucSeq) SameGpos(y NucLine) bool {
	return x.Fms == y.Fms && x.Chr == y.Gchr && x.PosStart == y.Gpos
}

func (x NucSeq) String() string {
	modString := ""
	// sort
	if len(x.Modify) > 0 {
		var kSlice []string
		for k := range x.Modify {
			kSlice = append(kSlice, k)
		}
		sort.Slice(kSlice, func(i, j int) bool {
			return kSlice[i][1:] < kSlice[j][1:]
		})
		for _, k := range kSlice {
			modString += fmt.Sprintf("%s%02d", k, x.Modify[k])
		}
	}
	rCnt := x.AllReadCount()
	return fmt.Sprintf("FMS:%d, CHR=%d, POS=%d~%d(%d), nuc=%d/%d, read=%d/%d\n\tnuc=%s\n\tmod=%s",
		x.Fms, x.Chr, x.PosStart, x.PosLast, x.PosLast[0]-x.PosStart+1,
		x.NucCount[0], x.NucCount[1],
		rCnt[0], rCnt[1],
		x.NucTB[0], modString)
}

func (x *NucSeq) Set(nl NucLine) {
	// prepare structure and reset all value
	x.Modify = make(map[string]int)
	x.Qid[0] = make(map[int]int)
	x.Qid[1] = make(map[int]int)
	x.ReadCount = [2]int{0, 0}
	x.NucCount = [2]int{0, 0}
	x.NucAll = ""
	x.NucTB = [2]string{"", ""}
	x.NucQualityTB = [2]string{"", ""}
	x.IsStarted = false

	x.Fms = nl.Fms
	x.Side = nl.Side
	x.Chr = nl.Gchr
	x.PosStart = nl.Gpos
	x.PosLast[0] = nl.Gpos
	x.PosLast[1] = nl.Gpos
	x.ReadCount[nl.Side] = nl.CountRead
	x.NucCount[nl.Side] = nl.CountRead
	if len(nl.Seq) > 1 {
		x.AddModify(nl)
	} else {
		x.NucTB[nl.Side] = strings.Repeat(nl.Seq, nl.CountRead)
	}
	x.NucAll = x.NucTB[0] + x.NucTB[1]
	x.AddQid(nl)
}

func (x NucSeq) IsValid() bool {
	if len(x.NucTB[0]) < 3 || len(x.NucTB[1]) < 3 {
		return false
	}
	nucTop, _ := Compress(x.NucTB[0])
	nucBottom, _ := Compress(x.NucTB[1])
	if nucTop != nucBottom {
		return false
	}
	return true
}

func (x *NucSeq) Add(nl NucLine) {
	if len(nl.Seq) > 1 {
		x.AddModify(nl)
	} else {
		if x.ReadCount[nl.Side] < nl.CountRead {
			x.ReadCount[nl.Side] = nl.CountRead
		}
		x.NucCount[nl.Side] += nl.CountRead
		x.NucTB[nl.Side] += strings.Repeat(nl.Seq, nl.CountRead)
		x.AddQid(nl)
	}
	x.NucAll = x.NucTB[0] + x.NucTB[1]
}

func (x *NucSeq) Append(y NucSeq) {
	for ix := 0; ix < 2; ix++ {
		if x.PosLast[ix] < y.PosLast[ix] {
			x.PosLast[ix] = y.PosLast[ix]
		}
		if x.ReadCount[ix] < y.ReadCount[ix] {
			x.ReadCount[ix] = y.ReadCount[ix]
		}
		x.NucCount[ix] += y.NucCount[ix]
		x.NucTB[ix] += "/" + y.NucTB[ix]
	}
	x.NucAll += "/" + y.NucTB[0] + y.NucTB[1]
	for k, v := range y.Modify {
		vx := x.Modify[k]
		vx += v
		x.Modify[k] = vx
	}
	for ix := 0; ix < 2; ix++ {
		for k, v := range y.Qid[ix] {
			vx := x.Qid[ix][k]
			if vx < v {
				vx = v
			}
			x.Qid[ix][k] = vx
		}
	}
}

func (x *NucSeq) AddModify(n1 NucLine) {
	tbStr := []string{"t", "b"}
	modCode := fmt.Sprintf("%s%03d%s", tbStr[n1.Side], n1.Gpos%1000, n1.Seq)
	v := x.Modify[modCode]
	v += n1.CountRead
	x.Modify[modCode] = v
}

func (x *NucSeq) AddQid(nl NucLine) {
	v := x.Qid[nl.Side][nl.Qid]
	if v < nl.CountRead {
		v = nl.CountRead
	}
	x.Qid[nl.Side][nl.Qid] = v
}

// AllReadCount return count of fastq lines
func (x NucSeq) AllReadCount() [2]int {
	vSum := [2]int{0, 0}
	for side := 0; side < 2; side++ {
		for _, v := range x.Qid[side] {
			vSum[side] += v
		}
	}
	return vSum
}

func (x *NucSeq) Finalize() (int, FlatFMS) {
	returnClass := 0
	var flatFMS FlatFMS
	flatFMS.RecKey = uuid.NewString()

	rCnt := x.AllReadCount() // 각 서열의 발생 건수를 다 더한다.
	if rCnt[0] > 1 && rCnt[1] > 1 {
		// 상위 하위 strand 2개 이상이 공존
		x.NucTB[0], x.NucQualityTB[0] = Compress(x.NucTB[0])
		x.NucTB[1], x.NucQualityTB[1] = Compress(x.NucTB[1])
		x.NucAll = nuc2.Nuc2String(x.NucTB[0], x.NucTB[1])
		x.NucAllQuality = nuc2.Nuc2DString(x.NucTB[0], x.NucTB[1])
		//x.NucAll, x.NucAllQuality = Compress(x.NucAll)

		// 상위와 하위 동일성 상관 없이 적합한 것으로 최종 결과에 반영
		// 일부 틀리는 부분은 틀리는 표시를 하도록 함
		returnClass = 0 // valid
		flatFMS.Set(*x)
	} else {
		// 상위 하위 수량이 1개 또는 없음
		if rCnt[0] == 0 || rCnt[1] == 0 {
			// 상위나 하위 하나만 존재
			returnClass = 2 // orphanCount++
		} else {
			// 상위 하나 하위 하나만 존재
			returnClass = 3 // tbOneCount++
		}
	}
	return returnClass, flatFMS
}

func Compress(w string) (string, string) {
	vMost := ""
	vQuality := ""
	tks := strings.Split(w, "/")
	for _, s := range tks {
		if len(s) > 0 {
			most, quality := MaxNuc(s)
			vMost += most
			vQuality += quality
		}
	}
	return vMost, vQuality
}

func MaxNuc(w string) (mostNuc string, quality string) {
	type ac struct {
		a byte
		c int
	}
	var count [128]ac
	for _, s := range []byte(w) {
		count[s].a = s
		count[s].c++
	}
	wc := count[:]
	sort.Slice(wc, func(i, j int) bool {
		return wc[i].c > wc[j].c
	})
	// check top 2 nuc count
	// 상위 2개가 동일한 개수이면 불명(*)으로 결정한다.
	if wc[0].c == wc[1].c {
		mostNuc = "*"
	} else {
		mostNuc = string(wc[0].a)
	}
	// quality
	quality = "."
	if wc[0].c < len(w) {
		quality = strconv.Itoa(int(float64(wc[0].c) / float64(len(w)) * 10.0))
	}

	return mostNuc, quality
}
