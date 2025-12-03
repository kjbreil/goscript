package main

import (
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
)

func genMain(filename string, modules map[string]string) (*jen.File, error) {
	f := jen.NewFile("main")
	f.Func().Id("main").Params().Block(
		jen.Id("ms").Op(":=").Index().
			Qual("github.com/kjbreil/goscript/pkg/module", "Module").
			ValuesFunc(func(group *jen.Group) {
				for k, v := range modules {
					group.Add(jen.Op("&")).Qual(v, strcase.ToCamel(k)).Block()
				}
			}),
		jen.List(jen.Id("config"), jen.Err()).Op(":=").
			Qual("github.com/kjbreil/goscript/pkg/core", "ParseConfig").
			Call(jen.Lit(filename), jen.Id("ms")),
		ifError(),
		jen.Line(),

		jen.List(jen.Id("gs"), jen.Err()).Op(":=").
			Qual("github.com/kjbreil/goscript/pkg/core", "New").
			Call(jen.List(jen.Id("config"), jen.Qual("github.com/kjbreil/goscript/pkg/core", "DefaultLogger").Call())),
		ifError(),
		jen.Line(),

		jen.Id("gs").Dot("UpdateModule").
			CallFunc(func(group *jen.Group) {
				for k, v := range modules {
					group.Add(jen.List(jen.Lit(k), jen.Qual("github.com/kjbreil/goscript/pkg/core", "GetModule").
						Types(jen.Op("*").Qual(v, strcase.ToCamel(k))).
						Call(jen.List(jen.Id("gs"), jen.Lit(k))),
					))
				}
			}),
		jen.Err().Op("=").Id("gs").Dot("Connect").Call(),
		ifError(),
		jen.Line(),

		jen.Id("done").Op(":=").Make(jen.List(jen.Chan().Qual("os", "Signal"), jen.Lit(1))),
		jen.Qual(
			"os/signal",
			"Notify",
		).Call(jen.Id("done"), jen.Qual("os", "Interrupt"), jen.Qual("syscall", "SIGINT"), jen.Qual("syscall", "SIGTERM")),
		jen.Id("gs").Dot("Logger").Call().Dot("Info").Call(jen.Lit("Everything is set up")),
		jen.Line(),

		jen.Op("<-").Id("done"),
		jen.Id("gs").Dot("Close").Call(),
	)
	return f, nil
}

func ifError() *jen.Statement {
	return jen.If(jen.Err().Op("!=").Nil()).Block(
		jen.Panic(jen.Err()),
	)
}
