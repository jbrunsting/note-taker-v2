package html

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/jbrunsting/note-taker-v2/utils"
	html2md "github.com/russross/blackfriday/v2"
)

const (
	notesDirKey = "$NOTES"
)

func getClass(tag string) string {
	o := ""
	for _, c := range tag {
		if ('a' <= c && c <= 'z') || ('0' <= c && c <= '9') {
			o += string(c)
		} else if c == ' ' {
			o += "_"
		}
	}
	return o
}

func getId(title string) string {
	return "id_n_" + strings.ReplaceAll(title, " ", "_")
}

func getToggles(tags []string) string {
	html := ""
	for _, tag := range tags {
		html += fmt.Sprintf(
			"<input id=\"id_c_%[1]s\" class=\"%[1]s\" type=\"checkbox\" checked/>",
			getClass(tag),
		)
	}
	html += "<input id=\"id_dark_mode\" type=\"checkbox\"/>"
	html += "<div class=\"tag-selector\">"
	for _, tag := range tags {
		html += fmt.Sprintf(
			"<label for=\"id_c_%s\">%s</label>",
			getClass(tag),
			tag,
		)
	}
	html += "<label id=\"id_dark_mode_toggle\" for=\"id_dark_mode\">☀</label>"
	html += "</div>"
	return html
}

type OrderedTag struct {
	Tag   string
	Count int
}

func WriteHtml(dir string) error {
	files, err := utils.GetNotesFiles(dir)
	if err != nil {
		return err
	}
	htmlCode, err := GenerateHTML(files, dir)
	if err != nil {
		return err
	}
	htmlPath := fmt.Sprintf("%v/index.html", dir)
	_, err = os.Stat(htmlPath)
	if os.IsNotExist(err) {
		file, err := os.Create(htmlPath)
		if err != nil {
			return err
		}
		file.Close()
	}
	return ioutil.WriteFile(htmlPath, []byte(htmlCode), 0644)
}

func GenerateHTML(files []os.FileInfo, dir string) (string, error) {
	oTags := make(map[string]*OrderedTag)
	html := ""
	for _, file := range files {
		filename := file.Name()
		path := fmt.Sprintf("%v/%v", dir, filename)

		bmd, err := ioutil.ReadFile(path)
		if err != nil {
			return html, err
		}
		md := strings.Replace(string(bmd), notesDirKey, dir, -1)

		tagsRegexp := regexp.MustCompile(`@[a-zA-Z0-9]+\s`)
		tags := tagsRegexp.FindAll([]byte(md), -1)
		tagNames := []string{}
		for _, tag := range tags {
			tagName := string(tag)[1:]
			tagNames = append(tagNames, tagName)
			oldCount := 0
			if oTag, ok := oTags[tagName]; ok {
				oldCount = oTag.Count
			}
			oTags[tagName] = &OrderedTag{tagName, oldCount + 1}
		}
    }

	for _, file := range files {
		filename := file.Name()
		path := fmt.Sprintf("%v/%v", dir, filename)
		name := strings.TrimSuffix(filename, filepath.Ext(filename))

		bmd, err := ioutil.ReadFile(path)
		if err != nil {
			return html, err
		}
		md := strings.Replace(string(bmd), notesDirKey, dir, -1)

		tagsRegexp := regexp.MustCompile(`@[a-zA-Z0-9]+\s`)
		tags := tagsRegexp.FindAll([]byte(md), -1)
		tagNames := []string{}
		for _, tag := range tags {
			tagName := string(tag)[1:]
			tagNames = append(tagNames, tagName)
		}

        oTagNames := []string{}
        for oTagName := range oTags {
            hasTag := false
            for _, tagName := range tagNames {
                if tagName == oTagName {
                    hasTag = true
                    break
                }
            }
            if !hasTag {
                oTagNames = append(oTagNames, oTagName)
            }
        }

		classes := ""
        for _, oTagName := range oTagNames {
			oTagName = strings.ToLower(oTagName)
			classes += " " + getClass(oTagName)
        }

		tagHtml := "<div class=\"tag\">"
		for _, tag := range tagNames {
			tag = strings.ToLower(tag)
			tagHtml += fmt.Sprintf("<p>%s</p>", tag)
		}
		tagHtml += "</div>"

		noteHtml := string(html2md.Run(
			[]byte(md),
			html2md.WithNoExtensions(),
		))

		html += fmt.Sprintf("<div class=\"__note__ %s\" id=\"%s\">", classes, getId(name))
		html += "<div class=\"header\">"
		html += fmt.Sprintf("<a href=\"#%s\" class=\"note-header\">%s</a>", getId(name), name)
		html += tagHtml
		html += "</div>"
		html += noteHtml
		html += "</div>"
	}

	vals := []*OrderedTag{}
	for _, v := range oTags {
		vals = append(vals, v)
	}

	sort.SliceStable(vals, func(i, j int) bool {
		return vals[i].Count > vals[j].Count
	})

	tags := []string{}
	for _, ot := range vals {
		tags = append(tags, ot.Tag)
	}

	html = getToggles(tags) + "<div id=\"id_body\"><div id=\"id_content\">" + html + "</div></div>"
	return "<html>" + getStyle(tags) + "<body>" + html + "</body></html>", nil
}

