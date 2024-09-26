// SPDX-License-Identifier: AGPL-3.0-only

// Code generated by client-gen. DO NOT EDIT.

package v0alpha1

import (
	"context"

	v0alpha1 "github.com/grafana/grafana/pkg/apis/alerting_notifications/v0alpha1"
	alertingnotificationsv0alpha1 "github.com/grafana/grafana/pkg/generated/applyconfiguration/alerting_notifications/v0alpha1"
	scheme "github.com/grafana/grafana/pkg/generated/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	gentype "k8s.io/client-go/gentype"
)

// TemplateGroupsGetter has a method to return a TemplateGroupInterface.
// A group's client should implement this interface.
type TemplateGroupsGetter interface {
	TemplateGroups(namespace string) TemplateGroupInterface
}

// TemplateGroupInterface has methods to work with TemplateGroup resources.
type TemplateGroupInterface interface {
	Create(ctx context.Context, templateGroup *v0alpha1.TemplateGroup, opts v1.CreateOptions) (*v0alpha1.TemplateGroup, error)
	Update(ctx context.Context, templateGroup *v0alpha1.TemplateGroup, opts v1.UpdateOptions) (*v0alpha1.TemplateGroup, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v0alpha1.TemplateGroup, error)
	List(ctx context.Context, opts v1.ListOptions) (*v0alpha1.TemplateGroupList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v0alpha1.TemplateGroup, err error)
	Apply(ctx context.Context, templateGroup *alertingnotificationsv0alpha1.TemplateGroupApplyConfiguration, opts v1.ApplyOptions) (result *v0alpha1.TemplateGroup, err error)
	TemplateGroupExpansion
}

// templateGroups implements TemplateGroupInterface
type templateGroups struct {
	*gentype.ClientWithListAndApply[*v0alpha1.TemplateGroup, *v0alpha1.TemplateGroupList, *alertingnotificationsv0alpha1.TemplateGroupApplyConfiguration]
}

// newTemplateGroups returns a TemplateGroups
func newTemplateGroups(c *NotificationsV0alpha1Client, namespace string) *templateGroups {
	return &templateGroups{
		gentype.NewClientWithListAndApply[*v0alpha1.TemplateGroup, *v0alpha1.TemplateGroupList, *alertingnotificationsv0alpha1.TemplateGroupApplyConfiguration](
			"templategroups",
			c.RESTClient(),
			scheme.ParameterCodec,
			namespace,
			func() *v0alpha1.TemplateGroup { return &v0alpha1.TemplateGroup{} },
			func() *v0alpha1.TemplateGroupList { return &v0alpha1.TemplateGroupList{} }),
	}
}
