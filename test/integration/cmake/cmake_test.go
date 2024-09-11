package cmake

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/objectbox/objectbox-generator/v4/test/assert"
	"github.com/objectbox/objectbox-generator/v4/test/cmake"
)

// excluded: cpp-multiple-targets,cpp-tree-multiple-targets (Visual Studio and XCode)
var projectsFlagString = flag.String("projects", "cpp-flat,cpp-tree,cpp-multiple-schema-dirs", "specify a subset of test projects (defaults to all)")
var singleGeneratorsFlagString = flag.String("singlegenerators", "", "specify comma-separated list of CMake single-config Generators (e.g. Unix Makefiles)")
var multiGeneratorsFlagString = flag.String("multigenerators", "", "specify comma-separated list of CMake multi-config Generators (e.g. Ninja Multi-Config)")
var autoDetectFlag = flag.Bool("autodetect", false, "auto-detect build tools (make, ninja)")
var noDefaultGeneratorFlag = flag.Bool("nodefaultgenerator", false, "do not run with default generator")
var noCleanUpFlag = flag.Bool("nocleanup", false, "leave temporary build directory in tmp folder for diagnostics")

func TestCMakeProjects(t *testing.T) {

	wd, err := os.Getwd()
	if err != nil {
		return
	}
	topdir := filepath.Join(wd, "..", "..", "..")

	var projects = []string{}
	var singleGenerators = []string{}
	var multiGenerators = []string{}

	if len(*projectsFlagString) > 0 {
		projects = strings.Split(*projectsFlagString, ",")
	}
	if len(*singleGeneratorsFlagString) > 0 {
		singleGenerators = strings.Split(*singleGeneratorsFlagString, ",")
	}
	if len(*multiGeneratorsFlagString) > 0 {
		multiGenerators = strings.Split(*multiGeneratorsFlagString, ",")
	}
	if *noDefaultGeneratorFlag == false {
		if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
			multiGenerators = append(multiGenerators, "default")
		} else {
			singleGenerators = append(singleGenerators, "default")
		}
	}
	if *autoDetectFlag {
		if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
			_, err := exec.LookPath("make")
			if err == nil {
				singleGenerators = append(singleGenerators, "Unix Makefiles")
			}
		}
		_, err = exec.LookPath("ninja")
		if err == nil {
			singleGenerators = append(singleGenerators, "Ninja")
			multiGenerators = append(multiGenerators, "Ninja Multi-Config")
		}
	}

	t.Logf("Projects: %s", projects)
	t.Logf("Single Generators: '%s'", singleGenerators)
	t.Logf("Multi Generators: '%s'", multiGenerators)

	tempRoot, err := ioutil.TempDir("", "objectbox-generator-test-integration-cmake")
	assert.NoErr(t, err)
	assert.True(t, len(tempRoot) > 0)
	if *noCleanUpFlag == false {
		defer func() {
			assert.NoErr(t, os.RemoveAll(tempRoot))
		}()
	} else {
		// Reprint the location on log at end of session on no cleanup
		defer func() {
			t.Logf("Temporary root build directory: %s", tempRoot)
		}()
	}
	t.Logf("Temporary root build directory: %s", tempRoot)

	var run = func(generator string, multiConfigBuild bool) {
		for _, project := range projects {
			var inSourceVariants = []bool{false, true}
			for _, inSource := range inSourceVariants {
				var variant = "default"
				var configureFlags = []string{fmt.Sprintf("-DTOPDIR=%s", topdir)}
				if inSource {
					variant = "insource"
					configureFlags = append(configureFlags, "-DDO_INSOURCE=TRUE")
				}
				t.Logf("Running CMake Test Project '%s' with CMake Generator '%s' (variant %s)", project, generator, variant)

				generatorLabel := strings.ReplaceAll(generator, " ", "_")

				var templateProjectDir = filepath.Join(wd, "projects", project)
				var testRootDir = filepath.Join(tempRoot, generatorLabel, project, variant)
				var confDir = filepath.Join(testRootDir, "src")
				var buildDir = filepath.Join(testRootDir, "build")

				t.Logf("TemplateProjectDir: %s", templateProjectDir)
				t.Logf("ConfDir: %s", confDir)
				t.Logf("BuildDir: %s", buildDir)
				t.Logf("Topdir: %s", topdir)
				t.Logf("Configure Flags: %s", configureFlags)

				cmake.CopyDir(wd, templateProjectDir, confDir)
				cmake.CopyFile(wd, "projects/common.cmake", testRootDir)
				var conf = &cmake.Cmake{
					Name:           "all",
					ConfDir:        confDir,
					BuildDir:       buildDir,
					Generator:      generator,
					ConfigureFlags: configureFlags,
				}
				err := os.MkdirAll(conf.BuildDir, 0750)
				if err != nil {
					t.Errorf("Failed to create build dir: %s", conf.BuildDir)
				} else {
					if stdOut, stdErr, err := conf.ConfigureRaw(); err != nil {
						t.Fatalf("cmake configuration failed: \n%s\n%s\n%s", stdOut, stdErr, err)
					} else {
						t.Logf("configuration output:\n%s", string(stdOut))
					}
					if multiConfigBuild {
						configs := []string{"Release", "Debug"}
						for _, config := range configs {
							if stdOut, stdErr, err := conf.BuildDefaultsWithConfig(config); err != nil {
								t.Fatalf("cmake build (configuration %s) failed: \n%s\n%s\n%s", config, stdOut, stdErr, err)
							} else {
								t.Logf("build (configuration %s) output:\n%s", config, string(stdOut))
							}
							if stdOut, stdErr, err := conf.BuildTargetWithConfig(config, "clean"); err != nil {
								t.Fatalf("cmake build clean (configuration %s) failed: \n%s\n%s\n%s", config, stdOut, stdErr, err)
							} else {
								t.Logf("clean (configuration %s) output:\n%s", config, string(stdOut))
							}
							if stdOut, stdErr, err := conf.BuildDefaultsWithConfig(config); err != nil {
								t.Fatalf("cmake build clean (configuration %s) failed: \n%s\n%s\n%s", config, stdOut, stdErr, err)
							} else {
								t.Logf("rebuild (configuration %s) output:\n%s", config, string(stdOut))
							}
						}
					} else {
						if stdOut, stdErr, err := conf.BuildDefaults(); err != nil {
							t.Fatalf("cmake build failed: \n%s\n%s\n%s", stdOut, stdErr, err)
						} else {
							t.Logf("build output:\n%s", string(stdOut))
						}
						if stdOut, stdErr, err := conf.BuildWithTarget("clean"); err != nil {
							t.Fatalf("cmake build clean failed: \n%s\n%s\n%s", stdOut, stdErr, err)
						} else {
							t.Logf("clean output:\n%s", string(stdOut))
						}
						if stdOut, stdErr, err := conf.BuildDefaults(); err != nil {
							t.Fatalf("cmake build clean failed: \n%s\n%s\n%s", stdOut, stdErr, err)
						} else {
							t.Logf("rebuild output:\n%s", string(stdOut))
						}
					}
					// Test updating for just "cpp-flat" using content from "update-cpp-flat" folder.
					if project == "cpp-flat" {
						t.Logf("**** Update project cpp-flat (task.fbs and main.cpp) ****")
						// depending on file system and build tools its possible the previous
						// rebuild and copying  happens too fast for the build system to detect changes,
						// -> delay after the last fbs generation before we 'create' the changed .fbs file
						time.Sleep(500 * time.Millisecond)
						cmake.CopyDir(wd, filepath.Join("update-cpp-flat", "."), filepath.Join(confDir))

						if multiConfigBuild {
							configs := []string{"Release", "Debug"}
							for _, config := range configs {
								if stdOut, stdErr, err := conf.BuildDefaultsWithConfig(config); err != nil {
									t.Fatalf("cmake build after update (configuration %s) failed: \n%s\n%s\n%s", config, stdOut, stdErr, err)
								} else {
									t.Logf("build after update (configuration %s) output:\n%s", config, string(stdOut))
								}
							}
						} else {
							if stdOut, stdErr, err := conf.BuildDefaults(); err != nil {
								t.Fatalf("cmake build after update failed: \n%s\n%s\n%s", stdOut, stdErr, err)
							} else {
								t.Logf("build after update output:\n%s", string(stdOut))
							}
						}
					}
				}
			}
		}

	}

	for _, singleGenerator := range singleGenerators {
		run(singleGenerator, false)
	}
	for _, multiGenerator := range multiGenerators {
		run(multiGenerator, true)
	}
}

func TestCMakeProjectModify(t *testing.T) {
}
