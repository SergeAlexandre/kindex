package cmd

import (
	"context"
	"fmt"
	"kindex/internal/handlers"
	"kindex/internal/httpsrv"
	"kindex/internal/misc"
	"net/http"
	"os"
	"path"

	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var flags struct {
	logConfig  misc.LogConfig
	httpConfig httpsrv.Config
	kubeconfig string
}

func init() {
	ServeCmd.PersistentFlags().StringVarP(&flags.logConfig.Mode, "logMode", "", "text", "Log mode ('text' or 'json')")
	ServeCmd.PersistentFlags().StringVarP(&flags.logConfig.Level, "logLevel", "l", "INFO", "Log level(DEBUG, INFO, WARN, ERROR)")

	ServeCmd.PersistentFlags().BoolVarP(&flags.httpConfig.Tls, "tls", "t", false, "enable TLS")
	ServeCmd.PersistentFlags().IntVar(&flags.httpConfig.DumpExchanges, "dumpExchanges", 0, "Dump http server req/resp (0, 1, 2 or 3")
	ServeCmd.PersistentFlags().StringVarP(&flags.httpConfig.BindAddr, "bindAddr", "a", "0.0.0.0", "Bind Address")
	ServeCmd.PersistentFlags().IntVarP(&flags.httpConfig.BindPort, "bindPort", "p", 7788, "Bind port")
	ServeCmd.PersistentFlags().StringVar(&flags.httpConfig.CertDir, "certDir", "", "Certificate Directory")
	ServeCmd.PersistentFlags().StringVar(&flags.httpConfig.CertName, "certName", "tls.crt", "Certificate Directory")
	ServeCmd.PersistentFlags().StringVar(&flags.httpConfig.KeyName, "keyName", "tls.key", "Certificate Directory")
	ServeCmd.PersistentFlags().StringVarP(&flags.kubeconfig, "kubeconfig", "k", "", "Kubeconfig file (overrides $KUBECONFIG and ~/.kube/config)")
	//ServeCmd.PersistentFlags().StringArrayVarP(&flags.oidcHttpConfig.AllowedOrigins, "allowedOrigins", "", []string{}, "Allowed Origins")
}

var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Runs an HTTP server that lists ingress links from the cluster",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		logger, err := misc.NewLogger(&flags.logConfig)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Unable to load logging configuration: %v\n", err)
			os.Exit(2)
		}

		logger.Info("Starting http server", "port", flags.httpConfig.BindPort)

		// Kubeconfig resolution (first match wins):
		// 1) --kubeconfig path
		// 2) else KUBECONFIG environment variable (standard merge of listed files)
		// 3) else ~/.kube/config
		restCfg, err := restConfigFromKubeconfig(flags.kubeconfig)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "kubeconfig: %v\n", err)
			os.Exit(2)
		}
		clientSet, err := kubernetes.NewForConfig(restCfg)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "kubernetes client: %v\n", err)
			os.Exit(2)
		}

		router := http.NewServeMux()

		// Setup server
		server := httpsrv.New("ingresses", &flags.httpConfig, router)
		ctx := logr.NewContextWithSlogLogger(context.Background(), logger)

		router.Handle("GET /", handlers.IngressesHandler(clientSet))
		router.Handle("GET /favicon.ico", handlers.FaviconHandler(path.Join("resources/static", "favicon.ico")))

		err = server.Start(ctx)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error on server launch: %v\n", err)
			os.Exit(1)
		}

	},
}

func restConfigFromKubeconfig(explicitPath string) (*rest.Config, error) {
	if explicitPath != "" {
		return clientcmd.BuildConfigFromFlags("", explicitPath)
	}
	// Do not use BuildConfigFromFlags("", ""): after a failed in-cluster attempt it uses
	// ClientConfigLoadingRules with an empty ExplicitPath and no Precedence, so ~/.kube/config
	// and KUBECONFIG are never read. Use the standard loading rules instead.
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{}).ClientConfig()
}
