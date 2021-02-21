package doxygen

import (
	"ScriptExecServer/pkg/xmlhelper"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

var (
	FunctionMember = "function"
	VariableMember = "variable"
	EnumMember     = "enum"
	DefineMember   = "define"
	TypedefMember  = "typedef"
	FriendMember   = "friend"

	ClassDoc     = "class"
	StructDoc    = "struct"
	FileDoc      = "file"
	DirDoc       = "dir"
	NamespaceDoc = "namespace"
	GroupDoc     = "group"
	UnionDoc     = "union"
	PageDoc      = "page"
)

type DocString struct {
	Content []interface{}
}

type ParameterItem struct {
	Name        string    `xml:"parameternamelist>parametername"`
	Description DocString `xml:"parameterdescription"`
}

type ParameterList struct {
	Kind  string          `xml:"kind,attr"`
	Items []ParameterItem `xml:"parameteritem"`
}

type Term struct {
	Content DocString
}

type VariableList struct {
	Items []DocString `xml:",any"`
}

type Ref struct {
	RefId   string    `xml:"refid,attr"`
	KindRef string    `xml:"kindref,attr"`
	Content DocString `xml:",any"`
}

type XRefSect struct {
	Id          string    `xml:"id,attr"`
	Title       string    `xml:"xreftitle"`
	Description DocString `xml:"xrefdescription"`
}

type ItemizedList struct {
	Items []DocString `xml:"listitem"`
}

type OrderedList struct {
	Items []DocString `xml:"listitem"`
}

type Bold struct {
	Content DocString
}

type Emphasis struct {
	Content DocString
}

type Verbatim struct {
	Content DocString
}

type Preformatted struct {
	Content DocString
}

type ComputerOutput struct {
	Content DocString
}

type Anchor struct {
	Id string `xml:"id,attr"`
}

type Linebreak struct{}

type ProgramListing struct {
	Content  DocString
	Filename string `xml:"filename,attr"`
}

type Paragraph struct {
	Content DocString `xml:",any"`
}

type Text struct {
	Content string `xml:",chardata"`
}

type Title struct {
	Content DocString `xml:",any"`
}

type Heading struct {
	Content DocString `xml:",any"`
	Level   int       `xml:"level,attr"`
}

type Section struct {
	Id      string    `xml:"id,attr"`
	Kind    string    `xml:"kind,attr"`
	Content DocString `xml:",any"`
}

type TableEntry struct {
	TableHead bool      `xml:"thead,attr"`
	Content   DocString `xml:",any"`
}

type TableRow struct {
	Columns []TableEntry `xml:"entry"`
}

type Table struct {
	RowCount    int `xml:"rows,attr"`
	ColumnCount int `xml:"cols,attr"`

	Rows []TableRow `xml:"row"`
}

type Image struct {
	Type        string `xml:"type,attr"`
	Name        string `xml:"name,attr"`
	Description string `xml:",chardata"`
}

type Descriptions struct {
	BriefDescription    DocString `xml:"briefdescription"`
	DetailedDescription DocString `xml:"detaileddescription"`
	InBodyDescription   DocString `xml:"inbodydescription"`
}

type Doxygen struct {
	CompoundDef CompoundDef `xml:"compounddef"`
}

type UnkownElement struct {
	xml.Name
}

type GraphChildNode struct {
	RefId     int      `xml:"refid,attr"`
	Relation  string   `xml:"relation,attr"`
	EdgeLabel []string `xml:"edgelabel"`
}

type GraphLink struct {
	RefId string `xml:"refid,attr"`
}

type GraphNode struct {
	Id       int              `xml:"id,attr"`
	Label    string           `xml:"label"`
	Link     GraphLink        `xml:"link"`
	Children []GraphChildNode `xml:"childnode"`
}

type Graph struct {
	Nodes []GraphNode `xml:"node"`
}

type CompoundDef struct {
	Descriptions

	CompoundName string `xml:"compoundname"`
	Title        string `xml:"title"`
	Id           string `xml:"id,attr"`
	Kind         string `xml:"kind,attr"`
	Language     string `xml:"language,attr"`
	Protection   string `xml:"prot,attr"`
	Virtual      string `xml:"virt,attr"`

	BaseCompoundRef BaseCompoundRef `xml:"basecompoundref"`
	Sections        []SectionDef    `xml:"sectiondef"`
	Location        Location        `xml:"location"`
	ProgramListing  *ProgramListing `xml:"programlisting"`

	Includes   []Include `xml:"includes"`
	IncludedBy []Include `xml:"includedby"`

	InnerClass      []InnerCompound `xml:"innerclass"`
	InnerFiles      []InnerCompound `xml:"innerfile"`
	InnerNamespaces []InnerCompound `xml:"innernamespace"`
	InnerGroups     []InnerCompound `xml:"innergroup"`
	InnerDirs       []InnerCompound `xml:"innerdir"`

	InheritanceGraph Graph `xml:"inheritancegraph"`
}

type Include struct {
	RefId string `xml:"refid,attr"`
	Local string `xml:"local,attr"`
	Value string `xml:",chardata"`
}

type InnerCompound struct {
	RefId string `xml:"refid,attr"`
	Prot  string `xml:"prot,attr"`
	Value string `xml:",chardata"`
}

type Location struct {
	File      string `xml:"file,attr"`
	Line      int    `xml:"line,attr"`
	Column    int    `xml:"column,attr"`
	BodyFile  string `xml:"bodyFile,attr"`
	BodyStart int    `xml:"bodyStart,attr"`
	BodyEnd   int    `xml:"bodyEnd,attr"`
}

type BaseCompoundRef struct {
	RefId string `xml:"refid,attr"`
	Prot  string `xml:"prot,attr"`
	Virt  string `xml:"virt,attr"`
	Value string `xml:",chardata"`
}

type SectionDef struct {
	Kind string `xml:"kind,attr"`
	Id   string `xml:"id,attr"`

	Header      string    `xml:"header"`
	Description DocString `xml:"description"`

	Members []*MemberDef `xml:"memberdef"`

	Functions []*FunctionMemberDef
	Enums     []*EnumMemberDef
	Variables []*VariableMemberDef
	Defines   []*DefineMemberDef
	Typedefs  []*TypedefMemberDef
	Friends   []*FriendMemberDef
}

type MemberDef struct {
	Kind     string `xml:"kind,attr"`
	Id       string `xml:"id,attr"`
	Prot     string `xml:"prot,attr"`
	Static   string `xml:"static,attr"`
	Const    string `xml:"const,attr"`
	Explicit string `xml:"explicit,attr"`
	Inline   string `xml:"inline,attr"`
	Strong   string `xml:"strong,attr"`
	Mutable  string `xml:"mutable,attr"`

	InnerXML []byte `xml:",innerxml"`
	// InnerXMLStr string `xml:",innerxml"`
}

type MemberDefContent struct {
	Type     string    `xml:"type"`
	Name     string    `xml:"name"`
	Location *Location `xml:"location"`
}

type FunctionParam struct {
	Type     DocString `xml:"type"`
	DeclName string    `xml:"declname"`
}

type ReferencedBy struct {
	RefId       string `xml:"refid,attr"`
	CompoundRef string `xml:"compoundref,attr"`
	Name        string `xml:",chardata"`

	StartLine int `xml:"startline,attr"`
	EndLine   int `xml:"endline,attr"`
}

type Reimplements struct {
	RefId string `xml:"refid,attr"`
	Name  string `xml:",chardata"`
}

type FunctionMemberDef struct {
	MemberDef
	Descriptions

	Type     DocString `xml:"type"`
	Name     string    `xml:"name"`
	Location Location  `xml:"location"`

	Definition string          `xml:"definition"`
	ArgsString string          `xml:"argsstring"`
	Params     []FunctionParam `xml:"param"`

	Reimplements    Reimplements   `xml:"reimplements"`
	ReimplementedBy []Reimplements `xml:"reimplementedby"`
}

type EnumMemberDef struct {
	MemberDef
	Descriptions

	Type     DocString `xml:"type"`
	Name     string    `xml:"name"`
	Location Location  `xml:"location"`

	Values []EnumValue `xml:"enumvalue"`
}

type EnumValue struct {
	Descriptions

	Id          string `xml:"id,attr"`
	Protection  string `xml:"prot,attr"`
	Name        string `xml:"name"`
	Initializer string `xml:"initializer"`
}

type VariableMemberDef struct {
	MemberDef
	Descriptions

	Type     DocString `xml:"type"`
	Name     string    `xml:"name"`
	Location Location  `xml:"location"`

	Definition string    `xml:"definition"`
	ArgsString DocString `xml:"argsstring"`
}

type DefineParam struct {
	Defname string `xml:"defname"`
}

type DefineMemberDef struct {
	MemberDef
	Descriptions

	Name        string        `xml:"name"`
	Params      []DefineParam `xml:"param"`
	Initializer string        `xml:"initializer"`
	Location    Location      `xml:"location"`
}

type TypedefMemberDef struct {
	MemberDef
	Descriptions

	Name       string    `xml:"name"`
	Type       DocString `xml:"type"`
	Definition string    `xml:"definition"`
	ArgsString DocString `xml:"argsstring"`
	Location   Location  `xml:"location"`
}

type FriendMemberDef struct {
	MemberDef
	Descriptions

	Name       string    `xml:"name"`
	Type       DocString `xml:"type"`
	Definition string    `xml:"definition"`
	Location   Location  `xml:"location"`

	ReferencedBy []ReferencedBy `xml:"referencedby"`
}

func (ty *Section) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "id":
			ty.Id = attr.Value
		case "kind":
			ty.Kind = attr.Value
		default:
			return errors.New(fmt.Sprintf("unknown section attribute: %s", attr.Name.Local))
		}
	}

	return ty.Content.UnmarshalXML(dec, start)
}

