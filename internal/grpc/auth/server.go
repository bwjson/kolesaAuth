package auth

import (
	"context"
	sso "github.com/bwjson/kolesa_proto/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type serverApi struct {
	sso.UnimplementedAuthServer
	auth Auth
}

type Auth interface {
	SendVerificationCode(
		ctx context.Context,
		phoneNumber string,
	) (err error)
	VerifyCode(
		ctx context.Context,
		phoneNumber string,
		code string,
	) (accessToken, refreshToken string, err error)
}

func Register(gRPCServer *grpc.Server, auth Auth) {
	sso.RegisterAuthServer(gRPCServer, &serverApi{auth: auth})
}

func (s *serverApi) SendVerificationCode(ctx context.Context, request *sso.SendVerificationCodeRequest) (*emptypb.Empty, error) {
	if len(request.PhoneNumber) != 12 {
		return nil, status.Error(codes.InvalidArgument, "Please provide phone number in appropriate format")
	}

	err := s.auth.SendVerificationCode(ctx, request.PhoneNumber)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *serverApi) VerifyCode(ctx context.Context, request *sso.VerifyCodeRequest) (*sso.VerifyCodeResponse, error) {
	if len(request.PhoneNumber) != 12 {
		return nil, status.Error(codes.InvalidArgument, "Please provide phone number in appropriate format")
	}

	if len(request.VerificationCode) != 4 {
		return nil, status.Error(codes.InvalidArgument, "Please provide correct verification code")
	}

	accessToken, refreshToken, err := s.auth.VerifyCode(ctx, request.PhoneNumber, request.VerificationCode)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &sso.VerifyCodeResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
