package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var (
	logger *log.Logger
	data   map[string]string
)

func main() {
	// Setup logger
	logger = log.New(os.Stdout, "env-config: ", log.LstdFlags)

	output := os.Stdout
	var err error

	// If there are not parameters used when calling the program report an error and produce the usage message
	if len(os.Args) == 1 {
		logger.Fatal("Template name is missing\nUsage: config <template_file_name> [<output_file_name>]")
	}

	// If the second parameter was produced when calling program will write the output into the file specified as a
	// second parameter
	if len(os.Args) > 2 {
		output, err = os.Create(os.Args[2])
		if err != nil {
			logger.Fatalf("Cannot open an output file: '%s'", os.Args[2])
		}

		defer output.Close()
	}

	// Prepare the data map from the environment variables
	data = make(map[string]string)

	env := os.Environ()
	for _, e := range env {
		pair := strings.SplitN(e, "=", 2)
		if pair[0] != "" {
			data[pair[0]] = pair[1]
		}
	}

	// Create a template from the file specified as a first parameter
	testTemplate, err := template.New(filepath.Base(os.Args[1])).Funcs(template.FuncMap{
		"env":         vars,
		"withDefault": varsWithDefault,
	}).ParseFiles(os.Args[1])
	if err != nil {
		logger.Fatal(err)
	}

	// Render the output from the template and the data map
	err = testTemplate.Execute(output, data)
	if err != nil {
		logger.Fatal(err)
	}

}

// this function extracts the variable value from the environment using the key
// If variable is absent it will return an empty string
func vars(data map[string]string, key string) string {
	return data[key]
}

// this function extracts the variable value from the environment using the key
// If variable is absent it will return the defaultValue instead
func varsWithDefault(data map[string]string, key string, defaultValue string) string {
	d := data[key]
	if d == "" {
		d = defaultValue
	}
	return d
}
