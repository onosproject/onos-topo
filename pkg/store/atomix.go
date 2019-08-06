package store

import (
	"github.com/atomix/atomix-go-client/pkg/client"
	"os"
)

const (
	atomixControllerEnv = "ATOMIX_CONTROLLER"
	atomixNamespaceEnv  = "ATOMIX_NAMESPACE"
	atomixAppEnv        = "ATOMIX_APP"
	atomixRaftGroup     = "ATOMIX_RAFT"
)

func getAtomixController() string {
	return os.Getenv(atomixControllerEnv)
}

func getAtomixNamespace() string {
	return os.Getenv(atomixNamespaceEnv)
}

func getAtomixApp() string {
	return os.Getenv(atomixAppEnv)
}

func getAtomixRaftGroup() string {
	return os.Getenv(atomixRaftGroup)
}

// getAtomixClient returns the Atomix client
func getAtomixClient() (*client.Client, error) {
	opts := []client.ClientOption{
		client.WithNamespace(getAtomixNamespace()),
		client.WithApplication(getAtomixApp()),
	}
	return client.NewClient(getAtomixController(), opts...)
}