func getStyle(tags []string) string {
	// We just make the CSS a big string so we can easily construct a single
	// html file that displays the notes, without relying on reading from an
	// external css file
	css := `
html {
    height: 100%;
}

body {
	margin: 0;
    font-family: Arial, Helvetica, sans-serif;
    height: 100%;
}

#id_body {
	width: 100%;
    padding: 0px;
	margin: 0px;
	margin-top: -100px;
	padding-top: 100px;
	min-height: 100%;
}

#id_content {
	margin: 0px auto;
	max-width: 800px;
}

p {
    margin: 5px 0px;
	font-size: 0.9em;
}

h1 {
    margin: 5px 0px;
    font-size: 1.6em;
}

h2 {
    margin: 5px 0px;
    font-size: 1.45em;
}

h3 {
    margin: 5px 0px;
    font-size: 1.3em;
}

h4 {
    margin: 5px 0px;
    font-size: 1.2em;
}

h5 {
    margin: 5px 0px;
    font-size: 1.1em;
}

h6 {
    margin: 5px 0px;
    font-size: 1em;
}

a {
	color: #6D9D99;
}

input {
    display: none;
}

label {
    margin: 0px 10px 0px 0px;
    padding: 3px 7px;
    border-radius: 3px;
    white-space: nowrap;
}

label:hover {
    cursor: pointer;
}

div.tag-selector {
	display: flex;
    overflow-x: auto;
    padding: 5px 10px;
	margin: 10px auto 0px auto;
	max-width: 800px;
}

#id_dark_mode_toggle {
    margin-right: 0px;
	margin-left: auto;
    box-shadow: 0px 0px 5px rgba(0, 0, 0, 0.5);
}

.__note__ {
    margin: 10px 0px;
    padding: 10px;
    border-radius: 3px;
    box-shadow: 0px 0px 5px rgba(0, 0, 0, 0.5);
}

.__note__ * {
	max-width: 100%;
}

.__note__ img {
	max-height: 450px;
	margin: auto;
	display: block;
	padding: 5px;
}

.header {
    overflow: auto;
    padding: 0px 5px 7px 0px;
	border-bottom: 1px solid #2E2E2E;
}

.note-header {
	cursor: pointer;
    font-weight: bold;
    margin: 1px 0px;
    font-size: 1em;
    display: inline-block;
	text-decoration: none;
}

.tag {
	display: inline-block;
	float: right;
    margin: 0px -5px;
    font-size: 0.8em;
}

.tag p {
    font-size: 1em;
    display: inline-block;
    margin: 0px 0px 0px 10px;
    padding: 1px 5px;
    border-radius: 3px;
}`
	// TODO: Should refactor this
	css += `
#id_body {
    color: #2E2E2E;
    background-color: #F4EFE5;
}

label {
    color: #F4EFE5;
    background-color: #6D9D99;
}

#id_dark_mode_toggle {
    color: #2E2E2E;
    background-color: #FAF8F3;
}

.__note__ {
    background-color: #FAF8F3;
}

.tag p {
	border: 2px solid;
	border-color: #6D9D99;
}

#id_dark_mode:checked ~ #id_body {
    color: #D1D1D1;
    background-color: #05070C;
}

#id_dark_mode:checked ~ .tag-selector label {
    color: #0B101A;
}

#id_dark_mode:checked ~ .tag-selector #id_dark_mode_toggle {
    color: #D1D1D1;
    background-color: #2E2E2E;
}

#id_dark_mode:checked ~ #id_body .__note__ {
    background-color: #2E2E2E;
}

#id_dark_mode:checked ~ #id_body .header {
	border-bottom: 1px solid #FAF8F3;
}
`
	for _, tag := range tags {
		css += fmt.Sprintf(`
input.%[1]s ~ #id_body div.%[1]s {
	display: block;
}
input.%[1]s:not(:checked) ~ #id_body div.%[1]s {
    display: none
}
input.%[1]s:checked ~ div > label[for=id_c_%[1]s] {
	background-color: #BFC9BC;
}
input.%[1]s:checked ~ #id_dark_mode:checked ~ div > label[for=id_c_%[1]s] {
	background-color: #556b69;
}
`,
			getClass(tag),
		)
	}
	return "<style>" + css + "</style>"
}
