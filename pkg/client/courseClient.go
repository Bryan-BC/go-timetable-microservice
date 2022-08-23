package client

import (
	"context"
	"log"

	"github.com/Bryan-BC/go-timetable-microservice/pkg/pb"
	"google.golang.org/grpc"
)

type CourseClient struct {
	Client pb.CourseServiceClient
}

func NewCourseClient(url string) *CourseClient {
	cc, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		log.Panicf("Error dialing gRPC server, %s \n", err)
	}

	return &CourseClient{Client: pb.NewCourseServiceClient(cc)}
}

func (svc *CourseClient) GetCourse(id string) (*pb.GetCourseResponse, error) {
	return svc.Client.GetCourse(
		context.Background(),
		&pb.GetCourseRequest{Id: id},
	)
}
