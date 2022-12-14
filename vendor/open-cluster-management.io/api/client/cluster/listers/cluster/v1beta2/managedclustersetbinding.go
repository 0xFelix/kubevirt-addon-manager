// Code generated by lister-gen. DO NOT EDIT.

package v1beta2

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
	v1beta2 "open-cluster-management.io/api/cluster/v1beta2"
)

// ManagedClusterSetBindingLister helps list ManagedClusterSetBindings.
// All objects returned here must be treated as read-only.
type ManagedClusterSetBindingLister interface {
	// List lists all ManagedClusterSetBindings in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta2.ManagedClusterSetBinding, err error)
	// ManagedClusterSetBindings returns an object that can list and get ManagedClusterSetBindings.
	ManagedClusterSetBindings(namespace string) ManagedClusterSetBindingNamespaceLister
	ManagedClusterSetBindingListerExpansion
}

// managedClusterSetBindingLister implements the ManagedClusterSetBindingLister interface.
type managedClusterSetBindingLister struct {
	indexer cache.Indexer
}

// NewManagedClusterSetBindingLister returns a new ManagedClusterSetBindingLister.
func NewManagedClusterSetBindingLister(indexer cache.Indexer) ManagedClusterSetBindingLister {
	return &managedClusterSetBindingLister{indexer: indexer}
}

// List lists all ManagedClusterSetBindings in the indexer.
func (s *managedClusterSetBindingLister) List(selector labels.Selector) (ret []*v1beta2.ManagedClusterSetBinding, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta2.ManagedClusterSetBinding))
	})
	return ret, err
}

// ManagedClusterSetBindings returns an object that can list and get ManagedClusterSetBindings.
func (s *managedClusterSetBindingLister) ManagedClusterSetBindings(namespace string) ManagedClusterSetBindingNamespaceLister {
	return managedClusterSetBindingNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ManagedClusterSetBindingNamespaceLister helps list and get ManagedClusterSetBindings.
// All objects returned here must be treated as read-only.
type ManagedClusterSetBindingNamespaceLister interface {
	// List lists all ManagedClusterSetBindings in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta2.ManagedClusterSetBinding, err error)
	// Get retrieves the ManagedClusterSetBinding from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1beta2.ManagedClusterSetBinding, error)
	ManagedClusterSetBindingNamespaceListerExpansion
}

// managedClusterSetBindingNamespaceLister implements the ManagedClusterSetBindingNamespaceLister
// interface.
type managedClusterSetBindingNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ManagedClusterSetBindings in the indexer for a given namespace.
func (s managedClusterSetBindingNamespaceLister) List(selector labels.Selector) (ret []*v1beta2.ManagedClusterSetBinding, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta2.ManagedClusterSetBinding))
	})
	return ret, err
}

// Get retrieves the ManagedClusterSetBinding from the indexer for a given namespace and name.
func (s managedClusterSetBindingNamespaceLister) Get(name string) (*v1beta2.ManagedClusterSetBinding, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1beta2.Resource("managedclustersetbinding"), name)
	}
	return obj.(*v1beta2.ManagedClusterSetBinding), nil
}
