package konfig

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/golgoth31/multiShellKonfig/pkg/shell"
	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
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

	if localContext.Context.Namespace == "" {
		localContext.Context.Namespace = "default"
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

	// copy
	_, err = os.Stat(path.Dir(outputFileName) + "/last-known.yaml")
	if err == nil {
		// copy last-known context
		fin, err := os.Open(path.Dir(outputFileName) + "/last-known.yaml")
		if err != nil {
			return "", err
		}
		defer fin.Close()

		fout, err := os.Create(outputFileName)
		if err != nil {
			return "", err
		}
		defer fout.Close()

		_, err = io.Copy(fout, fin)

		if err != nil {
			return "", err
		}
	} else {
		// save context for namespace
		if err := os.WriteFile(outputFileName, outContext, 0640); err != nil {
			return "", err
		}

		if err := os.WriteFile(path.Dir(outputFileName)+"/last-known.yaml", outContext, 0640); err != nil {
			return "", err
		}
	}

	return outputFileName, nil
}

func GetNSList(currentKonfigFile string) ([]string, error) {
	// read current context
	currentKonfig, err := clientcmd.LoadFromFile(currentKonfigFile)
	if err != nil {
		log.Debug().Err(err).Msg("unable to get current context")

		return []string{}, err
	}

	config := clientcmd.NewDefaultClientConfig(*currentKonfig, &clientcmd.ConfigOverrides{})

	// generate the rest client for the current context
	restClient, errRestClient := config.ClientConfig()
	if errRestClient != nil {
		log.Debug().Err(errRestClient).Msg("unable to create kube client")

		return []string{}, errRestClient
	}

	// creates the clientset
	clientset, errKube := kubernetes.NewForConfig(restClient)
	if errKube != nil {
		log.Debug().Err(errKube).Msg("unable to create kube client")

		return []string{}, errKube
	}

	namespaceList, errNs := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if errNs != nil {
		log.Debug().Err(errNs).Msg("unable to get ns")

		return []string{}, errNs
	}

	out := []string{}
	for _, namespace := range namespaceList.Items {
		out = append(out, namespace.GetName())
	}

	return out, nil
}
