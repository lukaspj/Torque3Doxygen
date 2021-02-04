package formatter

import (
	"ScriptExecServer/pkg/goxy"
	"fmt"
	"os"
	"text/template"
)

type Hugo struct {

}

func NewHugoFormatter() *Hugo {
	return &Hugo{}
}

func (h Hugo) WriteCompound(compound *goxy.CompoundDoc, path string) error {
	var err error

	err = os.MkdirAll(fmt.Sprintf("%s", path), 0644)
	if err != nil {
		return err
	}

	var mdType string
	switch compound.Kind {
	case goxy.Dir:
		mdType = "dir"
	default:
		mdType = "compound"
	}

	mdContent := fmt.Sprintf(`---
goxygen_id: "%s"
goxygen_key: "coding"
GeekdocFlatSection: true
title: "%s"
type: "%s"
---
`, compound.Id, compound.Title, mdType)

	t, err := template.New("compound").
		Parse(mdContent)

	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	err = t.ExecuteTemplate(f, compound.Name, compound)
	if err != nil {
		return err
	}
	/*
	err = ioutil.WriteFile(fmt.Sprintf("hugo/content/coding/%s/%s.md", compound.Kind, compound.Id), []byte(mdContent), 0644)
	if err != nil {
		return err
	}
*/
	return nil
}