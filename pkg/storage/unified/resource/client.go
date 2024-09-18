package resource

import (
	"github.com/fullstorydev/grpchan"
	"github.com/fullstorydev/grpchan/inprocgrpc"
	grpcAuth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"google.golang.org/grpc"

	grpcUtils "github.com/grafana/grafana/pkg/storage/unified/resource/grpc"
)

type ResourceClient interface {
	ResourceStoreClient
	ResourceIndexClient
	BlobStoreClient
	DiagnosticsClient
}

// Internal implementation
type resourceClient struct {
	ResourceStoreClient
	ResourceIndexClient
	BlobStoreClient
	DiagnosticsClient
}

func NewResourceClient(channel *grpc.ClientConn) ResourceClient {
	cc := grpchan.InterceptClientConn(channel, grpcUtils.UnaryClientInterceptor, grpcUtils.StreamClientInterceptor)
	return &resourceClient{
		ResourceStoreClient: NewResourceStoreClient(cc),
		ResourceIndexClient: NewResourceIndexClient(cc),
		BlobStoreClient:     NewBlobStoreClient(cc),
		DiagnosticsClient:   NewDiagnosticsClient(cc),
	}
}

func NewLocalResourceClient(server ResourceServer) ResourceClient {
	channel := &inprocgrpc.Channel{}

	auth := &grpcUtils.Authenticator{}

	channel.RegisterService(
		grpchan.InterceptServer(
			&ResourceStore_ServiceDesc,
			grpcAuth.UnaryServerInterceptor(auth.Authenticate),
			grpcAuth.StreamServerInterceptor(auth.Authenticate),
		),
		server,
	)

	channel.RegisterService(
		grpchan.InterceptServer(
			&ResourceIndex_ServiceDesc,
			grpcAuth.UnaryServerInterceptor(auth.Authenticate),
			grpcAuth.StreamServerInterceptor(auth.Authenticate),
		),
		server,
	)

	channel.RegisterService(
		grpchan.InterceptServer(
			&BlobStore_ServiceDesc,
			grpcAuth.UnaryServerInterceptor(auth.Authenticate),
			grpcAuth.StreamServerInterceptor(auth.Authenticate),
		),
		server,
	)

	channel.RegisterService(
		grpchan.InterceptServer(
			&Diagnostics_ServiceDesc,
			grpcAuth.UnaryServerInterceptor(auth.Authenticate),
			grpcAuth.StreamServerInterceptor(auth.Authenticate),
		),
		server,
	)

	cc := grpchan.InterceptClientConn(channel, grpcUtils.UnaryClientInterceptor, grpcUtils.StreamClientInterceptor)
	return &resourceClient{
		ResourceStoreClient: NewResourceStoreClient(cc),
		ResourceIndexClient: NewResourceIndexClient(cc),
		BlobStoreClient:     NewBlobStoreClient(cc),
		DiagnosticsClient:   NewDiagnosticsClient(cc),
	}
}
