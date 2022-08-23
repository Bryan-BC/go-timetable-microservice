package main

import (
	"log"
	"net"

	"github.com/Bryan-BC/go-timetable-microservice/pkg/client"
	"github.com/Bryan-BC/go-timetable-microservice/pkg/config"
	"github.com/Bryan-BC/go-timetable-microservice/pkg/db"
	"github.com/Bryan-BC/go-timetable-microservice/pkg/pb"
	"github.com/Bryan-BC/go-timetable-microservice/pkg/services"
	"google.golang.org/grpc"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Panicf("Error loading config, %s \n", err)
	}

	db := db.Init(c.DBURL)

	courseClient := client.NewCourseClient("localhost:5000")

	listener, err := net.Listen("tcp", c.Port)
	if err != nil {
		log.Panicf("Error listening, %s \n", err)
	}

	log.Printf("Listening on port %s \n", c.Port)

	s := services.Server{
		DBPointer:    &db,
		CourseClient: courseClient,
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTimetableServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(listener); err != nil {
		log.Panicf("Error serving timetable microservice, %s \n", err)
	}
}
