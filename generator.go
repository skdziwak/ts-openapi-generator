package main

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/sirupsen/logrus"
)

type TsGenerator struct {
}

func NewTsGenerator() *TsGenerator {
	return &TsGenerator{}
}

func (g *TsGenerator) HandleSchema(name string, schema *openapi3.Schema, writer CodeFileWriter, usedRefs map[string]bool) {
	interfaceName := name + "Props"
	className := name

	interfaceContent := g.generateInterface(interfaceName, schema, usedRefs)
	classContent := g.generateClass(className, interfaceName, schema)

	imports := g.generateImports(usedRefs)
	content := fmt.Sprintf("%s\n\n%s\n\n%s\n\nexport default %s;\n", imports, interfaceContent, classContent, className)
	trimmedContent := strings.TrimSpace(content)
	writer.Write("models/"+name+".ts", trimmedContent)
}

func (g *TsGenerator) generateInterface(name string, schema *openapi3.Schema, usedRefs map[string]bool) string {
	var properties []string
	for propName, propSchema := range schema.Properties {
		ref := propSchema.Ref
		propType := g.getTypeScriptType(propSchema.Value, &ref, usedRefs)
		optional := !g.isRequired(propName, schema)
		property := fmt.Sprintf("    %s%s: %s;", propName, g.optionalSuffix(optional), propType)
		properties = append(properties, property)
	}

	return fmt.Sprintf("export interface %s {\n%s\n}", name, strings.Join(properties, "\n"))
}

func (g *TsGenerator) generateClass(className, interfaceName string, schema *openapi3.Schema) string {
	var properties, constructorParams, constructorAssignments []string

	for propName := range schema.Properties {

		properties = append(properties, fmt.Sprintf("    readonly %s: %s['%s'];", propName, interfaceName, propName))
		constructorParams = append(constructorParams, propName)
		constructorAssignments = append(constructorAssignments, fmt.Sprintf("        this.%s = %s;", propName, propName))
	}

	return fmt.Sprintf(`class %s implements %s {
%s

    constructor({%s}: %s) {
%s
    }
}`, className, interfaceName, strings.Join(properties, "\n"), strings.Join(constructorParams, ", "), interfaceName, strings.Join(constructorAssignments, "\n"))
}

func (g *TsGenerator) generateImports(usedRefs map[string]bool) string {
	var imports []string
	for ref := range usedRefs {
		lastPart := strings.Split(ref, "/")
		importName := lastPart[len(lastPart)-1] + "Props"
		imports = append(imports, fmt.Sprintf("import { %s } from \"./%s\";", importName, lastPart[len(lastPart)-1]))
	}
	return strings.Join(imports, "\n")
}

func (g *TsGenerator) getTypeScriptType(schema *openapi3.Schema, ref *string, usedRefs map[string]bool) string {
	switch {
	case schema.Type.Is(openapi3.TypeString):
		return "string"
	case schema.Type.Is(openapi3.TypeInteger), schema.Type.Is(openapi3.TypeNumber):
		return "number"
	case schema.Type.Is(openapi3.TypeBoolean):
		return "boolean"
	case schema.Type.Is(openapi3.TypeArray):
		if schema.Items != nil {
			itemType := g.getTypeScriptType(schema.Items.Value, &schema.Items.Ref, usedRefs)
			return fmt.Sprintf("%s[]", itemType)
		} else {
			logrus.Warnf("Schema items are nil for %s", schema.Title)
			return "any[]"
		}
	case schema.Type.Is(openapi3.TypeObject):
		if ref != nil && *ref != "" {
			lastPart := strings.Split(*ref, "/")
			refName := lastPart[len(lastPart)-1]
			usedRefs[refName] = true
			return refName + "Props"
		} else {
			logrus.Warnf("Schema title is empty for an object type")
			return "Record<string, any>"
		}
	case schema.Type.Is(openapi3.TypeNull):
		return "null"
	default:
		return "any"
	}
}

func (g *TsGenerator) isRequired(propName string, schema *openapi3.Schema) bool {
	for _, required := range schema.Required {
		if required == propName {
			return true
		}
	}
	return false
}

func (g *TsGenerator) optionalSuffix(optional bool) string {
	if optional {
		return "?"
	}
	return ""
}
