package konfig

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	kubeClientConfig "k8s.io/client-go/tools/clientcmd/api/v1"
)

func (k *Konfig) readExternalData(fileDataPath string) ([]byte, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Debug().Err(err).Msg("unable to get home dir")

		return []byte{}, err
	}

	filePath := k.FilePath
	if strings.HasPrefix(filePath, "~/") {
		filePath = filepath.Join(homedir, filePath[2:])
	}

	filePath = filepath.Dir(filePath) + "/" + fileDataPath

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		log.Debug().Err(err).Msg("unable to read file")

		return []byte{}, err
	}

	return fileContent, nil
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

			if auth.AuthInfo.ClientCertificate != "" {
				var err error

				auth.AuthInfo.ClientCertificateData, err = k.readExternalData(auth.AuthInfo.ClientCertificate)
				if err != nil {
					return "", []byte{}, err
				}

				auth.AuthInfo.ClientCertificate = ""
			}

			if auth.AuthInfo.ClientKey != "" {
				var err error

				auth.AuthInfo.ClientKeyData, err = k.readExternalData(auth.AuthInfo.ClientKey)
				if err != nil {
					return "", []byte{}, err
				}

				auth.AuthInfo.ClientKey = ""
			}

			break
		}
	}

	cluster := kubeClientConfig.NamedCluster{}
	for index := range k.Content.Clusters {
		if k.Content.Clusters[index].Name == localContext.Context.Cluster {
			cluster = k.Content.Clusters[index]

			if cluster.Cluster.CertificateAuthority != "" {
				var err error

				cluster.Cluster.CertificateAuthorityData, err = k.readExternalData(cluster.Cluster.CertificateAuthority)
				if err != nil {
					return "", []byte{}, err
				}

				cluster.Cluster.CertificateAuthority = ""
			}

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