func (ty *TableEntry) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "thead":
			if attr.Value == "yes" {
				ty.TableHead = true
			} else if attr.Value == "no" {
				ty.TableHead = false
			} else {
				return errors.New(fmt.Sprintf("unknown boolean format: %s", attr.Value))
			}
		default:
			return errors.New(fmt.Sprintf("unknown section attribute: %s", attr.Name.Local))
		}
	}

	return ty.Content.UnmarshalXML(dec, start)
}

func (ty *Ref) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "refid":
			ty.RefId = attr.Value
		case "kindref":
			ty.KindRef = attr.Value
		default:
			return errors.New(fmt.Sprintf("unknown ref attribute: %s", attr.Name.Local))
		}
	}

	return ty.Content.UnmarshalXML(dec, start)
}

func (ty *Paragraph) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		default:
			return errors.New(fmt.Sprintf("unknown paragraph attribute: %s", attr.Name.Local))
		}
	}

	return ty.Content.UnmarshalXML(dec, start)
}

func (ty *Title) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		default:
			return errors.New(fmt.Sprintf("unknown title attribute: %s", attr.Name.Local))
		}
	}

	return ty.Content.UnmarshalXML(dec, start)
}

func (ty *Heading) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	var err error
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "level":
			ty.Level, err = strconv.Atoi(attr.Value)
			if err != nil {
				return err
			}
		default:
			return errors.New(fmt.Sprintf("unknown heading attribute: %s", attr.Name.Local))
		}
	}

	return ty.Content.UnmarshalXML(dec, start)
}

