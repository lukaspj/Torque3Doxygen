package main

import (
	"ScriptExecServer/pkg/doxygen"
	"ScriptExecServer/pkg/goxy"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type GoxygenEntity struct {
	Kind   string
	Entity interface{}
}

type GoxygenRef struct {
	Kind      string
	Name      string
	RefId     string
	ParentRef string
}

type GoxygenData struct {
	Entities map[string]GoxygenEntity
	Refs     map[string]GoxygenRef
}

type GeekdocBundleMenuItem struct {
	Name string                  `yaml:"name,omitempty"`
	Ref  string                  `yaml:"ref,omitempty"`
	Icon string                  `yaml:"icon,omitempty"`
	Sub  []GeekdocBundleMenuItem `yaml:"sub,omitempty"`
}

type GeekdocBundleMenu struct {
	Menus map[string][]GeekdocBundleMenuItem
}

func main() {
	docs := doxygen.ParseDoxygenFolder()

	compounds := make([]*goxy.CompoundDoc, 0)

	data := GoxygenData{
		Entities: make(map[string]GoxygenEntity),
		Refs:     make(map[string]GoxygenRef),
	}
	for _, doc := range docs {
		compound, err := goxy.CompoundFromDoxygen(doc)
		if err != nil {
			fmt.Println(fmt.Sprintf("unable to parse doxygen compound doc: %v, due to: %v", doc.CompoundDef.CompoundName, err))
		} else {
			compounds = append(compounds, compound)
		}
	}

	for _, compound := range compounds {
		for _, class := range compound.InnerClasses {
			for _, inner := range compounds {
				if inner.Id == class.RefId {
					inner.Parent = compound.Id
				}
			}
		}
		for _, namespace := range compound.InnerNamespaces {
			for _, inner := range compounds {
				if inner.Id == namespace.RefId {
					inner.Parent = compound.Id
				}
			}
		}
		for _, group := range compound.InnerGroups {
			for _, inner := range compounds {
				if inner.Id == group.RefId {
					inner.Parent = compound.Id
				}
			}
		}
		for _, file := range compound.InnerFiles {
			for _, inner := range compounds {
				if inner.Id == file.RefId {
					inner.Parent = compound.Id
				}
			}
		}
		for _, dir := range compound.InnerDirs {
			for _, inner := range compounds {
				if inner.Id == dir.RefId {
					inner.Parent = compound.Id
				}
			}
		}
	}

	// COMPOUNDS
	for _, compound := range compounds {
		mdContent := fmt.Sprintf(`---
goxygen_id: "%s"
GeekdocFlatSection: true
title: "%s"
type: "%s"
---
`, compound.Id, compound.Title, "compound")

		os.MkdirAll(fmt.Sprintf("hugo/content/%s", compound.Kind), 0644)
		ioutil.WriteFile(fmt.Sprintf("hugo/content/%s/%s.md", compound.Kind, compound.Id), []byte(mdContent), 0644)
		data.Entities[compound.Id] = GoxygenEntity{
			Kind:   string(compound.Kind),
			Entity: compound,
		}

		data.Refs[compound.Id] = GoxygenRef{
			Kind:      string(compound.Kind),
			Name:      compound.Name,
			ParentRef: "N/D",
			RefId:     compound.Id,
		}

		AddRefsFromDescriptions(data.Refs, compound.Id, compound.Descriptions)

		for _, section := range compound.Sections {
			AddRefsFromDocstring(data.Refs, compound.Id, section.Description)

			for _, function := range section.Functions {
				data.Refs[function.Id] = GoxygenRef{
					Kind:      "function",
					Name:      function.Name,
					ParentRef: compound.Id,
					RefId:     function.Id,
				}

				AddRefsFromDescriptions(data.Refs, compound.Id, function.Descriptions)
			}
			for _, enum := range section.Enums {
				data.Refs[enum.Id] = GoxygenRef{
					Kind:      "enum",
					Name:      enum.Name,
					ParentRef: compound.Id,
					RefId:     enum.Id,
				}

				for _, value := range enum.Values {
					data.Refs[value.Id] = GoxygenRef{
						Kind:      "enumvalue",
						Name:      value.Name,
						ParentRef: compound.Id,
						RefId:     value.Id,
					}
				}

				AddRefsFromDescriptions(data.Refs, compound.Id, enum.Descriptions)
			}
			for _, attr := range section.Attributes {
				data.Refs[attr.Id] = GoxygenRef{
					Kind:      "attribute",
					Name:      attr.Name,
					ParentRef: compound.Id,
					RefId:     attr.Id,
				}

				AddRefsFromDescriptions(data.Refs, compound.Id, attr.Descriptions)
			}
			for _, def := range section.Defines {
				data.Refs[def.Id] = GoxygenRef{
					Kind:      "define",
					Name:      def.Name,
					ParentRef: compound.Id,
					RefId:     def.Id,
				}

				AddRefsFromDescriptions(data.Refs, compound.Id, def.Descriptions)
			}
			for _, typedef := range section.Typedefs {
				data.Refs[typedef.Id] = GoxygenRef{
					Kind:      "typedef",
					Name:      typedef.Name,
					ParentRef: compound.Id,
					RefId:     typedef.Id,
				}

				AddRefsFromDescriptions(data.Refs, compound.Id, typedef.Descriptions)
			}
			for _, friend := range section.Friends {
				data.Refs[friend.Id] = GoxygenRef{
					Kind:      "friend",
					Name:      friend.Name,
					ParentRef: compound.Id,
					RefId:     friend.Id,
				}

				AddRefsFromDescriptions(data.Refs, compound.Id, friend.Descriptions)
			}
		}
	}

	os.MkdirAll("hugo/data", 0644)
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error: %v", errors.WithStack(err))
	}
	err = ioutil.WriteFile("hugo/data/goxygen.json", bytes, 0644)
	if err != nil {
		log.Fatalf("Error: %v", errors.WithStack(err))
	}

	bytes, err = json.Marshal(compounds[0:10])
	if err != nil {
		log.Fatalf("Error: %v", errors.WithStack(err))
	}
	err = ioutil.WriteFile("goxygen_excerpt_compounds.json", bytes, 0644)
	if err != nil {
		log.Fatalf("Error: %v", errors.WithStack(err))
	}

	pages := make([]GeekdocBundleMenuItem, 0)
	for _, compound := range compounds {
		if compound.Kind == goxy.Page {
			pages = append(pages, GeekdocBundleMenuItem{
				Name: compound.Title,
				Ref:  fmt.Sprintf("page/%s", compound.Id),
			})
		}
	}

	groups := make([]GeekdocBundleMenuItem, 0)
	for _, compound := range compounds {
		if compound.Kind == goxy.Group && compound.Parent == "" {
			group := GeekdocBundleMenuItem{
				Name: compound.Title,
				Ref:  fmt.Sprintf("group/%s", compound.Id),
				Sub:  []GeekdocBundleMenuItem{},
			}

			for _, inner := range compounds {
				if inner.Kind == goxy.Group && inner.Parent == compound.Id {
					group.Sub = append(group.Sub, GeekdocBundleMenuItem{
						Name: inner.Title,
						Ref:  fmt.Sprintf("group/%s", inner.Id),
					},
					)
				}
			}

			groups = append(groups, group)
		}
	}

	menu := map[string][]GeekdocBundleMenuItem{
		"main": {
			{
				Name: "Coding",
				Sub: []GeekdocBundleMenuItem{
					{
						Name: "Classes",
						Ref:  "/class",
					},
					{
						Name: "Files",
						Ref:  "/file",
					},
					{
						Name: "Dirs",
						Ref:  "/dir",
					},
					{
						Name: "Groups",
						Ref:  "/group",
					},
					{
						Name: "Namespaces",
						Ref:  "/namespace",
					},
					{
						Name: "Pages",
						Ref:  "/page",
						Sub:  pages,
					},
					{
						Name: "Unions",
						Ref:  "/union",
					},
				},
			},
		},
	}

	os.MkdirAll("hugo/data/menu", 0644)
	bytes, err = yaml.Marshal(menu)
	if err != nil {
		log.Fatalf("Error: %v", errors.WithStack(err))
	}
	err = ioutil.WriteFile("hugo/data/menu/main.yml", bytes, 0644)
	if err != nil {
		log.Fatalf("Error: %v", errors.WithStack(err))
	}

	return
}

