package goxy

import (
	"ScriptExecServer/pkg/doxygen"
	"errors"
	"fmt"
	"log"
	"strings"
)

func SectionsFromDoxygen(d *doxygen.Doxygen) ([]*SectionDoc, error) {
	sections := make([]*SectionDoc, 0)
	for _, section := range d.CompoundDef.Sections {
		s, err := SectionFromDoxygen(section)
		if err != nil {
			return nil, err
		}

		sections = append(sections, s)
	}

	return sections, nil
}

func CompoundFromDoxygen(d *doxygen.Doxygen) (*CompoundDoc, error) {
	var err error
	compound := &CompoundDoc{}

	compound.Name = d.CompoundDef.CompoundName
	compound.Title = d.CompoundDef.Title
	if compound.Title == "" {
		compound.Title = compound.Name
	}
	compound.Id = strings.ToLower(d.CompoundDef.Id)
	compound.Kind, err = KindFromDoxygen(d.CompoundDef.Kind)
	if err != nil {
		return nil, err
	}
	compound.Protection, err = ProtectionFromDoxygen(d.CompoundDef.Protection)
	if err != nil {
		return nil, err
	}
	compound.Location = LocationFromDoxygen(d.CompoundDef.Location)
	if d.CompoundDef.ProgramListing != nil {
		compound.ProgramListing, err = DocStringFromDoxygen(d.CompoundDef.ProgramListing.Content)
		if err != nil {
			return nil, err
		}
	}
	compound.Sections = make([]*SectionDoc, 0)
	for _, section := range d.CompoundDef.Sections {
		s, err := SectionFromDoxygen(section)
		if err != nil {
			return nil, err
		}

		if s.Header == "" {
			InferSectionHeader(compound.Kind, s)
		}

		compound.Sections = append(compound.Sections, s)
	}

	compound.InnerClasses = make([]InnerCompoundRef, 0)
	for _, ic := range d.CompoundDef.InnerClass {
		ic, err := InnerCompoundRefFromDoxygen(ic)
		if err != nil {
			return nil, err
		}

		compound.InnerClasses = append(compound.InnerClasses, ic)
	}

	compound.InnerFiles = make([]InnerCompoundRef, 0)
	for _, ic := range d.CompoundDef.InnerFiles {
		ic, err := InnerCompoundRefFromDoxygen(ic)
		if err != nil {
			return nil, err
		}

		compound.InnerFiles = append(compound.InnerFiles, ic)
	}

	compound.InnerDirs = make([]InnerCompoundRef, 0)
	for _, ic := range d.CompoundDef.InnerDirs {
		ic, err := InnerCompoundRefFromDoxygen(ic)
		if err != nil {
			return nil, err
		}

		compound.InnerDirs = append(compound.InnerDirs, ic)
	}

	compound.InnerGroups = make([]InnerCompoundRef, 0)
	for _, ic := range d.CompoundDef.InnerGroups {
		ic, err := InnerCompoundRefFromDoxygen(ic)
		if err != nil {
			return nil, err
		}

		compound.InnerGroups = append(compound.InnerGroups, ic)
	}

	compound.InnerNamespaces = make([]InnerCompoundRef, 0)
	for _, ic := range d.CompoundDef.InnerNamespaces {
		ic, err := InnerCompoundRefFromDoxygen(ic)
		if err != nil {
			return nil, err
		}

		compound.InnerNamespaces = append(compound.InnerNamespaces, ic)
	}

	compound.Descriptions, err = DescriptionsFromDoxygen(d.CompoundDef.Descriptions)
	if err != nil {
		return nil, err
	}

	return compound, nil
}

func InferSectionHeader(kind Kind, s *SectionDoc) {
	kindHeader := "Unknown"
	switch s.Kind {
	case Friends:
		kindHeader = "Friends"
	case Attributes:
		kindHeader = "Attributes"
	case StaticAttributes:
		kindHeader = "Static Attributes"
	case Functions:
		kindHeader = "Functions"
	case StaticFunctions:
		kindHeader = "Static Functions"
	case Types:
		kindHeader = "Types"
	case UserDefined:
		kindHeader = "User Defined"
	case Defines:
		kindHeader = "Defines"
	case Typedefs:
		kindHeader = "Typedefs"
	case Variables:
		kindHeader = "Variables"
	case Enums:
		kindHeader = "Enumerations"
	case Related:
		kindHeader = "Related"
	default:
		log.Fatal("unknown kind: ", s.Kind)
	}

	switch kind {
	case Group:
		s.Header = kindHeader
	case Class:
		fallthrough
	default:
		var protHeader string
		switch s.Protection {
		case Public:
			protHeader = "Public"
		case Protected:
			protHeader = "Protected"
		case Private:
			protHeader = "Private"
		}
		s.Header = fmt.Sprintf("%s %s", protHeader, kindHeader)
	}
}

