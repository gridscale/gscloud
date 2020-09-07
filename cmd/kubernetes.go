package cmd

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gridscale/gscloud/runtime"
	"github.com/gridscale/gscloud/utils"
	"github.com/kardianos/osext"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientauth "k8s.io/client-go/pkg/apis/clientauthentication/v1beta1"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func executablePath() string {
	filePath, err := osext.Executable()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return filePath
}

// clusterCmd represents the cluster command
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Actions on a Kubernetes cluster",
	Long:  "Actions on a Kubernetes cluster",
}

// kubernetesCmd represents the Kubernetes command
var kubernetesCmd = &cobra.Command{
	Use:   "kubernetes",
	Short: "Operate managed Kubernetes clusters",
	Long:  "Operate managed Kubernetes clusters",
}

// saveKubeconfigCmd represents the kubeconfig command
var saveKubeconfigCmd = &cobra.Command{
	Use:   "save-kubeconfig",
	Short: "Saves configuration of the given cluster into a kubeconfig",
	Long:  "Saves configuration of the given cluster into a kubeconfig or KUBECONFIG environment variable.",
	Run: func(cmd *cobra.Command, args []string) {
		kubeConfigFile, _ := cmd.Flags().GetString("kubeconfig")
		clusterID, _ := cmd.Flags().GetString("cluster")
		credentialPlugin, _ := cmd.Flags().GetBool("credential-plugin")
		kubeConfigEnv := os.Getenv("KUBECONFIG")

		pathOptions := clientcmd.NewDefaultPathOptions()
		if kubeConfigFile != "" {
			kubeConfigEnv = kubeConfigFile
			pathOptions.GlobalFile = kubeConfigFile
		}

		if kubeConfigEnv != "" && !utils.FileExists(kubeConfigEnv) {
			_, err := os.Create(kubeConfigEnv)
			if err != nil {
				log.Fatalln(err)
			}
		}

		currentKubeConfig, err := pathOptions.GetStartingConfig()
		if err != nil {
			log.Fatalf("Couldn't get starting config: %s", err)
		}

		op := rt.KubernetesOperator()
		newKubeConfig, _, err := fetchKubeConfigFromProvider(op, clusterID)
		if err != nil {
			log.Fatalf("Invalid kubeconfig: %s", err)
		}
		c := newKubeConfig.Clusters[0]
		u := newKubeConfig.Users[0]

		certificateAuthorityData, err := b64.StdEncoding.DecodeString(c.Cluster.CertificateAuthorityData)
		if err != nil {
			log.Fatalln(err)
		}

		currentKubeConfig.Clusters[c.Name] = &clientcmdapi.Cluster{
			Server:                   c.Cluster.Server,
			CertificateAuthorityData: certificateAuthorityData,
		}
		currentKubeConfig.AuthInfos[u.Name] = &clientcmdapi.AuthInfo{
			ClientCertificate: u.User.ClientKeyData,
			ClientKey:         u.User.ClientCertificateData,
		}
		if credentialPlugin {
			currentKubeConfig.AuthInfos[u.Name] = &clientcmdapi.AuthInfo{
				Exec: &clientcmdapi.ExecConfig{
					APIVersion: clientauth.SchemeGroupVersion.String(),
					Command:    executablePath(),
					Args: []string{
						"--config",
						runtime.ConfigPath(),
						"--account",
						rt.Account(),
						"kubernetes",
						"cluster",
						"exec-credential",
						"--cluster",
						clusterID,
					},
					Env: []clientcmdapi.ExecEnvVar{},
				},
			}
		} else {
			clientCertificateData, err := b64.StdEncoding.DecodeString(u.User.ClientCertificateData)
			if err != nil {
				log.Fatalln(err)
			}

			clientKeyData, err := b64.StdEncoding.DecodeString(u.User.ClientKeyData)
			if err != nil {
				log.Fatalln(err)
			}

			currentKubeConfig.AuthInfos[u.Name] = &clientcmdapi.AuthInfo{
				ClientCertificateData: clientCertificateData,
				ClientKeyData:         clientKeyData,
			}
		}

		currentKubeConfig.Contexts[newKubeConfig.CurrentContext] = &clientcmdapi.Context{
			Cluster:  c.Name,
			AuthInfo: u.Name,
		}
		currentKubeConfig.CurrentContext = newKubeConfig.CurrentContext

		err = clientcmd.ModifyConfig(pathOptions, *currentKubeConfig, true)
		if err != nil {
			log.Fatalln(err)
		}
	},
}

