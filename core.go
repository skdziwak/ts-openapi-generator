package main

import (
	"os"
	"path/filepath"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/sirupsen/logrus"
)

func panicOnError(err error) {
	if err != nil {
		logrus.Error(err)
		panic(err)
	}
}

type CodeFileWriter interface {
	Write(path string, content string) error
}

type CodeFileWriterImpl struct {
	folder string
}

func NewCodeFileWriter(folder string) CodeFileWriter {
	return &CodeFileWriterImpl{folder: folder}
}

func (w *CodeFileWriterImpl) Write(path string, content string) error {
	filePath := filepath.Join(w.folder, path)
	if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		return err
	}
	err := os.WriteFile(filePath, []byte(content), 0o644)
	if err != nil {
		return err
	}
	return nil
}

type CodeGenerator interface {
	HandleSchema(name string, schema *openapi3.Schema, writer CodeFileWriter, usedRefs map[string]bool)
}

func Generate(generator CodeGenerator, spec *OpenAPISpec, writer CodeFileWriter) {
	log := logrus.WithField("component", "generator")

	log.Info("Starting code generation")
	logrus.WithFields(logrus.Fields{
		"info":    spec.GetInfo(),
		"schemas": spec.GetSchemas(),
	}).Debug("Full OpenAPI specification")

	for name, schemaRef := range spec.GetSchemas() {
		log.WithField("schema", name).Debug("Handling schema")
		generator.HandleSchema(name, schemaRef.Value, writer, make(map[string]bool))
	}

	log.Info("Code generation complete")
}
