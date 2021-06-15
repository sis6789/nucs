package config

import (
	"fmt"
	"runtime"
	"testing"
)

func TestAny(t *testing.T) {
	if runtime.GOOS == "windows" {
		ReadConfig(`D:\keyomics\Projects\evaluateNucs\pctest.xml`)
	} else {
		ReadConfig("linux.xml")
	}

	mRec := MongoReport()
	fmt.Println(mRec)

	fmt.Println("BlastDb=", BlastDb())
	fmt.Println("BlastDbAbs=", BlastDbAbs())
	fmt.Println("BlastTaskCount=", BlastTaskCount())
	fmt.Println("BlastTaskCpuCount=", BlastTaskCpuCount())
	fmt.Println("ChromosomeDataDir=", ChromosomeDataDir())
	fmt.Println("ChromosomeDataDirAbs=", ChromosomeDataDirAbs())
	fmt.Println("FastqDir=", FastqDir())
	fmt.Println("FastqDirAbs=", FastqDirAbs())
	fmt.Println("FastqFilePattern=", FastqFilePattern())
	fmt.Println("FastqQueryExamine=", FastqQueryExamine())
	fmt.Println("FastqQueryTerminator=", FastqQueryTerminator())
	fmt.Println("FastqQueryTerminatorLength=", FastqQueryTerminatorLength())
	fmt.Println("FastqQueryTerminatorMismatch=", FastqQueryTerminatorMismatch())
	fmt.Println("JobTitle=", JobTitle())
	fmt.Println("JobTitleRun=", JobTitleRun())
	fmt.Println("L1prm=", L1prm())
	fmt.Println("LogDir=", LogDir())
	fmt.Println("LogDirAbs=", LogDirAbs())
	fmt.Println("MinimumQueryLength=", MinimumQueryLength())
	fmt.Println("MongodbAccess=", MongodbAccess())
	fmt.Println("RootDir=", RootDir())
	fmt.Println("RootDirAbs=", RootDirAbs())
	fmt.Println("RunName=", RunName())
	fmt.Println("SaveDir=", SaveDir())
	fmt.Println("SaveDirAbs=", SaveDirAbs())
	fmt.Println("Shift=", Shift())
	fmt.Println("ShiftDecode=", ShiftDecode())
	fmt.Println("TempDir=", TempDir())
	fmt.Println("TempDirAbs=", TempDirAbs())
	fmt.Println("TestLimit=", TestLimit())
	fmt.Println("TopBottomMinimumCount=", TopBottomMinimumCount())
	fmt.Println("TopBottomMinimumRatio=", TopBottomMinimumRatio())
	fmt.Println("TopBottomMinimumSum=", TopBottomMinimumSum())

	LogReport()

	LogAppend("append.fastq")
	LogNew("new.fastq")
	LogOpen("log.fastq")

	TempAppend("append.fastq")
	TempNew("new.fastq")
	TempOpen("log.fastq")

	SaveAppend("append.fastq")
	SaveNew("new.fastq")
	SaveOpen("log.fastq")

	fmt.Println(SaveName("하나님.만세"))
	fmt.Println(TempName("하나님.만세"))
	fmt.Println(LogName("하나님.만세"))
}
