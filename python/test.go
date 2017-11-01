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

package python

import (
	"android/soong/android"
)

// This file contains the module types for building Python test.

func init() {
	android.RegisterModuleType("python_test_host", PythonTestHostFactory)
}

type testDecorator struct {
	*binaryDecorator
}

func (test *testDecorator) install(ctx android.ModuleContext, file android.Path) {
	test.binaryDecorator.baseInstaller.install(ctx, file)
}

func NewTest(hod android.HostOrDeviceSupported) *Module {
	module, binary := NewBinary(hod)

	binary.baseInstaller = NewPythonInstaller("nativetest")

	test := &testDecorator{binaryDecorator: binary}

	module.bootstrapper = test
	module.installer = test

	return module
}

func PythonTestHostFactory() android.Module {
	module := NewTest(android.HostSupportedNoCross)

	return module.Init()
}
