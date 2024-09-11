package accesscontrol

import (
	"context"

	openfgav1 "github.com/openfga/api/proto/openfga/v1"

	"github.com/grafana/grafana/pkg/apimachinery/identity"
	"github.com/grafana/grafana/pkg/services/accesscontrol"
)

type recordingAccessControlFake struct {
	Disabled           bool
	EvaluateRecordings []struct {
		Permissions map[string][]string
		Evaluator   accesscontrol.Evaluator
	}
	Callback func(user identity.Requester, evaluator accesscontrol.Evaluator) (bool, error)
}

func (a *recordingAccessControlFake) Evaluate(_ context.Context, ur identity.Requester, evaluator accesscontrol.Evaluator) (bool, error) {
	a.EvaluateRecordings = append(a.EvaluateRecordings, struct {
		Permissions map[string][]string
		Evaluator   accesscontrol.Evaluator
	}{Permissions: ur.GetPermissions(), Evaluator: evaluator})
	if a.Callback == nil {
		return evaluator.Evaluate(ur.GetPermissions()), nil
	}
	return a.Callback(ur, evaluator)
}

func (a *recordingAccessControlFake) RegisterScopeAttributeResolver(prefix string, resolver accesscontrol.ScopeAttributeResolver) {
	// TODO implement me
	panic("implement me")
}

func (a *recordingAccessControlFake) IsDisabled() bool {
	return a.Disabled
}

func (a *recordingAccessControlFake) Check(ctx context.Context, in *openfgav1.CheckRequest) (*openfgav1.CheckResponse, error) {
	return nil, nil
}

func (a *recordingAccessControlFake) ListObjects(ctx context.Context, in *openfgav1.ListObjectsRequest) (*openfgav1.ListObjectsResponse, error) {
	return nil, nil
}

var _ accesscontrol.AccessControl = &recordingAccessControlFake{}
