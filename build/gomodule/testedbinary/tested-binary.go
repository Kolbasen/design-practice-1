package testedbinary

import (
	"fmt"
	"path"
	"regexp"

	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
)

var (
	// Package context used to define Ninja build rules.
	pctx = blueprint.NewPackageContext("github.com/Kolbasen/design-practice-1/build/gomodule/testedbinary")

	// Ninja rule to execute go build.
	goBuild = pctx.StaticRule("binaryBuild", blueprint.RuleParams{
		Command:     "cd $workDir && go build -o $outputPath $pkg",
		Description: "build go command $pkg",
	}, "workDir", "outputPath", "pkg")

	// Ninja rule to execute go mod vendor.
	goVendor = pctx.StaticRule("vendor", blueprint.RuleParams{
		Command:     "cd $workDir && go mod vendor",
		Description: "vendor dependencies of $name",
	}, "workDir", "name")

	// Ninja rule to execute go test.
	goTest = pctx.StaticRule("test", blueprint.RuleParams{
		Command:     "cd ${workDir} && go test -v ${testPkg} > ${outPath}",
		Description: "test ${testPkg}",
	}, "workDir", "outPath", "testPkg")
)

type testedBinaryModule struct {
	blueprint.SimpleName

	properties struct {
		// Go package name to build as a command with "go build".
		Pkg string
		// Go package name to test as a command with "go test".
		TestPkg string
		// List of source files.
		Srcs []string
		// Exclude patterns.
		SrcsExclude []string
		// If to call vendor command.
		VendorFirst bool
	}
}

func (tb *testedBinaryModule) GenerateBuildActions(ctx blueprint.ModuleContext) {
	name := ctx.ModuleName()
	config := bood.ExtractConfig(ctx)
	config.Debug.Printf("Adding build actions for go binary module '%s'", name)

	outputPath := path.Join(config.BaseOutputDir, "bin", name)

	testOutPath := path.Join(config.BaseOutputDir, "out.txt")

	var srcInputs []string
	var testInputs []string

	inputErors := false

	for _, src := range tb.properties.Srcs {
		if matches, err := ctx.GlobWithDeps(src, tb.properties.SrcsExclude); err == nil {
			for _, path := range matches {
				isTestFile, err := regexp.MatchString("^.*_test.go$", path)
				if err != nil {
					ctx.PropertyErrorf("srcs", "Error matching string %s", path)
					inputErors = true
					break
				}

				if isTestFile {
					testInputs = append(testInputs, path)
				} else {
					srcInputs = append(srcInputs, path)
				}
			}
		} else {
			ctx.PropertyErrorf("srcs", "Cannot resolve files that match pattern %s", src)
			inputErors = true
		}
	}

	if inputErors {
		return
	}

	if tb.properties.VendorFirst {
		vendorDirPath := path.Join(ctx.ModuleDir(), "vendor")
		ctx.Build(pctx, blueprint.BuildParams{
			Description: fmt.Sprintf("Vendor dependencies of %s", name),
			Rule:        goVendor,
			Outputs:     []string{vendorDirPath},
			Implicits:   []string{path.Join(ctx.ModuleDir(), "go.mod")},
			Optional:    true,
			Args: map[string]string{
				"workDir": ctx.ModuleDir(),
				"name":    name,
			},
		})
		srcInputs = append(srcInputs, vendorDirPath)
	}

	ctx.Build(pctx, blueprint.BuildParams{
		Description: fmt.Sprintf("Test module %s", tb.properties.TestPkg),
		Rule:        goTest,
		Outputs:     []string{testOutPath},
		Implicits:   append(srcInputs, testInputs...),
		Args: map[string]string{
			"outPath": testOutPath,
			"workDir": ctx.ModuleDir(),
			"testPkg": tb.properties.TestPkg,
		},
	})

	ctx.Build(pctx, blueprint.BuildParams{
		Description: fmt.Sprintf("Build %s as Go binary", name),
		Rule:        goBuild,
		Outputs:     []string{outputPath},
		Implicits:   srcInputs,
		Args: map[string]string{
			"outputPath": outputPath,
			"workDir":    ctx.ModuleDir(),
			"pkg":        tb.properties.Pkg,
		},
	})
}

func TestedBinaryFactory() (blueprint.Module, []interface{}) {
	mType := &testedBinaryModule{}
	return mType, []interface{}{&mType.SimpleName.Properties, &mType.properties}
}
