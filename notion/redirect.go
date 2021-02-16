package notion

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kjk/u"
)

var articleRedirectsTxt = ``

var redirects = [][]string{
	{"/index.html", "/"},
	{"/blog", "/"},
	{"/blog/", "/"},
}

var articleRedirects = make(map[string]string)

func readRedirects(store *Articles) {
	d := []byte(articleRedirectsTxt)
	lines := bytes.Split(d, []byte{'\n'})
	for _, l := range lines {
		if len(l) == 0 {
			continue
		}
		parts := strings.Split(string(l), "|")
		panicIf(len(parts) != 2, "malformed article_redirects.txt, len(parts) = %d (!2)", len(parts))
		idStr := parts[0]
		url := strings.TrimSpace(parts[1])
		idNum, err := strconv.Atoi(idStr)
		panicIf(err != nil, "malformed line in article_redirects.txt. Line:\n%s\nError: %s\n", l, err)
		id := u.EncodeBase64(idNum)
		a := store.idToArticle[id]
		if a != nil {
			articleRedirects[url] = id
			continue
		}
		//verbose("skipping redirect '%s' because article with id %d no longer present\n", string(l), id)
	}
}

var (
	netlifyRedirects []*netlifyRedirect
)

type netlifyRedirect struct {
	from string
	to   string
	// valid code is 301, 302, 200, 404
	code int
}

func netlifyAddRedirect(from, to string, code int) {
	r := netlifyRedirect{
		from: from,
		to:   to,
		code: code,
	}
	netlifyRedirects = append(netlifyRedirects, &r)
}

func netlifyAddRewrite(from, to string) {
	netlifyAddRedirect(from, to, 200)
}

func netflifyAddTempRedirect(from, to string) {
	netlifyAddRedirect(from, to, 302)
}

func netlifyAddStaticRedirects() {
	for _, redirect := range redirects {
		from := redirect[0]
		to := redirect[1]
		netflifyAddTempRedirect(from, to)
	}
}

func netlifyAddArticleRedirects(store *Articles) {
	for from, articleID := range articleRedirects {
		from = "/" + from
		article := store.idToArticle[articleID]
		panicIf(article == nil, "didn't find article for id '%s'", articleID)
		to := article.URL()
		netflifyAddTempRedirect(from, to) // TODO: change to permanent
	}

}

// redirect /article/:id/* => /article/:id/pretty-title
const netlifyRedirectsProlog = `/article/:id/*	/article/:id.html	200
`

func netlifyWriteRedirects() {
	buf := bytes.NewBufferString(netlifyRedirectsProlog)
	for _, r := range netlifyRedirects {
		s := fmt.Sprintf("%s\t%s\t%d\n", r.from, r.to, r.code)
		buf.WriteString(s)
	}
	netlifyWriteFile("_redirects", buf.Bytes())
}

// https://caddyserver.com/tutorial/caddyfile
// redirect /article/:id/* => /article/:id/pretty-title
var caddyProlog = `localhost:8080
root netlify_static
errors stdout
log stdout
rewrite / {
	r  ^/article/(.*)/.*$
	to /article/{1}.html
}
`

func isRewrite(r *netlifyRedirect) bool {
	return (r.code == 200) || strings.HasSuffix(r.from, "*")
}

func genCaddyRedir(r *netlifyRedirect) string {
	if r.from == "/" {
		return fmt.Sprintf("rewrite / %s\n", r.to)
	}
	if isRewrite(r) {
		// hack: caddy doesn't like `++` in from
		if strings.Contains(r.from, "++") {
			return ""
		}
		if strings.HasSuffix(r.from, "*") {
			base := strings.TrimSuffix(r.from, "*")
			to := strings.Replace(r.to, ":splat", "{1}", -1)
			return fmt.Sprintf(`
rewrite "%s" {
    regexp (.*)
    to %s
}
`, base, to)
		}
		return fmt.Sprintf(`
rewrite "^%s$" {
    to %s
}
`, r.from, r.to)
	}

	return fmt.Sprintf("redir \"%s\" \"%s\" %d\n", r.from, r.to, r.code)
}

func writeCaddyConfig() {
	path := filepath.Join("Caddyfile")
	f, err := os.Create(path)
	must(err)
	defer f.Close()

	_, err = f.Write([]byte(caddyProlog))
	must(err)
	for _, r := range netlifyRedirects {
		s := genCaddyRedir(r)
		_, err = io.WriteString(f, s)
		must(err)
	}
}
