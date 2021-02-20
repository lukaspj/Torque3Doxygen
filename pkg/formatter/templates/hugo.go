package templates

const DirCompound = `---
GeekdocFlatSection: true
title: "{{ .Compound.Title }}"
type: "{{ .Type }}"
url: "/{{ .Section }}/{{ .Compound.Kind }}/{{ .Compound.Id }}"

goxygen:
  kind: "{{ .Compound.Kind }}"
  section: "{{ .Section }}"
---
<div id="file-tree-container"></div>

<script defer>
	function visitFile(file) {
		const node = new TreeNode(file.Name);
		node.on('click', function () {
			window.location.href = "/{{ .Section }}/file/" + file.Id;
		})
		return node;
	}

	function visitDir(dir) {
		var node = new TreeNode(dir.Name);
		for (let i = 0; i < dir.Dirs.length; i++) {
			const subNode = visitDir(dir.Dirs[i]);
			node.addChild(subNode)
		}
		for (let i = 0; i < dir.Files.length; i++) {
			const subNode = visitFile(dir.Files[i]);
			node.addChild(subNode)
		}

		return node;
	}

	document.addEventListener("DOMContentLoaded", function () {
		var dirTree = {{ $.H.RenderDirJson .Compound }};

		const treeNode = visitDir(dirTree);

		var tree = new TreeView(treeNode, "#file-tree-container", TreeConfig);
		tree.collapseAllNodes();
		tree.reload();
	}, false);
</script>
`

const Compound = `---
GeekdocFlatSection: true
title: "{{ .Compound.Title }}"
type: "{{ .Type }}"
url: "/{{ .Section }}/{{ .Compound.Kind }}/{{ .Compound.Id }}"

goxygen:
  kind: "{{ .Compound.Kind }}"
  section: "{{ .Section }}"

GeekdocSearchKeywords:
  {{- range .Compound.InnerClasses }}
  - "{{ (index $.H.CompoundIdMap .RefId).Title }}"
  {{- end }}
  {{- range .Compound.InnerGroups }}
  - "{{ (index $.H.CompoundIdMap .RefId).Title }}"
  {{- end }}
  {{- range .Compound.InnerFiles }}
  - "{{ (index $.H.CompoundIdMap .RefId).Title }}"
  {{- end }}
  {{- range .Compound.InnerDirs }}
  - "{{ (index $.H.CompoundIdMap .RefId).Title }}"
  {{- end }}
  {{- range .Compound.Sections }}
  - "{{ .Header }}"
    {{- range .Enums }}
  -  "{{ .Name }}"
      {{- range .Values }}
  - "{{ .Name }}"
      {{- end }}
    {{- end }}
    {{- range .Functions }}
  - "{{ .Name }}"
    {{- end }}
    {{- range .Attributes }}
  - "{{ .Name }}"
    {{- end }}
    {{- range .Defines }}
  - "{{ .Name }}"
    {{- end }}
    {{- range .Typedefs }}
  - "{{ .Name }}"
    {{- end }}
    {{- range .Friends }}
  - "{{ .Name }}"
    {{- end }}
  {{- end }}
---

{{ if .Compound.Location.File }}
<p>
	<a href="/{{.Section}}/file/{{ .Compound.Location.FileRefId }}">{{ .Compound.Location.File }}</a>
</p>
{{ end }}

{{ if .Compound.BriefDescription }}
{{ $.H.RenderDocstring .Compound.BriefDescription }}
{{ end }}

<p>
	<a href="#detailed_description">More...</a>
</p>

{{ if (len .Compound.InnerClasses) }}
<h2>Classes:</h2>
<div class="inner-compound-briefs">
{{ range $compound := .Compound.InnerClasses }}
{{ $.H.RenderInnerCompound $compound }}
{{ end }}
</div>
{{ end }}

{{ if (len .Compound.InnerNamespaces) }}
<h2>Namespaces:</h2>
<div class="inner-compound-briefs">
{{ range $compound := .Compound.InnerNamespaces }}
{{ $.H.RenderInnerCompound $compound }}
{{ end }}
</div>
{{ end }}

{{ if (len .Compound.InnerGroups) }}
<h2>Groups:</h2>
<div class="inner-compound-briefs">
{{ range $compound := .Compound.InnerGroups }}
{{ $.H.RenderInnerCompound $compound }}
{{ end }}
</div>
{{ end }}

{{ if (len .Compound.InnerFiles) }}
<h2>Files:</h2>
<div class="inner-compound-briefs">
{{ range $compound := .Compound.InnerFiles }}
{{ $.H.RenderInnerCompound $compound }}
{{ end }}
</div>
{{ end }}

{{ if (len .Compound.InnerDirs) }}
<h2>Dirs:</h2>
<div class="inner-compound-briefs">
{{ range .Compound.InnerDirs }}
{{ $.H.RenderInnerCompound . }}
{{ end }}
</div>
{{ end }}

{{ range .Compound.Sections }}
{{ $.H.RenderSectionBrief . }}
{{ end }} 

<a id="detailed_description"></a>
<h2>Detailed Description</h2>
{{ if .Compound.BriefDescription }}
{{ $.H.RenderDocstring .Compound.BriefDescription }}
{{ end }}
{{ if .Compound.DetailedDescription }}
{{ $.H.RenderDocstring .Compound.DetailedDescription }}
{{ end }}

{{ range .Compound.Sections }}
{{ $.H.RenderSection . }}
{{ end }}

{{ if .Compound.ProgramListing.Content }}
{{ $.H.RenderHighlightWithLineNos "C++" ($.H.RenderDocstring .Compound.ProgramListing) }}
{{ end }}`

