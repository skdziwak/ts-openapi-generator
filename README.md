# TypeScript OpenAPI Generator

This project is a Go-based tool that generates TypeScript code from OpenAPI 3.0 specifications. It parses an OpenAPI YAML file and generates TypeScript interfaces and classes based on the schemas defined in the specification.

## Features

- Parses OpenAPI 3.0 YAML specifications
- Generates TypeScript interfaces and classes for each schema
- Supports basic data types (string, number, boolean, array, object)
- Handles optional properties
- Outputs generated code to specified folder

## Installation

Ensure you have Go 1.22.1 or later installed on your system. Then, clone this repository and build the project:

    git clone https://github.com/skdziwak/ts-openapi-generator.git
    cd ts-openapi-generator
    go build

## Usage

Run the generator with the following command:

    ./ts-openapi-generator -i <input-yaml-file> -o <output-folder> [-l <log-level>]

Arguments:
- `-i`: Path to the input OpenAPI YAML file (required)
- `-o`: Path to the output folder for generated code (required)
- `-l`: Log level (optional, default: info)

Example:

    ./ts-openapi-generator -i api-spec.yaml -o ./generated -l debug

## Generated Code Structure

For each schema in the OpenAPI specification, the generator creates a TypeScript file in the `models` subdirectory of the output folder. Each file contains:

1. An interface representing the schema properties
2. A class implementing the interface
3. A default export of the class

## Project Structure

- `main.go`: Entry point of the application
- `core.go`: Core interfaces and utilities
- `generator.go`: TypeScript code generation logic
- `openapi.go`: OpenAPI specification parsing and handling

## Key Components

### OpenAPISpec

The `OpenAPISpec` struct in `openapi.go` represents the parsed OpenAPI specification. It provides methods to access paths and schemas defined in the spec.

### CodeGenerator and CodeFileWriter

These interfaces, defined in `core.go`, are central to the code generation process:

- `CodeGenerator`: Defines the `HandleSchema` method for generating code from a schema
- `CodeFileWriter`: Provides a `Write` method for outputting generated code to files

### TsGenerator

The `TsGenerator` struct in `generator.go` implements the `CodeGenerator` interface. It contains the logic for generating TypeScript interfaces and classes from OpenAPI schemas.

## Dependencies

Main external dependencies:
- github.com/getkin/kin-openapi: For parsing OpenAPI specifications
- github.com/sirupsen/logrus: For logging

See `go.mod` for a complete list of dependencies.

## Contributing

Contributions are welcome! Please submit pull requests with any improvements or bug fixes.

## License

This project is licensed under the MIT License - see below for details:

MIT License

Copyright (c) 2024 Szymon Dziwak

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.