package main

import (
	"fmt"
	"strings"
	"os"
)

func main() {
	outputDir := "src/pkg/internal/ast"
	baseName := "expr"
	types := []string{"Binary : left Expr, operator Token, right Expr", "Grouping : expression Expr", "Literal : value Object", "Unary : operator Token, right Expr"}
	err := defineAst(outputDir, baseName, types)
	if err != nil {
		panic(err)
	}
}

func defineAst(outputDir string, baseName string, types []string) (err error) {
	path := outputDir + "/" + baseName + ".go"
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	dirs := strings.Split(outputDir, "/")
	fmt.Fprintf(file, "package " + dirs[len(dirs) - 1] + "\n\n")
	fmt.Fprintf(file, "type " + baseName + " interface {\n")
	fmt.Fprintf(file, "\tAccept(visitor Visitor[any]) any\n")
	fmt.Fprintf(file, "}\n\n")

	// Define classes
	for _, typ := range types {
		className := strings.TrimSpace(strings.Split(typ, ":")[0])
		classType := strings.TrimSpace(strings.Split(typ, ":")[1])
		fields := strings.Split(strings.TrimSpace(strings.Split(typ, ":")[1]), ", ")
		defineType(file, baseName, className, classType, fields)
	}

	return nil
}

func defineType(file *os.File, baseName string, className string, classType string, fields []string) {
	// Define struct type
	fmt.Fprintf(file, "type " + className + " struct {\n")

	for _, field := range fields {
		name, ftype := strings.Split(field, " ")[0], strings.Split(field, " ")[1]
		fmt.Fprintf(file, "\t%v %v\n", name, ftype)
	}

	fmt.Fprintf(file, "}\n")
	
	fmt.Fprintln(file, "")

	// Define constructor
	fmt.Fprintf(file, "func New%v(%v) *%v {\n", className, classType, className)		
	fmt.Fprintf(file, "\treturn &%v{", className)
	for index, field := range fields {
		name, ftype := strings.Split(field, " ")[0], strings.Split(field, " ")[1]
		fmt.Fprintf(file, "%v: %v", name, ftype)
		if index < len(fields) - 1 {
			fmt.Fprint(file, ", ")
		}
	}
	fmt.Fprintf(file, "}\n")
	fmt.Fprintf(file, "}\n\n")

	// Define visitor method
	fmt.Fprintf(file, "func (%c %v) Accept(visitor Visitor[any]) any {\n", strings.ToLower(className)[0], className)
	fmt.Fprintf(file, "\treturn visitor.visit%v%v(%c)\n", className, baseName, strings.ToLower(className)[0])
	fmt.Fprintf(file, "}\n\n")
}

func defineVisitor(file *os.File, baseName string, types []string) (err error) {
	fmt.Fprintf(file, "type Visitor[R any] interface {\n")
	
	for _, typ := range types {
		typeName := strings.TrimSpace(strings.Split(typ, ":")[0])
		fmt.Fprintf(file, "\tVisit%v%v(%v %v) R\n", typeName, baseName, strings.ToLower(baseName), typeName)
	}

	fmt.Fprintf(file, "}\n\n")
	return nil
}
