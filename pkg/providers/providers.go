package providers

import (
	"context"
	"fmt"

	"lab.nardis.io/kontinuum/api/v1alpha1"
)

// CloudProviderClient defines the interface that all cloud providers must implement
type CloudProviderClient interface {
	// VerifyClusterConnection checks if the cluster is accessible
	VerifyClusterConnection(ctx context.Context, spec v1alpha1.MemberClusterSpec) (bool, error)

	// GetKubeconfig retrieves the kubeconfig for the remote cluster
	GetKubeconfig(ctx context.Context, spec v1alpha1.MemberClusterSpec) (string, error)

	// GetClusterInfo retrieves information about the cluster
	GetClusterInfo(ctx context.Context, spec v1alpha1.MemberClusterSpec) (*ClusterInfo, error)
}

// ClusterInfo contains common cluster information
type ClusterInfo struct {
	Version   string
	NodeCount int
	Region    string
	NodePools []NodePoolInfo
}

type NodePoolInfo struct {
	Name        string
	NodeCount   int
	MachineType string
}

// GetCloudProviderClient returns the appropriate cloud provider client
func GetCloudProviderClient(provider v1alpha1.CloudProvider) (CloudProviderClient, error) {
	switch provider {
	// case v1alpha1.ProviderAKS:
	// 	return NewAKSClient()
	case v1alpha1.ProviderEKS:
		return NewEKSClient()
	// case v1alpha1.ProviderGKE:
	// 	return NewGKEClient()
	default:
		return nil, fmt.Errorf("unsupported cloud provider: %s", provider)
	}
}
