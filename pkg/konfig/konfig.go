package konfig

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/golgoth31/multiShellKonfig/pkg/shell"
	"github.com/rs/zerolog/log"
	kubeClientConfig "k8s.io/client-go/tools/clientcmd/api/v1"
	"sigs.k8s.io/yaml"
)

func Load(path string, homeDir string) (*kubeClientConfig.Config, error) {
	kubeConfig := kubeClientConfig.Config{}

	if strings.HasPrefix(path, "~/") {
		path = filepath.Join(homeDir, path[2:])
	}

	fileContent, err := os.ReadFile(path)
	if err != nil {
		log.Debug().Err(err).Msg("unable to read file")

		return &kubeConfig, err
	}

	if err := yaml.Unmarshal(fileContent, &kubeConfig); err != nil {
		log.Debug().Err(err).Msg("unable to load konfig")

		return &kubeConfig, err
	}

	return &kubeConfig, nil
}

func Generate(context *shell.ContextDef, kubeConfig *kubeClientConfig.Config, contextsPath string) (string, error) {
	localContext := kubeClientConfig.NamedContext{}
	for _, c := range kubeConfig.Contexts {
		if c.Name == context.Name {
			localContext = c

			break
		}
	}

	auth := kubeClientConfig.NamedAuthInfo{}
	for _, authInfo := range kubeConfig.AuthInfos {
		if authInfo.Name == localContext.Context.AuthInfo {
			auth = authInfo

			break
		}
	}

	cluster := kubeClientConfig.NamedCluster{}
	for _, clusterInfo := range kubeConfig.Clusters {
		if clusterInfo.Name == localContext.Context.Cluster {
			cluster = clusterInfo

			break
		}
	}

	//build config file for context
	newFile := kubeClientConfig.Config{
		APIVersion:  kubeConfig.APIVersion,
		Kind:        kubeConfig.Kind,
		Preferences: kubeConfig.Preferences,
		Clusters: []kubeClientConfig.NamedCluster{
			cluster,
		},
		AuthInfos: []kubeClientConfig.NamedAuthInfo{
			auth,
		},
	}

	// compute the sha of cluster/auth tupple
	contextShaData, err := json.Marshal(newFile)
	if err != nil {
		log.Debug().Err(err).Msg("cannot marshall config")

		return "", err
	}
	contextSha := sha256.Sum256(contextShaData)

	// generate new current context
	newFile.Contexts = []kubeClientConfig.NamedContext{
		localContext,
	}
	newFile.CurrentContext = context.Name

	outContext, err := json.Marshal(newFile)
	if err != nil {
		log.Debug().Err(err).Msg("cannot marshall config")

		return "", err
	}

	outputFileName := contextsPath + "/" + context.FileID + "/" + fmt.Sprintf("%x", contextSha) + "/" + localContext.Context.Namespace + ".yaml"
	// create directory for cluster:auth tupple
	if _, err := os.Stat(path.Dir(outputFileName)); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			if errMkdir := os.MkdirAll(path.Dir(outputFileName), 0755); errMkdir != nil {
				return "", errMkdir
			}
		}
	}

	// save context for namespace
	if err := os.WriteFile(outputFileName, outContext, 0640); err != nil {
		return "", err
	}

	return outputFileName, nil
}
