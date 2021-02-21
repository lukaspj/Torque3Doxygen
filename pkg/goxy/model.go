package goxy

const (
	Class     Kind = "class"
	File      Kind = "file"
	Struct    Kind = "struct"
	Namespace Kind = "namespace"
	Group     Kind = "group"
	Dir       Kind = "dir"
	Union     Kind = "union"
	Page      Kind = "page"
)

const (
	Public Protection = iota
	Protected
	Private
)

const (
	Functions SectionKind = iota
	StaticFunctions
	Attributes
	StaticAttributes
	Types
	Friends
	UserDefined
	Defines
	Typedefs
	Variables
	Enums
	Related
)

const (
	Member KindRef = iota
	Compound
)

const (
	Section        DocStringType = "section"
	Paragraph      DocStringType = "paragraph"
	Anchor         DocStringType = "anchor"
	Image          DocStringType = "image"
	Text           DocStringType = "text"
	Ref            DocStringType = "ref"
	Title          DocStringType = "title"
	Heading        DocStringType = "heading"
	XRefSect       DocStringType = "xrefsect"
	Table          DocStringType = "table"
	ParameterList  DocStringType = "parameterlist"
	ItemizedList   DocStringType = "itemizedlist"
	OrderedList    DocStringType = "orderedlist"
	VariableList   DocStringType = "variablelist"
	Bold           DocStringType = "bold"
	Emphasis       DocStringType = "emphasis"
	Verbatim       DocStringType = "verbatim"
	Preformatted   DocStringType = "preformatted"
	Term           DocStringType = "term"
	ComputerOutput DocStringType = "computeroutput"
	LineBreak      DocStringType = "linebreak"
	Highlight      DocStringType = "highlight"
)

type Kind string
type Protection int
type SectionKind int
type KindRef int
type DocStringType string

type SourceLocation struct {
	File          string
	FileRefId     string
	Line          int
	Column        int
	BodyFile      string
	BodyFileRefId string
	BodyStart     int
	BodyEnd       int
}

type DocStringSection struct {
	Id      string
	Kind    string
	Content DocString
}

type DocStringParagraph struct {
	Content DocString
}

type DocStringText struct {
	Content string
}

type DocStringTitle struct {
	Content DocString
}

type DocStringHeading struct {
	Content DocString
	Level   int
}

type DocStringParameterItem struct {
	Name        string
	Description DocString
}

type DocStringParameterList struct {
	Kind  string
	Items []DocStringParameterItem
}

type DocStringTableEntry struct {
	Head    bool
	Content DocString
}

type DocStringTable struct {
	Rows [][]DocStringTableEntry
}

type DocStringXRefSect struct {
	Id          string
	Title       string
	Description DocString
}

type DocStringOrderedList struct {
	Items []DocString
}

type DocStringItemizedList struct {
	Items []DocString
}

type DocStringVariableList struct {
	Items []DocString
}

type DocStringBold struct {
	Content DocString
}

type DocStringEmphasis struct {
	Content DocString
}

type DocStringVerbatim struct {
	Content DocString
}

type DocStringPreformatted struct {
	Content DocString
}

type DocStringTerm struct {
	Content DocString
}

type DocStringComputerOutput struct {
	Content DocString
}

type DocStringLinebreak struct{}

type DocStringHighlight struct {
	Content DocString
}

type DocStringRef struct {
	RefId   string
	KindRef string
	Content DocString
}

type DocStringAnchor struct {
	Id string
}

type DocStringImage struct {
	Name        string
	Type        string
	Description string
}

type DocStringElement struct {
	Type  DocStringType
	Value interface{}
}

type DocString struct {
	Content []DocStringElement
}

type Descriptions struct {
	BriefDescription    DocString `json:"brief_description,omitempty"`
	DetailedDescription DocString `json:"detailed_description,omitempty"`
	InBodyDescription   DocString `json:"in_body_description,omitempty"`
}

type EnumValue struct {
	Descriptions

	Id          string
	Initializer string
	Protection  Protection
	Name        string
}

type ClassAttributeDoc struct {
	Descriptions

	Id         string
	Name       string
	Protection Protection
	Type       DocString
	Location   SourceLocation

	Definition string
	ArgsString DocString
}

type InnerCompoundRef struct {
	RefId      string
	Protection Protection
	Value      string
}

type EnumDoc struct {
	Descriptions

	Id         string
	Name       string
	Protection Protection
	Type       DocString
	Location   SourceLocation

	Values []EnumValue
}

type FriendDoc struct {
	Descriptions

	Id         string
	Name       string
	Protection Protection
	Type       DocString
	Definition string
	Location   SourceLocation
}

type FunctionParam struct {
	Type     DocString
	DeclName string
}

type Reimplements struct {
	RefId    string
	MemberId string
	ParentId string
}

type FunctionDoc struct {
	Descriptions

	Id         string
	Name       string
	Protection Protection
	Type       DocString
	Location   SourceLocation

	Definition string
	ArgsString string
	Params     []FunctionParam

	Reimplements    Reimplements
	ReimplementedBy []Reimplements
}

type SectionDoc struct {
	Id         string
	Kind       SectionKind
	Protection Protection

	Header      string
	Description DocString

	Attributes []*ClassAttributeDoc
	Enums      []*EnumDoc
	Functions  []*FunctionDoc
	Defines    []*DefineDoc
	Typedefs   []*TypedefDoc
	Friends    []*FriendDoc
}

type DefineParam struct {
	Defname string
}

type DefineDoc struct {
	Descriptions

	Id          string
	Name        string
	Initializer string
	Params      []DefineParam
}

type TypedefDoc struct {
	Descriptions

	Id         string
	Name       string
	Type       DocString
	Definition string
	ArgsString DocString
}

type CompoundRef struct {
	Kind      string
	Name      string
	RefId     string
	ParentRef string
}

type GraphNode struct {
	Id    int
	Label string
	RefId string
}

type GraphEdge struct {
	FromId    int
	ToId      int
	Relation  string
	EdgeLabel string
}

type Graph struct {
	Nodes []GraphNode
	Edges []GraphEdge
}

func (g Graph) ResolveId(id int) *GraphNode {
	for _, node := range g.Nodes {
		if node.Id == id {
			return &node
		}
	}
	return nil
}

type CompoundDoc struct {
	Descriptions

	Parent string

	Id   string
	Kind Kind

	Name       string
	Title      string
	Sections   []*SectionDoc
	Protection Protection

	InnerClasses    []InnerCompoundRef
	InnerFiles      []InnerCompoundRef
	InnerDirs       []InnerCompoundRef
	InnerGroups     []InnerCompoundRef
	InnerNamespaces []InnerCompoundRef

	Location       SourceLocation
	ProgramListing DocString

	InheritanceGraph Graph
}