func AddRefsFromDescriptions(refs map[string]GoxygenRef, id string, d goxy.Descriptions) {
	AddRefsFromDocstring(refs, id, d.DetailedDescription)
	AddRefsFromDocstring(refs, id, d.BriefDescription)
	AddRefsFromDocstring(refs, id, d.InBodyDescription)
}

func AddRefsFromDocstring(refs map[string]GoxygenRef, id string, doc goxy.DocString) {
	for _, element := range doc.Content {
		switch element.Type {
		case goxy.Anchor:
			a := element.Value.(goxy.DocStringAnchor)
			refs[a.Id] = GoxygenRef{
				Kind:      string(goxy.Anchor),
				Name:      "N/A",
				RefId:     a.Id,
				ParentRef: id,
			}
		case goxy.Section:
			s := element.Value.(goxy.DocStringSection)
			if s.Id != "" {
				refs[s.Id] = GoxygenRef{
					Kind:      s.Kind,
					Name:      "N/A",
					RefId:     s.Id,
					ParentRef: id,
				}
			}
			AddRefsFromDocstring(refs, id, s.Content)
		case goxy.Paragraph:
			p := element.Value.(goxy.DocStringParagraph)
			AddRefsFromDocstring(refs, id, p.Content)
		case goxy.Title:
			v := element.Value.(goxy.DocStringTitle)
			AddRefsFromDocstring(refs, id, v.Content)
		case goxy.Heading:
			v := element.Value.(goxy.DocStringHeading)
			AddRefsFromDocstring(refs, id, v.Content)
		case goxy.ParameterList:
			v := element.Value.(goxy.DocStringParameterList)
			for _, item := range v.Items {
				AddRefsFromDocstring(refs, id, item.Description)
			}
		case goxy.XRefSect:
			v := element.Value.(goxy.DocStringXRefSect)
			AddRefsFromDocstring(refs, id, v.Description)
		case goxy.OrderedList:
			v := element.Value.(goxy.DocStringOrderedList)
			for _, item := range v.Items {
				AddRefsFromDocstring(refs, id, item)
			}
		case goxy.ItemizedList:
			v := element.Value.(goxy.DocStringItemizedList)
			for _, item := range v.Items {
				AddRefsFromDocstring(refs, id, item)
			}
		case goxy.Bold:
			v := element.Value.(goxy.DocStringBold)
			AddRefsFromDocstring(refs, id, v.Content)
		case goxy.Emphasis:
			v := element.Value.(goxy.DocStringEmphasis)
			AddRefsFromDocstring(refs, id, v.Content)
		case goxy.Verbatim:
			v := element.Value.(goxy.DocStringVerbatim)
			AddRefsFromDocstring(refs, id, v.Content)
		case goxy.Term:
			v := element.Value.(goxy.DocStringTerm)
			AddRefsFromDocstring(refs, id, v.Content)
		case goxy.ComputerOutput:
			v := element.Value.(goxy.DocStringComputerOutput)
			AddRefsFromDocstring(refs, id, v.Content)
		case goxy.Highlight:
			v := element.Value.(goxy.DocStringHighlight)
			AddRefsFromDocstring(refs, id, v.Content)
		case goxy.Ref:
			v := element.Value.(goxy.DocStringRef)
			AddRefsFromDocstring(refs, id, v.Content)
		}
	}
}
