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
	"strings"

	"github.com/maksim-paskal/k8s-images-cli/pkg/api"
	"github.com/maksim-paskal/k8s-images-cli/pkg/config"
	"github.com/maksim-paskal/k8s-images-cli/pkg/imageignore"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

var version = flag.Bool("version", false, "show version")

func main() { //nolint:cyclop
	flag.Parse()

	logLevel, err := log.ParseLevel(*config.Get().LogLevel)
	if err != nil {
		log.WithError(err).Fatal("error parse level")
	}

	log.SetLevel(logLevel)

	if *version {
		os.Stdout.WriteString(config.GetVersion())
		os.Exit(0)
	}

	imageignore, err := imageignore.New(*config.Get().ImageIgnoreFile)
	if err != nil {
		log.Debug(err)
	}

	kubeconfigs := strings.Split(*config.Get().KubeConfigFile, ",")
	images := make(map[string]v1.Pod)

	for _, kubeconfig := range kubeconfigs {
		kubeconfigImages, err := api.GetPodsImages(kubeconfig, imageignore)
		if err != nil {
			log.Fatal(err)
		}

		for k, v := range kubeconfigImages {
			if _, ok := images[k]; !ok {
				images[k] = v
			}
		}
	}

	if len(*config.Get().Image) > 0 || len(*config.Get().ImagePullPolicy) > 0 {
		log.Debug("close app when Image or ImagePullPolicy argument found")

		return
	}

	result := []string{}

	for k := range images {
		result = append(result, k)
	}

	sort.Strings(result)

	printResults(imageignore, result)
}

func printResults(imageignore *imageignore.Type, result []string) {
	for _, image := range result {
		isIgnored := imageignore.Match(image)

		if !isIgnored {
			//nolint:forbidigo
			fmt.Println(image)
		} else {
			log.Debugf("ignored by imageignore %s", image)
		}
	}
}
