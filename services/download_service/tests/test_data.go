package tests

import (
	"time"

	pb "github.com/NesterovYehor/TextNest/services/download_service/api"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	key            = "test-key"
	userId         = "test-userid"
	expirationDate = timestamppb.New(time.Now().Add(time.Minute))
	title          = "test-title"
)

var data = pb.Metadata{
	Key:         key,
	ExpiredDate: expirationDate,
	Title:       title,
	CreatedAt:   timestamppb.New(time.Now()),
}
