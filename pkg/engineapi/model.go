package engineapi

type EngineFunctionArgument struct {
	Name string `xml:"name,attr"`
}

type EngineFunction struct {
	Name       string `xml:"name,attr"`
	ReturnType string `xml:"returnType,attr"`
	Symbol     string `xml:"symbol,attr"`
	IsCallback string `xml:"isCallback,attr"`
	IsVariadic string `xml:"isVariadic,attr"`
	Docs       string `xml:"docs,attr"`

	Arguments []EngineFunctionArgument `xml:"arguments>EngineFunctionArgument"`
}

type EngineEnum struct {
}

type EngineEnumType struct {
	Name           string `xml:"name,attr"`
	IsAbstract     string `xml:"isAbstract,attr"`
	IsInstantiable string `xml:"isInstantiable,attr"`
	IsDisposable   string `xml:"isDisposable,attr"`
	IsSingleton    string `xml:"isSingleton,attr"`
	Docs           string `xml:"docs,attr"`

	Enums []EngineEnum `xml:"enums>EngineEnum"`
}

type Exports struct {
	Functions []EngineFunction `xml:"EngineFunction"`
	Enums     []EngineEnumType `xml:"EngineEnumType"`
}

type EngineExportScope struct {
	Exports Exports `xml:"exports"`
}