func InnerCompoundRefFromDoxygen(ic doxygen.InnerCompound) (InnerCompoundRef, error) {
	var res InnerCompoundRef
	var err error
	res.RefId = strings.ToLower(ic.RefId)
	res.Value = ic.Value
	res.Protection, err = ProtectionFromDoxygen(ic.Prot)
	return res, err
}

func KindFromDoxygen(kind string) (Kind, error) {
	switch kind {
	case "class":
		return Class, nil
	case "file":
		return File, nil
	case "struct":
		return Class, nil
	case "namespace":
		return Namespace, nil
	case "group":
		return Group, nil
	case "union":
		return Union, nil
	case "dir":
		return Dir, nil
	case "page":
		return Page, nil
	default:
		return "", errors.New(fmt.Sprintf("unable to convert kind string, unknown value: %s", kind))
	}
}

func SectionFromDoxygen(section doxygen.SectionDef) (*SectionDoc, error) {
	var err error
	s := &SectionDoc{}

	s.Header = section.Header
	s.Description, err = DocStringFromDoxygen(section.Description)
	if err != nil {
		return nil, err
	}

	s.Kind, err = SectionKindFromDoxygen(section.Kind)
	if err != nil {
		firstDash := strings.Index(section.Kind, "-")
		if firstDash > 0 {
			s.Protection, err = ProtectionFromDoxygen(section.Kind[:firstDash])
			if err != nil {
				return nil, err
			}
			s.Kind, err = SectionKindFromDoxygen(section.Kind[firstDash+1:])
			if err != nil {
				return nil, err
			}
		}
	}

	s.Id = strings.ToLower(section.Id)

	for _, function := range section.Functions {
		f := &FunctionDoc{}
		f.Id = strings.ToLower(function.Id)
		f.Name = function.Name
		f.Protection, err = ProtectionFromDoxygen(function.Prot)
		if err != nil {
			return nil, err
		}
		f.Location = LocationFromDoxygen(function.Location)
		f.Type, err = DocStringFromDoxygen(function.Type)
		if err != nil {
			return nil, err
		}

		f.Type, err = DocStringFromDoxygen(function.Type)
		if err != nil {
			return nil, err
		}
		f.Definition = function.Definition
		f.ArgsString = function.ArgsString
		f.Params, err = FunctionParamsFromDoxygen(function.Params)
		if err != nil {
			return nil, err
		}
		f.Descriptions, err = DescriptionsFromDoxygen(function.Descriptions)
		if err != nil {
			return nil, err
		}

		if function.Reimplements.RefId != "" {
			fRefId := strings.ToLower(function.Reimplements.RefId)
			f.Reimplements = Reimplements{
				RefId:    fRefId,
				MemberId: strings.Split(fRefId, "_")[0],
				ParentId: strings.Split(fRefId, "_")[1],
			}
		}
		f.ReimplementedBy = make([]Reimplements, 0)
		for _, r := range function.ReimplementedBy {
			rRefId := strings.ToLower(r.RefId)
			f.ReimplementedBy = append(f.ReimplementedBy, Reimplements{
				RefId:    rRefId,
				MemberId: strings.Split(rRefId, "_")[0],
				ParentId: strings.Split(rRefId, "_")[1],
			})
		}

		s.Functions = append(s.Functions, f)
	}

	for _, enum := range section.Enums {
		e := &EnumDoc{}
		e.Id = strings.ToLower(enum.Id)
		e.Name = enum.Name
		e.Protection, err = ProtectionFromDoxygen(enum.Prot)
		if err != nil {
			return nil, err
		}
		e.Location = LocationFromDoxygen(enum.Location)
		e.Type, err = DocStringFromDoxygen(enum.Type)
		if err != nil {
			return nil, err
		}
		e.Values, _ = EnumValuesFromDoxygen(enum.Values)
		e.Descriptions, err = DescriptionsFromDoxygen(enum.Descriptions)
		if err != nil {
			return nil, err
		}

		s.Enums = append(s.Enums, e)
	}

	for _, variable := range section.Variables {
		a := &ClassAttributeDoc{}
		a.Id = strings.ToLower(variable.Id)
		a.Name = variable.Name
		a.Protection, err = ProtectionFromDoxygen(variable.Prot)
		if err != nil {
			return nil, err
		}
		a.Location = LocationFromDoxygen(variable.Location)
		a.Type, err = DocStringFromDoxygen(variable.Type)
		if err != nil {
			return nil, err
		}
		a.Definition = variable.Definition
		a.Descriptions, err = DescriptionsFromDoxygen(variable.Descriptions)
		if err != nil {
			return nil, err
		}
		a.ArgsString, err = DocStringFromDoxygen(variable.ArgsString)
		if err != nil {
			return nil, err
		}

		s.Attributes = append(s.Attributes, a)
	}

	for _, define := range section.Defines {
		d := &DefineDoc{}
		d.Id = strings.ToLower(define.Id)
		d.Name = define.Name
		d.Initializer = define.Initializer
		d.Params, err = DefineParamsFromDoxygen(define.Params)
		if err != nil {
			return nil, err
		}
		d.Descriptions, err = DescriptionsFromDoxygen(define.Descriptions)
		if err != nil {
			return nil, err
		}

		s.Defines = append(s.Defines, d)
	}

	for _, typedef := range section.Typedefs {
		d := &TypedefDoc{}
		d.Id = strings.ToLower(typedef.Id)
		d.Name = typedef.Name
		d.Definition = typedef.Definition
		d.Type, err = DocStringFromDoxygen(typedef.Type)
		if err != nil {
			return nil, err
		}
		d.Descriptions, err = DescriptionsFromDoxygen(typedef.Descriptions)
		if err != nil {
			return nil, err
		}
		d.ArgsString, err = DocStringFromDoxygen(typedef.ArgsString)
		if err != nil {
			return nil, err
		}

		s.Typedefs = append(s.Typedefs, d)
	}

	for _, friend := range section.Friends {
		f := &FriendDoc{}
		f.Id = strings.ToLower(friend.Id)
		f.Name = friend.Name
		f.Definition = friend.Definition
		f.Type, err = DocStringFromDoxygen(friend.Type)
		if err != nil {
			return nil, err
		}
		f.Descriptions, err = DescriptionsFromDoxygen(friend.Descriptions)
		if err != nil {
			return nil, err
		}

		s.Friends = append(s.Friends, f)
	}

	return s, nil
}

