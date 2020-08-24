/*
Copyright 2019 Red Hat Inc.

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

package plan

import (
	"context"
	"errors"
	"github.com/konveyor/controller/pkg/logging"
	libref "github.com/konveyor/controller/pkg/ref"
	api "github.com/konveyor/virt-controller/pkg/apis/virt/v1alpha1"
	"github.com/konveyor/virt-controller/pkg/settings"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"time"
)

var log = logging.WithName("plan")

//
// Application settings.
var (
	Settings = &settings.Settings
	FastReQ  = time.Second * 3
)

//
// Creates a new Plan Controller and adds it to the Manager.
func Add(mgr manager.Manager) error {
	reconciler := &Reconciler{
		Client: mgr.GetClient(),
		scheme: mgr.GetScheme(),
	}
	cnt, err := controller.New(
		"plan-controller",
		mgr,
		controller.Options{
			Reconciler: reconciler,
		})
	if err != nil {
		log.Trace(err)
		return err
	}
	// Primary CR.
	err = cnt.Watch(
		&source.Kind{Type: &api.Plan{}},
		&handler.EnqueueRequestForObject{},
		&PlanPredicate{})
	if err != nil {
		log.Trace(err)
		return err
	}
	// References.
	err = cnt.Watch(
		&source.Kind{
			Type: &api.Provider{},
		},
		libref.Handler(),
		&ProviderPredicate{})
	if err != nil {
		log.Trace(err)
		return err
	}
	err = cnt.Watch(
		&source.Kind{
			Type: &api.Map{},
		},
		libref.Handler(),
		&MapPredicate{})
	if err != nil {
		log.Trace(err)
		return err
	}
	err = cnt.Watch(
		&source.Kind{
			Type: &api.Host{},
		},
		libref.Handler(),
		&MapPredicate{})
	if err != nil {
		log.Trace(err)
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &Reconciler{}

//
// Reconciles a Plan object.
type Reconciler struct {
	client.Client
	scheme *runtime.Scheme
}

//
// Reconcile a Plan CR.
func (r *Reconciler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	var err error

	// Reset the logger.
	log.Reset()

	// Fetch the CR.
	plan := &api.Plan{}
	err = r.Get(context.TODO(), request.NamespacedName, plan)
	if err != nil {
		if k8serr.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		log.Trace(err)
		return reconcile.Result{}, err
	}

	// Begin staging conditions.
	plan.Status.BeginStagingConditions()

	// Validations.
	err = r.validate(plan)
	if err != nil {
		if errors.Is(err, ProviderInvNotReady) {
			return reconcile.Result{
				RequeueAfter: FastReQ,
			}, nil
		}
		log.Trace(err)
		return reconcile.Result{RequeueAfter: FastReQ}, nil
	}

	// Ready condition.
	if !plan.Status.HasBlockerCondition() {
		plan.Status.SetReady(true, ReadyMessage)
	}

	// End staging conditions.
	plan.Status.EndStagingConditions()

	// Apply changes.
	plan.Status.ObservedGeneration = plan.Generation
	err = r.Status().Update(context.TODO(), plan)
	if err != nil {
		log.Trace(err)
		return reconcile.Result{RequeueAfter: FastReQ}, nil
	}

	// Done
	return reconcile.Result{}, nil
}
