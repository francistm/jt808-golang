package main

import (
	"bytes"
	"io/ioutil"

	. "github.com/dave/jennifer/jen"
	"github.com/francistm/jt808-golang/internal/generator"
)

func main() {
	f := NewFile("jt808")
	f.HeaderComment("Code generated by generator/message, DO NOT MODIFY MANUALLY")
	messageStructs, err := generator.GetAllMessageStructs()

	if err != nil {
		panic(err)
	}

	for _, messageStruct := range messageStructs {
		f.Func().Params(Id("m").Id("*PackBody")).Id("As"+messageStruct.Name).
			Params().
			Op("*").Qual("github.com/francistm/jt808-golang/message", messageStruct.Name).
			Block(
				Return(Id("m").Dot("body").Assert(Op("*").Qual("github.com/francistm/jt808-golang/message", messageStruct.Name))),
			)

		f.Line()
	}

	buf := new(bytes.Buffer)

	if err := f.Render(buf); err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile("messages.gen.go", buf.Bytes(), 0644); err != nil {
		panic(err)
	}
}