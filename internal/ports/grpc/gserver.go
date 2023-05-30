package grpc

import (
	"context"
	"ads/internal/app"
	"log"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type gRPCServer interface {
	AdServiceServer
}

type gRPCServerStruct struct {
	A app.App
	UnimplementedAdServiceServer
}

func (g *gRPCServerStruct) CreateAd(ctx context.Context, req *CreateAdRequest) (*AdResponse, error) {
	if err := g.A.CheckUser(ctx, req.GetUserId()); err != nil {
		log.Println("not found user in db for create ad ", err)
		return nil, status.Error(codes.NotFound, "User not found")
	}
	ad, err := g.A.CreateAd(ctx, req.GetTitle(), req.GetText(), req.GetUserId())
	if err != nil {
		log.Println("error in create ad ", err)
		return nil, status.Error(codes.InvalidArgument, "error create ad")
	}
	log.Printf("user %v && create ad %v \n", ad.AuthorID, ad.ID)
	return &AdResponse{AuthorId: ad.AuthorID, Id: ad.ID, Published: ad.Published, Title: ad.Title, Text: ad.Text}, nil
}

func (g *gRPCServerStruct) ChangeAdStatus(ctx context.Context, req *ChangeAdStatusRequest) (*AdResponse, error) {
	ad, err := g.A.ChangeAdStatus(ctx, req.GetAdId(), req.GetPublished(), req.GetUserId())
	if err != nil { 
		log.Println("error in change status: ", err)
		return nil, status.Error(codes.InvalidArgument, "error change status")
	}
	log.Println("change ad status: adID ", ad.ID, " published: ", ad.Published)
	return &AdResponse{AuthorId: ad.AuthorID, Id: ad.ID, Published: ad.Published, Title: ad.Title, Text: ad.Text}, nil
}

func (g *gRPCServerStruct) UpdateAd(ctx context.Context, req *UpdateAdRequest) (*AdResponse, error) {
	ad, err := g.A.UpdateAd(ctx, req.GetUserId(), req.GetTitle(), req.GetText(), req.GetAdId())
	if err != nil {
		log.Println("error in update ad ", err) 
		return nil, status.Error(codes.InvalidArgument, "error update ad")
	}
	log.Println("update ad ", ad)
	return &AdResponse{AuthorId: ad.AuthorID, Id: ad.ID, Published: ad.Published, Title: ad.Title, Text: ad.Text}, nil
}

func (g *gRPCServerStruct) ListAds(ctx context.Context, req *emptypb.Empty) (*ListAdResponse, error) {
	ads, err := g.A.ListAds(ctx)
	if err != nil {
		log.Println("error in list ads ", err) 
		return nil, status.Error(codes.InvalidArgument, "error list ads")
	}
	var adsResponse []*AdResponse
	for _, ad := range ads {
		adResponse := &AdResponse{
			AuthorId: ad.AuthorID,
			Id: ad.ID,
			Published: ad.Published,
			Title: ad.Title,
			Text: ad.Text,
		}
		adsResponse = append(adsResponse, adResponse)
	}
	log.Println("got ads list")
	return &ListAdResponse{List: adsResponse}, nil
}

func (g *gRPCServerStruct) CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
	u, err := g.A.CreateUser(ctx, req.GetName(), "ourGopher@mai.com")
	if err != nil {
		log.Println("error in create user", err)
		return nil, status.Error(codes.InvalidArgument, "error create user")
	}
	log.Println("create user", u)
	return &UserResponse{Name: u.NickName, Id: u.UserID}, nil
}

func (g *gRPCServerStruct) GetUser(ctx context.Context, req *GetUserRequest) (*UserResponse, error) {
	u, err := g.A.GetUser(ctx, req.GetId())
	if err != nil {
		log.Println("error in get user", err)
		return nil, status.Error(codes.InvalidArgument, "error get user")
	}
	log.Println("got user", u)
	return &UserResponse{Name: u.NickName, Id: u.UserID}, nil
}

func (g *gRPCServerStruct) DeleteUser(ctx context.Context, req *DeleteUserRequest) (*emptypb.Empty, error) {
	if err := g.A.DeleteUser(ctx, req.GetId()); err != nil {
		log.Println("error :: ", err)
		return nil, status.Error(codes.NotFound, "User not found")
	}
	log.Println("deleted user", req.GetId())
	return &emptypb.Empty{}, nil
}

func (g *gRPCServerStruct) DeleteAd(ctx context.Context, req *DeleteAdRequest) (*emptypb.Empty, error) {
	ad, err := g.A.DeleteAd(ctx, req.GetAuthorId(), req.GetAdId())
	if err != nil {
		log.Println("err in deleteAd :: ", ad.ID)
		return nil, status.Error(codes.NotFound, "ad not found")
	}
	log.Println("deleted ad : ", ad)
	return &emptypb.Empty{}, nil
}

func UnaryServerInterceptorLogMethod(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	log.Println("<<<--- interceptor SERVER --->>> NAME METHOD: ", info.FullMethod)
	h, err := handler(ctx, req)
	if err != nil {
		log.Printf("err %v handler %v \n", err, h)
		return nil, err
	}
	return h, nil
}

func UnaryServerInterceptorPanicMethod(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Panicln("!!! Interceprtor Panic !!!", r, "||| NAME METHOD : ", info.FullMethod)
		}
	}()
	return handler(ctx, req)
}

func UnaryClientInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	log.Println("<<<--- interceptor CLIENT --->>> NAME METHOD: ", method)

	return invoker(ctx, method, req, reply, cc, opts...)
}

func NewService(a app.App) gRPCServer {
	return &gRPCServerStruct{A: a}
}
