/*
Copyright 2017 Jetstack Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// This file was automatically generated by informer-gen

package internalversion

import (
	internalinterfaces "github.com/jetstack/navigator/pkg/client/informers/internalversion/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// CassandraClusters returns a CassandraClusterInformer.
	CassandraClusters() CassandraClusterInformer
	// ElasticsearchClusters returns a ElasticsearchClusterInformer.
	ElasticsearchClusters() ElasticsearchClusterInformer
	// Pilots returns a PilotInformer.
	Pilots() PilotInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// CassandraClusters returns a CassandraClusterInformer.
func (v *version) CassandraClusters() CassandraClusterInformer {
	return &cassandraClusterInformer{factory: v.factory, filter: v.filter}
}

// ElasticsearchClusters returns a ElasticsearchClusterInformer.
func (v *version) ElasticsearchClusters() ElasticsearchClusterInformer {
	return &elasticsearchClusterInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// Pilots returns a PilotInformer.
func (v *version) Pilots() PilotInformer {
	return &pilotInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}
