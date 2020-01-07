package cmd

type contextInfo struct {
	User string `yaml:"user"`
}
type clusterInfo struct {
	CertificateAuthorityData string `yaml:"certificate-authority-data"`
	Server                   string `yaml:"server"`
}

type contextEntry struct {
	Name    string      `yaml:"name"`
	Context contextInfo `yaml:"context"`
}

type clusterEntry struct {
	Name    string      `yaml:"name"`
	Cluster clusterInfo `yaml:"cluster"`
}

type userEntry struct {
	Name string   `yaml:"name"`
	User userInfo `yaml:"user"`
}

type authProviderConfig struct {
	AccessToken string `yaml:"access-token"`
}

type authProviderInfo struct {
	Config authProviderConfig `yaml:"config"`
}

type userInfo struct {
	ClientCertificateData string `yaml:"client-certificate-data"`
	ClientKeyData         string `yaml:"client-key-data"`
}

type kubeConfig struct {
	Kind           string         `yaml:"kind,omitempty"`
	Contexts       []contextEntry `yaml:"contexts"`
	CurrentContext string         `yaml:"current-context"`
	Users          []userEntry    `yaml:"users"`
	Clusters       []clusterEntry `yaml:"clusters"`
}