func (ty *Bold) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		default:
			return errors.New(fmt.Sprintf("unknown bold attribute: %s", attr.Name.Local))
		}
	}

	return ty.Content.UnmarshalXML(dec, start)
}

func (ty *Emphasis) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		default:
			return errors.New(fmt.Sprintf("unknown emphasis attribute: %s", attr.Name.Local))
		}
	}

	return ty.Content.UnmarshalXML(dec, start)
}

func (ty *Verbatim) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		default:
			return errors.New(fmt.Sprintf("unknown verbatim attribute: %s", attr.Name.Local))
		}
	}

	return ty.Content.UnmarshalXML(dec, start)
}

func (ty *Preformatted) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		default:
			return errors.New(fmt.Sprintf("unknown preformatted attribute: %s", attr.Name.Local))
		}
	}

	return ty.Content.UnmarshalXML(dec, start)
}

func (ty *ComputerOutput) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		default:
			return errors.New(fmt.Sprintf("unknown computeroutput attribute: %s", attr.Name.Local))
		}
	}

	return ty.Content.UnmarshalXML(dec, start)
}

func (ty *Term) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		default:
			return errors.New(fmt.Sprintf("unknown term attribute: %s", attr.Name.Local))
		}
	}

	return ty.Content.UnmarshalXML(dec, start)
}

func (ty *ProgramListing) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "filename":
			ty.Filename = attr.Value
		default:
			return errors.New(fmt.Sprintf("unknown type attribute: %s", attr.Name.Local))
		}
	}

	return ty.Content.UnmarshalXML(dec, start)
}

