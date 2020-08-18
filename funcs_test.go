package ecsmeta_test

import (
	"os"
	"testing"
	"log"

	"github.com/mashiike/ecsmeta"
)

func TestMustFuncMap(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	os.Setenv("ECS_CONTAINER_METADATA_URI", server.URL + "/v4")
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
