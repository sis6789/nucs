package config

import (
	_ "embed"
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/mitchellh/go-homedir"
)

type jobConfig struct {
	BlastDb                      string `xml:"blast_db"`
	BlastDbAbs                   string `xml:"blast_db_abs"`
	BlastTaskCount               int    `xml:"blast_task_count"`
	BlastTaskCpuCount            int    `xml:"blast_task_cpu_count"`
	ChromosomeDataDir            string `xml:"chromosome_data_dir"`
	ChromosomeDataDirAbs         string `xml:"chromosome_data_dir_abs"`
	FastqDir                     string `xml:"fastq_dir"`
	FastqDirAbs                  string `xml:"fastq_dir_abs"`
	FastqFilePattern             string `xml:"fastq_file_pattern"`
	FastqQueryExamine            string `xml:"fastq_query_examine"`
	FastqQueryTerminator         string `xml:"fastq_query_terminator"`
	FastqQueryTerminatorLength   int    `xml:"fastq_query_terminator_length"`
	FastqQueryTerminatorMismatch int    `xml:"fastq_query_terminator_mismatch"`
	JobTitle                     string `xml:"job_title"`
	JobTitleRun                  string `xml:"job_title_run"`
	LogDir                       string `xml:"log_dir"`
	LogDirAbs                    string `xml:"log_dir_abs"`
	MinimumQueryLength           int    `xml:"minimum_query_length"`
	MongodbAccess                string `xml:"mongodb_access"`
	RootDir                      string `xml:"root_dir"`
	RootDirAbs                   string `xml:"root_dir_abs"`
	RunName                      string `xml:"run_name"`
	SaveDir                      string `xml:"save_dir"`
	SaveDirAbs                   string `xml:"save_dir_abs"`
	Shift                        string `xml:"shift"`
	ShiftDecode                  string `xml:"shift_decode"`
	TempDir                      string `xml:"temp_dir"`
	TempDirAbs                   string `xml:"temp_dir_abs"`
	TestLimit                    int    `xml:"test_limit"`
}

var obj jobConfig

//go:embed linux.xml
var linuxXML []byte

//go:embed windows.xml
var windowsXML []byte

var separator string
var readConfig = false