func FunctionParamsFromDoxygen(params []doxygen.FunctionParam) ([]FunctionParam, error) {
	var err error
	r := make([]FunctionParam, len(params))

	for i, param := range params {
		r[i] = FunctionParam{
			DeclName: param.DeclName,
		}
		r[i].Type, err = DocStringFromDoxygen(param.Type)
		if err != nil {
			return nil, err
		}
	}
	return r, nil
}

func DefineParamsFromDoxygen(params []doxygen.DefineParam) ([]DefineParam, error) {
	r := make([]DefineParam, len(params))

	for i, param := range params {
		r[i] = DefineParam{
			Defname: param.Defname,
		}
	}
	return r, nil
}

func EnumValuesFromDoxygen(values []doxygen.EnumValue) ([]EnumValue, error) {
	var err error
	r := make([]EnumValue, len(values))

	for i, value := range values {
		r[i] = EnumValue{
			Name:        value.Name,
			Id:          strings.ToLower(value.Id),
			Initializer: value.Initializer,
		}
		r[i].Descriptions, err = DescriptionsFromDoxygen(value.Descriptions)
		if err != nil {
			return nil, err
		}
		r[i].Protection, err = ProtectionFromDoxygen(value.Protection)
		if err != nil {
			return nil, err
		}
	}

	return r, nil
}

func DescriptionsFromDoxygen(d doxygen.Descriptions) (Descriptions, error) {
	var err error

	r := Descriptions{}
	r.BriefDescription, err = DocStringFromDoxygen(d.BriefDescription)
	if err != nil {
		return Descriptions{}, err
	}
	r.DetailedDescription, err = DocStringFromDoxygen(d.DetailedDescription)
	if err != nil {
		return Descriptions{}, err
	}
	r.InBodyDescription, err = DocStringFromDoxygen(d.InBodyDescription)
	if err != nil {
		return Descriptions{}, err
	}

	return r, nil
}

