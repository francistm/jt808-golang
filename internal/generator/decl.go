package generator

import (
	"fmt"
	"sort"
	"strconv"
)

type MesgDecl struct {
	MesgId   uint16
	Versions []MesgDeclVersion
}

type MesgDeclVersion struct {
	Version    int
	StructName string
}

func parseDeclName(structName string) (uint16, int, error) {
	var (
		version int
		matched = structNameRegex.FindStringSubmatch(structName)
	)

	if len(matched) != 4 {
		return 0, 0, fmt.Errorf("invalid name %s", structName)
	}

	mesgIdParsed, err := strconv.ParseUint(matched[1], 16, 16)

	if err != nil {
		return 0, 0, fmt.Errorf("invalid name %s, %w", structName, err)
	}

	if matched[3] != "" {
		versionParsed, err := strconv.ParseUint(matched[3], 10, 8)

		if err != nil {
			return 0, 0, fmt.Errorf("invalid name %s, %w", structName, err)
		}

		version = int(versionParsed)
	}

	return uint16(mesgIdParsed), version, nil
}

func buildMesgDecls(declNames []string) ([]*MesgDecl, error) {
	var (
		mesgDecls    = make([]*MesgDecl, 0, 100)
		mesgDeclsMap = make(map[uint16]*MesgDecl, 100)
	)

	for _, typeName := range declNames {
		mesgId, version, err := parseDeclName(typeName)

		if err != nil {
			return nil, err
		}

		mesgDecl, ok := mesgDeclsMap[mesgId]

		if !ok {
			mesgDecl = &MesgDecl{
				MesgId:   mesgId,
				Versions: make([]MesgDeclVersion, 0, 2),
			}

			mesgDecl.Versions = append(mesgDecl.Versions, MesgDeclVersion{
				Version:    version,
				StructName: typeName,
			})

			mesgDeclsMap[mesgId] = mesgDecl
			mesgDecls = append(mesgDecls, mesgDecl)
		} else {
			mesgDecl.Versions = append(mesgDecl.Versions, MesgDeclVersion{
				Version:    version,
				StructName: typeName,
			})
		}
	}

	sort.Slice(mesgDecls, func(i, j int) bool {
		return mesgDecls[i].MesgId < mesgDecls[j].MesgId
	})

	for _, mesgDecl := range mesgDecls {
		sort.Slice(mesgDecl.Versions, func(i, j int) bool {
			return mesgDecl.Versions[i].Version < mesgDecl.Versions[j].Version
		})
	}

	return mesgDecls, nil
}