const InnerCompound = `<div class="inner-compound-briefs__item">
	<div class="inner-compound-briefs__item__kind">
		{{ (index $.H.CompoundRefs .RefId).Kind }}
	</div>
	<div class="inner-compound-briefs__item__description">
		<div class="inner-compound-briefs__item__description__name">
			{{ $.H.RenderRef .RefId .Value}}
		</div>
		<div class="inner-compound-briefs__item__description__brief">
			{{ $.H.RenderDocstring (index $.H.CompoundIdMap .RefId).BriefDescription }}
		</div>
	</div>
</div>
`

const SectionBrief = `<div class="compound-section">
{{ with .Section.Header }}
<h2>{{ . }}</h2>
{{ end }}

{{ with .Section.Description }}
{{ $.H.RenderDocstring . }}
{{ end }}

{{ with .Section.Enums }}
<div class="section-briefs">
{{ range . }}
    <div class="section-briefs__item">
        <div class="section-briefs__item__kind">
            enum
        </div>
        <div class="section-briefs__item__description">
            <div class="section-briefs__item__description__name">
                {{ $name := .Name}}
                {{ if HasPrefix $name "@" }}
                {{ $name = "_Anonymous_" }}
                {{ end }}
				{{ $.H.RenderHighlight "C++" (printf "<a href=\"#%s\">%s</a> %s" .Id $name ($.H.RenderEnumBody .Values)) }}
            </div>
            <div class="section-briefs__item__description__brief">
				{{ $.H.RenderDocstring .BriefDescription }}
            </div>
        </div>
    </div>
{{ end }}
</div>
{{ end }}

{{ with .Section.Functions }}
<div class="section-briefs">
{{ range . }}
    <div class="section-briefs__item">
        <div class="section-briefs__item__kind">
            {{ $.H.RenderHighlight "C++" ($.H.RenderDocstring .Type) }}
        </div>
        <div class="section-briefs__item__description">
            <div class="section-briefs__item__description__name">
            	{{ $.H.RenderHighlight "C++" ($.H.RenderBriefFunctionDecl .) }}
            </div>
            <div class="section-briefs__item__description__brief">
				{{ $.H.RenderDocstring .BriefDescription }}
            </div>
        </div>
    </div>
{{ end }}
</div>
{{ end }}

{{ with .Section.Attributes }}
<div class="section-briefs">
{{ range . }}
    <div class="section-briefs__item">
        <div class="section-briefs__item__kind">
            {{ $.H.RenderHighlight "C++" ($.H.RenderDocstring .Type) }}
        </div>
        <div class="section-briefs__item__description">
            <div class="section-briefs__item__description__name">
            	{{ $.H.RenderHighlight "C++" (printf "<a href=\"#%s\">%s</a> %s" .Id .Name ($.H.RenderDocstring .ArgsString)) }}
            </div>
            <div class="section-briefs__item__description__brief">
				{{ $.H.RenderDocstring .BriefDescription }}
            </div>
        </div>
    </div>
{{ end }}
</div>
{{ end }}

{{ with .Section.Defines }}
<div class="section-briefs">
{{ range . }}
    <div class="section-briefs__item">
        <div class="section-briefs__item__kind">
            define
        </div>
        <div class="section-briefs__item__description">
            <div class="section-briefs__item__description__name">
            	{{ $.H.RenderHighlight "C++" ($.H.RenderBriefDefineDecl .) }}
            </div>
            <div class="section-briefs__item__description__brief">
				{{ $.H.RenderDocstring .BriefDescription }}
            </div>
        </div>
    </div>
{{ end }}
</div>
{{ end }}

{{ with .Section.Typedefs }}
<div class="section-briefs">
{{ range . }}
    <div class="section-briefs__item">
        <div class="section-briefs__item__kind">
            {{ $.H.RenderHighlight "C++" ($.H.RenderDocstring .Type) }}
        </div>
        <div class="section-briefs__item__description">
            <div class="section-briefs__item__description__name">
            	{{ $.H.RenderHighlight "C++" (printf "%s %s" .Name ($.H.RenderDocstring .ArgsString)) }}
            </div>
            <div class="section-briefs__item__description__brief">
				{{ $.H.RenderDocstring .BriefDescription }}
            </div>
        </div>
    </div>
{{ end }}
</div>
{{ end }}

{{ with .Section.Friends }}
<div class="section-briefs">
{{ range . }}
    <div class="section-briefs__item">
        <div class="section-briefs__item__kind">
            {{ $.H.RenderHighlight "C++" ($.H.RenderDocstring .Type) }}
        </div>
        <div class="section-briefs__item__description">
            <div class="section-briefs__item__description__name">
			{{ $.H.RenderHighlight "C++" (printf "<a href=\"#%s\">%s</a>" .Id .Name) }}
            </div>
            <div class="section-briefs__item__description__brief">
				{{ $.H.RenderDocstring .BriefDescription }}
            </div>
        </div>
    </div>
{{ end }}
</div>
{{ end }}
</div>
`

