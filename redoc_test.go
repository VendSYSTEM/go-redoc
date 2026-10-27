package redoc_test

import (
	"embed"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/VendSYSTEM/go-redoc/v3"
	"github.com/stretchr/testify/assert"
)

func TestBodyRedocModuleAndSafeValues(t *testing.T) {
	url := `/spec?value="</script><script>alert(1)</script>&unicode=\` + "\u2028"
	r := redoc.Redoc{SpecPath: url, Title: `<title>&"`, Description: `<desc>&"`}
	body, err := r.Body()
	if err != nil {
		t.Fatal(err)
	}
	html := string(body)
	if !strings.Contains(redoc.JavaScript, "Version: 3.0.1") {
		t.Error("embedded bundle is not Redoc 3.0.1")
	}
	if !strings.Contains(html, `<script type="module">`) || strings.Contains(html, "Redoc.init(") {
		t.Error("page must initialize the Redoc module without the removed global")
	}
	if strings.Count(strings.ToLower(html), "</script>") != 1 {
		t.Error("inline scripts must not be terminated by the bundle or user values")
	}
	match := regexp.MustCompile(`init\(("[^\n]*?"), \{\}, document\.getElementById\("main"\)\)`).FindStringSubmatch(html)
	if len(match) != 2 {
		t.Fatal("module initialization with a JSON-encoded URL not found")
	}
	var got string
	if err := json.Unmarshal([]byte(match[1]), &got); err != nil || got != url {
		t.Errorf("spec URL round trip = %q, %v; want %q", got, err, url)
	}
	if !strings.Contains(html, `&lt;title&gt;&amp;&#34;`) || !strings.Contains(html, `&lt;desc&gt;&amp;&#34;`) {
		t.Error("title and description must be HTML-escaped")
	}
}

//go:embed testdata/spec.json
var spec embed.FS

func TestRedoc(t *testing.T) {
	r := redoc.Redoc{
		SpecFile: "testdata/spec.json",
		SpecFS:   &spec,
		SpecPath: "/openapi.json", // "/openapi.yaml"
		Title:    "Test API",
	}

	t.Run("Body", func(t *testing.T) {
		body, err := r.Body()
		assert.NoError(t, err)
		assert.Contains(t, string(body), r.Title)
	})

	t.Run("Handler", func(t *testing.T) {
		handler := r.Handler()

		t.Run("Spec", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
			w := httptest.NewRecorder()
			handler(w, req)

			resp := w.Result()
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Contains(t, string(body), `"swagger":"2.0"`)
		})

		t.Run("Docs", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()
			handler(w, req)

			resp := w.Result()
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.Equal(t, "text/html", resp.Header.Get("Content-Type"))

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Contains(t, string(body), r.Title)
		})
	})
}
