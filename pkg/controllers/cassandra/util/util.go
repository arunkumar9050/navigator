package util

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/selection"

	"github.com/jetstack/navigator/pkg/apis/navigator"
	v1alpha1 "github.com/jetstack/navigator/pkg/apis/navigator/v1alpha1"
)

const (
	typeName             = "cass"
	kindName             = "CassandraCluster"
	ClusterNameLabelKey  = "navigator.jetstack.io/cassandra-cluster-name"
	NodePoolNameLabelKey = "navigator.jetstack.io/cassandra-node-pool-name"
	DefaultCqlPort       = 9042
)

func NewControllerRef(c *v1alpha1.CassandraCluster) metav1.OwnerReference {
	return *metav1.NewControllerRef(c, schema.GroupVersionKind{
		Group:   navigator.GroupName,
		Version: "v1alpha1",
		Kind:    kindName,
	})
}

func ResourceBaseName(c *v1alpha1.CassandraCluster) string {
	return typeName + "-" + c.Name
}

func NodePoolResourceName(c *v1alpha1.CassandraCluster, np *v1alpha1.CassandraClusterNodePool) string {
	return fmt.Sprintf("%s-%s", ResourceBaseName(c), np.Name)
}

func SeedProviderServiceName(c *v1alpha1.CassandraCluster) string {
	return fmt.Sprintf("%s-seedprovider", ResourceBaseName(c))
}

func CqlServiceName(c *v1alpha1.CassandraCluster) string {
	return fmt.Sprintf("%s-cql", ResourceBaseName(c))
}

func ServiceAccountName(c *v1alpha1.CassandraCluster) string {
	return ResourceBaseName(c)
}

func PilotRBACRoleName(c *v1alpha1.CassandraCluster) string {
	return fmt.Sprintf("%s-pilot", ResourceBaseName(c))
}

func ClusterLabels(c *v1alpha1.CassandraCluster) map[string]string {
	return map[string]string{
		"app":               "cassandracluster",
		ClusterNameLabelKey: c.Name,
	}
}

func SelectorForCluster(c *v1alpha1.CassandraCluster) (labels.Selector, error) {
	clusterNameReq, err := labels.NewRequirement(
		ClusterNameLabelKey,
		selection.Equals,
		[]string{c.Name},
	)
	if err != nil {
		return nil, err
	}
	return labels.NewSelector().Add(*clusterNameReq), nil
}

func NodePoolLabels(
	c *v1alpha1.CassandraCluster,
	poolName string,
) map[string]string {
	labels := ClusterLabels(c)
	if poolName != "" {
		labels[NodePoolNameLabelKey] = poolName
	}
	return labels
}

func Int32Ptr(i int32) *int32 {
	return &i
}

func OwnerCheck(
	obj metav1.Object,
	owner metav1.Object,
) error {
	if !metav1.IsControlledBy(obj, owner) {
		ownerRef := metav1.GetControllerOf(obj)
		return fmt.Errorf(
			"'%s/%s' is foreign owned: "+
				"it is owned by '%v', not '%s/%s'.",
			obj.GetNamespace(), obj.GetName(),
			ownerRef,
			owner.GetNamespace(), owner.GetName(),
		)
	}
	return nil
}
