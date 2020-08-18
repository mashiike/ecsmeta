package ecsmeta_test

import (
	"bytes"
	"log"
	"os"
	"testing"
	"text/template"

	"github.com/mashiike/ecsmeta"
)

func TestMustFuncMap(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	os.Setenv("ECS_CONTAINER_METADATA_URI", server.URL+"/v4")
	funcMap := ecsmeta.MustFuncMapWithName(
		"myfunc",
		ecsmeta.WithLogger(log.New(os.Stderr, "", log.LstdFlags)),
	)
	fn := funcMap["myfunc"].(func(string) string)
	if fn == nil {
		t.Error("no function")
	}
	if val := fn(".Containers[0].DockerName"); val != "query-metadata" {
		t.Errorf("unexpected DockerName: %s", val)
	}
	defer func() {
		err := recover()
		if err == nil {
			t.Error("must be panic")
		}
	}()
	fn("invalid query")
}

func TestTemplate(t *testing.T) {
	server := newTestServer()
	defer server.Close()
	os.Setenv("ECS_CONTAINER_METADATA_URI", server.URL+"/v4")
	const templateText = `Image: {{ecsmeta ".Containers[0].Image"}}`
	funcMap := ecsmeta.MustFuncMap()

	tmpl, err := template.New("txt").Funcs(funcMap).Parse(templateText)
	if err != nil {
		t.Fatalf("parsing: %s", err)
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, nil)
	if err != nil {
		t.Fatalf("execution: %s", err)
	}
	if buf.String() != "Image: mreferre/eksutils" {
		t.Fatalf("result: %s", buf.String())
	}
}
