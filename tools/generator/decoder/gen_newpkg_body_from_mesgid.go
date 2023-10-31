package main

import (
	. "github.com/dave/jennifer/jen"
	"github.com/francistm/jt808-golang/internal/generator"
)

func genNewPkgBodyFromMesgId(f *File, mesgDecls []*generator.MessageDecl) {
	f.Func().
		Params(
			Id("m").Op("*").Id("MessagePack").Types(Id("T")),
		).
		Id("NewPackBodyFromMesgId").
		Params().
		Parens(List(
			Any(),
			Error(),
		)).
		Block(
			If(Id("m").Dot("PackHeader").Dot("Package").Op("!=").Nil()).Block(
				Return(New(Id("PartialPackBody")), Nil()),
			).Else().Block(
				Switch(Id("m").Dot("PackHeader").Dot("MessageID")).BlockFunc(func(g *Group) {
					for _, mesgDecl := range mesgDecls {
						g.Case(Lit(mesgDecl.MesgId)).Block(Return(New(Id(mesgDecl.StructName)), Nil()))
					}

					g.Default().Return(
						Nil(),
						Qual("fmt", "Errorf").Call(
							Lit("unsupported messageId: 0x%.4X"),
							Id("m").Dot("PackHeader").Dot("MessageID"),
						),
					)
				}),
			),
		)
}
