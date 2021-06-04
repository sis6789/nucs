package file_category

import (
	"path/filepath"
	"sort"
	"strings"
)

var FNumCat []int
var FCatNames []string

func PrepareFileCategory(fileNames []string) {
	var fileNumberCategory []int

	// "-" 앞 이름에서 고유한 것 결정
	var fileCatNames []string
	cats := make(map[string]int)
	for _, n := range fileNames {
		cats[strings.Split(filepath.Base(n), "-")[0]] = 1
	}
	// 정렬
	for k := range cats {
		fileCatNames = append(fileCatNames, k)
	}
	sort.Slice(fileCatNames, func(i, j int) bool {
		return fileCatNames[i] < fileCatNames[j]
	})
	// 고유이름에 순번 부여
	for ix, n := range fileCatNames {
		cats[n] = ix
	}
	// 각 경로에 대한 카테고리 번호 부여
	for _, n := range fileNames {
		fileNumberCategory = append(fileNumberCategory, cats[strings.Split(filepath.Base(n), "-")[0]])
	}

	FNumCat = fileNumberCategory
	FCatNames = fileCatNames
}

func FileCategory(fms uint32) int {
	return FNumCat[fms>>17]
}

func FileCategoryFmss(fmss uint32) int {
	return FNumCat[fmss>>18]
}
