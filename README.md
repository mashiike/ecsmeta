# ecsmeta
This is a utility module for getting ECS ​​metadata.
Mainly used to embed ECS metadata in text/template (or http://github.com/kayac/go-config ).

See below for ECS metadata.
https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-metadata-endpoint.html

This respects https://github.com/fujiwara/tfstate-lookup.

## Usage (Go package)

See details in [godoc](https://godoc.org/github.com/mashiike/ecsmeta).

```go
package main

import (
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/mashiike/ecsmeta"
)

func main() {
    const templateText = `Image: {{ecsmeta ".Containers[0].Image"}}`
	funcMap := ecsmeta.MustFuncMap()

	tmpl, err := template.New("txt").Funcs(funcMap).Parse(templateText)
	if err != nil {
		log.Fatalf("parsing: %s", err)
	}
	err = tmpl.Execute(os.Stdout, nil)
	if err != nil {
		log.Fatalf("execution: %s", err)
	}
}
```

## LICENSE

[Mozilla Public License Version 2.0](LICENSE)
