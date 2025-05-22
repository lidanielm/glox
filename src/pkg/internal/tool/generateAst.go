package main

import (
	"fmt"
	"strings"
	"os"
)

func main() {
	outputDir := "out"
	baseName := "Expr"
	types := []string{"Binary : Expr left, Token operator, Expr right"}
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

	fmt.Fprintf(file, "package " + outputDir + "\n\n")
	fmt.Fprintf(file, "type " + baseName + " interface {\n")

	fmt.Fprintf(file, "}\n\n")

	for _, typ := range types {
		className := strings.TrimSpace(strings.Split(typ, ":")[0])
		fields := strings.Split(strings.TrimSpace(strings.Split(typ, ":")[1]), ", ")
		fmt.Fprintf(file, "type " + className + " struct {\n")
	
		for _, field := range fields {
			name, ftype := strings.Split(field, " ")[0], strings.Split(field, " ")[1]
			fmt.Fprintf(file, "\t%v %v\n", name, ftype)
		}

		fmt.Fprintf(file, "}\n")
		
		fmt.Fprintln(file, "")

		fmt.Fprintf(file, "func New%v(%v) *%v {\n", baseName, strings.TrimSpace(strings.Split(typ, ":")[1]), className)		
		fmt.Fprintf(file, "\treturn &%v{", className)
		for index, field := range fields {
			name, ftype := strings.Split(field, " ")[0], strings.Split(field, " ")[1]
			fmt.Fprintf(file, "%v: %v", name, ftype)
			if index < len(fields) - 1 {
				fmt.Fprint(file, ", ")
			}
		}
		fmt.Fprintf(file, "}\n")
		fmt.Fprintf(file, "}\n")
	}

	return nil
}
