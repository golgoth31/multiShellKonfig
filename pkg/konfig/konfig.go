package konfig

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	kubeClientConfig "k8s.io/client-go/tools/clientcmd/api/v1"
	"sigs.k8s.io/yaml"
)

const (
	filePerm = 0600
	dirPerm  = 0700
)

func Load(filePath, homeDir string) (*kubeClientConfig.Config, error) {
	kubeConfig := kubeClientConfig.Config{}

	if strings.HasPrefix(filePath, "~/") {
		filePath = filepath.Join(homeDir, filePath[2:])
	}

	fileContent, err := os.ReadFile(filePath)
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

func (k *Konfig) Generate(contextName, contextsPath string) (string, []byte, error) {
	localContext := kubeClientConfig.NamedContext{}

	for _, c := range k.Content.Contexts {
		if c.Name == contextName {
			localContext = c

			break
		}
	}

	localContext.Context.Namespace = ""

	auth := kubeClientConfig.NamedAuthInfo{}
	for index := range k.Content.AuthInfos {
		if k.Content.AuthInfos[index].Name == localContext.Context.AuthInfo {
			auth = k.Content.AuthInfos[index]

			break
		}
	}

	cluster := kubeClientConfig.NamedCluster{}
	for index := range k.Content.Clusters {
		if k.Content.Clusters[index].Name == localContext.Context.Cluster {
			cluster = k.Content.Clusters[index]

			break
		}
	}

	// build config file for context
	newFile := kubeClientConfig.Config{
		APIVersion:  k.Content.APIVersion,
		Kind:        k.Content.Kind,
		Preferences: k.Content.Preferences,
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

		return "", []byte{}, err
	}
	contextSha := sha256.Sum256(contextShaData)

	if localContext.Context.Namespace == "" {
		localContext.Context.Namespace = "default"
	}

	lastNS, err := os.ReadFile(contextsPath + "/" + k.FileID + "/" + fmt.Sprintf("%x", contextSha) + "/last-namespace")
	if err != nil {
		log.Debug().Err(err).Msg("cannot read last-namespace file")
	}

	log.Debug().Msgf("last namespace: %s", lastNS)

	if lastNS != nil {
		localContext.Context.Namespace = string(lastNS)
	}

	// generate new current context
	newFile.Contexts = []kubeClientConfig.NamedContext{
		localContext,
	}
	newFile.CurrentContext = contextName

	outFileData, err := json.Marshal(newFile)
	if err != nil {
		log.Debug().Err(err).Msg("cannot marshall config")

		return "", []byte{}, err
	}

	outFileName := contextsPath + "/" + k.FileID + "/" + fmt.Sprintf("%x", contextSha) + "/" + localContext.Context.Namespace + ".yaml"

	return outFileName, outFileData, nil
}

func SaveContextFile(fileName string, fileData []byte) error {
	// create directory for cluster:auth tupple
	if _, err := os.Stat(path.Dir(fileName)); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			if errMkdir := os.MkdirAll(path.Dir(fileName), dirPerm); errMkdir != nil {
				return errMkdir
			}
		}
	}

	lastNS := strings.TrimSuffix(path.Base(fileName), filepath.Ext(fileName))

	if err := os.WriteFile(fileName, fileData, filePerm); err != nil {
		return err
	}

	if err := os.WriteFile(path.Dir(fileName)+"/last-namespace", []byte(lastNS), filePerm); err != nil {
		return err
	}

	return nil
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
	for index := range namespaceList.Items {
		out = append(out, namespaceList.Items[index].GetName())
	}

	return out, nil
}
