package konfig

import kubeClientConfig "k8s.io/client-go/tools/clientcmd/api/v1"

type Konfig struct {
	FileID   string
	FilePath string
	Content  *kubeClientConfig.Config
}
