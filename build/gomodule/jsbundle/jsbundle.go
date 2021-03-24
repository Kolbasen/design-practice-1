package jsbundle

import (
	"fmt"
	"path"
	"strconv"

	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
)

var (
	// Package context used to define Ninja build rules.
	pctx = blueprint.NewPackageContext("github.com/Kolbasen/design-practice-1/build/gomodule/jsbundle")

	// Ninja rule to execute build
	jsBuild = pctx.StaticRule("build", blueprint.RuleParams{
		Command:     "cd $workDir/js && npx webpack --env ENTRY=${entry} --env SHOULD_OBFUSCATE=${shouldObfuscate} FILENAME=${name} --config=webpack.config.js",
		Description: "build js bundle",
	}, "entry", "shouldObfuscate", "name", "workDir")
)

type testedBinaryModule struct {
	blueprint.SimpleName

	properties struct {
		// List of source files.
		Srcs []string
		// If to call obfuscation.
		Obfuscate bool
		// Path to srcs
		Path string
	}
}

func (tb *testedBinaryModule) GenerateBuildActions(ctx blueprint.ModuleContext) {
	name := ctx.ModuleName()
	config := bood.ExtractConfig(ctx)
	config.Debug.Printf("Adding build actions for js bundle module '%s'", name)

	outPath := path.Join(config.BaseOutputDir, tb.properties.Path)

	var resultFiles = ""

	inputErors := false

	for _, src := range tb.properties.Srcs {
		if _, err := ctx.GlobWithDeps(src, []string{}); err == nil {
			addStr := "," + src
			if len(resultFiles) == 0 {
				addStr = src
			}
			resultFiles += addStr

		} else {
			ctx.PropertyErrorf("srcs", "Cannot resolve files that match pattern %s", src)
			inputErors = true
		}
	}

	if inputErors {
		return
	}

	ctx.Build(pctx, blueprint.BuildParams{
		Description: fmt.Sprintf("Js bundle of %s", name),
		Rule:        jsBuild,
		Outputs:     []string{outPath},
		Args: map[string]string{
			"workDir":         ctx.ModuleDir(),
			"name":            name,
			"entry":           resultFiles,
			"shouldObfuscate": strconv.FormatBool(tb.properties.Obfuscate),
		},
	})
}

func JsBundleFactory() (blueprint.Module, []interface{}) {
	mType := &testedBinaryModule{}
	return mType, []interface{}{&mType.SimpleName.Properties, &mType.properties}
}
