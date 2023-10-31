package generator

import (
	"go/ast"
	"regexp"
)

// regexp to match message struct name
// e.g. Body0001, Body0001_11, Body0001_19
var structNameRegex = regexp.MustCompile(`^Body(\d{4})(_(\d{2}))?$`)

func parseMesgDeclNames(pkgMap map[string]*ast.Package) []string {
	declNames := make([]string, 0, 100)

	for _, pkg := range pkgMap {
		for _, astFile := range pkg.Files {
			if astFile.Scope == nil {
				continue
			}

			for _, obj := range astFile.Scope.Objects {
				if obj.Kind != ast.Typ {
					continue
				}

				decl := obj.Decl.(*ast.TypeSpec)

				if _, ok := decl.Type.(*ast.StructType); !ok {
					continue
				}

				if !structNameRegex.MatchString(obj.Name) {
					continue
				}

				declNames = append(declNames, obj.Name)
			}
		}
	}

	return declNames
}
