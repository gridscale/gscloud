package cmd

import (
	"context"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientauth "k8s.io/client-go/pkg/apis/clientauthentication/v1beta1"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// k8sOperator is used for testing purpose,
// we can mock data return from the gsclient via interface.
type k8sOperator interface {
	RenewK8sCredentials(ctx context.Context, id string) error
	GetPaaSService(ctx context.Context, id string) (gsclient.PaaSService, error)
}

// k8s action enums
const (
	k8sSaveConfigAction = iota
	k8sExecCredAction
)

// produceK8SCmdRunFunc takes an instance of a struct that implements `k8sOperator`
// returns a `cmdRunFunc`
func produceK8SCmdRunFunc(o k8sOperator, action int) cmdRunFunc {
	switch action {
	case k8sSaveConfigAction:
		return func(cmd *cobra.Command, args []string) {
			kubeConfigFile, _ := cmd.Flags().GetString("kubeconfig")
			clusterID, _ := cmd.Flags().GetString("cluster")
			credentialPlugin, _ := cmd.Flags().GetBool("credential-plugin")
			kubeConfigEnv := os.Getenv("KUBECONFIG")

			pathOptions := clientcmd.NewDefaultPathOptions()
			if kubeConfigFile != "" {
				kubeConfigEnv = kubeConfigFile
				pathOptions.GlobalFile = kubeConfigFile
			}

			if kubeConfigEnv != "" && !fileExists(kubeConfigEnv) {
				_, err := os.Create(kubeConfigEnv)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
			}

			currentKubeConfig, err := pathOptions.GetStartingConfig()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			newKubeConfig := fetchKubeConfigFromProvider(o, clusterID)
			if len(newKubeConfig.Clusters) == 0 || len(newKubeConfig.Users) == 0 {
				fmt.Fprintln(os.Stderr, "Error: Invalid kubeconfig")
				os.Exit(1)
			}
			c := newKubeConfig.Clusters[0]
			u := newKubeConfig.Users[0]

			certificateAuthorityData, err := b64.StdEncoding.DecodeString(c.Cluster.CertificateAuthorityData)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
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
						Command:    cliPath(),
						Args: []string{
							"--config",
							cliConfigPath(),
							"--account",
							account,
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
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}

				clientKeyData, err := b64.StdEncoding.DecodeString(u.User.ClientKeyData)
				if err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
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
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

		}

	case k8sExecCredAction:
		return func(cmd *cobra.Command, args []string) {
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
				fmt.Fprintln(os.Stderr, err)
			}

			if execCredential == nil {
				newKubeConfig := fetchKubeConfigFromProvider(o, clusterID)
				if len(newKubeConfig.Users) != 0 {
					u := newKubeConfig.Users[0]
					clientKeyData, err := b64.StdEncoding.DecodeString(u.User.ClientKeyData)
					if err != nil {
						fmt.Println(err)
					}
					clientCertificateData, err := b64.StdEncoding.DecodeString(u.User.ClientCertificateData)
					if err != nil {
						fmt.Println(err)
					}

					execCredential = &clientauth.ExecCredential{
						TypeMeta: metav1.TypeMeta{
							Kind:       "ExecCredential",
							APIVersion: clientauth.SchemeGroupVersion.String(),
						},
						Status: &clientauth.ExecCredentialStatus{
							ClientKeyData:         string(clientKeyData),
							ClientCertificateData: string(clientCertificateData),
							ExpirationTimestamp:   &metav1.Time{Time: time.Now().Add(time.Hour)},
						},
					}

					if err := cacheKubeConfig(clusterID, execCredential); err != nil {
						fmt.Fprintln(os.Stderr, err)
					}
				}

			}
			if execCredential == nil {
				fmt.Println("Error: Could not retrieve kubeconfig from provider for account: ", account)
				return
			}
			execCredentialJSON, err := json.MarshalIndent(execCredential, "", "    ")
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			// this output will be used by kubectl
			fmt.Println(string(execCredentialJSON))
		}

	default:
	}
	return nil
}

func initK8SCmd() {
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
		Run:   produceK8SCmdRunFunc(client, k8sSaveConfigAction),
	}

	// execCredentialCmd represents the getCertificate command
	var execCredentialCmd = &cobra.Command{
		Use:   "exec-credential",
		Short: "Provides client credentials to kubectl command",
		Long:  "exec-credential provides client credentials to kubectl command.",
		Run:   produceK8SCmdRunFunc(client, k8sExecCredAction),
	}

	rootCmd.AddCommand(kubernetesCmd)
	kubernetesCmd.AddCommand(clusterCmd)
	clusterCmd.AddCommand(execCredentialCmd)
	clusterCmd.AddCommand(saveKubeconfigCmd)
	saveKubeconfigCmd.Flags().String("kubeconfig", "", "(optional) absolute path to the kubeconfig file")
	saveKubeconfigCmd.Flags().String("cluster", "", "The cluster's uuid")
	saveKubeconfigCmd.MarkFlagRequired("cluster")
	saveKubeconfigCmd.Flags().Bool("credential-plugin", false, "Enables credential plugin authentication method (exec-credential)")
	execCredentialCmd.Flags().String("kubeconfig", "", "(optional) absolute path to the kubeconfig file")
	execCredentialCmd.Flags().String("cluster", "", "The cluster's uuid")
	execCredentialCmd.MarkFlagRequired("cluster")
}

func fetchKubeConfigFromProvider(o k8sOperator, id string) *kubeConfig {
	var kc kubeConfig

	if err := o.RenewK8sCredentials(context.Background(), id); err != nil {
		os.Exit(1)
	}

	paaSService, err := o.GetPaaSService(context.Background(), id)
	if err != nil {
		os.Exit(1)
	}

	if len(paaSService.Properties.Credentials) != 0 {
		err := yaml.Unmarshal([]byte(paaSService.Properties.Credentials[0].KubeConfig), &kc)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	return &kc
}

func kubeConfigCachePath() string {
	return filepath.Join(cliCachePath(), "exec-credential")
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
