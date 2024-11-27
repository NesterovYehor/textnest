package models

import (
	"time"

	download_service "github.com/NesterovYehor/TextNest/services/download_service/internal/grpc_server"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Metadata struct {
	Key         string
	CreatedAt   time.Time
	ExpiredDate time.Time
}

// Convert models.Metadata to Protobuf Metadata
func (m *Metadata) ToProto() *download_service.Metadata {
	return &download_service.Metadata{
		Key:         m.Key,
		CreatedAt:   timestamppb.New(m.CreatedAt),
		ExpiredDate: timestamppb.New(m.ExpiredDate),
	}
}

// Convert Protobuf Metadata to models.Metadata
func FromProto(proto *download_service.Metadata) *Metadata {
	return &Metadata{
		Key:         proto.Key,
		CreatedAt:   proto.CreatedAt.AsTime(),
		ExpiredDate: proto.ExpiredDate.AsTime(),
	}
}
