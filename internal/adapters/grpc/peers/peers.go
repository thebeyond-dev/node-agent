package amneziawg

import (
	"context"

	"github.com/thebeyond-net/node-agent/internal/core/ports"
	v1 "github.com/thebeyond-net/node-agent/pkg/amneziawg/v1"
	"github.com/thebeyond-net/node-agent/pkg/amneziawg/v1/amneziawgv1connect"
)

type AmneziaWGServiceServer struct {
	amneziawgv1connect.UnimplementedAmneziaWGServiceHandler
	usecase ports.PeerUseCase
}

func NewAmneziaWGServiceServer(usecase ports.PeerUseCase) *AmneziaWGServiceServer {
	return &AmneziaWGServiceServer{usecase: usecase}
}

func (s *AmneziaWGServiceServer) CreatePeer(ctx context.Context, req *v1.CreatePeerRequest) (*v1.CreatePeerResponse, error) {
	cfg, pubKey, err := s.usecase.CreatePeer(ctx)
	return &v1.CreatePeerResponse{
		Success: err == nil,
		Config:  cfg,
		Pubkey:  pubKey,
	}, err
}

func (s *AmneziaWGServiceServer) DeletePeer(ctx context.Context, req *v1.DeletePeerRequest) (*v1.DeletePeerResponse, error) {
	err := s.usecase.DeletePeer(ctx, req.GetPubkey())
	return &v1.DeletePeerResponse{
		Success: err == nil,
	}, err
}
