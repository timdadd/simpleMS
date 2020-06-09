package routeguide_test

import (
	"context"
	"fmt"
	"lib/grpc_test"
	rgmock "routeguide/mox"
	pb "routeguide/pb"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"google.golang.org/protobuf/proto"
)

type s struct {
	grpc_test.Tester
}

func Test(t *testing.T) {
	grpc_test.RunSubTests(t, s{})
}

var msg = pb.RouteNote{
	Location: &pb.Point{Latitude: 17, Longitude: 29},
	Message:  "Taxi-cab",
}

func (s) TestRouteChat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock for the stream returned by RouteChat
	stream := rgmock.NewMockRouteGuide_RouteChatClient(ctrl)
	// set expectation on sending.
	stream.EXPECT().Send(
		gomock.Any(),
	).Return(nil)
	// Set expectation on receiving.
	stream.EXPECT().Recv().Return(&msg, nil)
	stream.EXPECT().CloseSend().Return(nil)
	// Create mock for the client interface.
	rgclient := rgmock.NewMockRouteGuideClient(ctrl)
	// Set expectation on RouteChat
	rgclient.EXPECT().RouteChat(
		gomock.Any(),
	).Return(stream, nil)
	if err := testRouteChat(rgclient); err != nil {
		t.Fatalf("Test failed: %v", err)
	}
}

func testRouteChat(client pb.RouteGuideClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.RouteChat(ctx)
	if err != nil {
		return err
	}
	if err := stream.Send(&msg); err != nil {
		return err
	}
	if err := stream.CloseSend(); err != nil {
		return err
	}
	got, err := stream.Recv()
	if err != nil {
		return err
	}
	if !proto.Equal(got, &msg) {
		return fmt.Errorf("stream.Recv() = %v, want %v", got, &msg)
	}
	return nil
}
