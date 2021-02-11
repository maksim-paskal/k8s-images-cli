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
package main

import (
	"bufio"
	"os"
	"regexp"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type ImageIgnore struct {
	rules    []string
	isLoaded bool
}

func getImageIgnore() (*ImageIgnore, error) {
	ig := ImageIgnore{
		rules:    make([]string, 0),
		isLoaded: false,
	}

	file, err := os.Open(*appConfig.ImageIgnoreFile)
	if err != nil {
		return &ig, errors.Wrap(err, "cant open file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ig.rules = append(ig.rules, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.WithError(err).Debug("some error in scaner")
	}

	ig.isLoaded = true

	return &ig, nil
}

func (ig *ImageIgnore) match(image string) bool {
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
