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
			Const().Id("unsupportErrMesg").Op("=").Lit("unsupport protocol version: %d"),
			Line(),

			If(Id("m").Dot("PackHeader").Dot("Package").Op("!=").Nil()).Block(
				Return(New(Id("PartialPackBody")), Nil()),
			).Else().Block(
				Switch(Id("m").Dot("PackHeader").Dot("MessageID")).BlockFunc(func(g *Group) {
					for _, mesgDecl := range mesgDecls {
						g.Case(Lit(mesgDecl.MesgId)).BlockFunc(func(g *Group) {
							if mesgDecl.Versions[0] != nil {
								g.Return(New(Id(mesgDecl.Versions[0].StructName)), Nil())
							} else {
								if mesgDecl.Versions[1] != nil {
									g.If(Id("m").Dot("PackHeader").Dot("Property").Dot("Version").Op("==").Id("Version2013")).
										Block(Return(
											New(Id(mesgDecl.Versions[1].StructName)),
											Nil(),
										))
								}

								if mesgDecl.Versions[2] != nil {
									g.If(Id("m").Dot("PackHeader").Dot("Property").Dot("Version").Op("==").Id("Version2019")).
										Block(Return(
											New(Id(mesgDecl.Versions[2].StructName)),
											Nil(),
										))
								}

								g.Return(
									Nil(),
									Qual("fmt", "Errorf").Call(
										Id("unsupportErrMesg"),
										Id("m").Dot("PackHeader").Dot("Property").Dot("Version"),
									),
								)
							}
						})
					}

					g.Default().Return(
						Nil(),
						Id("ErrMesgNotSupport"),
					)
				}),
			),
		)
}