const Section = `<div class="compound-section">
{{ with .Section.Id }}
<a class="anchor" id="{{ . }}"></a>
{{ end }}

{{ with .Section.Header }}
<h2>{{ . }}</h2>
{{ end }}

{{ with .Section.Description }}
{{ $.H.RenderDocstring . }}
{{ end }}

{{ with .Section.Enums }}
{{ range . }}
<a class="anchor" id="{{ .Id }}"></a>
{{ $.H.RenderHighlight "C++" .Name }}
<h3>Enumerator</h3>
<dl class="enumerator">
	{{ range .Values }}
	<dt>{{ .Name }} <i>{{ .Initializer }}</i></dt>
	<dd>
		{{ $.H.RenderDocstring .BriefDescription }}
		{{ $.H.RenderDocstring .DetailedDescription }}
	</dd>
	{{ end }}
</dl>
{{ $.H.RenderDocstring .BriefDescription }}
{{ $.H.RenderDocstring .DetailedDescription }}
{{ end }}
{{ end }}

{{ with .Section.Functions }}
{{ range . }}
	<a class="anchor" id="{{ .Id }}"></a>
	{{ $.H.RenderHighlight "C++" ($.H.RenderBriefFunctionDecl .) }}
	
	<p>
	{{ if .Reimplements.RefId }}
		Reimplemented from: {{ $.H.RenderReimplementedFrom . }}
	{{ else }}
	{{ $.H.RenderDocstring .BriefDescription }}
	{{ $.H.RenderDocstring .DetailedDescription }}
	{{ end }}
	</p>
	
	<p>
	{{ if .ReimplementedBy }}
		Reimplemented by: {{ $.H.RenderReimplementedBy . }}
	{{ end }}
	</p>
{{ end }}
{{ end }}

{{ with .Section.Attributes }}
{{ range . }}
	<a class="anchor" id="{{ .Id }}"></a>
	{{ $.H.RenderHighlight "C++" (printf "%s %s %s" ($.H.RenderDocstring .Type) .Name ($.H.RenderDocstring .ArgsString)) }}

	{{ $.H.RenderDocstring .BriefDescription }}
	{{ $.H.RenderDocstring .DetailedDescription }}
{{ end }}
{{ end }}

{{ with .Section.Defines }}
{{ range . }}
	<a class="anchor" id="{{ .Id }}"></a>
	{{ $.H.RenderHighlight "C++" ($.H.RenderBriefDefineDecl .) }}

	{{ $.H.RenderDocstring .BriefDescription }}
	{{ $.H.RenderDocstring .DetailedDescription }}
{{ end }}
{{ end }}

{{ with .Section.Typedefs }}
{{ range . }}
	<a class="anchor" id="{{ .Id }}"></a>
	{{ $.H.RenderHighlight "C++" (printf "typedef %s %s %s" ($.H.RenderDocstring .Type) .Name ($.H.RenderDocstring .ArgsString)) }}

	{{ $.H.RenderDocstring .BriefDescription }}
	{{ $.H.RenderDocstring .DetailedDescription }}
{{ end }}
{{ end }}
`

const DocstringSection = `<div class="docstring-section">
{{ with .Id }}
<a class="anchor" id="{{ . }}"></a>
{{ end }}
{{ if .Kind }}
<div class="docstring-section__title docstring-section__title--{{ .Kind }}">
	{{ .Kind }}:
</div>
<div class="docstring-section__content">
	{{ $.H.RenderDocstring .Content}}
</div>
{{ else }}
{{ $.H.RenderDocstring .Content}}
{{ end }}
</div>`

const DocstringParameterList = `<b>Parameters:</b>
<table class="goxy-parameterlist">
	<tbody>
	{{ range .Items }}
		<tr>
			<td>{{ .Name }}</td>
			<td>{{ $.H.RenderDocstring .Description }}</td>
		</tr>
	{{ end }}
	</tbody>
</table>`

const DocstringTable = `<table>
	<tbody>
	{{ range .Rows }}
		<tr>
		{{ range . }}
			{{ if .Head }}
			<th>
			{{ else }}
			<td> 
			{{ end }}
			<td>{{ $.H.RenderDocstring .Content }}</td>
			{{ if .Head }}
			</th>
			{{ else }}
			</td> 
			{{ end }}
		{{ end }}
		</tr>
	{{ end }}
	</tbody>
</table>`
