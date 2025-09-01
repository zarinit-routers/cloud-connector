package grpc

import (
	"context"

	"github.com/zarinit-routers/cloud-connector/models"
	"github.com/zarinit-routers/cloud-connector/storage/repository"
	pb "github.com/zarinit-routers/connector-rpc/gen/connector"
)

type nodesService struct {
	pb.UnimplementedNodesServer
}

func newNodesService() *nodesService {
	return &nodesService{}
}

func (s *nodesService) NodesByGroup(ctx context.Context, req *pb.NodesByGroupRequest) (*pb.NodesResponse, error) {
	groupId, err := models.UUIDFromString(req.GroupId)
	if err != nil {
		return nil, err
	}

	data, err := repository.GetQueries().GetNodes(ctx, groupId)
	if err != nil {
		return nil, err
	}
	nodes := []*pb.Node{}
	for _, d := range data {

		tags, err := repository.GetQueries().GetTags(ctx, d.Id)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, &pb.Node{
			Id:   d.Id.String(),
			Name: d.Name.String,
			Tags: tags,
		})
	}
	return &pb.NodesResponse{
		Clients: nodes,
	}, nil
}

func (s *nodesService) AddTag(ctx context.Context, req *pb.TagRequest) (*pb.Node, error) {
	id, err := models.UUIDFromString(req.ModeId)
	if err != nil {
		return nil, err
	}
	err = repository.GetQueries().AddTag(ctx, repository.AddTagParams{
		Tag:    req.Tag,
		NodeId: id,
	})
	if err != nil {
		return nil, err
	}

	node, err := repository.GetQueries().GetNode(ctx, id)
	if err != nil {
		return nil, err
	}
	tags, err := repository.GetQueries().GetTags(ctx, id)
	if err != nil {
		return nil, err
	}
	return &pb.Node{
		Id:   req.ModeId,
		Name: node.Name.String,
		Tags: tags,
	}, nil
}

func (s *nodesService) RemoveTag(ctx context.Context, req *pb.TagRequest) (*pb.Node, error) {
	id, err := models.UUIDFromString(req.ModeId)
	if err != nil {
		return nil, err
	}
	err = repository.GetQueries().RemoveTag(ctx, repository.RemoveTagParams{
		Tag:    req.Tag,
		NodeId: id,
	})
	if err != nil {
		return nil, err
	}

	node, err := repository.GetQueries().GetNode(ctx, id)
	if err != nil {
		return nil, err
	}
	tags, err := repository.GetQueries().GetTags(ctx, id)
	if err != nil {
		return nil, err
	}
	return &pb.Node{
		Id:   req.ModeId,
		Name: node.Name.String,
		Tags: tags,
	}, nil
}
