package main

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/micro/micro/v3/service/build"
	"github.com/micro/micro/v3/service/build/client"
)

const testDir = "test"

func main() {
	client := client.NewBuilder()

	// Create a buffer to write our archive to and the zip archive.
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	// walkFn zips each file in the directory
	walkFn := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		relpath, err := filepath.Rel(testDir, path)
		if err != nil {
			return err
		}

		f, err := w.Create(relpath)
		if err != nil {
			return err
		}

		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		_, err = f.Write(bytes)
		return err
	}

	// Add some files to the archive.
	if err := filepath.Walk("./test", walkFn); err != nil {
		w.Close()
		log.Fatalf("Error archiving test directory: %v", err)
	}

	// Close the zip
	if err := w.Close(); err != nil {
		log.Fatalf("Error closing zip writer: %v", err)
	}

	// build the source using the client
	res, err := client.Build(buf, build.Archive("zip"))
	if err != nil {
		log.Fatalf("Error building source: %v", err)
	}

	// get the current directory
	path, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting working directory: %v", err)
	}

	// if a previous build exists, delete if
	resultPath := filepath.Join(path, "output")
	if _, err := os.Stat(resultPath); err == nil {
		os.Remove(resultPath)
	}

	// create a file to write the output to
	file, err := os.Create(resultPath)
	if err != nil {
		log.Fatalf("Error creating result file: %v", err)
	}
	defer file.Close()

	// copy the result to the file
	if _, err := io.Copy(file, res); err != nil {
		log.Fatalf("Error copying result to file: %v", err)
	}

	// set the file permissons so the binary is executable
	if err := os.Chmod(file.Name(), 0333); err != nil {
		log.Fatalf("Error setting file permissions: %v", err)
	}

	// print the result
	log.Printf("Result written to ./output")
}
