package config

import (
	"fmt"
)

// MongoReport - return config description for mongo DBMS storage
func MongoReport() []string {
	var result []string
	result = append(result, fmt.Sprint("BlastDb=", obj.BlastDb))
	result = append(result, fmt.Sprint("BlastDbAbs=", obj.BlastDbAbs))
	result = append(result, fmt.Sprint("BlastTaskCount=", obj.BlastTaskCount))
	result = append(result, fmt.Sprint("BlastTaskCpuCount=", obj.BlastTaskCpuCount))
	result = append(result, fmt.Sprint("ChromosomeDataDir=", obj.ChromosomeDataDir))
	result = append(result, fmt.Sprint("ChromosomeDataDirAbs=", obj.ChromosomeDataDirAbs))
	result = append(result, fmt.Sprint("FastqDir=", obj.FastqDir))
	result = append(result, fmt.Sprint("FastqDirAbs=", obj.FastqDirAbs))
	result = append(result, fmt.Sprint("FastqFilePattern=", obj.FastqFilePattern))
	result = append(result, fmt.Sprint("FastqQueryExamine=", obj.FastqQueryExamine))
	result = append(result, fmt.Sprint("FastqQueryTerminator=", obj.FastqQueryTerminator))
	result = append(result, fmt.Sprint("FastqQueryTerminatorLength=", obj.FastqQueryTerminatorLength))
	result = append(result, fmt.Sprint("FastqQueryTerminatorMismatch=", obj.FastqQueryTerminatorMismatch))
	result = append(result, fmt.Sprint("JobTitle=", obj.JobTitle))
	result = append(result, fmt.Sprint("LogDir=", obj.LogDir))
	result = append(result, fmt.Sprint("LogDirAbs=", obj.LogDirAbs))
	result = append(result, fmt.Sprint("MinimumQueryLength=", obj.MinimumQueryLength))
	result = append(result, fmt.Sprint("MongodbAccess=", obj.MongodbAccess))
	result = append(result, fmt.Sprint("RootDir=", obj.RootDir))
	result = append(result, fmt.Sprint("RootDirAbs=", obj.RootDirAbs))
	result = append(result, fmt.Sprint("RunName=", obj.RunName))
	result = append(result, fmt.Sprint("SaveDir=", obj.SaveDir))
	result = append(result, fmt.Sprint("SaveDirAbs=", obj.SaveDirAbs))
	result = append(result, fmt.Sprint("Shift=", obj.Shift))
	result = append(result, fmt.Sprint("ShiftDecode=", obj.ShiftDecode))
	result = append(result, fmt.Sprint("TempDir=", obj.TempDir))
	result = append(result, fmt.Sprint("TempDirAbs=", obj.TempDirAbs))
	result = append(result, fmt.Sprint("TestLimit=", obj.TestLimit))
	return result
}
