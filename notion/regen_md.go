package notion

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/kjk/u"
)

var (
	mdOutWhitelist = make(map[string]bool)
)

func isMarkdownFile(path string) bool {
	path = strings.ToLower(path)
	return strings.HasSuffix(path, ".md")
}

func getFilesRecur(dir string, shouldInclude func(s string) bool) ([]string, error) {
	var res []string
	dirsToVisit := []string{dir}
	for len(dirsToVisit) > 0 {
		dir := dirsToVisit[0]
		dirsToVisit = dirsToVisit[1:]
		fileInfos, err := ioutil.ReadDir(dir)
		if err != nil {
			return res, err
		}
		for _, fi := range fileInfos {
			path := filepath.Join(dir, fi.Name())
			if fi.IsDir() {
				dirsToVisit = append(dirsToVisit, path)
				continue
			}
			if fi.Mode().IsRegular() {
				if (shouldInclude == nil) || shouldInclude(path) {
					res = append(res, path)
				}
			}
		}

	}
	return res, nil
}

func parseMd(d []byte) ([]byte, map[string]string) {
	meta := make(map[string]string)
	d = u.NormalizeNewlines(d)
	lines := bytes.Split(d, []byte{10})
	for len(lines) > 0 {
		l := lines[0]
		parts := bytes.SplitN(l, []byte{':'}, 2)
		if len(parts) != 2 {
			break
		}
		k := string(bytes.TrimSpace(parts[0]))
		k = capitalize(k)
		v := strings.TrimSpace(string(parts[1]))
		meta[k] = v
		lines = lines[1:]
	}

	// remove empty lines at the top
	for len(lines) > 0 {
		l := lines[0]
		if len(l) > 0 {
			break
		}
		lines = lines[1:]
	}

	d = bytes.Join(lines, []byte{10})
	return d, meta
}

const keyTitle = "Title"

func mdToHTML(mdFile, templateFile, htmlFile string) {
	md, err := ioutil.ReadFile(mdFile)
	must(err)
	md, meta := parseMd(md)
	body := markdownToHTML(md, "")

	model := make(map[string]interface{})
	title := meta[keyTitle]
	model[keyTitle] = title
	model["BodyHTML"] = template.HTML(body)

	templateName := filepath.Base(templateFile)
	templates = template.Must(template.ParseFiles(templateFile))
	var buf bytes.Buffer
	err = templates.ExecuteTemplate(&buf, templateName, model)
	must(err)
	err = ioutil.WriteFile(htmlFile, buf.Bytes(), 0644)
	must(err)
	verbose("%s => %s\n", mdFile, htmlFile)
}

// template is a _md.tmpl.html file in the same directory
func findMdTemplate(mdFile string) string {
	dir := filepath.Dir(mdFile)
	path := filepath.Join(dir, "_md.tmpl.html")
	_, err := os.Stat(path)
	if err != nil {
		return ""
	}
	return path
}

func regenMd() {
	mdFiles, err := getFilesRecur("www", isMarkdownFile)
	must(err)
	for _, mdFile := range mdFiles {
		htmlFile := replaceExt(mdFile, ".html")
		templateFile := findMdTemplate(mdFile)
		if templateFile == "" {
			verbose("%s : skipping because no template file\n", mdFile)
			continue
		}
		verbose("%s\n", mdFile)
		mdToHTML(mdFile, templateFile, htmlFile)
		verbose("Whitelisted: %s\n", htmlFile)
		mdOutWhitelist[htmlFile] = true
	}
}
