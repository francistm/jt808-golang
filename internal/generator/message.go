package generator

import (
	"go/parser"
	"go/token"
	"os"
	"path"
)

func ParseMesgDecls() ([]*MesgDecl, error) {
	var (
		srcPath string
		fileSet = token.NewFileSet()
	)

	cwd, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	srcPath = path.Join(cwd, "message")
	pkgMap, err := parser.ParseDir(fileSet, srcPath, nil, parser.DeclarationErrors)

	if err != nil {
		return nil, err
	}

	declNames := parseMesgDeclNames(pkgMap)
	mesgDecls, err := buildMesgDecls(declNames)

	if err != nil {
		return nil, err
	}

	return mesgDecls, nil
}
