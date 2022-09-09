package main

import (
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	_ "embed"

	"golang.org/x/tools/go/packages"
)

//go:embed generated_typegen.template
var typegenTemplate string

func loadAllImportedTypes(path string) (map[string]string, []string) {
	imports := map[string]string{}
	types := map[string]struct{}{}

	pkgs, err := packages.Load(
		&packages.Config{
			Mode: packages.NeedTypes | packages.NeedTypesInfo | packages.NeedImports | packages.NeedDeps | packages.NeedSyntax,
			Fset: token.NewFileSet(),
		},
		"file="+path,
	)
	if err != nil {
		panic(err)
	}

	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			for _, imp := range file.Imports {
				if strings.Contains(imp.Comment.Text(), "typereg:") {
					suggested := strings.TrimPrefix(strings.TrimSpace(imp.Comment.Text()), "typereg:")

					name := strings.Trim(imp.Path.Value, "\"")

					if suggested != "" {
						imports[name] = suggested
					} else {
						imports[name] = imp.Name.String()
					}
				}
			}
		}

		for name, imp := range pkg.Imports {
			if as, ok := imports[name]; ok {
				for _, typ := range imp.TypesInfo.Types {
					if typ.IsType() && !typ.IsBuiltin() && unicode.IsUpper(rune(trueTypeName(typ.Type.String(), name)[0])) {
						if strings.HasPrefix(typ.Type.String(), name) {
							types[transformFullyQualifiedType(typ.Type.String(), name, as)] = struct{}{}
						}

						if strings.HasPrefix(typ.Type.String(), "*"+name) {
							types[transformFullyQualifiedType(typ.Type.String(), "*"+name, as)] = struct{}{}
						}
					}
				}
			}
		}
	}

	return imports, setToSlice(types)
}

func transformFullyQualifiedType(typ, name, as string) string {
	return strings.Replace(typ, name, as, 1)
}

func trueTypeName(typ, name string) string {
	typ = strings.TrimPrefix(typ, "*")
	return strings.TrimPrefix(typ, name+".")
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	imports, types := loadAllImportedTypes(filepath.Join(cwd, os.Getenv("GOFILE")))

	temp, err := template.New("type").Parse(typegenTemplate)
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile(generateOutputFileName(), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}

	temp.Execute(
		file,
		map[string]interface{}{
			"imports": imports,
			"types":   types,
			"package": os.Getenv("GOPACKAGE"),
		},
	)
}

func setToSlice(set map[string]struct{}) []string {
	slice := []string{}
	for k := range set {
		slice = append(slice, k)
	}

	return slice
}

func generateOutputFileName() string {
	return strings.Replace(os.Getenv("GOFILE"), ".go", ".typereg.generated.go", 1)
}
