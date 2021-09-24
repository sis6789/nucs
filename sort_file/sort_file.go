package sort_file

import (
	"github.com/google/uuid"
	"github.com/sis6789/nucs/caller"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func SortFile(fileName string, keyList ...string) {
	if len(keyList) == 0 {
		return
	}
	var err error
	var sortInFile, sortOutFile, sortErrFile *os.File
	folder := filepath.Dir(fileName)
	sortInFn := fileName
	sortOutFn := filepath.Join(folder, "so_"+uuid.NewString()+".txt")
	sortErrFn := filepath.Join(folder, "se_"+uuid.NewString()+".txt")
	if sortInFile, err = os.Open(sortInFn); err != nil {
		log.Fatalln(caller.Caller(), err)
	}
	if sortOutFile, err = os.Create(sortOutFn); err != nil {
		log.Fatalln(caller.Caller(), err)
	}
	if sortErrFile, err = os.Create(sortErrFn); err != nil {
		log.Fatalln(caller.Caller(), err)
	}

	var cmdPath string
	var sortOrder []string
	if runtime.GOOS == "windows" {
		cmdPath, err = exec.LookPath("wsl")
		if err != nil {
			log.Fatalln(caller.Caller(), err)
			return
		}
		sortOrder = []string{cmdPath, "sort"}
	} else {
		cmdPath, err = exec.LookPath("sort")
		if err != nil {
			log.Fatalln(caller.Caller(), err)
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
		log.Println(startErr)
		_ = sortInFile.Close()
		_ = sortOutFile.Close()
		_ = sortErrFile.Close()
		log.Fatalln(caller.Caller(), fileName)
	}
	// wait end of sort
	if endErr := cmd.Wait(); endErr != nil {
		log.Println(endErr)
		_ = sortInFile.Close()
		_ = sortOutFile.Close()
		_ = sortErrFile.Close()
		log.Fatalln(caller.Caller(), fileName)
	}
	// close
	if err = sortInFile.Close(); err != nil {
		log.Fatalln(caller.Caller(), err)
	}
	if err = sortOutFile.Sync(); err != nil {
		log.Fatalln(caller.Caller(), err)
	} else if err = sortOutFile.Close(); err != nil {
		log.Fatalln(caller.Caller(), err)
	}
	if err = sortErrFile.Sync(); err != nil {
		log.Fatalln(caller.Caller(), err)
	} else if err = sortErrFile.Close(); err != nil {
		log.Fatalln(caller.Caller(), err)
	}
	// remove error file
	if err = os.Remove(sortErrFn); err != nil {
		log.Fatalln(caller.Caller(), err)
	}
	if err = os.Remove(sortInFn); err != nil {
		log.Fatalln(caller.Caller(), err)
	}
	if err = os.Rename(sortOutFn, sortInFn); err != nil {
		log.Fatalln(caller.Caller(), err)
	}
}
