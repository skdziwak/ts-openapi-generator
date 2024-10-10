package main

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/sirupsen/logrus"
)

type TsGenerator struct{}

func (g *TsGenerator) HandleSchema(name string, schema *openapi3.Schema, writer CodeFileWriter) {
	interfaceName := name + "Props"
	className := name

	interfaceContent := g.generateInterface(interfaceName, schema)
	classContent := g.generateClass(className, interfaceName, schema)

	content := fmt.Sprintf("%s\n\n%s\n\nexport default %s;\n", interfaceContent, classContent, className)
	writer.Write("models/"+name+".ts", content)
}

func (g *TsGenerator) generateInterface(name string, schema *openapi3.Schema) string {
	var properties []string
	for propName, propSchema := range schema.Properties {
		ref := propSchema.Ref
		propType := g.getTypeScriptType(propSchema.Value, &ref)
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

func (g *TsGenerator) getTypeScriptType(schema *openapi3.Schema, ref *string) string {
	switch {
	case schema.Type.Is(openapi3.TypeString):
		return "string"
	case schema.Type.Is(openapi3.TypeInteger), schema.Type.Is(openapi3.TypeNumber):
		return "number"
	case schema.Type.Is(openapi3.TypeBoolean):
		return "boolean"
	case schema.Type.Is(openapi3.TypeArray):
		if schema.Items != nil {
			itemType := g.getTypeScriptType(schema.Items.Value, &schema.Items.Ref)
			return fmt.Sprintf("%s[]", itemType)
		} else {
			logrus.Warnf("Schema items are nil for %s", schema.Title)
			return "any[]"
		}
	case schema.Type.Is(openapi3.TypeObject):
		if ref != nil && *ref != "" {
			lastPart := strings.Split(*ref, "/")
			return lastPart[len(lastPart)-1] + "Props"
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
