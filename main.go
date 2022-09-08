// Copyright Nitric Pty Ltd.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"golang.org/x/exp/slices"
	"golang.org/x/mod/modfile"
)

type stringCommaList struct {
	value []string
}

func (s *stringCommaList) String() string {
	return strings.Join(s.value, ",")
}

func (s *stringCommaList) Set(v string) error {
	s.value = strings.Split(v, ",")
	return nil
}

func main() {
	var ignoreList = &stringCommaList{}

	flag.Var(ignoreList, "i", "module to ignore")

	flag.Parse()

	b, err := os.ReadFile("go.mod")
	if err != nil {
		panic(err)
	}

	mf, err := modfile.Parse("go.mod", b, nil)
	if err != nil {
		panic(err)
	}

	mods := []string{}

	for _, r := range mf.Require {
		if slices.Contains(ignoreList.value, r.Mod.Path) {
			continue
		}

		if r.Indirect {
			continue // only update directly required modules
		}

		mods = append(mods, r.Mod.Path)
	}

	fmt.Print(strings.Join(mods, " "))
}
