package main

import (
	. "github.com/dave/jennifer/jen"
	"github.com/francistm/jt808-golang/internal/generator"
)

func genNewPkgBodyFromMesgId(f *File, mesgDecls []*generator.MesgDecl) {
	f.Func().
		Params(
			Id("m").Op("*").Id("MessagePack").Types(Id("T")),
		).
		Id("NewPackBodyFromMesgId").
		Params().
		Parens(List(
			Id("MesgBody"),
			Error(),
		)).
		Block(
			If(Id("m").Dot("PackHeader").Dot("Package").Op("!=").Nil()).Block(
				Return(New(Id("PartialPackBody")), Nil()),
			).Else().Block(
				Switch(Id("m").Dot("PackHeader").Dot("MessageID")).BlockFunc(func(g *Group) {
					for _, mesgDecl := range mesgDecls {
						if len(mesgDecl.Versions) == 1 {
							g.Case(Lit(mesgDecl.MesgId)).Block(Return(New(Id(mesgDecl.Versions[0].StructName)), Nil()))
						} else if len(mesgDecl.Versions) == 2 {
							g.Case(Lit(mesgDecl.MesgId)).BlockFunc(func(g *Group) {
								g.If(Id("m").Dot("PackHeader").Dot("Property").Dot("Version").Op("==").Id("Version2013")).
									Block(Return(
										New(Id(mesgDecl.Versions[0].StructName)),
										Nil(),
									)).
									Else().If(Id("m").Dot("PackHeader").Dot("Property").Dot("Version").Op("==").Id("Version2019")).
									Block(Return(
										New(Id(mesgDecl.Versions[1].StructName)),
										Nil(),
									)).
									Else().
									Block(Return(
										Nil(),
										Qual("fmt", "Errorf").Call(
											Lit("unsupport protocol version: %d"),
											Id("m").Dot("PackHeader").Dot("Property").Dot("Version"),
										),
									))
							})
						}
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
