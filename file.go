package main

import (
	"bufio"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func WriteFstabFileContentToTempFile(entries []*FstabLine, dst string) error {
	file, err := os.OpenFile(dst, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)

	defer file.Close()
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(file)
	for _, cnf := range entries {
		_, err := writer.WriteString(fmt.Sprintf("%s\n", cnf.GenerateFstabEntryString()))
		if err != nil {
			log.Fatalf("failed creating file: %s", err)
		}
	}
	writer.Flush()

	return nil
}

func ReadConfigFromXmlFile(path string) ([]*Config, error) {
	// read yml file
	ymlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})

	err = yaml.Unmarshal(ymlFile, &data)

	if err != nil {
		return nil, err
	}

	// get data and put into configs
	content := data["fstab"]
	if content == nil {
		return nil, errors.New("fstab key not found")
	}
	fconfig, ok := content.(map[string]interface{})
	if !ok {
		return nil, errors.New("can't parse fstab value")
	}

	configs, err := NewConfigs(fconfig)
	if err != nil {
		return nil, err
	}
	return configs, nil
}

func CopyFile(src, dst string) error {

	// check if file exists
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	// check if file is regular file
	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	//open source file
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	//open dst file
	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	//copy file
	_, err = io.Copy(destination, source)
	return err
}
