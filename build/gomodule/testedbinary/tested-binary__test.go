package testedbinary

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
)

func TestBinFactory(t *testing.T) {
	ctx := blueprint.NewContext()

	ctx.MockFileSystem(map[string][]byte{
		"Blueprints": []byte(`
			go_binary_tested {
			  name: "test-out",
			  srcs: ["test-src.go"],
			  pkg: ".",
			  testPkg: ".",
	          vendorFirst: true,
			  srcsExclude: ["**/*_test.go"],
			}
		`),
		"test-src.go": nil,
	})

	ctx.RegisterModuleType("go_binary_tested", TestedBinaryFactory)

	cfg := bood.NewConfig()

	_, errs := ctx.ParseBlueprintsFiles(".", cfg)
	if len(errs) != 0 {
		t.Fatalf("Syntax errors in the test blueprint file: %s", errs)
	}

	_, errs = ctx.PrepareBuildActions(cfg)
	if len(errs) != 0 {
		t.Errorf("Unexpected errors while preparing build actions: %s", errs)
	}
	buffer := new(bytes.Buffer)
	if err := ctx.WriteBuildFile(buffer); err != nil {
		t.Errorf("Error writing ninja file: %s", err)
	} else {
		text := buffer.String()
		t.Logf("Gennerated ninja build file:\n%s", text)
		if !strings.Contains(text, "out/bin/test-out: ") {
			t.Errorf("Generated ninja file does not have build of the test module")
		}
		if !strings.Contains(text, "test-src.go") {
			t.Errorf("Generated ninja file does not have source dependency")
		}
		if !strings.Contains(text, "build vendor: g.testedbinary.vendor | go.mod") {
			t.Errorf("Generated ninja file does not have vendor build rule")
		}
	}
}
