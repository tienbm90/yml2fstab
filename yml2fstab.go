package main

import (
	"flag"
	"log"
)

var (
	inFile  = flag.String("in", "input.yml", "Path to configuration file. Default is ./input.yml")
	outFile = flag.String("out", "/etc/fstab", "Path to output file. Default is /etc/fstab")
	tmpFile = flag.String("tmp-file", "/tmp/fstab.temp", "Path to temp file. Default is /tmp/fstab.temp")
)

func main() {

	//parse agrument
	flag.Parse()
	//read config from file

	configs, err := ReadConfigFromXmlFile(*inFile)

	if err != nil {
		log.Printf("Parser error: %s", err.Error())
		return
	}
	// create fstab entries
	entries := make([]*FstabLine, 0)

	for _, cnf := range configs {
		ent := NewFstabLineFromConfig(*cnf)
		entries = append(entries, ent)
	}

	// write fstab file
	err = WriteFstabFileContentToTempFile(entries, *tmpFile)
	if err != nil {
		log.Printf("Write file error: %s", err.Error())
		return
	}

	// copy file to /etc/fstab
	err = CopyFile(*tmpFile, *outFile)
	if err != nil {
		log.Printf(	"Copy file error: %s", err.Error())
	}
}
