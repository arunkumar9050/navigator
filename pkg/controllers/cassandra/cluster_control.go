package cassandra

import (
	"github.com/golang/glog"
	v1alpha1 "github.com/jetstack-experimental/navigator/pkg/apis/navigator/v1alpha1"
)

const (
	errorSync = "ErrSync"

	successSync = "SuccessSync"

	messageErrorSyncServiceAccount = "Error syncing service account: %s"
	messageErrorSyncConfigMap      = "Error syncing config map: %s"
	messageErrorSyncService        = "Error syncing service: %s"
	messageErrorSyncNodePools      = "Error syncing node pools: %s"
	messageSuccessSync             = "Successfully synced CassandraCluster"
)

type ControlInterface interface {
	Sync(*v1alpha1.CassandraCluster) (v1alpha1.CassandraClusterStatus, error)
}

var _ ControlInterface = &defaultCassandraClusterControl{}

type defaultCassandraClusterControl struct{}

func NewController() ControlInterface {
	return &defaultCassandraClusterControl{}
}
func (e *defaultCassandraClusterControl) Sync(
	c *v1alpha1.CassandraCluster,
) (v1alpha1.CassandraClusterStatus, error) {
	glog.V(4).Infof("defaultCassandraClusterControl.Sync")
	return c.Status, nil
}
