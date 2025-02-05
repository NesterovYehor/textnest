package models

import (
	"time"

	pb "github.com/NesterovYehor/TextNest/services/download_service/api"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Metadata struct {
	Key            string
	CreatedAt      time.Time
	ExpirationDate time.Time // Renamed from ExpiredDate
}

// Convert models.Metadata to Protobuf Metadata
func (m *Metadata) ToProto() *pb.Metadata {
	return &pb.Metadata{
		Key:         m.Key,
		CreatedAt:   timestamppb.New(m.CreatedAt),
		ExpiredDate: timestamppb.New(m.ExpirationDate),
	}
}

// Convert Protobuf Metadata to models.Metadata
func FromProto(proto *pb.Metadata) *Metadata {
	return &Metadata{
		Key:         proto.Key,
		CreatedAt:   proto.CreatedAt.AsTime(),
		ExpirationDate: proto.ExpiredDate.AsTime(),
	}
}
