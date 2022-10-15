package biz

import (
	"context"

	pb "moredoc/api/v1"
	"moredoc/middleware/auth"
	"moredoc/model"
	"moredoc/util"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type GroupAPIService struct {
	pb.UnimplementedGroupAPIServer
	dbModel *model.DBModel
	logger  *zap.Logger
}

func NewGroupAPIService(dbModel *model.DBModel, logger *zap.Logger) (service *GroupAPIService) {
	return &GroupAPIService{dbModel: dbModel, logger: logger.Named("GroupAPIService")}
}

// CreateGroup 创建用户组
// 0. 检查用户权限
// 1. 检查用户组是否存在
// 2. 创建用户组
func (s *GroupAPIService) CreateGroup(ctx context.Context, req *pb.Group) (*pb.Group, error) {
	s.logger.Debug("CreateGroup", zap.Any("req", req))

	userClaims, ok := ctx.Value(auth.CtxKeyUserClaims).(*auth.UserClaims)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, ErrorMessageInvalidToken)
	}

	fullMethod, _ := ctx.Value(auth.CtxKeyFullMethod).(string)
	if yes := s.dbModel.CheckPermissionByUserId(userClaims.UserId, fullMethod); !yes {
		return nil, status.Errorf(codes.PermissionDenied, ErrorMessagePermissionDenied)
	}

	return &pb.Group{}, nil
}
func (s *GroupAPIService) UpdateGroup(ctx context.Context, req *pb.Group) (*pb.Group, error) {
	return &pb.Group{}, nil
}
func (s *GroupAPIService) DeleteGroup(ctx context.Context, req *pb.DeleteGroupRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
func (s *GroupAPIService) GetGroup(ctx context.Context, req *pb.GetGroupRequest) (*pb.Group, error) {
	return &pb.Group{}, nil
}

// ListGroup 列出用户组。所有用户都可以查询
func (s *GroupAPIService) ListGroup(ctx context.Context, req *pb.ListGroupRequest) (*pb.ListGroupReply, error) {
	s.logger.Debug("ListGroup", zap.Any("req", req))
	groups, _, err := s.dbModel.GetGroupList(model.OptionGetGroupList{Page: 1, Size: 100000, SelectFields: req.Field, WithCount: false})
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	var pbGroups []*pb.Group
	util.CopyStruct(&groups, &pbGroups)
	return &pb.ListGroupReply{Group: pbGroups}, nil
}