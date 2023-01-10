package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "app/calculatorpb"
)

type server struct {
	*pb.UnimplementedCalculatorServiseServer
}

func (s *server) SquareRoot(ctx context.Context, req *pb.SquareRootRequest) (*pb.SquareRootResponse, error) {

	var sr float64 = float64(req.Number) / 2
	var temp float64
	for {
		temp = sr
		sr = (temp + (float64(req.Number) / temp)) / 2
		if (temp - sr) == 0 {
			break
		}
	}
	return &pb.SquareRootResponse{
		SqrNumber: float64(sr),
	}, nil
}

func (s *server) PerfectNumber(req *pb.PerfectNumberRequest, stream pb.CalculatorServise_PerfectNumberServer) error {

	fmt.Println("Req:", req)

	number := req.GetNumber()
	nums := int64(1)
	for nums <= number {
		if findPerfectNumber(nums) {
			stream.Send(&pb.PerfectNumberResponse{
				PerfectNumber: nums,
			})
			nums++
		} else {
			nums++
		}
	}

	return nil
}

func findPerfectNumber(n int64) bool {
	var total int64
	for i := int64(1); i < n; i++ {
		if n%i == 0 {
			total += i
		}
	}

	if total == n {
		return true
	} else {
		return false
	}
}

func (s *server) TotalNumber(stream pb.CalculatorServise_TotalNumberServer) error {

	var (
		total float64
	)

	for {

		req, err := stream.Recv()

		if err == io.EOF {

			err = stream.SendAndClose(&pb.TotalNumberResponse{
				TotalNumber: total,
			})

			if err != nil {
				log.Println("error while TotalNumber Recv:", err)
				return err
			}

			return nil
		}

		if err != nil {
			log.Println("error while TotalNumber Recv:", err)
			return err
		}

		total += req.Number
	}
}

func (s *server) FindMinimum(stream pb.CalculatorServise_FindMinimumServer) error {
	boolean := true
	minimum := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if boolean {
			minimum = req.Number
			boolean = false
			err = stream.Send(&pb.FindMinimumResponse{Minimum: minimum})
			if err != nil {
				log.Println("error while FindMinimum Send:", err)
				return err
			}
		}

		if err != nil {
			log.Println("error while FindMinimum Recv:", err)
			return err
		}

		fmt.Println(req)

		if req.Number < minimum {
			minimum = req.Number
			err = stream.Send(&pb.FindMinimumResponse{Minimum: minimum})
			if err != nil {
				log.Println("error while FindMinimum Send:", err)
				return err
			}
		}
	}
}

func main() {

	lis, err := net.Listen("tcp", ":9002")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	pb.RegisterCalculatorServiseServer(s, &server{})

	fmt.Println("Listening :9002...")
	if err = s.Serve(lis); err != nil {
		panic(err)
	}
}