func (ty *DocString) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	for {
		t, err := dec.Token()
		if err != nil {
			return err
		}
		switch tt := t.(type) {
		case xml.StartElement:
			switch tt.Name.Local {
			case "ref":
				var r Ref
				err = dec.DecodeElement(&r, &tt)
				if err != nil {
					return err
				}
				ty.Content = append(ty.Content, r)
			case "para":
				var p Paragraph
				err = dec.DecodeElement(&p, &tt)
				if err != nil {
					return err
				}
				ty.Content = append(ty.Content, p)
			case "title":
				var t Title
				err = dec.DecodeElement(&t, &tt)
				if err != nil {
					return err
				}
				ty.Content = append(ty.Content, t)
			case "heading":
				var t Heading
				err = dec.DecodeElement(&t, &tt)
				if err != nil {
					return err
				}
				ty.Content = append(ty.Content, t)
			// TODO: Different sectnumbers?
			case "sect1":
				fallthrough
			case "sect2":
				fallthrough
			case "simplesect":
				var s Section
				err = dec.DecodeElement(&s, &tt)
				if err != nil {
					return err
				}
				ty.Content = append(ty.Content, s)
			case "table":
				var t Table
				err = dec.DecodeElement(&t, &tt)
				if err != nil {
					return err
				}
				ty.Content = append(ty.Content, t)
			case "parameterlist":
				var p ParameterList
				err = dec.DecodeElement(&p, &tt)
				if err != nil {
					return err
				}
				ty.Content = append(ty.Content, p)
			case "variablelist":
				var p VariableList
				err = dec.DecodeElement(&p, &tt)
				if err != nil {
					return err
				}
				ty.Content = append(ty.Content, p)
			case "xrefsect":
				var x XRefSect
				err = dec.DecodeElement(&x, &tt)
				if err != nil {
					return err
				}
				ty.Content = append(ty.Content, x)
			case "itemizedlist":
				var l ItemizedList
				err = dec.DecodeElement(&l, &tt)
				if err != nil {
					return err
				}
				ty.Content = append(ty.Content, l)
			case "orderedlist":
				var l OrderedList
				err = dec.DecodeElement(&l, &tt)
				if err != nil {
					return err
				}
				ty.Content = append(ty.Content, l)
			case "bold":
				var b Bold
				err = dec.DecodeElement(&b, &tt)
				if err != nil {
					return err
				}
				ty.Content = append(ty.Content, b)
			case "emphasis":
				var b Emphasis
				err = dec.DecodeElement(&b, &tt)
				if err != nil {
					return err
				}
				ty.Content = append(ty.Content, b)
			case "verbatim":
				var v Verbatim
				err = dec.DecodeElement(&v, &tt)
				if err != nil {
					return err
				}
				ty.Content = append(ty.Content, v)
			case "preformatted":
				var p Preformatted
				err = dec.DecodeElement(&p, &tt)
				if err != nil {
					return err
				}
				ty.Content = append(ty.Content, p)
			case "computeroutput":
				var c ComputerOutput
				err = dec.DecodeElement(&c, &tt)
				if err != nil {
					return err
				}
				ty.Content = append(ty.Content, c)
			case "term":
				var c Term
				err = dec.DecodeElement(&c, &tt)
				if err != nil {
					return err
				}
				ty.Content = append(ty.Content, c)
			case "linebreak":
				var l Linebreak
				err = dec.DecodeElement(&l, &tt)
				if err != nil {
					return err
				}
				ty.Content = append(ty.Content, l)
			case "anchor":
				var l Anchor
				err = dec.DecodeElement(&l, &tt)
				if err != nil {
					return err
				}
				ty.Content = append(ty.Content, l)
			case "image":
				var i Image
				err = dec.DecodeElement(&i, &tt)
				if err != nil {
					return err
				}
				ty.Content = append(ty.Content, i)
			case "programlisting":
				var p ProgramListing
				err = dec.DecodeElement(&p, &tt)
				if err != nil {
					return err
				}
				ty.Content = append(ty.Content, p)
			case "ndash":
				var t Text
				t.Content = "\u2013"
				ty.Content = append(ty.Content, t)
			case "hruler":
				//TODO
				var t Text
				t.Content = "---------------------------------"
				ty.Content = append(ty.Content, t)
			case "ulink":
				//TODO
				var t Text
				t.Content = "<href>"
				ty.Content = append(ty.Content, t)
			case "zwj":
				var t Text
				t.Content = "\u200D"
				ty.Content = append(ty.Content, t)
			case "sp":
				var t Text
				t.Content = " "
				ty.Content = append(ty.Content, t)
			case "codeline":
			case "highlight":
			default:
				return errors.New(fmt.Sprintf("unknown token `%s` in docstring element", tt.Name.Local))
			}
		case xml.CharData:
			var t Text
			t.Content = string(tt)
			ty.Content = append(ty.Content, t)
		case xml.EndElement:
			if tt == start.End() {
				return nil
			}
		default:
			return errors.New(fmt.Sprintf("unknown token type %v in docstring element", t))
		}
	}
}