func DocStringFromDoxygen(t doxygen.DocString) (DocString, error) {
	var err error
	parts := make([]DocStringElement, len(t.Content))
	for i, c := range t.Content {
		switch cc := c.(type) {
		case doxygen.Ref:
			e := DocStringRef{
				RefId:   strings.ToLower(cc.RefId),
				KindRef: cc.KindRef,
			}
			e.Content, err = DocStringFromDoxygen(cc.Content)
			if err != nil {
				return DocString{}, err
			}

			parts[i] = DocStringElement{
				Type:  Ref,
				Value: e,
			}
		case doxygen.Anchor:
			parts[i] = DocStringElement{
				Type:  Anchor,
				Value: DocStringAnchor{strings.ToLower(cc.Id)},
			}
		case doxygen.Image:
			parts[i] = DocStringElement{
				Type:  Image,
				Value: DocStringImage{
					Name:        cc.Name,
					Type:        cc.Type,
					Description: cc.Description,
				},
			}
		case doxygen.Text:
			parts[i] = DocStringElement{
				Type:  Text,
				Value: DocStringText{cc.Content},
			}
		case doxygen.Section:
			s := DocStringSection{
				Id:   strings.ToLower(cc.Id),
				Kind: cc.Kind,
			}
			s.Content, err = DocStringFromDoxygen(cc.Content)
			if err != nil {
				return DocString{}, err
			}
			parts[i] = DocStringElement{
				Type:  Section,
				Value: s,
			}
		case doxygen.Paragraph:
			p := DocStringParagraph{}
			p.Content, err = DocStringFromDoxygen(cc.Content)
			if err != nil {
				return DocString{}, err
			}
			parts[i] = DocStringElement{
				Type:  Paragraph,
				Value: p,
			}
		case doxygen.Title:
			t := DocStringTitle{}
			t.Content, err = DocStringFromDoxygen(cc.Content)
			if err != nil {
				return DocString{}, err
			}
			parts[i] = DocStringElement{
				Type:  Title,
				Value: t,
			}
		case doxygen.Heading:
			t := DocStringHeading{}
			t.Content, err = DocStringFromDoxygen(cc.Content)
			if err != nil {
				return DocString{}, err
			}
			t.Level = cc.Level
			parts[i] = DocStringElement{
				Type:  Heading,
				Value: t,
			}
		case doxygen.XRefSect:
			t := DocStringXRefSect{
				Id:    strings.ToLower(cc.Id),
				Title: cc.Title,
			}
			t.Description, err = DocStringFromDoxygen(cc.Description)
			if err != nil {
				return DocString{}, err
			}
			parts[i] = DocStringElement{
				Type:  XRefSect,
				Value: t,
			}
		case doxygen.ItemizedList:
			t := DocStringItemizedList{
				Items: make([]DocString, len(cc.Items)),
			}
			for i, item := range cc.Items {
				t.Items[i], err = DocStringFromDoxygen(item)
				if err != nil {
					return DocString{}, err
				}
			}
			parts[i] = DocStringElement{
				Type:  ItemizedList,
				Value: t,
			}
		case doxygen.OrderedList:
			t := DocStringOrderedList{
				Items: make([]DocString, len(cc.Items)),
			}
			for i, item := range cc.Items {
				t.Items[i], err = DocStringFromDoxygen(item)
				if err != nil {
					return DocString{}, err
				}
			}
			parts[i] = DocStringElement{
				Type:  OrderedList,
				Value: t,
			}
		case doxygen.VariableList:
			t := DocStringVariableList{
				Items: make([]DocString, len(cc.Items)),
			}
			for i, item := range cc.Items {
				t.Items[i], err = DocStringFromDoxygen(item)
				if err != nil {
					return DocString{}, err
				}
			}
			parts[i] = DocStringElement{
				Type:  VariableList,
				Value: t,
			}
		case doxygen.ParameterList:
			l := DocStringParameterList{
				Kind:  cc.Kind,
				Items: make([]DocStringParameterItem, len(cc.Items)),
			}

			for i, item := range cc.Items {
				l.Items[i] = DocStringParameterItem{
					Name: item.Name,
				}
				l.Items[i].Description, err = DocStringFromDoxygen(item.Description)
				if err != nil {
					return DocString{}, err
				}
			}
			parts[i] = DocStringElement{
				Type:  ParameterList,
				Value: l,
			}
		case doxygen.Table:
			t := DocStringTable{}

			t.Rows = make([][]DocStringTableEntry, len(cc.Rows))
			for idx, row := range cc.Rows {
				t.Rows[idx] = make([]DocStringTableEntry, len(row.Columns))
				for jdx, col := range row.Columns {
					t.Rows[idx][jdx] = DocStringTableEntry{
						Head:    col.TableHead,
					}
					t.Rows[idx][jdx].Content, err = DocStringFromDoxygen(col.Content)
					if err != nil {
						return DocString{}, err
					}
				}
			}

			parts[i] = DocStringElement{
				Type:  Table,
				Value: t,
			}
		case doxygen.Bold:
			t := DocStringBold{}
			t.Content, err = DocStringFromDoxygen(cc.Content)
			if err != nil {
				return DocString{}, err
			}
			parts[i] = DocStringElement{
				Type:  Bold,
				Value: t,
			}
		case doxygen.Verbatim:
			t := DocStringVerbatim{}
			t.Content, err = DocStringFromDoxygen(cc.Content)
			if err != nil {
				return DocString{}, err
			}
			parts[i] = DocStringElement{
				Type:  Verbatim,
				Value: t,
			}
		case doxygen.Emphasis:
			t := DocStringEmphasis{}
			t.Content, err = DocStringFromDoxygen(cc.Content)
			if err != nil {
				return DocString{}, err
			}
			parts[i] = DocStringElement{
				Type:  Emphasis,
				Value: t,
			}
		case doxygen.Term:
			t := DocStringTerm{}
			t.Content, err = DocStringFromDoxygen(cc.Content)
			if err != nil {
				return DocString{}, err
			}
			parts[i] = DocStringElement{
				Type:  Term,
				Value: t,
			}
		case doxygen.Preformatted:
			t := DocStringPreformatted{}
			t.Content, err = DocStringFromDoxygen(cc.Content)
			if err != nil {
				return DocString{}, err
			}
			parts[i] = DocStringElement{
				Type:  Preformatted,
				Value: t,
			}
		case doxygen.ComputerOutput:
			t := DocStringComputerOutput{}
			t.Content, err = DocStringFromDoxygen(cc.Content)
			if err != nil {
				return DocString{}, err
			}
			parts[i] = DocStringElement{
				Type:  ComputerOutput,
				Value: t,
			}
		case doxygen.ProgramListing:
			h := DocStringHighlight{}
			h.Content, err = DocStringFromDoxygen(cc.Content)
			if err != nil {
				return DocString{}, err
			}
			parts[i] = DocStringElement{
				Type:  Highlight,
				Value: h,
			}
		case doxygen.Linebreak:
			parts[i] = DocStringElement{
				Type:  LineBreak,
				Value: DocStringLinebreak{},
			}
		default:
			return DocString{}, errors.New(fmt.Sprintf("unable to convert doc string, unknown type: %T", cc))
		}
	}

	return DocString{
		Content: parts,
	}, nil
}

