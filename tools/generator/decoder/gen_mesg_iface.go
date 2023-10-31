package main

import (
	. "github.com/dave/jennifer/jen"
	"github.com/francistm/jt808-golang/internal/generator"
)

func genMesgBodyIface(f *File, mesgDecls []*generator.MessageDecl) {
	for _, mesgDecl := range mesgDecls {
		f.Func().
			Params(
				Op("*").Id(mesgDecl.StructName),
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
