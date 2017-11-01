// Copyright 2017 Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package build

import (
	"android/soong/finder"
	"android/soong/fs"
	"android/soong/ui/logger"
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// This file provides an interface to the Finder for use in Soong UI
// This file stores configuration information about which files to find

// NewSourceFinder returns a new Finder configured to search for source files.
// Callers of NewSourceFinder should call <f.Shutdown()> when done
func NewSourceFinder(ctx Context, config Config) (f *finder.Finder) {
	ctx.BeginTrace("find modules")
	defer ctx.EndTrace()

	dir, err := os.Getwd()
	if err != nil {
		ctx.Fatalf("No working directory for module-finder: %v", err.Error())
	}
	cacheParams := finder.CacheParams{
		WorkingDirectory: dir,
		RootDirs:         []string{"."},
		ExcludeDirs:      []string{".git", ".repo"},
		PruneFiles:       []string{".out-dir", ".find-ignore"},
		IncludeFiles:     []string{"Android.mk", "Android.bp", "Blueprints", "CleanSpec.mk"},
	}
	dumpDir := config.FileListDir()
	f, err = finder.New(cacheParams, fs.OsFs, logger.New(ioutil.Discard),
		filepath.Join(dumpDir, "files.db"))
	if err != nil {
		ctx.Fatalf("Could not create module-finder: %v", err)
	}
	return f
}

// FindSources searches for source files known to <f> and writes them to the filesystem for
// use later.
func FindSources(ctx Context, config Config, f *finder.Finder) {
	// note that dumpDir in FindSources may be different than dumpDir in NewSourceFinder
	// if a caller such as multiproduct_kati wants to share one Finder among several builds
	dumpDir := config.FileListDir()
	os.MkdirAll(dumpDir, 0777)

	androidMks := f.FindFirstNamedAt(".", "Android.mk")
	err := dumpListToFile(androidMks, filepath.Join(dumpDir, "Android.mk.list"))
	if err != nil {
		ctx.Fatalf("Could not export module list: %v", err)
	}

	cleanSpecs := f.FindFirstNamedAt(".", "CleanSpec.mk")
	dumpListToFile(cleanSpecs, filepath.Join(dumpDir, "CleanSpec.mk.list"))
	if err != nil {
		ctx.Fatalf("Could not export module list: %v", err)
	}

	isBlueprintFile := func(dir finder.DirEntries) (dirs []string, files []string) {
		files = []string{}
		for _, file := range dir.FileNames {
			if file == "Android.bp" || file == "Blueprints" {
				files = append(files, file)
			}
		}

		return dir.DirNames, files
	}
	androidBps := f.FindMatching(".", isBlueprintFile)
	err = dumpListToFile(androidBps, filepath.Join(dumpDir, "Android.bp.list"))
	if err != nil {
		ctx.Fatalf("Could not find modules: %v", err)
	}
}

func dumpListToFile(list []string, filePath string) (err error) {
	desiredText := strings.Join(list, "\n")
	desiredBytes := []byte(desiredText)
	actualBytes, readErr := ioutil.ReadFile(filePath)
	if readErr != nil || !bytes.Equal(desiredBytes, actualBytes) {
		err = ioutil.WriteFile(filePath, desiredBytes, 0777)
	}
	return err
}
