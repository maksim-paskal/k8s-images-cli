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
package imageignore

import (
	"bufio"
	"os"
	"regexp"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Type struct {
	rules    []string
	isLoaded bool
}

func New(imageIgnoreFile string) (*Type, error) {
	imageIgnore := Type{
		rules:    make([]string, 0),
		isLoaded: false,
	}

	file, err := os.Open(imageIgnoreFile)
	if err != nil {
		return &imageIgnore, errors.Wrap(err, "cant open file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		imageIgnore.rules = append(imageIgnore.rules, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.WithError(err).Debug("some error in scaner")
	}

	imageIgnore.isLoaded = true

	return &imageIgnore, nil
}

func (ig *Type) Match(image string) bool {
	if ig.isLoaded && len(ig.rules) > 0 {
		for _, ignoreRule := range ig.rules {
			isMatch, err := regexp.MatchString(ignoreRule, image)
			if err != nil {
				log.WithError(err).Error("error in regexp.MatchString")
			}

			if isMatch {
				return true
			}
		}
	}

	return false
}
