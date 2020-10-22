package main

import (
	"context"

	pb "github.com/m3o/services/infrastructure/proto"
	"github.com/micro/micro/v3/service/errors"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	k8s "github.com/scaleway/scaleway-sdk-go/api/k8s/v1beta4"
	"github.com/scaleway/scaleway-sdk-go/api/lb/v1"
)

type infrastructure struct{}

func (i *infrastructure) Check(ctx context.Context, request *pb.CheckRequest, response *pb.CheckResponse) error {
	issues, err := checkInfraUsage()
	if err != nil {
		return err
	}
	response.Issues = make([]string, len(issues))
	for i, iss := range issues {
		response.Issues[i] = iss
	}
	return nil
}

func (i *infrastructure) Summary(ctx context.Context, req *pb.SummaryRequest, rsp *pb.SummaryResponse) error {
	clRsp, err := k8sAPI.ListClusters(&k8s.ListClustersRequest{Region: scalewayRegion})
	if err != nil {
		return errors.InternalServerError("infrastructure.Summary", "Error listing clusters: %v", err)
	}
	rsp.Clusters = int32(clRsp.TotalCount)

	lbRsp, err := lbAPI.ListLbs(&lb.ListLbsRequest{Region: scalewayRegion})
	if err != nil {
		return errors.InternalServerError("infrastructure.Summary", "Error listing load balancers: %v", err)
	}
	rsp.LoadBalancers = int32(lbRsp.TotalCount)

	svrRsp, err := inAPI.ListServers(&instance.ListServersRequest{Zone: scalewayZone})
	if err != nil {
		return errors.InternalServerError("infrastructure.Summary", "Error listing servers: %v", err)
	}
	rsp.Servers = int32(svrRsp.TotalCount)

	volRsp, err := inAPI.ListVolumes(&instance.ListVolumesRequest{Zone: scalewayZone})
	if err != nil {
		return errors.InternalServerError("infrastructure.Summary", "Error listing volumes: %v", err)
	}
	rsp.Volumes = int32(volRsp.TotalCount)

	return nil
}
