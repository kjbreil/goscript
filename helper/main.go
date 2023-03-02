package main

import (
	"fmt"
)

func main() {

	devices := DevicesInit()

	fmt.Println(devices)
	//
	//external := make(map[string]*jen.File)
	//models := make(map[string]*jen.File)
	//
	//fileList := []string{
	//	"types",
	//}
	//
	//modelsFileList := []string{
	//	"entity",
	//}
	//
	//outFolder := "entities"
	//modelFolder := "model"
	//
	//for _, v := range append(DeviceNames, fileList...) {
	//	external[v] = jen.NewFilePathName(fmt.Sprintf("./%s/%s.go", outFolder, v), outFolder)
	//	external[v].Comment("////////////////////////////////////////////////////////////////////////////////")
	//	external[v].Comment("Do not modify this file, it is automatically generated")
	//	external[v].Comment("////////////////////////////////////////////////////////////////////////////////")
	//}
	//
	//for _, v := range modelsFileList {
	//	models[v] = jen.NewFilePathName(fmt.Sprintf("./%s/%s.go", modelFolder, v), modelFolder)
	//	models[v].Comment("////////////////////////////////////////////////////////////////////////////////")
	//	models[v].Comment("Do not modify this file, it is automatically generated")
	//	models[v].Comment("////////////////////////////////////////////////////////////////////////////////")
	//}
	//
	//fileList = []string{"types"}
	//
	//external["types"].Type().Id("Entity").Interface()
	//
	//for _, d := range devices {
	//
	//	st := make(map[string][]*jen.Statement)
	//	for _, key := range d.Attributes {
	//		if key.DataType == "" {
	//			continue
	//		}
	//		st[key.Name] = append(st[key.Name], FieldAdder(key.Name, key.DataType))
	//	}
	//
	//	for _, v := range []string{
	//		"friendly_name",
	//	} {
	//		st[v] = append(st[v], jen.Id(strcase.ToCamel(v)).Op("*").String().Tag(map[string]string{"json": v + ",omitempty"}))
	//	}
	//
	//	for _, v := range []string{
	//		"additional",
	//	} {
	//		st[v] = append(st[v], jen.Id(strcase.ToCamel(v)).Map(jen.String()).Interface().Tag(map[string]string{"json": v + ",omitempty"}))
	//	}
	//
	//	sortedKeys := []string{}
	//	for key := range st {
	//		sortedKeys = append(sortedKeys, key)
	//	}
	//	sort.Strings(sortedKeys)
	//	external[d.Name].Type().Id(strcase.ToCamel(d.Name)).StructFunc(
	//		func(g *jen.Group) {
	//			for _, key := range sortedKeys {
	//				v := st[key]
	//				for _, val := range v {
	//					g.Add(val)
	//				}
	//			}
	//		},
	//	)
	//
	//	external[d.Name].Func().
	//		Id(fmt.Sprintf("Get%s", strcase.ToCamel(d.Name))).
	//		Params(jen.Id("attributes").Map(jen.String()).Interface()).
	//		Op("*").Id(strcase.ToCamel(d.Name)).
	//		Block(
	//			jen.Var().Id(string(strcase.ToLowerCamel(d.Name)[0])).Id(strcase.ToCamel(d.Name)),
	//			jen.Id("fillFields").Params(jen.Op("&").Id(string(strcase.ToLowerCamel(d.Name)[0])), jen.Id("attributes")),
	//			jen.Id(string(strcase.ToLowerCamel(d.Name)[0])).Dot("Additional").Op("=").Id("attributes"),
	//			jen.Return(jen.Op("&").Id(string(strcase.ToLowerCamel(d.Name)[0]))),
	//		)
	//
	//}
	//
	//for _, v := range modelsFileList {
	//	sortedKeys := []string{}
	//	for _, key := range devices {
	//		sortedKeys = append(sortedKeys, key.Name)
	//	}
	//	sort.Strings(sortedKeys)
	//
	//	for i := 0; i < len(sortedKeys); i++ {
	//		camelName := strcase.ToCamel(sortedKeys[i])
	//		models[v].Type().Id(fmt.Sprintf("On%sHandler", camelName)).Func().Params(jen.Id("message").Op("*").Id("Message"), jen.Id("newAttrs"), jen.Id("oldAttrs").Op("*").Qual("github.com/kjbreil/hass-ws/entities", camelName))
	//	}
	//
	//	models[v].Type().Id("OnTypeHandlers").StructFunc(func(g *jen.Group) {
	//		for i := 0; i < len(sortedKeys); i++ {
	//			camelName := strcase.ToCamel(sortedKeys[i])
	//			g.Add(jen.Id(fmt.Sprintf("On%s", camelName)).Id(fmt.Sprintf("On%sHandler", camelName)))
	//		}
	//	})
	//
	//	models[v].Func().Params(
	//		jen.Id("o").Op("*").Id("OnTypeHandlers"),
	//	).Id("Run").Params(
	//		jen.Id("message").Op("*").Id("Message"),
	//	).Block(
	//		jen.If(jen.Id("message").Dot("Event").Op("==").Id("nil").Op("||").
	//			Id("message").Dot("Event").Dot("Data").Op("==").Id("nil").Op("||").
	//			Id("message").Dot("Event").Dot("Data").Dot("EntityId").Op("==").Id("nil")).
	//			Block(jen.Return()),
	//		jen.Id("entityType").Op(":=").Qual("strings", "Split").Params(
	//			jen.Op("*").Id("message").Dot("Event").Dot("Data").Dot("EntityId"),
	//			jen.Lit("."),
	//		).Id("[0]"),
	//
	//		jen.Switch(jen.Id("entityType")).BlockFunc(func(g *jen.Group) {
	//			for i := 0; i < len(sortedKeys); i++ {
	//				camelName := strcase.ToCamel(sortedKeys[i])
	//				underName := strcase.ToSnake(sortedKeys[i])
	//				g.Add(jen.Case(jen.Lit(underName)).Block(
	//					jen.If(jen.Id("o").Dot(fmt.Sprintf("On%s", camelName)).Op("==").Id("nil")).Block(jen.Return()),
	//					jen.Id("newAttrs").Op(":=").Qual("github.com/kjbreil/hass-ws/entities", fmt.Sprintf("Get%s", camelName)).Params(jen.Id("message").Dot("Event").Dot("Data").Dot("NewState").Dot("Attributes")),
	//					jen.Id("oldAttrs").Op(":=").Qual("github.com/kjbreil/hass-ws/entities", fmt.Sprintf("Get%s", camelName)).Params(jen.Id("message").Dot("Event").Dot("Data").Dot("OldState").Dot("Attributes")),
	//					jen.Id("o").Dot(fmt.Sprintf("On%s", camelName)).Params(jen.Id("message"), jen.Id("newAttrs"), jen.Id("oldAttrs")),
	//				))
	//				// .Params(jen.Id("message").Id("newAttrs").Id("oldAttrs"))
	//			}
	//		}),
	//	)
	//
	//	models[v].Func().Params(
	//		jen.Id("o").Op("*").Id("OnTypeHandlers"),
	//	).Id("RunStates").Params(
	//		jen.Id("message").Op("*").Id("Message"),
	//	).Block(
	//		jen.If(jen.Id("message").Dot("Result").Op("==").Id("nil")).
	//			Block(jen.Return()),
	//		jen.For(jen.Id("_").Op(",").Id("r").Op(":=").Range().Id("message").Dot("Result")).Block(
	//
	//			jen.Id("entityType").Op(":=").Qual("strings", "Split").Params(
	//				jen.Op("*").Id("r").Dot("EntityId"),
	//				jen.Lit("."),
	//			).Id("[0]"),
	//			jen.Switch(jen.Id("entityType")).BlockFunc(func(g *jen.Group) {
	//				for i := 0; i < len(sortedKeys); i++ {
	//					camelName := strcase.ToCamel(sortedKeys[i])
	//					underName := strcase.ToSnake(sortedKeys[i])
	//					g.Add(jen.Case(jen.Lit(underName)).Block(
	//						jen.If(jen.Id("o").Dot(fmt.Sprintf("On%s", camelName)).Op("==").Id("nil")).Block(jen.Continue()),
	//						jen.Id("newAttrs").Op(":=").Qual("github.com/kjbreil/hass-ws/entities", fmt.Sprintf("Get%s", camelName)).Params(jen.Id("r").Dot("Attributes")),
	//						jen.Id("newMsg").Op(":=").Op("&").Id("Message").Values(jen.Dict{
	//
	//							jen.Id("ID"):   jen.Id("message").Dot("ID"),
	//							jen.Id("Type"): jen.Id("MessageTypeEvent"),
	//							jen.Id("Event"): jen.Op("&").Id("Event").Values(
	//								jen.Dict{
	//									jen.Id("Data"): jen.Op("&").Id("Data").Values(jen.Dict{
	//										jen.Id("EntityId"): jen.Id("r").Dot("EntityId"),
	//										jen.Id("NewState"): jen.Op("&").Id("State").Values(
	//											jen.Dict{
	//												jen.Id("EntityId"):    jen.Id("r").Dot("EntityId"),
	//												jen.Id("LastChanged"): jen.Id("r").Dot("LastChanged"),
	//												jen.Id("State"):       jen.Id("r").Dot("State"),
	//												jen.Id("Attributes"):  jen.Id("r").Dot("Attributes"),
	//												jen.Id("LastUpdated"): jen.Id("r").Dot("LastUpdated"),
	//												jen.Id("Context"):     jen.Id("r").Dot("Context"),
	//											},
	//										),
	//										jen.Id("OldState"): jen.Id("nil"),
	//									}),
	//									jen.Id("Context"):   jen.Id("r").Dot("Context"),
	//									jen.Id("EventType"): jen.Id("EventTypeState"),
	//								},
	//							),
	//						},
	//						),
	//
	//						jen.Id("o").Dot(fmt.Sprintf("On%s", camelName)).Params(jen.Id("newMsg"), jen.Id("newAttrs"), jen.Id("nil")),
	//					))
	//					// .Params(jen.Id("message").Id("newAttrs").Id("oldAttrs"))
	//				}
	//			}),
	//		),
	//	)
	//
	//	models[v].Type().Id("OnEntityHandlers").Map(jen.String()).Func().Params(jen.Id("message").Op("*").Id("Message"))
	//
	//	models[v].Func().Params(
	//		jen.Id("o").Id("OnEntityHandlers"),
	//	).Id("Run").Params(
	//		jen.Id("message").Op("*").Id("Message"),
	//	).Block(
	//		jen.If(jen.Id("message").Dot("Event").Op("==").Id("nil").Op("||").
	//			Id("message").Dot("Event").Dot("Data").Op("==").Id("nil").Op("||").
	//			Id("message").Dot("Event").Dot("Data").Dot("EntityId").Op("==").Id("nil")).
	//			Block(jen.Return()),
	//		jen.For(jen.Id("k").Op(",").Id("v").Op(":=").Range().Id("o")).Block(
	//			jen.If(jen.Id("k").Op("==").Op("*").Id("message").Dot("Event").Dot("Data").Dot("EntityId").Block(
	//				jen.Id("v").Call(jen.Id("message")),
	//			)),
	//		),
	//	)
	//
	//	models[v].Func().Params(
	//		jen.Id("o").Id("OnEntityHandlers"),
	//	).Id("RunStates").Params(
	//		jen.Id("message").Op("*").Id("Message"),
	//	).Block(
	//		jen.If(jen.Id("message").Dot("Result").Op("==").Id("nil")).
	//			Block(jen.Return()),
	//		jen.For(jen.Id("_").Op(",").Id("r").Op(":=").Range().Id("message").Dot("Result")).Block(
	//			jen.For(jen.Id("k").Op(",").Id("v").Op(":=").Range().Id("o")).Block(
	//
	//				jen.If(jen.Id("k").Op("==").Op("*").Id("r").Dot("EntityId").Block(
	//
	//					jen.Id("newMsg").Op(":=").Op("&").Id("Message").Values(jen.Dict{
	//
	//						jen.Id("ID"):   jen.Id("message").Dot("ID"),
	//						jen.Id("Type"): jen.Id("MessageTypeEvent"),
	//						jen.Id("Event"): jen.Op("&").Id("Event").Values(
	//							jen.Dict{
	//								jen.Id("Data"): jen.Op("&").Id("Data").Values(jen.Dict{
	//									jen.Id("EntityId"): jen.Id("r").Dot("EntityId"),
	//									jen.Id("NewState"): jen.Op("&").Id("State").Values(
	//										jen.Dict{
	//											jen.Id("EntityId"):    jen.Id("r").Dot("EntityId"),
	//											jen.Id("LastChanged"): jen.Id("r").Dot("LastChanged"),
	//											jen.Id("State"):       jen.Id("r").Dot("State"),
	//											jen.Id("Attributes"):  jen.Id("r").Dot("Attributes"),
	//											jen.Id("LastUpdated"): jen.Id("r").Dot("LastUpdated"),
	//											jen.Id("Context"):     jen.Id("r").Dot("Context"),
	//										},
	//									),
	//									jen.Id("OldState"): jen.Id("nil"),
	//								}),
	//								jen.Id("Context"):   jen.Id("r").Dot("Context"),
	//								jen.Id("EventType"): jen.Id("EventTypeState"),
	//							},
	//						),
	//					},
	//					),
	//
	//					jen.Id("v").Call(jen.Id("newMsg")),
	//				)),
	//			),
	//		),
	//	)
	//}
	//
	//for k, v := range external {
	//	v.Save(fmt.Sprintf("./%s/%s.go", outFolder, k))
	//}
	//for k, v := range models {
	//	err := v.Save(fmt.Sprintf("./%s/%s.go", modelFolder, k))
	//	if err != nil {
	//		log.Panicln(err)
	//	}
	//}

}