func LocationFromDoxygen(location doxygen.Location) SourceLocation {
	return SourceLocation{
		File:      location.File,
		Line:      location.Line,
		Column:    location.Column,
		BodyFile:  location.BodyFile,
		BodyStart: location.BodyStart,
		BodyEnd:   location.BodyEnd,
	}
}

func SectionKindFromDoxygen(kind string) (SectionKind, error) {
	switch kind {
	case "attrib":
		return Attributes, nil
	case "var":
		return Variables, nil
	case "function":
		fallthrough
	case "func":
		return Functions, nil
	case "static-func":
		return StaticFunctions, nil
	case "static-attrib":
		return StaticAttributes, nil
	case "type":
		return Types, nil
	case "friend":
		return Friends, nil
	case "user-defined":
		return UserDefined, nil
	case "define":
		return Defines, nil
	case "typedef":
		return Typedefs, nil
	case "enum":
		return Enums, nil
	case "related":
		return Related, nil
	default:
		return -1, errors.New(fmt.Sprintf("unable to convert sectionkind string, unknown value: %s", kind))
	}
}

func ProtectionFromDoxygen(protection string) (Protection, error) {
	switch protection {
	case "public":
		return Public, nil
	case "private":
		return Private, nil
	case "protected":
		return Protected, nil
	case "":
		return -1, nil
	default:
		return -1, errors.New(fmt.Sprintf("unable to convert protection string, unknown value: %s", protection))
	}
}