func init() {
	readConfig = false
	if runtime.GOOS == "windows" {
		if err := xml.Unmarshal(windowsXML, &obj); err != nil {
			log.Fatal(err)
		}
		separator = `\`
	} else {
		if err := xml.Unmarshal(linuxXML, &obj); err != nil {
			log.Fatal(err)
		}
		separator = `/`
	}
	obj.RunName = time.Now().Format("20060102-150405.999999")
	obj.FastqQueryTerminatorLength = len(obj.FastqQueryTerminator)
	obj.MongodbAccess = "mongodb://localhost:27017"
}

func n(p string) string {
	var expand string
	expand, err := homedir.Expand(p)
	if err != nil {
		log.Fatalln(err)
	}
	abs, err := filepath.Abs(expand)
	if err != nil {
		log.Fatalln(err)
	}
	return abs
}

func checkBlank() {
	if obj.BlastTaskCount == 0 {
		obj.BlastTaskCount = 10
	}
	if obj.BlastTaskCpuCount == 0 {
		obj.BlastTaskCpuCount = 1
	}
	if obj.BlastDb == "" && readConfig {
		log.Fatalln("no BlastDb")
	}
	if obj.ChromosomeDataDir == "" && readConfig {
		log.Fatalln("no ChromosomeDataDir")
	}
	if obj.FastqDir == "" && readConfig {
		log.Fatalln("no FastqDir")
	}
	if obj.JobTitle == "" {
		log.Fatalln("no JobTitle")
	}
	if obj.MongodbAccess == "" {
		log.Fatalln("no MongodbAccess")
	}
	if obj.RootDir == "" {
		log.Fatalln("no RootDir")
	}
}

func setAbs() {
	var err error

	if obj.RunName == "" {
		obj.RunName = time.Now().Format("20060102-150405.999999")
	}
	obj.JobTitleRun = obj.JobTitle + separator + obj.RunName
	if obj.BlastDbAbs, err = filepath.Abs(n(obj.BlastDb)); err != nil {
		log.Fatalln(err)
	}
	if obj.ChromosomeDataDirAbs, err = filepath.Abs(n(obj.ChromosomeDataDir)); err != nil {
		log.Fatalln(err)
	}
	if obj.FastqDirAbs, err = filepath.Abs(n(obj.FastqDir)); err != nil {
		log.Fatalln(err)
	}
	if obj.RootDirAbs, err = filepath.Abs(n(obj.RootDir)); err != nil {
		log.Fatalln(err)
	}
	if obj.LogDir == "" {
		obj.LogDir = obj.RootDir + separator + obj.JobTitleRun + separator + "log"
	}
	if obj.LogDirAbs, err = filepath.Abs(n(obj.LogDir)); err != nil {
		log.Fatalln(err)
	}
	if obj.SaveDir == "" {
		obj.SaveDir = obj.RootDir + separator + obj.JobTitleRun + separator + "save"
	}
	if obj.SaveDirAbs, err = filepath.Abs(n(obj.SaveDir)); err != nil {
		log.Fatalln(err)
	}
	if obj.TempDir == "" {
		obj.TempDir = obj.RootDir + separator + obj.JobTitleRun + separator + "temp"
	}
	if obj.TempDirAbs, err = filepath.Abs(n(obj.TempDir)); err != nil {
		log.Fatalln(err)
	}
}

func makeDirectory() {
	if !readConfig {
		return
	}
	if err := os.MkdirAll(obj.RootDirAbs, 0777); err != nil {
		log.Fatalln(err)
	}
	if err := os.MkdirAll(obj.LogDirAbs, 0777); err != nil {
		log.Fatalln(err)
	}
	if err := os.MkdirAll(obj.SaveDirAbs, 0777); err != nil {
		log.Fatalln(err)
	}
	if err := os.MkdirAll(obj.TempDirAbs, 0777); err != nil {
		log.Fatalln(err)
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ReadConfig(configFileName string) {
	readConfig = true
	configFile, err := os.Open(configFileName)
	checkErr(err)
	doc, err := ioutil.ReadAll(configFile)
	checkErr(err)
	_ = configFile.Close()
	err = xml.Unmarshal(doc, &obj)
	checkErr(err)

	checkBlank()
	setAbs()
	makeDirectory()

	obj.FastqQueryTerminatorLength = len(obj.FastqQueryTerminator)

	readConfig = false
}

func LogReport() {
	log.Println("BlastDb=", obj.BlastDb)
	log.Println("BlastDbAbs=", obj.BlastDbAbs)
	log.Println("BlastTaskCount=", obj.BlastTaskCount)
	log.Println("BlastTaskCpuCount=", obj.BlastTaskCpuCount)
	log.Println("ChromosomeDataDir=", obj.ChromosomeDataDir)
	log.Println("ChromosomeDataDirAbs=", obj.ChromosomeDataDirAbs)
	log.Println("FastqDir=", obj.FastqDir)
	log.Println("FastqDirAbs=", obj.FastqDirAbs)
	log.Println("FastqFilePattern=", obj.FastqFilePattern)
	log.Println("FastqQueryExamine=", obj.FastqQueryExamine)
	log.Println("FastqQueryTerminator=", obj.FastqQueryTerminator)
	log.Println("FastqQueryTerminatorLength=", obj.FastqQueryTerminatorLength)
	log.Println("FastqQueryTerminatorMismatch=", obj.FastqQueryTerminatorMismatch)
	log.Println("JobTitle=", obj.JobTitle)
	log.Println("LogDir=", obj.LogDir)
	log.Println("LogDirAbs=", obj.LogDirAbs)
	log.Println("MinimumQueryLength=", obj.MinimumQueryLength)
	log.Println("MongodbAccess=", obj.MongodbAccess)
	log.Println("RootDir=", obj.RootDir)
	log.Println("RootDirAbs=", obj.RootDirAbs)
	log.Println("RunName=", obj.RunName)
	log.Println("SaveDir=", obj.SaveDir)
	log.Println("SaveDirAbs=", obj.SaveDirAbs)
	log.Println("Shift=", obj.Shift)
	log.Println("ShiftDecode=", obj.ShiftDecode)
	log.Println("TempDir=", obj.TempDir)
	log.Println("TempDirAbs=", obj.TempDirAbs)
	log.Println("TestLimit=", obj.TestLimit)
}

func BlastDb() string                   { return obj.BlastDb }
func BlastDbAbs() string                { return obj.BlastDbAbs }
func BlastTaskCount() int               { return obj.BlastTaskCount }
func BlastTaskCpuCount() int            { return obj.BlastTaskCpuCount }
func ChromosomeDataDir() string         { return obj.ChromosomeDataDir }
func ChromosomeDataDirAbs() string      { return obj.ChromosomeDataDirAbs }
func FastqDir() string                  { return obj.FastqDir }
func FastqDirAbs() string               { return obj.FastqDirAbs }
func FastqFilePattern() string          { return obj.FastqFilePattern }
func FastqQueryExamine() string         { return obj.FastqQueryExamine }
func FastqQueryTerminator() string      { return obj.FastqQueryTerminator }
func FastqQueryTerminatorLength() int   { return obj.FastqQueryTerminatorLength }
func FastqQueryTerminatorMismatch() int { return obj.FastqQueryTerminatorMismatch }
func JobTitle() string                  { return obj.JobTitle }
func LogDir() string                    { return obj.LogDir }
func LogDirAbs() string                 { return obj.LogDirAbs }
func MinimumQueryLength() int           { return obj.MinimumQueryLength }
func MongodbAccess() string             { return obj.MongodbAccess }
func RootDir() string                   { return obj.RootDir }
func RootDirAbs() string                { return obj.RootDirAbs }
func RunName() string                   { return obj.RunName }
func SaveDir() string                   { return obj.SaveDir }
func SaveDirAbs() string                { return obj.SaveDirAbs }
func Shift() string                     { return obj.Shift }
func ShiftDecode() string               { return obj.ShiftDecode }
func TempDir() string                   { return obj.TempDir }
func TempDirAbs() string                { return obj.TempDirAbs }
func TestLimit() int                    { return obj.TestLimit }

func SaveName(fn string) string {
	abs, err := filepath.Abs(obj.SaveDirAbs + separator + fn)
	if err != nil {
		log.Fatalln(err)
	}
	return abs
}

func TempName(fn string) string {
	abs, err := filepath.Abs(obj.TempDirAbs + separator + fn)
	if err != nil {
		log.Fatalln(err)
	}
	return abs
}

func LogName(fn string) string {
	abs, err := filepath.Abs(obj.LogDirAbs + separator + fn)
	if err != nil {
		log.Fatalln(err)
	}
	return abs
}

func SaveNew(fn string) *os.File {
	fullPath := obj.SaveDirAbs + separator + fn
	file, err := os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatalln(err)
	}
	return file
}

func TempNew(fn string) *os.File {
	fullPath := obj.TempDirAbs + separator + fn
	file, err := os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatalln(err)
	}
	return file
}

func LogNew(fn string) *os.File {
	fullPath := obj.LogDirAbs + separator + fn
	file, err := os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatalln(err)
	}
	return file
}

func SaveAppend(fn string) *os.File {
	fullPath := obj.SaveDirAbs + separator + fn
	file, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalln(err)
	}
	return file
}

func TempAppend(fn string) *os.File {
	fullPath := obj.TempDirAbs + separator + fn
	file, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalln(err)
	}
	return file
}

func LogAppend(fn string) *os.File {
	fullPath := obj.LogDirAbs + separator + fn
	file, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalln(err)
	}
	return file
}

func SaveOpen(fn string) *os.File {
	fullPath := obj.SaveDirAbs + separator + fn
	file, err := os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalln(err)
	}
	return file
}

func TempOpen(fn string) *os.File {
	fullPath := obj.TempDirAbs + separator + fn
	file, err := os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalln(err)
	}
	return file
}

func LogOpen(fn string) *os.File {
	fullPath := obj.LogDirAbs + separator + fn
	file, err := os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalln(err)
	}
	return file
}
