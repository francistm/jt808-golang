package main

import (
	. "github.com/dave/jennifer/jen"
	"github.com/francistm/jt808-golang/internal/generator"
)

func genMesgBodyIface(f *File, mesgDecls []*generator.MesgDecl) {
	for _, mesgDecl := range mesgDecls {
		for _, version := range mesgDecl.Versions {
			if version == nil {
				continue
			}

			f.Func().
				Params(
					Op("*").Id(version.StructName),
				).
				Id("MesgId").
				Params().
				Uint16().
				Block(
					Return(Lit(mesgDecl.MesgId)),
				)

			f.Line()
		}
	}
}
