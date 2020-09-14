package ocp

import (
	"context"
	net "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	liberr "github.com/konveyor/controller/pkg/error"
	libocp "github.com/konveyor/controller/pkg/inventory/container/ocp"
	libref "github.com/konveyor/controller/pkg/ref"
	model "github.com/konveyor/virt-controller/pkg/controller/provider/model/ocp"
	storage "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

//
// StorageClass
type StorageClass struct {
	libocp.BaseCollection
}

//
// Get the kubernetes object being collected.
func (r *StorageClass) Object() runtime.Object {
	return &storage.StorageClass{}
}

//
// Reconcile.
// Achieve initial consistency.
func (r *StorageClass) Reconcile(ctx context.Context) (err error) {
	pClient := r.Reconciler.Client()
	list := &storage.StorageClassList{}
	err = pClient.List(context.TODO(), nil, list)
	if err != nil {
		err = liberr.Wrap(err)
		return
	}
	db := r.Reconciler.DB()
	tx, err := db.Begin()
	if err != nil {
		err = liberr.Wrap(err)
		return
	}
	defer tx.End()
	for _, resource := range list.Items {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		m := &model.StorageClass{}
		m.With(&resource)
		r.Reconciler.UpdateThreshold(m)
		Log.Info("Create", libref.ToKind(m), m.String())
		err = db.Insert(m)
		if err != nil {
			err = liberr.Wrap(err)
			return
		}
	}
	err = tx.Commit()
	if err != nil {
		err = liberr.Wrap(err)
		return
	}

	return
}

//
// Resource created watch event.
func (r *StorageClass) Create(e event.CreateEvent) bool {
	object, cast := e.Object.(*storage.StorageClass)
	if !cast {
		return false
	}
	m := &model.StorageClass{}
	m.With(object)
	r.Reconciler.Create(m)

	return false
}

//
// Resource updated watch event.
func (r *StorageClass) Update(e event.UpdateEvent) bool {
	object, cast := e.ObjectNew.(*storage.StorageClass)
	if !cast {
		return false
	}
	m := &model.StorageClass{}
	m.With(object)
	r.Reconciler.Update(m)

	return false
}

//
// Resource deleted watch event.
func (r *StorageClass) Delete(e event.DeleteEvent) bool {
	object, cast := e.Object.(*storage.StorageClass)
	if !cast {
		return false
	}
	m := &model.StorageClass{}
	m.With(object)
	r.Reconciler.Delete(m)

	return false
}

//
// Ignored.
func (r *StorageClass) Generic(e event.GenericEvent) bool {
	return false
}

//
// NetworkAttachmentDefinition
type NetworkAttachmentDefinition struct {
	libocp.BaseCollection
}

//
// Get the kubernetes object being collected.
func (r *NetworkAttachmentDefinition) Object() runtime.Object {
	return &net.NetworkAttachmentDefinition{}
}

//
// Reconcile.
// Achieve initial consistency.
func (r *NetworkAttachmentDefinition) Reconcile(ctx context.Context) (err error) {
	pClient := r.Reconciler.Client()
	list := &net.NetworkAttachmentDefinitionList{}
	err = pClient.List(context.TODO(), nil, list)
	if err != nil {
		err = liberr.Wrap(err)
		return
	}
	db := r.Reconciler.DB()
	tx, err := db.Begin()
	if err != nil {
		err = liberr.Wrap(err)
		return
	}
	defer tx.End()
	for _, resource := range list.Items {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		m := &model.NetworkAttachmentDefinition{}
		m.With(&resource)
		r.Reconciler.UpdateThreshold(m)
		Log.Info("Create", libref.ToKind(m), m.String())
		err = db.Insert(m)
		if err != nil {
			err = liberr.Wrap(err)
			return
		}
	}
	err = tx.Commit()
	if err != nil {
		err = liberr.Wrap(err)
		return
	}

	return
}

//
// Resource created watch event.
func (r *NetworkAttachmentDefinition) Create(e event.CreateEvent) bool {
	object, cast := e.Object.(*net.NetworkAttachmentDefinition)
	if !cast {
		return false
	}
	m := &model.NetworkAttachmentDefinition{}
	m.With(object)
	r.Reconciler.Create(m)

	return false
}

//
// Resource updated watch event.
func (r *NetworkAttachmentDefinition) Update(e event.UpdateEvent) bool {
	object, cast := e.ObjectNew.(*net.NetworkAttachmentDefinition)
	if !cast {
		return false
	}
	m := &model.NetworkAttachmentDefinition{}
	m.With(object)
	r.Reconciler.Update(m)

	return false
}

//
// Resource deleted watch event.
func (r *NetworkAttachmentDefinition) Delete(e event.DeleteEvent) bool {
	object, cast := e.Object.(*net.NetworkAttachmentDefinition)
	if !cast {
		return false
	}
	m := &model.NetworkAttachmentDefinition{}
	m.With(object)
	r.Reconciler.Delete(m)

	return false
}

//
// Ignored.
func (r *NetworkAttachmentDefinition) Generic(e event.GenericEvent) bool {
	return false
}