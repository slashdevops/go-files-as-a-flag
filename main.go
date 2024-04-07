package main

import (
	"flag"
	"fmt"
	"os"
)

// FileVar is a custom flag type for files
// This should implement the Value interface of the flag package
// Reference: https://pkg.go.dev/gg-scm.io/tool/internal/flag#FlagSet.Var
type FileVar struct {
	*os.File
}

// String presents the current value as a string.
func (f *FileVar) String() string {
	if f.File == nil {
		return ""
	}

	return f.File.Name()
}

// Set is called once, in command line order, for each flag present.
func (f *FileVar) Set(value string) error {
	file, err := os.OpenFile(value, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}

	f.File = file
	return nil
}

// Get returns the contents of the Value.
func (f *FileVar) Get() interface{} {
	return f.File
}

// IsBoolFlag returns true if the flag is a boolean flag
func (f *FileVar) IsBoolFlag() bool {
	return false
}

func main() {
	// Create a new flag set
	fs := flag.NewFlagSet("File as Flag CLI", flag.ExitOnError)

	// Add a flag to get some content
	var content string
	fs.StringVar(&content, "file.content", "", "content to write to the file")

	// Add a custom file flag
	file := &FileVar{os.Stdout}
	defer file.Close()
	fs.Var(file, "output.file", "output file")

	// Parse the command line arguments
	fs.Parse(os.Args[1:])

	// Check if the content is empty (required)
	if content == "" {
		fs.PrintDefaults()
		fmt.Println("error: '-file.content' is required")
		os.Exit(1)
	}

	// Write the content to the file
	file.Write([]byte(content))
}
