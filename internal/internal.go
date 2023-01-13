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
package internal

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/maksim-paskal/k8s-images-cli/pkg/api"
	"github.com/maksim-paskal/k8s-images-cli/pkg/config"
	"github.com/maksim-paskal/k8s-images-cli/pkg/imageignore"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func Run(ctx context.Context) error {
	imageignore, err := imageignore.New(*config.Get().ImageIgnoreFile)
	if err != nil {
		log.Debug(err)
	}

	kubeconfigs := strings.Split(*config.Get().KubeConfigFile, ",")
	images := make(map[string]*api.Image)

	for _, kubeconfig := range kubeconfigs {
		kubeconfigImages, err := api.GetPodsImages(ctx, kubeconfig, imageignore)
		if err != nil {
			return errors.Wrap(err, "errors getting pods images")
		}

		for k, v := range kubeconfigImages {
			if _, ok := images[k]; !ok {
				images[k] = v
			}
		}
	}

	if len(*config.Get().Image) > 0 || len(*config.Get().ImagePullPolicy) > 0 || len(*config.Get().ImageArch) > 0 {
		log.Debug("close app when Image or ImagePullPolicy argument found")

		return nil
	}

	printResults(imageignore, images)

	return nil
}

func printResults(imageignore *imageignore.Type, images map[string]*api.Image) {
	result := []string{}

	for k, v := range images {
		image := k

		if *config.Get().ShowArch {
			image = fmt.Sprintf("%s (%s)", k, v.ImageArch)
		}

		result = append(result, image)
	}

	sort.Strings(result)

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
