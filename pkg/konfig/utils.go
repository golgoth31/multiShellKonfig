package konfig

import (
	"context"
	"errors"
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

func GetNSList(currentKonfigFile string) (string, []string, error) {
	// read current context
	currentKonfig, err := clientcmd.LoadFromFile(currentKonfigFile)
	if err != nil {
		log.Debug().Err(err).Msg("unable to get current context")

		return "", []string{}, err
	}

	config := clientcmd.NewDefaultClientConfig(*currentKonfig, &clientcmd.ConfigOverrides{})

	// generate the rest client for the current context
	restClient, errRestClient := config.ClientConfig()
	if errRestClient != nil {
		log.Debug().Err(errRestClient).Msg("unable to create kube client")

		return "", []string{}, errRestClient
	}

	// creates the clientset
	clientset, errKube := kubernetes.NewForConfig(restClient)
	if errKube != nil {
		log.Debug().Err(errKube).Msg("unable to create kube client")

		return "", []string{}, errKube
	}

	namespaceList, errNs := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if errNs != nil {
		log.Debug().Err(errNs).Msg("unable to get ns")

		return "", []string{}, errNs
	}

	out := []string{}
	for index := range namespaceList.Items {
		out = append(out, namespaceList.Items[index].GetName())
	}

	return currentKonfig.Contexts[currentKonfig.CurrentContext].Namespace, out, nil
}
