package main

import (
	"fmt"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
)

// OpenAPISpec represents the parsed OpenAPI specification
type OpenAPISpec struct {
	Spec *openapi3.T
}

// ParseOpenAPIFile parses the OpenAPI YAML file and returns an OpenAPISpec
func ParseOpenAPIFile(filePath string) (*OpenAPISpec, error) {
	// Read the YAML file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	// Parse the YAML into an OpenAPI 3.0 document
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(data)
	if err != nil {
		return nil, fmt.Errorf("error parsing OpenAPI spec: %v", err)
	}

	// Validate the document
	err = doc.Validate(loader.Context)
	if err != nil {
		return nil, fmt.Errorf("invalid OpenAPI spec: %v", err)
	}

	return &OpenAPISpec{
		Spec: doc,
	}, nil
}

// GetPaths returns all paths defined in the OpenAPI specification
func (o *OpenAPISpec) GetPaths() openapi3.Paths {
	return *o.Spec.Paths
}

// GetSchemas returns all schemas defined in the OpenAPI specification
func (o *OpenAPISpec) GetSchemas() map[string]*openapi3.SchemaRef {
	return o.Spec.Components.Schemas
}

// GetInfo returns the API info from the OpenAPI specification
func (o *OpenAPISpec) GetInfo() *openapi3.Info {
	return o.Spec.Info
}
