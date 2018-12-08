package main

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/nuclio/nuclio/cmd/dlx/app"
	"github.com/nuclio/nuclio/pkg/errors"
)

func getNamespace(namespaceArgument string) string {

	// if the namespace was passed in the arguments, use that
	if namespaceArgument != "" {
		return namespaceArgument
	}

	// if the namespace exists in env, use that
	if namespaceEnv := os.Getenv("NUCLIO_PROXIER_NAMESPACE"); namespaceEnv != "" {
		return namespaceEnv
	}

	// if nothing was passed, assume "this" namespace
	return "@nuclio.selfNamespace"
}

func main() {
	listenAddress := flag.String("listen-addr", ":8080", "IP/port on which dlx listens on")
	kubeconfigPath := flag.String("kubeconfig-path", "", "Path of kubeconfig file")
	namespace := flag.String("namespace", "", "Namespace to listen on, or * for all")
	flag.Parse()

	// get the namespace from args -> env -> default (*)
	resolvedNamespace := getNamespace(*namespace)

	// if the namespace is set to @nuclio.selfNamespace, use the namespace we're in right now
	if resolvedNamespace == "@nuclio.selfNamespace" {

		// get namespace from within the pod. if found, return that
		if namespacePod, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace"); err == nil {
			resolvedNamespace = string(namespacePod)
		}
	}

	if err := app.Run(*kubeconfigPath, *listenAddress, resolvedNamespace); err != nil {
		errors.PrintErrorStack(os.Stderr, err, 5)

		os.Exit(1)
	}
}