func (sec *SectionDef) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		switch attr.Name.Local {
		case "kind":
			sec.Kind = attr.Value
		case "id":
			sec.Kind = attr.Value
		default:
			return errors.New("unknown section attribute")
		}
	}

	for {
		t, err := dec.Token()
		if err != nil {
			return err
		}

		switch tt := t.(type) {
		case xml.StartElement:
			switch tt.Name.Local {
			case "memberdef":
				var kind string
				for _, attr := range tt.Attr {
					if attr.Name.Local == "kind" {
						kind = attr.Value
					}
				}
				if kind == "" {
					return errors.New("missing kind on memberdef")
				}

				switch kind {
				case FunctionMember:
					var f FunctionMemberDef
					err = dec.DecodeElement(&f, &tt)
					if err != nil {
						return err
					}
					sec.Functions = append(sec.Functions, &f)
				case EnumMember:
					var e EnumMemberDef
					err = dec.DecodeElement(&e, &tt)
					if err != nil {
						return err
					}
					sec.Enums = append(sec.Enums, &e)
				case VariableMember:
					var v VariableMemberDef
					err = dec.DecodeElement(&v, &tt)
					if err != nil {
						return err
					}
					sec.Variables = append(sec.Variables, &v)
				case DefineMember:
					var d DefineMemberDef
					err = dec.DecodeElement(&d, &tt)
					if err != nil {
						return err
					}
					sec.Defines = append(sec.Defines, &d)
				case TypedefMember:
					var t TypedefMemberDef
					err = dec.DecodeElement(&t, &tt)
					if err != nil {
						return err
					}
					sec.Typedefs = append(sec.Typedefs, &t)
				case FriendMember:
					var f FriendMemberDef
					err = dec.DecodeElement(&f, &tt)
					if err != nil {
						return err
					}
					sec.Friends = append(sec.Friends, &f)
				default:
					return errors.New(fmt.Sprintf("unknown member kind: %s", kind))
				}
			case "header":
				var h string
				err = dec.DecodeElement(&h, &tt)
				if err != nil {
					return err
				}
				sec.Header = h
			case "description":
				var d DocString
				err = dec.DecodeElement(&d, &tt)
				if err != nil {
					return err
				}
				sec.Description = d
			default:
				return errors.New(fmt.Sprintf("unknown section element: %s", tt.Name.Local))
			}
		case xml.EndElement:
			if tt == start.End() {
				return nil
			}
		}
	}
}

func ParseDoxygenFolder(path string) []*Doxygen {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	doxygenDocs := make([]*Doxygen, 0, len(files))

	for _, f := range files {
		d, err := readFile(fmt.Sprintf("%s/%s", path, f.Name()))

		if err != nil {
			log.Printf("Error reading file (%s): %v\n", f.Name(), err)
		}
		if d != nil {
			doxygenDocs = append(doxygenDocs, d)
		}
	}
	return doxygenDocs
}

func readFile(path string) (*Doxygen, error) {
	if !strings.HasSuffix(path, ".xml") {
		return nil, nil
	}
	if strings.HasSuffix(path, "index.xml") {
		return nil, nil
	}

	data, err := xmlhelper.ReadFileWithBadUTF8(fmt.Sprintf("%s", path))
	if err != nil {
		log.Fatal(fmt.Sprintf("Error parsing file (%s): %v\n", path, err))
	}

	doxygen := &Doxygen{}
	err = xml.Unmarshal(data, doxygen)
	if err != nil {
		log.Printf("%s", string(data))
		log.Fatal(fmt.Sprintf("Error parsing file (%s): %v\n", path, err))
	}
	return doxygen, nil
}
