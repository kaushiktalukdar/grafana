// SPDX-License-Identifier: AGPL-3.0-only

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package applyconfiguration

import (
	v0alpha1 "github.com/grafana/grafana/pkg/apis/alerting_notifications/v0alpha1"
	servicev0alpha1 "github.com/grafana/grafana/pkg/apis/service/v0alpha1"
	alertingnotificationsv0alpha1 "github.com/grafana/grafana/pkg/generated/applyconfiguration/alerting_notifications/v0alpha1"
	internal "github.com/grafana/grafana/pkg/generated/applyconfiguration/internal"
	applyconfigurationservicev0alpha1 "github.com/grafana/grafana/pkg/generated/applyconfiguration/service/v0alpha1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	testing "k8s.io/client-go/testing"
)

// ForKind returns an apply configuration type for the given GroupVersionKind, or nil if no
// apply configuration type exists for the given GroupVersionKind.
func ForKind(kind schema.GroupVersionKind) interface{} {
	switch kind {
	// Group=notifications.alerting.grafana.app, Version=v0alpha1
	case v0alpha1.SchemeGroupVersion.WithKind("Integration"):
		return &alertingnotificationsv0alpha1.IntegrationApplyConfiguration{}
	case v0alpha1.SchemeGroupVersion.WithKind("Interval"):
		return &alertingnotificationsv0alpha1.IntervalApplyConfiguration{}
	case v0alpha1.SchemeGroupVersion.WithKind("Matcher"):
		return &alertingnotificationsv0alpha1.MatcherApplyConfiguration{}
	case v0alpha1.SchemeGroupVersion.WithKind("Receiver"):
		return &alertingnotificationsv0alpha1.ReceiverApplyConfiguration{}
	case v0alpha1.SchemeGroupVersion.WithKind("ReceiverSpec"):
		return &alertingnotificationsv0alpha1.ReceiverSpecApplyConfiguration{}
	case v0alpha1.SchemeGroupVersion.WithKind("Route"):
		return &alertingnotificationsv0alpha1.RouteApplyConfiguration{}
	case v0alpha1.SchemeGroupVersion.WithKind("RouteDefaults"):
		return &alertingnotificationsv0alpha1.RouteDefaultsApplyConfiguration{}
	case v0alpha1.SchemeGroupVersion.WithKind("RoutingTree"):
		return &alertingnotificationsv0alpha1.RoutingTreeApplyConfiguration{}
	case v0alpha1.SchemeGroupVersion.WithKind("RoutingTreeSpec"):
		return &alertingnotificationsv0alpha1.RoutingTreeSpecApplyConfiguration{}
	case v0alpha1.SchemeGroupVersion.WithKind("TemplateGroup"):
		return &alertingnotificationsv0alpha1.TemplateGroupApplyConfiguration{}
	case v0alpha1.SchemeGroupVersion.WithKind("TemplateGroupSpec"):
		return &alertingnotificationsv0alpha1.TemplateGroupSpecApplyConfiguration{}
	case v0alpha1.SchemeGroupVersion.WithKind("TimeInterval"):
		return &alertingnotificationsv0alpha1.TimeIntervalApplyConfiguration{}
	case v0alpha1.SchemeGroupVersion.WithKind("TimeIntervalSpec"):
		return &alertingnotificationsv0alpha1.TimeIntervalSpecApplyConfiguration{}
	case v0alpha1.SchemeGroupVersion.WithKind("TimeRange"):
		return &alertingnotificationsv0alpha1.TimeRangeApplyConfiguration{}

		// Group=service.grafana.app, Version=v0alpha1
	case servicev0alpha1.SchemeGroupVersion.WithKind("ExternalName"):
		return &applyconfigurationservicev0alpha1.ExternalNameApplyConfiguration{}
	case servicev0alpha1.SchemeGroupVersion.WithKind("ExternalNameSpec"):
		return &applyconfigurationservicev0alpha1.ExternalNameSpecApplyConfiguration{}

	}
	return nil
}

func NewTypeConverter(scheme *runtime.Scheme) *testing.TypeConverter {
	return &testing.TypeConverter{Scheme: scheme, TypeResolver: internal.Parser()}
}
