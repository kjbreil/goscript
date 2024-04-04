package main

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
)

func genMain(filename string, modules map[string]string) ([]byte, error) {
	f := NewFile("main")
	f.Func().Id("main").Params().Block(
		Id("ms").Op(":=").Index().
			Qual("github.com/kjbreil/goscript/pkg/core", "Module").
			ValuesFunc(func(group *Group) {
				for k, v := range modules {
					group.Add(Op("&")).Qual(v, k).Block()
				}
			}),
		List(Id("config"), Error()).Op(":=").
			Qual("github.com/kjbreil/goscript/pkg/core", "ParseConfig").
			Call(Lit(filename)),
		ifError(),
		List(Id("gs"), Error()).Op(":=").
			Qual("github.com/kjbreil/goscript/pkg/core", "New").
			Call(List(Id("config"), Qual("github.com/kjbreil/goscript/pkg/core", "DefaultLogger").Call())),
		ifError(),

		Id("gs").Dot("UpdateModule").
			CallFunc(func(group *Group) {
				for k, v := range modules {
					group.Add(List(Lit(k), Qual("github.com/kjbreil/goscript/pkg/core", "GetModule").
						Types(Op("*").Qual(v, strcase.ToCamel(k))).
						Call(List(Id("gs"), Lit(k))),
					))
				}
			}),
		Id("err").Op(":=").Id("gs").Dot("Connect"),
		ifError(),
		Id("done").Op(":=").Make(List(Chan().Qual("os", "Signal"), Lit(1))),
		Qual("signal", "Notify").Call(Id("done"), Qual("os", "Interrupt"), Qual("syscall", "SIGINT"), Qual("syscall", "SIGTERM")),
		Id("gs").Dot("Logger").Call().Dot("Info").Call(Lit("Everything is set up")),
		Op("<-").Id("done"),
		Id("gs").Dot("Close").Call(),
	)
	fmt.Printf("%#v", f)
	return nil, nil
}

func ifError() *Statement {
	return If(Err().Op("!=").Nil()).Block(
		Panic(Err()),
	)
}