// execCredentialCmd represents the getCertificate command
var execCredentialCmd = &cobra.Command{
	Use:   "exec-credential",
	Short: "Provides client credentials to kubectl command",
	Long:  "exec-credential provides client credentials to kubectl command.",
	Run: func(cmd *cobra.Command, args []string) {
		kubeConfigFile, _ := cmd.Flags().GetString("kubeconfig")
		clusterID, _ := cmd.Flags().GetString("cluster")

		kubectlDefaults := clientcmd.NewDefaultPathOptions()
		if kubeConfigFile != "" {
			kubectlDefaults.GlobalFile = kubeConfigFile
		}

		_, err := kubectlDefaults.GetStartingConfig()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		execCredential, err := loadCachedKubeConfig(clusterID)
		if err != nil {
			log.Fatalf("Couldn't load cached kubeconfig: %s", err)
		}

		op := rt.KubernetesOperator()

		if execCredential == nil {
			newKubeConfig, expirationTime, err := fetchKubeConfigFromProvider(op, clusterID)
			if err != nil {
				log.Fatalf("Couldn't fetch kubeconfig: %s", err)
			}

			u := newKubeConfig.Users[0]
			clientKeyData, err := b64.StdEncoding.DecodeString(u.User.ClientKeyData)
			if err != nil {
				fmt.Println(err)
			}
			clientCertificateData, err := b64.StdEncoding.DecodeString(u.User.ClientCertificateData)
			if err != nil {
				fmt.Println(err)
			}

			if expirationTime.IsZero() {
				expirationTime = time.Now().Add(time.Hour)
			}

			execCredential = &clientauth.ExecCredential{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ExecCredential",
					APIVersion: clientauth.SchemeGroupVersion.String(),
				},
				Status: &clientauth.ExecCredentialStatus{
					ClientKeyData:         string(clientKeyData),
					ClientCertificateData: string(clientCertificateData),
					ExpirationTimestamp:   &metav1.Time{Time: expirationTime},
				},
			}

			if err := cacheKubeConfig(clusterID, execCredential); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
		execCredentialJSON, err := json.MarshalIndent(execCredential, "", "    ")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		// this output will be used by kubectl
		fmt.Println(string(execCredentialJSON))
	},
}

func init() {
	saveKubeconfigCmd.Flags().String("kubeconfig", "", "(optional) absolute path to the kubeconfig file")
	saveKubeconfigCmd.Flags().String("cluster", "", "The cluster's uuid")
	saveKubeconfigCmd.MarkFlagRequired("cluster")
	saveKubeconfigCmd.Flags().Bool("credential-plugin", false, "Enables credential plugin authentication method (exec-credential)")
	clusterCmd.AddCommand(saveKubeconfigCmd)

	execCredentialCmd.Flags().String("kubeconfig", "", "(optional) absolute path to the kubeconfig file")
	execCredentialCmd.Flags().String("cluster", "", "The cluster's uuid")
	execCredentialCmd.MarkFlagRequired("cluster")
	clusterCmd.AddCommand(execCredentialCmd)

	kubernetesCmd.AddCommand(clusterCmd)
	rootCmd.AddCommand(kubernetesCmd)
}

func fetchKubeConfigFromProvider(op runtime.KubernetesOperator, id string) (kubeConfig, time.Time, error) {
	var kc kubeConfig
	var expirationTime time.Time

	if err := op.RenewK8sCredentials(context.Background(), id); err != nil {
		return kubeConfig{}, time.Time{}, err
	}

	platformService, err := op.GetPaaSService(context.Background(), id)
	if err != nil {
		return kubeConfig{}, time.Time{}, err
	}

	if len(platformService.Properties.Credentials) != 0 {
		err := yaml.Unmarshal([]byte(platformService.Properties.Credentials[0].KubeConfig), &kc)
		if err != nil {
			return kubeConfig{}, time.Time{}, err
		}
		expirationTime = platformService.Properties.Credentials[0].ExpirationTime.Time
	}

	return kc, expirationTime, nil
}

func kubeConfigCachePath() string {
	return filepath.Join(runtime.CachePath(), "exec-credential")
}

func cachedKubeConfigPath(id string) string {
	return filepath.Join(kubeConfigCachePath(), id+".json")
}

func cacheKubeConfig(id string, execCredential *clientauth.ExecCredential) error {
	if execCredential.Status.ExpirationTimestamp.IsZero() {
		return nil
	}

	cachePath := kubeConfigCachePath()
	if err := os.MkdirAll(cachePath, os.FileMode(0700)); err != nil {
		return err
	}

	path := cachedKubeConfigPath(id)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(0600))
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(execCredential)
}

func loadCachedKubeConfig(id string) (*clientauth.ExecCredential, error) {
	kubeConfigPath := cachedKubeConfigPath(id)
	f, err := os.Open(kubeConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}

	defer f.Close()

	var execCredential *clientauth.ExecCredential
	if err := json.NewDecoder(f).Decode(&execCredential); err != nil {
		return nil, err
	}

	timeStamp := execCredential.Status.ExpirationTimestamp

	if execCredential.Status == nil || timeStamp.IsZero() || timeStamp.Time.Before(time.Now()) {
		err = os.Remove(kubeConfigPath)
		return nil, err
	}

	return execCredential, nil
}
