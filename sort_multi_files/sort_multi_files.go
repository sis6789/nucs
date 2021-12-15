package sort_multi_files

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sync"

	"github.com/google/uuid"
)

func SortFiles(folder string, pattern string, keyList ...string) {
	var guard = make(chan struct{}, 5) // 동시수행을 5개로 제한
	var wg sync.WaitGroup

	targets := targetList(folder, pattern)
	for _, fn := range targets {
		guard <- struct{}{}
		wg.Add(1)
		go doSort(&wg, &guard, folder, fn, keyList...)
	}
	wg.Wait()
	close(guard)
}

func targetList(folder string, pattern string) []string {
	// 대상 파일 목록
	regPattern := regexp.MustCompile(pattern)
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		log.Fatalln(err)
	}
	var fnList []string
	for _, fi := range files {
		if fi.IsDir() {
			continue
		}
		if regPattern.MatchString(fi.Name()) {
			fnList = append(fnList, fi.Name())
		}
	}
	return fnList
}

func doSort(wg *sync.WaitGroup, guard *chan struct{}, folder, fileName string, keyList ...string) {
	defer func() {
		wg.Done()
		<-*guard
	}()

	sortInFn := filepath.Join(folder, fileName)
	sortOutFn := filepath.Join(folder, "so_"+uuid.NewString()+".txt")
	sortErrFn := filepath.Join(folder, "se_"+uuid.NewString()+".txt")

	sortInFile, err := os.Open(sortInFn)
	if err != nil {
		log.Fatalln(err)
	}
	sortOutFile, err := os.Create(sortOutFn)
	if err != nil {
		log.Fatalln(err)
	}
	sortErrFile, err := os.Create(sortErrFn)
	if err != nil {
		log.Fatalln(err)
	}

	var cmdPath string
	var sortOrder []string
	if runtime.GOOS == "windows" {
		cmdPath, err = exec.LookPath("wsl")
		if err != nil {
			log.Fatalln(err)
			return
		}
		sortOrder = []string{cmdPath, "sort"}
	} else {
		cmdPath, err = exec.LookPath("sort")
		if err != nil {
			log.Fatalln(err)
			return
		}
		sortOrder = []string{cmdPath}
	}
	sortOrder = append(sortOrder, keyList...)
	cmd := exec.Cmd{
		Path:   cmdPath,
		Args:   sortOrder,
		Dir:    folder,
		Stdin:  sortInFile,
		Stdout: sortOutFile,
		Stderr: sortErrFile,
	}
	// invoke sort
	if startErr := cmd.Start(); startErr != nil {
		log.Print(startErr)
		_ = sortInFile.Close()
		_ = sortOutFile.Close()
		_ = sortErrFile.Close()
		log.Fatalln(fileName)
	}
	// wait end of sort
	if endErr := cmd.Wait(); endErr != nil {
		log.Print(endErr)
		_ = sortInFile.Close()
		_ = sortOutFile.Close()
		_ = sortErrFile.Close()
		log.Fatalln(fileName)
	}
	// close
	if err = sortInFile.Close(); err != nil {
		log.Fatalln(err)
	}
	if err = sortOutFile.Sync(); err != nil {
		log.Fatalln(err)
	} else if err = sortOutFile.Close(); err != nil {
		log.Fatalln(err)
	}
	if err = sortErrFile.Sync(); err != nil {
		log.Fatalln(err)
	} else if err = sortErrFile.Close(); err != nil {
		log.Fatalln(err)
	}
	// remove error file
	if err = os.Remove(sortErrFn); err != nil {
		log.Fatalln(err)
	}
	if err = os.Remove(sortInFn); err != nil {
		log.Fatalln(err)
	}
	if err = os.Rename(sortOutFn, sortInFn); err != nil {
		log.Fatalln(err)
	}
}
