/*
Copyright paskal.maksim@gmail.com
Licensed under the Apache License, Version 2.0 (the "License")
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package imageignore_test

import (
	"testing"

	"github.com/maksim-paskal/k8s-images-cli/pkg/imageignore"
)

func TestGetImageIgnore(t *testing.T) {
	t.Parallel()

	imageignore, err := imageignore.New("./testdata/imageignore.test")
	if err != nil {
		t.Fatal(err)
	}

	mustMatch := []string{
		"test",
		"testtest",
		"docker.io/test",
		"asddsfdsfsdf-end",
	}

	for _, v := range mustMatch {
		if !imageignore.Match(v) {
			t.Fatalf("image must match %s", v)
		}
	}

	mustNotMatch := []string{
		"asdffasdfasfasf",
		"papapap.odpdpd/asdasdsd.pds",
		"papapap.s",
		"99393993#ssdsds",
	}

	for _, v := range mustNotMatch {
		if imageignore.Match(v) {
			t.Fatalf("image must not match %s", v)
		}
	}
}
