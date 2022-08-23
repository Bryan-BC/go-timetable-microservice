package services

import (
	"context"
	"net/http"
	"sort"

	"github.com/Bryan-BC/go-timetable-microservice/pkg/client"
	"github.com/Bryan-BC/go-timetable-microservice/pkg/db"
	"github.com/Bryan-BC/go-timetable-microservice/pkg/models"
	"github.com/Bryan-BC/go-timetable-microservice/pkg/pb"
)

type Server struct {
	DBPointer    *db.DB
	CourseClient *client.CourseClient
}

func (s *Server) GetTimetable(ctx context.Context, req *pb.GetTimetableRequest) (*pb.GetTimetableResponse, error) {
	var timetable models.Timetable

	if result := s.DBPointer.DataBase.Where(&models.Timetable{Id: req.Id}).First(&timetable); result.Error != nil {
		return &pb.GetTimetableResponse{
			Status: http.StatusNotFound,
			Error:  "Timetable not found",
		}, nil
	}

	return &pb.GetTimetableResponse{
		Status: http.StatusOK,
		Timetable: &pb.Timetable{
			Id:      timetable.Id,
			Courses: timetable.Courses,
		},
	}, nil
}

func (s *Server) CreateTimetable(ctx context.Context, req *pb.CreateTimetableRequest) (*pb.CreateTimetableResponse, error) {
	var timetable models.Timetable
	var courses []*pb.Course
	for _, courseId := range req.Courses {
		resp, _ := s.CourseClient.GetCourse(courseId)
		if resp.Status == http.StatusNotFound {
			return &pb.CreateTimetableResponse{
				Status: http.StatusNotFound,
				Error:  resp.Error,
			}, nil
		}
		courses = append(courses, resp.Course)
	}
	GenerateTimetable(&timetable, courses)

	if result := s.DBPointer.DataBase.Create(&timetable); result.Error != nil {
		return &pb.CreateTimetableResponse{
			Status: http.StatusInternalServerError,
			Error:  "Error creating timetable",
		}, nil
	}

	return &pb.CreateTimetableResponse{
		Status: http.StatusCreated,
		Timetable: &pb.Timetable{
			Id:      timetable.Id,
			Courses: timetable.Courses,
		},
	}, nil
}

func GenerateTimetable(timetable *models.Timetable, courses []*pb.Course) {
	courseLengths := make([]int, len(courses))
	courseToIdx := make(map[string]int)
	courseIndices := make([][]*pb.CourseIndex, len(courses))
	for i, course := range courses {
		courseIndices[i] = course.Schedule
		courseToIdx[course.Name] = i
		courseLengths[i] = len(course.Schedule)
	}
	idxPermutations := [][]int{}
	GeneratePerm([]int{}, courseLengths, make([]int, len(courses)), len(courses), &idxPermutations)
	for _, idxs := range idxPermutations {
		if clashes := CheckClashes(courseIndices, idxs); len(clashes) == 0 {
			timetable.Courses = []*pb.TimetableSchedule{}
			for i, idx := range idxs {
				timetable.Courses = append(timetable.Courses, &pb.TimetableSchedule{
					CourseName: courseIndices[i][idx].Name,
					Days:       ConvertCourseDayToTimetableDay(courseIndices[i][idx].Days),
				})
			}
			return
		}
	}
}

func GeneratePerm(curr, courseLengths, currIdxs []int, length int, idxPermutations *[][]int) {
	if len(curr) == length {
		*idxPermutations = append(*idxPermutations, curr)
		i := length - 1
		currIdxs[i]++
		for i >= 0 && currIdxs[i] == courseLengths[i] {
			currIdxs[i] = 0
			i--
			if i == -1 {
				return
			}
			currIdxs[i]++
		}
		GeneratePerm([]int{}, courseLengths, currIdxs, length, idxPermutations)
	} else {
		curr = append(curr, currIdxs[len(curr)])
		GeneratePerm(curr, courseLengths, currIdxs, length, idxPermutations)
	}
}

func CheckClashes(courses [][]*pb.CourseIndex, idxs []int) map[string][][]int {
	clashes := make(map[string][][]int)
	days := make(map[string][][]int)
	for i, idx := range idxs {
		for _, day := range courses[i][idx].Days {
			for _, tuple := range day.Timings {
				if _, ok := days[day.Day]; !ok {
					days[day.Day] = [][]int{}
				}
				days[day.Day] = append(days[day.Day], []int{int(tuple.Start), int(tuple.End)})
			}
		}
	}

	for _, timeslots := range days {
		sort.Slice(timeslots, func(i, j int) bool {
			return timeslots[i][0] < timeslots[j][0]
		})
	}

	for day, timeslots := range days {
		for i := 0; i < len(timeslots); i++ {
			for j := i + 1; j < len(timeslots); j++ {
				if timeslots[i][1] > timeslots[j][0] {
					clashes[day] = append(clashes[day], []int{timeslots[i][0], timeslots[j][0]})
				} else {
					break
				}
			}
		}
	}
	return clashes
}

func ConvertCourseDayToTimetableDay(courseDays []*pb.CourseDay) []*pb.TimetableDay {
	days := []*pb.TimetableDay{}
	for _, day := range courseDays {
		days = append(days, &pb.TimetableDay{
			Day:     day.Day,
			Timings: ConvertCourseTimingsToTimetableTimings(day.Timings),
		})
	}
	return days
}

func ConvertCourseTimingsToTimetableTimings(courseTimings []*pb.CourseIntTuple) []*pb.TimetableIntTuple {
	timings := []*pb.TimetableIntTuple{}
	for _, timing := range courseTimings {
		timings = append(timings, &pb.TimetableIntTuple{
			Start: timing.Start,
			End:   timing.End,
		})
	}
	return timings
}
