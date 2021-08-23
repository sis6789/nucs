//go:build windows

package json_config

import (
	"path/filepath"
)

// K0 : For Keyomics
func K0() {

	jsonConfig["blast_db"] = "D:/Utility/NCBI/db/hg38"
	jsonConfig["chromosome_data_dir"] = "D:/HumanGenome"
	jsonConfig["fastq_dir"] = "D:/keyomics_test/fastq"
	jsonConfig["root_dir"] = "D:/keyomics_test"

	jsonConfig["blast_task_count"] = 200
	jsonConfig["blast_task_cpu_count"] = 4
	jsonConfig["fastq_file_pattern"] = "^.+fq$"
	jsonConfig["fastq_query_examine"] = "^(\\w{7,12})(CAGCTG([AC])CGTCAGTCT)(\\w+)$"
	jsonConfig["fastq_query_terminator"] = "GATCGGAAGAGCACACGTCTGAACTCCAGTCAC"
	jsonConfig["fastq_query_terminator_mismatch"] = 3
	jsonConfig["job_title"] = "test-bbb"
	jsonConfig["l1prm"] = "GATCGGAAGAGCACACGTCTGAACTCCAGTCAC"
	jsonConfig["minimum_query_length"] = 26
	jsonConfig["mongodb_access"] = "mongodb://localhost:27017"
	jsonConfig["run_name"] = "run01"
	jsonConfig["shift"] = ",T,GT,AGT,CAGT"
	jsonConfig["shift_decode"] = "----,T---,GT--,AGT-,CAGT"
	jsonConfig["test_limit"] = 10000
	jsonConfig["top_bottom_minimum_count"] = 2
	jsonConfig["top_bottom_minimum_ratio"] = 0.2
	jsonConfig["top_bottom_minimum_sum"] = 4

	K1()
}

func K1() {

	jsonConfig["fastq_query_terminator_length"] = len(jsonConfig["fastq_query_terminator"].(string))
	jsonConfig["work_dir"] = filepath.Join(jsonConfig["root_dir"].(string),
		jsonConfig["job_title"].(string), jsonConfig["run_name"].(string))
	jsonConfig["log_dir"] = filepath.Join(jsonConfig["work_dir"].(string), "log")
	jsonConfig["save_dir"] = filepath.Join(jsonConfig["work_dir"].(string), "save")
	jsonConfig["temp_dir"] = filepath.Join(jsonConfig["work_dir"].(string), "temp")
}
