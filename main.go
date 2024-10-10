package main

import (
	"flag"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logrus.SetLevel(logrus.InfoLevel)
	// Parse command-line arguments
	inputYAML := flag.String("i", "", "Path to the input OpenAPI YAML file (required)")
	outputFolder := flag.String("o", "", "Path to the output folder for generated code (required)")
	logLevel := flag.String("l", "info", "Log level (optional, default: info)")
	flag.Parse()

	// Validate required arguments
	if *inputYAML == "" || *outputFolder == "" {
		flag.Usage()
		logrus.Fatal("Both input_yaml and output_folder are required")
	}

	// Set log level based on the provided argument
	level, err := logrus.ParseLevel(*logLevel)
	if err != nil {
		logrus.Warnf("Invalid log level '%s', defaulting to info", *logLevel)
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	logrus.Infof("Input YAML: %s", *inputYAML)
	logrus.Infof("Output folder: %s", *outputFolder)
	logrus.Infof("Log level: %s", level)

	spec, err := ParseOpenAPIFile(*inputYAML)
	if err != nil {
		logrus.Fatalf("Error parsing OpenAPI file: %v", err)
	}

	writer := NewCodeFileWriter(*outputFolder)

	generator := NewTsGenerator()
	Generate(generator, spec, writer)

}
