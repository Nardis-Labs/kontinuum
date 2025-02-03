package providers

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"

// 	container "cloud.google.com/go/container/apiv1"
// 	containerpb "cloud.google.com/go/container/apiv1/containerpb"
// 	"google.golang.org/api/option"
// 	"lab.nardis.io/kontinuum/api/v1alpha1"
// )

// type GKEClient struct {
// 	client *container.ClusterManagerClient
// }

// func NewGKEClient() (*GKEClient, error) {
// 	ctx := context.Background()
// 	client, err := container.NewClusterManagerClient(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create GKE client: %v", err)
// 	}

// 	return &GKEClient{
// 		client: client,
// 	}, nil
// }

// // parseGKECredentials parses the credentials from the RemoteCluster spec
// func parseGKECredentials(spec v1alpha1.MemberClusterSpec) (projectID string, location string, err error) {
// 	var creds struct {
// 		ProjectID string `json:"projectId"`
// 		Location  string `json:"location"` // Can be a zone or region
// 	}

// 	if err := json.Unmarshal([]byte(spec.Credentials.Data), &creds); err != nil {
// 		return "", "", fmt.Errorf("failed to parse GKE credentials: %v", err)
// 	}

// 	if creds.ProjectID == "" || creds.Location == "" {
// 		return "", "", fmt.Errorf("projectId and location are required in GKE credentials")
// 	}

// 	return creds.ProjectID, creds.Location, nil
// }

// func (g *GKEClient) VerifyClusterConnection(ctx context.Context, spec v1alpha1.MemberClusterSpec) (bool, error) {
// 	projectID, location, err := parseGKECredentials(spec)
// 	if err != nil {
// 		return false, err
// 	}

// 	req := &containerpb.GetClusterRequest{
// 		Name: fmt.Sprintf("projects/%s/locations/%s/clusters/%s",
// 			projectID, location, spec.ClusterName),
// 	}

// 	cluster, err := g.client.GetCluster(ctx, req)
// 	if err != nil {
// 		return false, fmt.Errorf("failed to get cluster: %v", err)
// 	}

// 	// Check if cluster is running
// 	return cluster.Status == containerpb.Cluster_RUNNING, nil
// }

// func (g *GKEClient) GetKubeconfig(ctx context.Context, spec v1alpha1.MemberClusterSpec) (string, error) {
// 	projectID, location, err := parseGKECredentials(spec)
// 	if err != nil {
// 		return "", err
// 	}

// 	req := &containerpb.GetClusterRequest{
// 		Name: fmt.Sprintf("projects/%s/locations/%s/clusters/%s",
// 			projectID, location, spec.ClusterName),
// 	}

// 	cluster, err := g.client.GetCluster(ctx, req)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to get cluster: %v", err)
// 	}

// 	// Generate kubeconfig
// 	kubeconfig := map[string]interface{}{
// 		"apiVersion":      "v1",
// 		"kind":            "Config",
// 		"current-context": spec.ClusterName,
// 		"clusters": []map[string]interface{}{
// 			{
// 				"name": spec.ClusterName,
// 				"cluster": map[string]interface{}{
// 					"server":                     "https://" + cluster.Endpoint,
// 					"certificate-authority-data": cluster.MasterAuth.ClusterCaCertificate,
// 				},
// 			},
// 		},
// 		"contexts": []map[string]interface{}{
// 			{
// 				"name": spec.ClusterName,
// 				"context": map[string]interface{}{
// 					"cluster": spec.ClusterName,
// 					"user":    "gke-user",
// 				},
// 			},
// 		},
// 		"users": []map[string]interface{}{
// 			{
// 				"name": "gke-user",
// 				"user": map[string]interface{}{
// 					"auth-provider": map[string]interface{}{
// 						"name": "gcp",
// 					},
// 				},
// 			},
// 		},
// 	}

// 	kubeconfigBytes, err := json.Marshal(kubeconfig)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to marshal kubeconfig: %v", err)
// 	}

// 	return string(kubeconfigBytes), nil
// }

// func (g *GKEClient) GetClusterInfo(ctx context.Context, spec v1alpha1.MemberClusterSpec) (*ClusterInfo, error) {
// 	projectID, location, err := parseGKECredentials(spec)
// 	if err != nil {
// 		return nil, err
// 	}

// 	req := &containerpb.GetClusterRequest{
// 		Name: fmt.Sprintf("projects/%s/locations/%s/clusters/%s",
// 			projectID, location, spec.ClusterName),
// 	}

// 	cluster, err := g.client.GetCluster(ctx, req)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get cluster: %v", err)
// 	}

// 	info := &ClusterInfo{
// 		Version:   cluster.CurrentMasterVersion,
// 		Region:    location,
// 		NodeCount: 0,
// 		NodePools: []NodePoolInfo{},
// 	}

// 	// Get node pools
// 	npReq := &containerpb.ListNodePoolsRequest{
// 		Parent: fmt.Sprintf("projects/%s/locations/%s/clusters/%s",
// 			projectID, location, spec.ClusterName),
// 	}

// 	nodePools, err := g.client.ListNodePools(ctx, npReq)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to list node pools: %v", err)
// 	}

// 	for _, np := range nodePools.NodePools {
// 		nodePool := NodePoolInfo{
// 			Name:        np.Name,
// 			NodeCount:   int(np.InitialNodeCount),
// 			MachineType: np.Config.MachineType,
// 		}
// 		info.NodePools = append(info.NodePools, nodePool)
// 		info.NodeCount += nodePool.NodeCount
// 	}

// 	return info, nil
// }

// // Helper function to create GKE client with specific credentials
// func NewGKEClientWithCredentials(ctx context.Context, credentialsJSON []byte) (*GKEClient, error) {
// 	client, err := container.NewClusterManagerClient(
// 		ctx,
// 		option.WithCredentialsJSON(credentialsJSON),
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create GKE client with credentials: %v", err)
// 	}

// 	return &GKEClient{
// 		client: client,
// 	}, nil
// }

// // Example credential structure for GKE
// type GKECredentials struct {
// 	ProjectID      string `json:"projectId"`
// 	Location       string `json:"location"`
// 	ServiceAccount string `json:"serviceAccount"` // Base64 encoded service account key
// }
