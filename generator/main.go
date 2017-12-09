package main

import (
	"html/template"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/rancher/types/config"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"path/filepath"
	"fmt"
)

var underscoreRegexp = regexp.MustCompile(`([a-z])([A-Z])`)

func ToLowerCamelCase(input string) string {
	return (strings.ToLower(input[:1]) + input[1:])
}

func addUnderscore(input string) string {
	return strings.ToLower(underscoreRegexp.ReplaceAllString(input, `${1}_${2}`))
}

func capitalize(s string) string {
	if len(s) <= 1 {
		return strings.ToUpper(s)
	}

	return strings.ToUpper(s[:1]) + s[1:]
}

func main() {
	if err := filepath.Walk("../..", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.Contains(path, "vendor") {
			return filepath.SkipDir
		}

		if strings.HasPrefix(info.Name(), "zz_generated") {
			fmt.Println("Removing", path)
			if err := os.Remove(path); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		logrus.Fatal(err)
	}
	funcMap := template.FuncMap{
		"toLowerCamelCase":  ToLowerCamelCase,
		"toLowerUnderscore": addUnderscore,
		"capitalize":        capitalize,
		"upper":             strings.ToUpper,
	}
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		logrus.Fatal(err)
	}

	management, err := config.NewManagementContext(*kubeConfig)
	if err != nil {
		logrus.Fatal(err)
	}
	schemaClient := management.Management.DynamicSchemas("")
	schemas, err := schemaClient.List(metav1.ListOptions{})
	if err != nil {
		logrus.Fatal(err)
	}
	for _, schema := range schemas.Items {
		wd, err := os.Getwd()
		if err != nil {
			logrus.Fatal(err)
		}
		output, err := os.Create(path.Join(wd, "generator", strings.ToLower("zz_generated_"+addUnderscore(schema.Name))+".go"))
		if err != nil {
			logrus.Fatal(err)
		}
		data := map[string]interface{}{
			"schema": getTypeMap(schema.Spec),
			"Name":   capitalize(strings.TrimSuffix(schema.Name, "config")),
		}
		typeTemplate, err := template.New("type.template").Funcs(funcMap).ParseFiles(path.Join(wd, "generator", "type.template"))
		if err != nil {
			logrus.Fatal(err)
		}
		if err := typeTemplate.Execute(output, data); err != nil {
			logrus.Fatal(err)
		}
	}
}

func getTypeMap(schema v3.DynamicSchemaSpec) map[string]string {
	result := map[string]string{}
	for name, field := range schema.ResourceFields {
		if name == "id" {
			continue
		}

		fieldName := capitalize(name)

		if strings.HasPrefix(field.Type, "reference") || strings.HasPrefix(field.Type, "date") || strings.HasPrefix(field.Type, "enum") {
			result[fieldName] = "string"
		} else if strings.HasPrefix(field.Type, "array[reference[") {
			result[fieldName] = "[]string"
		} else if strings.HasPrefix(field.Type, "array") {
			switch field.Type {
			case "array[reference]":
				fallthrough
			case "array[date]":
				fallthrough
			case "array[enum]":
				fallthrough
			case "array[string]":
				result[fieldName] = "[]string"
			case "array[int]":
				result[fieldName] = "[]int64"
			case "array[float64]":
				result[fieldName] = "[]float64"
			case "array[json]":
				result[fieldName] = "[]interface{}"
			default:
				s := strings.TrimLeft(field.Type, "array[")
				s = strings.TrimRight(s, "]")
				result[fieldName] = "[]" + capitalize(s)
			}
		} else if strings.HasPrefix(field.Type, "map") {
			result[fieldName] = "map[string]interface{}"
		} else if strings.HasPrefix(field.Type, "json") {
			result[fieldName] = "interface{}"
		} else if strings.HasPrefix(field.Type, "boolean") {
			result[fieldName] = "bool"
		} else if strings.HasPrefix(field.Type, "extensionPoint") {
			result[fieldName] = "interface{}"
		} else if strings.HasPrefix(field.Type, "float") {
			result[fieldName] = "float64"
		} else if strings.HasPrefix(field.Type, "int") {
			result[fieldName] = "int64"
		} else if strings.HasPrefix(field.Type, "password") {
			result[fieldName] = "string"
		} else if strings.HasPrefix(field.Type, "bool") {
			result[fieldName] = "bool"
		} else if field.Nullable {
			result[fieldName] = "*" + capitalize(field.Type)
		} else {
			result[fieldName] = capitalize(field.Type)
		}
	}

	return result
}
