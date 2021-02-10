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
	"flag"
	"fmt"
	"os"
	"sort"

	log "github.com/sirupsen/logrus"
)

//nolint:gochecknoglobals
var gitVersion string = "dev"

func main() {
	flag.Parse()

	logLevel, err := log.ParseLevel(*appConfig.LogLevel)
	if err != nil {
		log.WithError(err).Fatal("error parse level")
	}

	log.SetLevel(logLevel)

	if *appConfig.showVersion {
		os.Stdout.WriteString(appConfig.Version)
		os.Exit(0)
	}

	clientset, err := getKubernetesClient()
	if err != nil {
		log.Fatal(err)
	}

	images, err := getPodsImages(clientset)
	if err != nil {
		log.Fatal(err)
	}

	result := []string{}

	for k := range images {
		result = append(result, k)
	}

	sort.Strings(result)

	imageignore, err := getImageIgnore()
	if err != nil {
		log.Debug(err)
	}

	for _, image := range result {
		isIgnored := imageignore.match(image)

		if !isIgnored {
			//nolint:forbidigo
			fmt.Println(image)
		}
	}
}
