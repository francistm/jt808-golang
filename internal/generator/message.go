package generator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strconv"
)

var structNameRegex = regexp.MustCompile(`^Body(\d+)$`)

type MessageDecl struct {
	Name     string
	HeaderID uint16
}

func buildMessageStructFromName(structName string) (*MessageDecl, error) {
	matched := structNameRegex.FindStringSubmatch(structName)

	if len(matched) != 2 {
		return nil, fmt.Errorf("messageName %s is invalid", structName)
	}

	headerID, err := strconv.ParseInt(matched[1], 16, 32)

	if err != nil {
		return nil, err
	}

	messageStruct := MessageDecl{
		Name:     structName,
		HeaderID: uint16(headerID),
	}

	return &messageStruct, nil
}

func GetAllMessageDecls() ([]*MessageDecl, error) {
	var messageStructs []*MessageDecl
	fileSet := token.NewFileSet()
	pkgMap, err := parser.ParseDir(fileSet, "message", nil, 0)

	if err != nil {
		return nil, err
	}

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

				s, err := buildMessageStructFromName(obj.Name)

				if err != nil {
					return nil, err
				}

				messageStructs = append(messageStructs, s)
			}
		}
	}

	return messageStructs, nil
}
