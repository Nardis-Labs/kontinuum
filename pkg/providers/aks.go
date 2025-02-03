package providers

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"

// 	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
// 	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice"

// 	"lab.nardis.io/kontinuum/api/v1alpha1"
// )

// type AKSClient struct {
// 	client *armcontainerservice.ManagedClustersClient
// }

// func NewAKSClient() (*AKSClient, error) {
// 	cred, err := azidentity.NewDefaultAzureCredential(nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	client, err := armcontainerservice.NewManagedClustersClient(cred, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &AKSClient{
// 		client: client,
// 	}, nil
// }

// func (a *AKSClient) VerifyClusterConnection(ctx context.Context, spec v1alpha1.MemberClusterSpec) (bool, error) {
// 	// Parse credentials from spec
// 	var creds struct {
// 		SubscriptionID string `json:"subscriptionId"`
// 		ResourceGroup  string `json:"resourceGroup"`
// 	}

// 	if err := json.Unmarshal([]byte(spec.Credentials.Data), &creds); err != nil {
// 		return false, err
// 	}

// 	result, err := a.client.Get(ctx, creds.ResourceGroup, spec.ClusterName, nil)
// 	if err != nil {
// 		return false, err
// 	}

// 	return result.Properties.ProvisioningState != nil && *result.Properties.ProvisioningState == "Succeeded", nil
// }

// func (a *AKSClient) GetKubeconfig(ctx context.Context, spec v1alpha1.MemberClusterSpec) (string, error) {
// 	var creds struct {
// 		SubscriptionID string `json:"subscriptionId"`
// 		ResourceGroup  string `json:"resourceGroup"`
// 	}

// 	if err := json.Unmarshal([]byte(spec.Credentials.Data), &creds); err != nil {
// 		return "", err
// 	}

// 	credential, err := a.client.ListClusterAdminCredentials(ctx, creds.ResourceGroup, spec.ClusterName, nil)
// 	if err != nil {
// 		return "", err
// 	}

// 	if credential.Kubeconfigs == nil || len(credential.Kubeconfigs) == 0 {
// 		return "", fmt.Errorf("no kubeconfig found for cluster %s", spec.ClusterName)
// 	}

// 	return string(credential.Kubeconfigs[0].Value), nil
// }

// func (a *AKSClient) GetClusterInfo(ctx context.Context, spec v1alpha1.MemberClusterSpec) (*ClusterInfo, error) {
// 	var creds struct {
// 		SubscriptionID string `json:"subscriptionId"`
// 		ResourceGroup  string `json:"resourceGroup"`
// 	}

// 	if err := json.Unmarshal([]byte(spec.Credentials.Data), &creds); err != nil {
// 		return nil, err
// 	}

// 	cluster, err := a.client.Get(ctx, creds.ResourceGroup, spec.ClusterName, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	info := &ClusterInfo{
// 		Version:   *cluster.Properties.KubernetesVersion,
// 		Region:    *cluster.Location,
// 		NodeCount: 0,
// 		NodePools: []NodePoolInfo{},
// 	}

// 	// Get node pools
// 	nodePools, err := a.client.ListAgentPools(ctx, creds.ResourceGroup, spec.ClusterName, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, np := range nodePools.Value {
// 		nodePool := NodePoolInfo{
// 			Name:        *np.Name,
// 			NodeCount:   int(*np.Properties.Count),
// 			MachineType: *np.Properties.VMSize,
// 		}
// 		info.NodePools = append(info.NodePools, nodePool)
// 		info.NodeCount += nodePool.NodeCount
// 	}

// 	return info, nil
// }
