package formatter

import (
	"ScriptExecServer/pkg/formatter/templates"
	"ScriptExecServer/pkg/goxy"
	"bufio"
	"bytes"
	"fmt"
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/pkg/errors"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

type Hugo struct {
	Section string

	CompoundIdMap map[string]*goxy.CompoundDoc
	CompoundRefs  map[string]goxy.CompoundRef
}

var funcMap = template.FuncMap{
	"HasPrefix": strings.HasPrefix,
}

func NewHugoFormatter(section string, idMap map[string]*goxy.CompoundDoc, refs map[string]goxy.CompoundRef) *Hugo {
	return &Hugo{
		Section:       section,
		CompoundIdMap: idMap,
		CompoundRefs:  refs,
	}
}

type CompoundTemplateModel struct {
	Section  string
	Type     string
	H        *Hugo
	Compound *goxy.CompoundDoc
}

type HighlightOpts struct {
	ShowLineNumbers bool
	Language        string
}

func (h *Hugo) RenderHighlight(language string, content string) string {
	buf := bytes.NewBufferString("")

	formatter := html.New(
		html.WithClasses(true),
	)

	style := styles.Get("swapoff")
	if style == nil {
		style = styles.Fallback
	}

	lexer := lexers.Get(language)
	lexer = chroma.Coalesce(lexer)
	iterator, _ := lexer.Tokenise(nil, content)
	_ = formatter.Format(buf, style, iterator)

	return ReplaceHrefInHighlight(buf.String())
}

func (h *Hugo) RenderHighlightWithLineNos(language string, content string) string {
	buf := bytes.NewBufferString("")

	formatter := html.New(
		html.WithClasses(true),
		html.WithLineNumbers(true),
	)

	style := styles.Get("swapoff")
	if style == nil {
		style = styles.Fallback
	}

	lexer := lexers.Get(language)
	lexer = chroma.Coalesce(lexer)
	iterator, _ := lexer.Tokenise(nil, content)
	_ = formatter.Format(buf, style, iterator)

	return ReplaceHrefInHighlight(buf.String())
}

func minInt(x, y int) int {
	if x > y {
		return y
	} else {
		return x
	}
}

func ReplaceHrefInHighlight(s string) string {
	// The simplest case
	re1, _ := regexp.Compile("<span [^>]+>((?:-|:|\\(|\\)|\\*|&amp;|&gt;|&lt;)*)&lt;</span>\\s*<span [^>]+>a</span>\\s*<span [^>]+>href</span>\\s*<span [^>]+>=</span>\\s*<span [^>]+>(?:&#34;|\")([^<]+)(?:&#34;|\")</span>\\s*<span [^>]+>&gt;((?:::|~)*)</span>\\s*((?:<span [^>]+>[^<]+</span>\\s*)*?)\\s*(<span [^>]+>[^<]+</span>)\\s*<span [^>]+>((?:&lt;|&gt;|=|\\*|!|-|\\||&amp;|\\^|\\+|~|/|:)*)&lt;/</span>\\s*<span [^>]+>a</span>\\s*<span [^>]+>&gt;((?:\\*|&amp;|&lt;|&gt;|-|\\(|\\)|:|\\+|=|!|/|\\?|~)*)</span>")
	// Handle :: namespacing
	// re2, _ := regexp.Compile("<span [^>]+>((?:-|\\(|\\)|\\*|&amp;|&gt;|&lt;)*)&lt;</span>\\s*<span [^>]+>a</span>\\s*<span [^>]+>href</span>\\s*<span [^>]+>=</span>\\s*<span [^>]+>(?:&#34;|\")([^<]+)(?:&#34;|\")</span>\\s*<span [^>]+>&gt;(~?)</span>\\s*(<span [^>]+>[^<]+</span>\\s*<span [^>]+>::</span>)?\\s*(<span [^>]+>[^<]+</span>)\\s*<span [^>]+>((?:&gt;|=)*)&lt;/</span>\\s*<span [^>]+>a</span>\\s*<span [^>]+>&gt;((?:\\*|&amp;|&lt;|&gt;|-|\\(|\\)|:)*)</span>")
	// Handle Generics
	// re3, _ := regexp.Compile("<span [^>]+>((?:-|\\(|\\)|\\*|&amp;|&gt;|&lt;)*)&lt;</span>\\s*<span [^>]+>a</span>\\s*<span [^>]+>href</span>\\s*<span [^>]+>=</span>\\s*<span [^>]+>(?:&#34;|\")([^<]+)(?:&#34;|\")</span>\\s*<span [^>]+>&gt;(~?)</span>\\s*(<span [^>]+>[^<]+</span>\\s*<span [^>]+>&lt;</span>)?\\s*(<span [^>]+>[^<]+</span>)\\s*<span [^>]+>((?:&gt;|=)*)&lt;/</span>\\s*<span [^>]+>a</span>\\s*<span [^>]+>&gt;((?:\\*|&amp;|&lt;|&gt;|-|\\(|\\)|:)*)</span>")

	ret := re1.ReplaceAllString(s, "<span class=\"o\">$1</span><a href=$2><span class=\"n\">$3</span>$4$5$6</a><span class=\"o\">$7</span>")
	tmp := ret
	for {
		tmp = re1.ReplaceAllString(ret, "<span class=\"o\">$1</span><a href=$2><span class=\"n\">$3</span>$4$5$6</a><span class=\"o\">$7</span>")
		if tmp == ret {
			break
		} else {
			ret = tmp
		}
	}

	failIdx := strings.Index(ret, "<span class=\"o\">&lt;</span><span class=\"n\">a</span> <span class=\"n\">href")
	if failIdx >= 0 {
		panic(fmt.Sprintf("unable to transform link in: %s,\n\n\nreturn: %s\n\n\ncontext: %s", s, ret, ret[failIdx:minInt(len(ret), failIdx + 2000)]))
	}
	return ret
	/*re3.ReplaceAllString(
		re2.ReplaceAllString(
			re1.ReplaceAllString(s, "<span class=\"o\">$1</span><a href=$2><span class=\"n\">$3</span>$4$5$6</a><span class=\"o\">$7</span>"),
			"<span class=\"o\">$1</span><a href=$2><span class=\"n\">$3</span>$4$5$6</a><span class=\"o\">$7</span>",
		),
		"<span class=\"o\">$1</span><a href=$2><span class=\"n\">$3</span>$4$5$6</a><span class=\"o\">&gt;$7</span>",
	)*/
}

func (h *Hugo) RenderBriefFunctionDecl(function goxy.FunctionDoc) string {
	buf := bytes.NewBufferString("")

	_, _ = fmt.Fprintf(buf, "<a href=\"#%s\">%s</a>(", function.Id, function.Name)

	paramStrings := make([]string, len(function.Params))
	for idx, param := range function.Params {
		paramStrings[idx] = fmt.Sprintf("%s %s", h.RenderDocstring(param.Type), param.DeclName)
	}

	_, _ = fmt.Fprint(buf, strings.Join(paramStrings, ", "))
	_, _ = fmt.Fprint(buf, ")")

	return buf.String()
}

func (h *Hugo) RenderBriefDefineDecl(define goxy.DefineDoc) string {
	buf := bytes.NewBufferString("")

	_, _ = fmt.Fprintf(buf, "<a href=\"#%s\">%s</a>(", define.Id, define.Name)

	paramStrings := make([]string, len(define.Params))
	for idx, param := range define.Params {
		paramStrings[idx] = fmt.Sprintf("%s", param.Defname)
	}

	_, _ = fmt.Fprint(buf, strings.Join(paramStrings, ", "))
	_, _ = fmt.Fprintf(buf, ") %s", define.Initializer)

	return buf.String()
}

func (h *Hugo) RenderReimplementedFrom(f goxy.FunctionDoc) string {
	ref, ok := h.CompoundRefs[f.Reimplements.RefId]
	if !ok {
		return "&lt;UNKNOWN TYPE&gt;"
	}

	pRef, ok := h.CompoundRefs[ref.ParentRef]
	if !ok {
		return "&lt;UNKNOWN PARENT TYPE&gt;"
	}

	return fmt.Sprintf("<a href=\"/%s/%s/%s#%s\">%s</a>", h.Section, pRef.Kind, strings.ToLower(pRef.RefId), ref.RefId, pRef.Name)
}

func (h *Hugo) RenderReimplementedBy(f goxy.FunctionDoc) string {
	buf := bytes.NewBufferString("")

	for i, reimplements := range f.ReimplementedBy {
		if i > 0 {
			_, _ = fmt.Fprint(buf, ", ")
		}

		ref, ok := h.CompoundRefs[reimplements.RefId]
		if !ok {
			_, _ = fmt.Fprint(buf, "&lt;UNKNOWN TYPE&gt;")
			continue

		}

		pRef, ok := h.CompoundRefs[ref.ParentRef]
		if !ok {
			_, _ = fmt.Fprint(buf, "&lt;UNKNOWN PARENT TYPE&gt;")
			continue
		}

		_, _ = fmt.Fprintf(buf, "<a href=\"/%s/%s/%s#%s\">%s</a>", h.Section, pRef.Kind, strings.ToLower(pRef.RefId), ref.RefId, pRef.Name)
	}

	return buf.String()
}

func (h *Hugo) RenderRef(refId, content string) string {
	if c, ok := h.CompoundRefs[refId]; !ok {
		log.Printf("error: %+v", fmt.Errorf("unknown ref: %s", refId))
		return content
	} else {
		if p, ok := h.CompoundRefs[c.ParentRef]; ok {
			return fmt.Sprintf("<a href=\"/%s/%s/%s#%s\">%s</a>", h.Section, p.Kind, strings.ToLower(p.RefId), c.RefId, content)
		} else {
			return fmt.Sprintf("<a href=\"/%s/%s/%s\">%s</a>", h.Section, c.Kind, strings.ToLower(c.RefId), content)
		}
	}
}

func (h *Hugo) RenderDocstring(docstring goxy.DocString) string {
	buf := bytes.NewBufferString("")

	for _, element := range docstring.Content {
		switch e := element.Value.(type) {
		case goxy.DocStringText:
			_, _ = fmt.Fprintf(buf, strings.ReplaceAll(e.Content, "{{", "££@$$"))
		case goxy.DocStringParagraph:
			_, _ = fmt.Fprintf(buf, "<p>%s</p>", h.RenderDocstring(e.Content))
		case goxy.DocStringEmphasis:
			_, _ = fmt.Fprintf(buf, "<em>%s</em>", h.RenderDocstring(e.Content))
		case goxy.DocStringBold:
			_, _ = fmt.Fprintf(buf, "<b>%s</b>", h.RenderDocstring(e.Content))
		case goxy.DocStringVerbatim:
			_, _ = fmt.Fprintf(buf, "<pre>%s</pre>", h.RenderDocstring(e.Content))
		case goxy.DocStringPreformatted:
			_, _ = fmt.Fprintf(buf, "<pre>%s</pre>", h.RenderDocstring(e.Content))
		case goxy.DocStringComputerOutput:
			_, _ = fmt.Fprintf(buf, "<pre>%s</pre>", h.RenderDocstring(e.Content))
		case goxy.DocStringItemizedList:
			_, _ = fmt.Fprintf(buf, "<ul>")
			for _, item := range e.Items {
				_, _ = fmt.Fprintf(buf, "<li>%s</li>", h.RenderDocstring(item))
			}
			_, _ = fmt.Fprintf(buf, "</ul>")
		case goxy.DocStringOrderedList:
			_, _ = fmt.Fprintf(buf, "<ol>")
			for _, item := range e.Items {
				_, _ = fmt.Fprintf(buf, "<li>%s</li>", h.RenderDocstring(item))
			}
			_, _ = fmt.Fprintf(buf, "</ol>")
		case goxy.DocStringVariableList:
			_, _ = fmt.Fprintf(buf, "<dl>")
			for _, item := range e.Items {
				if item.Content[0].Type == goxy.Term {
					_, _ = fmt.Fprintf(buf, "<dt>%s</dt>", h.RenderDocstring(item))
				} else {
					_, _ = fmt.Fprintf(buf, "<dd>%s</dd>", h.RenderDocstring(item))
				}
			}
			_, _ = fmt.Fprintf(buf, "</dl>")
		case goxy.DocStringTerm:
			_, _ = fmt.Fprintf(buf, "%s", h.RenderDocstring(e.Content))
		case goxy.DocStringHeading:
			_, _ = fmt.Fprintf(buf, "<h%d>%s</h%d>", e.Level, h.RenderDocstring(e.Content), e.Level)
		case goxy.DocStringXRefSect:
			_, _ = fmt.Fprintf(buf, h.RenderRef(e.Id, fmt.Sprintf("<b>%s</b>: %s", e.Title, e.Description)))
		case goxy.DocStringRef:
			_, _ = fmt.Fprintf(buf, h.RenderRef(e.RefId, h.RenderDocstring(e.Content)))
		case goxy.DocStringAnchor:
			_, _ = fmt.Fprintf(buf, "<a id=\"%s\"></a>", e.Id)
		case goxy.DocStringSection:
			t, err := template.New("docstringsection").
				Funcs(funcMap).
				Parse(templates.DocstringSection)
			if err != nil {
				log.Fatalf("error: %+v", errors.WithStack(err))
			}

			err = t.ExecuteTemplate(buf, "docstringsection", map[string]interface{}{
				"H":       h,
				"Id":      e.Id,
				"Kind":    e.Kind,
				"Content": e.Content,
			})
			if err != nil {
				log.Fatalf("error: %+v", errors.WithStack(err))
			}
		case goxy.DocStringTitle:
			_, _ = fmt.Fprintf(buf, "<h2>%s</h2>", h.RenderDocstring(e.Content))
		case goxy.DocStringParameterList:
			t, err := template.New("docstringparamlist").
				Funcs(funcMap).
				Parse(templates.DocstringParameterList)
			if err != nil {
				log.Fatalf("error: %+v", errors.WithStack(err))
			}

			err = t.ExecuteTemplate(buf, "docstringparamlist", map[string]interface{}{
				"H":     h,
				"Kind":  e.Kind,
				"Items": e.Items,
			})
			if err != nil {
				log.Fatalf("error: %+v", errors.WithStack(err))
			}
		case goxy.DocStringTable:
			t, err := template.New("docstringtable").
				Funcs(funcMap).
				Parse(templates.DocstringTable)
			if err != nil {
				log.Fatalf("error: %+v", errors.WithStack(err))
			}

			err = t.ExecuteTemplate(buf, "docstringtable", map[string]interface{}{
				"H":    h,
				"Rows": e.Rows,
			})
			if err != nil {
				log.Fatalf("error: %+v", errors.WithStack(err))
			}
		case goxy.DocStringImage:
			_, _ = fmt.Fprintf(buf, "<img src=\"%s\" alt=\"%s\" />", e.Name, e.Description)
		case goxy.DocStringHighlight:
			_, _ = fmt.Fprintf(buf, "%s", h.RenderHighlight("C++", h.RenderDocstring(e.Content)))
		case goxy.DocStringLinebreak:
			_, _ = fmt.Fprint(buf, "<br />")
		default:
			log.Fatalf("error: %+v", errors.WithStack(errors.New("unable to resolve docstring type: "+string(element.Type))))
		}
	}

	return buf.String()
}

func (h *Hugo) RenderEnumBody(values []goxy.EnumValue) string {
	buf := bytes.NewBufferString("{")
	if len(values) < 2 {
		_, _ = fmt.Fprint(buf, "\n")
	}
	for _, value := range values {
		_, _ = fmt.Fprintf(buf, "  %s %s\n", value.Name, value.Initializer)
	}
	_, _ = fmt.Fprintf(buf, "}")
	return buf.String()
}

func (h *Hugo) RenderInnerCompound(compound goxy.InnerCompoundRef) string {
	buf := bytes.NewBufferString("")

	t, err := template.New("innercompound").
		Funcs(funcMap).
		Parse(templates.InnerCompound)
	if err != nil {
		log.Fatalf("error: %+v", errors.WithStack(err))
	}

	err = t.ExecuteTemplate(buf, "innercompound", map[string]interface{}{
		"H":          h,
		"RefId":      compound.RefId,
		"Value":      compound.Value,
		"Protection": compound.Protection,
	})
	if err != nil {
		log.Fatalf("error: %+v", errors.WithStack(err))
	}

	return buf.String()
}

func (h *Hugo) RenderSectionBrief(section *goxy.SectionDoc) string {
	buf := bytes.NewBufferString("")

	t, err := template.New("sectionbrief").
		Funcs(funcMap).
		Parse(templates.SectionBrief)
	if err != nil {
		log.Fatalf("error: %+v", errors.WithStack(err))
	}

	err = t.ExecuteTemplate(buf, "sectionbrief", map[string]interface{}{
		"H":       h,
		"Section": section,
	})
	if err != nil {
		log.Fatalf("error: %+v", errors.WithStack(err))
	}

	return buf.String()
}

func (h *Hugo) RenderSection(section *goxy.SectionDoc) string {
	buf := bytes.NewBufferString("")

	t, err := template.New("section").
		Funcs(funcMap).
		Parse(templates.Section)
	if err != nil {
		log.Fatalf("error: %+v", errors.WithStack(err))
	}

	err = t.ExecuteTemplate(buf, "section", map[string]interface{}{
		"H":       h,
		"Section": section,
	})
	if err != nil {
		log.Fatalf("error: %+v", errors.WithStack(err))
	}

	return buf.String()
}

func (h *Hugo) RenderDirJson(compound *goxy.CompoundDoc) string {
	buf := bytes.NewBufferString("")

	_, _ = fmt.Fprintf(buf, "{\"Name\":\"%s\",\"Id\":\"%s\"", compound.Name, compound.Id)

	_, _ = fmt.Fprint(buf, ",\"Dirs\":[")
	for idx, dir := range compound.InnerDirs {
		if idx > 0 {
			_, _ = fmt.Fprint(buf, ",")
		}
		_, _ = fmt.Fprint(buf, h.RenderDirJson(h.CompoundIdMap[dir.RefId]))
	}
	_, _ = fmt.Fprint(buf, "],\"Files\":[")
	for idx, file := range compound.InnerFiles {
		if idx > 0 {
			_, _ = fmt.Fprint(buf, ",")
		}
		fileCompound := h.CompoundIdMap[file.RefId]
		_, _ = fmt.Fprintf(buf, "{\"Name\":\"%s\",\"Id\":\"%s\"}", fileCompound.Name, fileCompound.Id)
	}
	_, _ = fmt.Fprint(buf, "]}")

	return buf.String()
}

func (h *Hugo) WriteCompound(compound *goxy.CompoundDoc, path string) error {
	var err error

	err = os.MkdirAll(fmt.Sprintf("%s", filepath.Dir(path)), 0644)
	if err != nil {
		return errors.WithStack(err)
	}

	var mdType string
	switch compound.Kind {
	case goxy.Dir:
		mdType = "dir"
	default:
		mdType = "compound"
	}

	model := CompoundTemplateModel{
		H:        h,
		Section:  h.Section,
		Type:     mdType,
		Compound: compound,
	}

	compoundTemplate := templates.Compound
	if compound.Kind == goxy.Dir {
		compoundTemplate = templates.DirCompound
	}

	t, err := template.New("compound").
		Funcs(funcMap).
		Parse(compoundTemplate)

	if err != nil {
		return errors.WithStack(err)
	}

	f, err := os.Create(path)
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()

	buf := bytes.NewBufferString("")
	err = t.ExecuteTemplate(buf, "compound", model)
	if err != nil {
		return errors.WithStack(err)
	}
	w := bufio.NewWriter(f)
	_, _ = w.WriteString(strings.ReplaceAll(buf.String(), "££@$$", "{{\"{\"}}"))

	err = w.Flush()
	if err != nil {
		return errors.WithStack(err)
	}
	/*
		err = ioutil.WriteFile(fmt.Sprintf("hugo/content/coding/%s/%s.md", compound.Kind, compound.Id), []byte(mdContent), 0644)
		if err != nil {
			return err
		}
	*/
	return nil
}
