package providers

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"lab.nardis.io/kontinuum/api/v1alpha1"
	"lab.nardis.io/kontinuum/pkg/retry"

	ctrl "sigs.k8s.io/controller-runtime"
)

var log = ctrl.Log.WithName("eks")

type EKSClient struct {
	client *eks.Client
}

func NewEKSClient() (*EKSClient, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	client := eks.NewFromConfig(cfg)

	return &EKSClient{
		client: client,
	}, nil
}

func (e *EKSClient) VerifyClusterConnection(ctx context.Context, spec v1alpha1.HubClusterSpec) (bool, error) {
	var result *eks.DescribeClusterOutput

	op := func(ctx context.Context) error {
		var err error
		result, err = e.client.DescribeCluster(ctx, &eks.DescribeClusterInput{
			Name: &spec.ClusterName,
		})
		return err
	}

	retryOpts := retry.RetryOptions{
		MaxRetries:      3,
		InitialInterval: 1 * time.Second,
		MaxInterval:     10 * time.Second,
		BackoffFactor:   2.0,
	}

	if err := retry.WithRetry(ctx, op, retryOpts); err != nil {
		return false, fmt.Errorf("failed to describe EKS cluster after retries: %v", err)
	}

	return result.Cluster.Status == "ACTIVE", nil
}

func (e *EKSClient) GetKubeconfig(ctx context.Context, spec v1alpha1.HubClusterSpec) (string, error) {
	var result *eks.DescribeClusterOutput

	op := func(ctx context.Context) error {
		var err error
		result, err = e.client.DescribeCluster(ctx, &eks.DescribeClusterInput{
			Name: &spec.ClusterName,
		})
		return err
	}

	retryOpts := retry.RetryOptions{
		MaxRetries:      3,
		InitialInterval: 1 * time.Second,
		MaxInterval:     10 * time.Second,
		BackoffFactor:   2.0,
	}

	if err := retry.WithRetry(ctx, op, retryOpts); err != nil {
		return "", fmt.Errorf("failed to get EKS cluster info after retries: %v", err)
	}

	// Generate kubeconfig using cluster endpoint and certificate authority
	kubeconfig := generateEKSKubeconfig(
		spec.ClusterName,
		*result.Cluster.Endpoint,
		*result.Cluster.CertificateAuthority.Data,
	)

	return kubeconfig, nil
}

func (e *EKSClient) GetClusterInfo(ctx context.Context, spec v1alpha1.HubClusterSpec) (*ClusterInfo, error) {
	var cluster *eks.DescribeClusterOutput
	var nodeGroups *eks.ListNodegroupsOutput

	// Get cluster info with retry
	clusterOp := func(ctx context.Context) error {
		var err error
		cluster, err = e.client.DescribeCluster(ctx, &eks.DescribeClusterInput{
			Name: &spec.ClusterName,
		})
		return err
	}

	// Get node groups with retry
	nodeGroupsOp := func(ctx context.Context) error {
		var err error
		nodeGroups, err = e.client.ListNodegroups(ctx, &eks.ListNodegroupsInput{
			ClusterName: &spec.ClusterName,
		})
		return err
	}

	retryOpts := retry.RetryOptions{
		MaxRetries:      3,
		InitialInterval: 1 * time.Second,
		MaxInterval:     10 * time.Second,
		BackoffFactor:   2.0,
	}

	if err := retry.WithRetry(ctx, clusterOp, retryOpts); err != nil {
		return nil, fmt.Errorf("failed to get EKS cluster after retries: %v", err)
	}

	if err := retry.WithRetry(ctx, nodeGroupsOp, retryOpts); err != nil {
		return nil, fmt.Errorf("failed to list EKS node groups after retries: %v", err)
	}

	info := &ClusterInfo{
		Version:   *cluster.Cluster.Version,
		Region:    spec.Region,
		NodePools: []NodePoolInfo{},
	}

	// Get details for each node group with retry
	for _, ngName := range nodeGroups.Nodegroups {
		var ng *eks.DescribeNodegroupOutput

		nodeGroupOp := func(ctx context.Context) error {
			var err error
			ng, err = e.client.DescribeNodegroup(ctx, &eks.DescribeNodegroupInput{
				ClusterName:   &spec.ClusterName,
				NodegroupName: &ngName,
			})
			return err
		}

		if err := retry.WithRetry(ctx, nodeGroupOp, retryOpts); err != nil {

			log.Error(err, "failed to describe node group", "nodegroup", ngName)
			continue
		}

		nodePool := NodePoolInfo{
			Name:        ngName,
			NodeCount:   int(*ng.Nodegroup.ScalingConfig.DesiredSize),
			MachineType: ng.Nodegroup.InstanceTypes[0],
		}
		info.NodePools = append(info.NodePools, nodePool)
		info.NodeCount += nodePool.NodeCount
	}

	return info, nil
}

// generateEKSKubeconfig generates a kubeconfig for the EKS cluster.
func generateEKSKubeconfig(clusterName, endpoint, caData string) string {
	return fmt.Sprintf(`
apiVersion: v1
clusters:
- cluster:
	server: %s
	certificate-authority-data: %s
  name: %s
contexts:
- context:
	cluster: %s
	user: %s
  name: %s
current-context: %s
kind: Config
preferences: {}
users:
- name: %s
  user:
	exec:
	  apiVersion: client.authentication.k8s.io/v1alpha1
	  command: aws-iam-authenticator
	  args:
		- "token"
		- "-i"
		- "%s"
`, endpoint, caData, clusterName, clusterName, clusterName, clusterName, clusterName, clusterName, clusterName)
}
