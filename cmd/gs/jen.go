package main

import (
	. "github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
)

func genMain(filename string, modules map[string]string) (*File, error) {
	f := NewFile("main")
	f.Func().Id("main").Params().Block(
		Id("ms").Op(":=").Index().
			Qual("github.com/kjbreil/goscript/pkg/module", "Module").
			ValuesFunc(func(group *Group) {
				for k, v := range modules {
					group.Add(Op("&")).Qual(v, strcase.ToCamel(k)).Block()
				}
			}),
		List(Id("config"), Err()).Op(":=").
			Qual("github.com/kjbreil/goscript/pkg/core", "ParseConfig").
			Call(Lit(filename), Id("ms")),
		ifError(),
		Line(),

		List(Id("gs"), Err()).Op(":=").
			Qual("github.com/kjbreil/goscript/pkg/core", "New").
			Call(List(Id("config"), Qual("github.com/kjbreil/goscript/pkg/core", "DefaultLogger").Call())),
		ifError(),
		Line(),

		Id("gs").Dot("UpdateModule").
			CallFunc(func(group *Group) {
				for k, v := range modules {
					group.Add(List(Lit(k), Qual("github.com/kjbreil/goscript/pkg/core", "GetModule").
						Types(Op("*").Qual(v, strcase.ToCamel(k))).
						Call(List(Id("gs"), Lit(k))),
					))
				}
			}),
		Err().Op("=").Id("gs").Dot("Connect").Call(),
		ifError(),
		Line(),

		Id("done").Op(":=").Make(List(Chan().Qual("os", "Signal"), Lit(1))),
		Qual(
			"os/signal",
			"Notify",
		).Call(Id("done"), Qual("os", "Interrupt"), Qual("syscall", "SIGINT"), Qual("syscall", "SIGTERM")),
		Id("gs").Dot("Logger").Call().Dot("Info").Call(Lit("Everything is set up")),
		Line(),

		Op("<-").Id("done"),
		Id("gs").Dot("Close").Call(),
	)
	return f, nil
}

func ifError() *Statement {
	return If(Err().Op("!=").Nil()).Block(
		Panic(Err()),
	)
}
