package main

import (
	"context"
	"log"
	"net"

	"github.com/thargro5/cloudnativecourse/lab5/movieapi"
	"google.golang.org/grpc"
)

const (
	port = ":50050"
)

// server is used to implement movieapi.MovieInfoServer.
type server struct {
	movieData map[string]*movieapi.MovieReply
}

// GetMovieInfo implements movieapi.MovieInfoServer
func (s *server) GetMovieInfo(ctx context.Context, in *movieapi.MovieRequest) (*movieapi.MovieReply, error) {
	log.Printf("Received request for movie info: %v", in)
	movie, ok := s.movieData[in.Title]
	if !ok {
		return nil, grpc.Errorf(grpc.Code(grpc.NotFound), "Movie not found: %s", in.Title)
	}
	return movie, nil
}

// SetMovieInfo implements movieapi.MovieInfoServer
func (s *server) SetMovieInfo(ctx context.Context, in *movieapi.MovieData) (*movieapi.Status, error) {
	log.Printf("Received request to set movie info: %v", in)
	// Store the movie data
	s.movieData[in.Title] = &movieapi.MovieReply{
		Year:     in.Year,
		Director: in.Director,
		Cast:     in.Cast,
	}
	return &movieapi.Status{Code: "Success"}, nil
}

func newServer() *server {
	return &server{
		movieData: make(map[string]*movieapi.MovieReply),
	}
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	movieapi.RegisterMovieInfoServer(s, newServer())
	log.Printf("Server listening on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
