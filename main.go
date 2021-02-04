package main

import (
	"ScriptExecServer/pkg/doxygen"
	"ScriptExecServer/pkg/formatter"
	"ScriptExecServer/pkg/goxy"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type GoxygenData struct {
	Entities map[string]*goxy.CompoundDoc
	Refs     map[string]goxy.CompoundRef
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
	/*r := chi.NewRouter()

	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatalf("server crashed with error: %v", err)
	}
	 */

	scriptDocs := doxygen.ParseDoxygenFolder("script-doxygen/xml")

	scriptCompounds, scriptData := ExtractDoxygenMetadata(scriptDocs)

	docs := doxygen.ParseDoxygenFolder("doxygen/xml")

	compounds, data := ExtractDoxygenMetadata(docs)

	scriptingFormatter := formatter.NewHugoFormatter("scripting", scriptData.Entities, scriptData.Refs)
	codingFormatter := formatter.NewHugoFormatter("coding", data.Entities, data.Refs)

	for _, compound := range compounds {
		err := codingFormatter.WriteCompound(compound, fmt.Sprintf("hugo/content/coding/%s/%s.html", compound.Kind, compound.Id))

		if err != nil {
			log.Fatalf("Error: %+v", err)
		}
	}

	for _, compound := range scriptCompounds {
		err := scriptingFormatter.WriteCompound(compound, fmt.Sprintf("hugo/content/scripting/%s/%s.html", compound.Kind, compound.Id))

		if err != nil {
			log.Fatalf("Error: %+v", err)
		}
	}

	os.MkdirAll("hugo/data", 0644)
	bytes, err := json.Marshal(map[string]GoxygenData {
		"coding": data,
		"scripting": scriptData,
	})
	if err != nil {
		log.Fatalf("Error: %v", errors.WithStack(err))
	}
	err = ioutil.WriteFile("hugo/data/goxygen.json", bytes, 0644)
	if err != nil {
		log.Fatalf("Error: %v", errors.WithStack(err))
	}

	codingPages := make([]GeekdocBundleMenuItem, 0)
	for _, compound := range compounds {
		if compound.Kind == goxy.Page {
			codingPages = append(codingPages, GeekdocBundleMenuItem{
				Name: compound.Title,
				Ref:  fmt.Sprintf("coding/page/%s", compound.Id),
			})
		}
	}

	scriptingPages := make([]GeekdocBundleMenuItem, 0)
	for _, compound := range scriptCompounds {
		if compound.Kind == goxy.Page {
			scriptingPages = append(scriptingPages, GeekdocBundleMenuItem{
				Name: compound.Title,
				Ref:  fmt.Sprintf("scripting/page/%s", compound.Id),
			})
		}
	}

	engineGroups := make([]GeekdocBundleMenuItem, 0)
	var rootDirRefId string
	for _, compound := range compounds {
		if compound.Kind == goxy.Group && compound.Parent == "" {
			group := GeekdocBundleMenuItem{
				Name: compound.Title,
				Ref:  fmt.Sprintf("coding/group/%s", compound.Id),
				Sub:  []GeekdocBundleMenuItem{},
			}

			for _, inner := range compounds {
				if inner.Kind == goxy.Group && inner.Parent == compound.Id {
					group.Sub = append(group.Sub, GeekdocBundleMenuItem{
						Name: inner.Title,
						Ref:  fmt.Sprintf("coding/group/%s", inner.Id),
					},
					)
				}
			}

			engineGroups = append(engineGroups, group)
		}

		if compound.Kind == goxy.Dir && compound.Title == "Engine" {
			rootDirRefId = compound.Id
		}
	}

	scriptGroups := make([]GeekdocBundleMenuItem, 0)
	for _, compound := range scriptCompounds {
		if compound.Kind == goxy.Group && compound.Parent == "" {
			group := GeekdocBundleMenuItem{
				Name: compound.Title,
				Ref:  fmt.Sprintf("scripting/group/%s", compound.Id),
				Sub:  []GeekdocBundleMenuItem{},
			}

			for _, inner := range compounds {
				if inner.Kind == goxy.Group && inner.Parent == compound.Id {
					group.Sub = append(group.Sub, GeekdocBundleMenuItem{
						Name: inner.Title,
						Ref:  fmt.Sprintf("scripting/group/%s", inner.Id),
					},
					)
				}
			}

			scriptGroups = append(scriptGroups, group)
		}
	}

	menu := map[string][]GeekdocBundleMenuItem{
		"main": {
			{
				Name: "Coding Reference",
				Sub: []GeekdocBundleMenuItem{
					{
						Name: "Classes",
						Ref:  "coding/class",
					},
					{
						Name: "Files",
						Ref:  fmt.Sprintf("coding/dir/%s", rootDirRefId),
					},
					{
						Name: "Groups",
						Ref:  "coding/group",
					},
					{
						Name: "Namespaces",
						Ref:  "coding/namespace",
					},
					{
						Name: "Pages",
						Ref:  "coding/page",
						Sub:  codingPages,
					},
					{
						Name: "Unions",
						Ref:  "coding/union",
					},
				},
			},
			{
				Name: "Scripting Reference",
				Sub: []GeekdocBundleMenuItem{
					{
						Name: "Classes",
						Ref:  "scripting/class",
					},
					{
						Name: "Groups",
						Ref:  "scripting/group",
					},
					{
						Name: "Namespaces",
						Ref:  "scripting/namespace",
					},
					{
						Name: "Pages",
						Ref:  "scripting/page",
						Sub:  codingPages,
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

func ExtractDoxygenMetadata(docs []*doxygen.Doxygen) ([]*goxy.CompoundDoc, GoxygenData) {
	compounds := make([]*goxy.CompoundDoc, 0)

	data := GoxygenData{
		Entities: make(map[string]*goxy.CompoundDoc),
		Refs:     make(map[string]goxy.CompoundRef),
	}
	files := make(map[string]*goxy.CompoundDoc)
	for _, doc := range docs {
		compound, err := goxy.CompoundFromDoxygen(doc)
		if err != nil {
			fmt.Println(fmt.Sprintf("unable to parse doxygen compound doc: %v, due to: %v", doc.CompoundDef.CompoundName, err))
		} else {
			compounds = append(compounds, compound)
			if compound.Kind == goxy.File {
				files[strings.ToLower(compound.Location.File)] = compound
			}
		}
	}

	for _, compound := range compounds {
		if file, ok := files[strings.ToLower(compound.Location.File)]; ok {
			compound.Location.FileRefId = file.Id
		}
		if file, ok := files[strings.ToLower(compound.Location.BodyFile)]; ok {
			compound.Location.BodyFileRefId = file.Id
		}

		if compound.Kind == goxy.Class {
			for _, class := range compound.InnerClasses {
				for _, inner := range compounds {
					if inner.Id == class.RefId {
						inner.Parent = compound.Id
					}
				}
			}
		}

		if compound.Kind == goxy.Namespace {
			for _, namespace := range compound.InnerNamespaces {
				for _, inner := range compounds {
					if inner.Id == namespace.RefId {
						inner.Parent = compound.Id
					}
				}
			}
		}

		if compound.Kind == goxy.Group {
			for _, group := range compound.InnerGroups {
				for _, inner := range compounds {
					if inner.Id == group.RefId {
						inner.Parent = compound.Id
					}
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
		data.Entities[compound.Id] = compound

		data.Refs[compound.Id] = goxy.CompoundRef{
			Kind:      string(compound.Kind),
			Name:      compound.Name,
			ParentRef: "N/D",
			RefId:     compound.Id,
		}

		AddRefsFromDescriptions(data.Refs, compound.Id, compound.Descriptions)

		for _, section := range compound.Sections {
			AddRefsFromDocstring(data.Refs, compound.Id, section.Description)

			for _, function := range section.Functions {
				data.Refs[function.Id] = goxy.CompoundRef{
					Kind:      "function",
					Name:      function.Name,
					ParentRef: compound.Id,
					RefId:     function.Id,
				}

				AddRefsFromDescriptions(data.Refs, compound.Id, function.Descriptions)
			}
			for _, enum := range section.Enums {
				data.Refs[enum.Id] = goxy.CompoundRef{
					Kind:      "enum",
					Name:      enum.Name,
					ParentRef: compound.Id,
					RefId:     enum.Id,
				}

				for _, value := range enum.Values {
					data.Refs[value.Id] = goxy.CompoundRef{
						Kind:      "enumvalue",
						Name:      value.Name,
						ParentRef: compound.Id,
						RefId:     value.Id,
					}
				}

				AddRefsFromDescriptions(data.Refs, compound.Id, enum.Descriptions)
			}
			for _, attr := range section.Attributes {
				data.Refs[attr.Id] = goxy.CompoundRef{
					Kind:      "attribute",
					Name:      attr.Name,
					ParentRef: compound.Id,
					RefId:     attr.Id,
				}

				AddRefsFromDescriptions(data.Refs, compound.Id, attr.Descriptions)
			}
			for _, def := range section.Defines {
				data.Refs[def.Id] = goxy.CompoundRef{
					Kind:      "define",
					Name:      def.Name,
					ParentRef: compound.Id,
					RefId:     def.Id,
				}

				AddRefsFromDescriptions(data.Refs, compound.Id, def.Descriptions)
			}
			for _, typedef := range section.Typedefs {
				data.Refs[typedef.Id] = goxy.CompoundRef{
					Kind:      "typedef",
					Name:      typedef.Name,
					ParentRef: compound.Id,
					RefId:     typedef.Id,
				}

				AddRefsFromDescriptions(data.Refs, compound.Id, typedef.Descriptions)
			}
			for _, friend := range section.Friends {
				data.Refs[friend.Id] = goxy.CompoundRef{
					Kind:      "friend",
					Name:      friend.Name,
					ParentRef: compound.Id,
					RefId:     friend.Id,
				}

				AddRefsFromDescriptions(data.Refs, compound.Id, friend.Descriptions)
			}
		}
	}

	return compounds, data
}

func AddRefsFromDescriptions(refs map[string]goxy.CompoundRef, id string, d goxy.Descriptions) {
	AddRefsFromDocstring(refs, id, d.DetailedDescription)
	AddRefsFromDocstring(refs, id, d.BriefDescription)
	AddRefsFromDocstring(refs, id, d.InBodyDescription)
}

func AddRefsFromDocstring(refs map[string]goxy.CompoundRef, id string, doc goxy.DocString) {
	for _, element := range doc.Content {
		switch element.Type {
		case goxy.Anchor:
			a := element.Value.(goxy.DocStringAnchor)
			refs[a.Id] = goxy.CompoundRef{
				Kind:      string(goxy.Anchor),
				Name:      "N/A",
				RefId:     a.Id,
				ParentRef: id,
			}
		case goxy.Section:
			s := element.Value.(goxy.DocStringSection)
			if s.Id != "" {
				refs[s.Id] = goxy.CompoundRef{
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
		case goxy.VariableList:
			v := element.Value.(goxy.DocStringVariableList)
			for _, item := range v.Items {
				AddRefsFromDocstring(refs, id, item)
			}
		case goxy.Table:
			v := element.Value.(goxy.DocStringTable)
			for _, row := range v.Rows {
				for _, col := range row {
					AddRefsFromDocstring(refs, id, col.Content)
				}
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
		case goxy.Preformatted:
			v := element.Value.(goxy.DocStringPreformatted)
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
		case goxy.Image:
		case goxy.Text:
		case goxy.LineBreak:
		default:
			fmt.Printf("unhandled docstring ref: %s\n", element.Type)
		}
	}
}
