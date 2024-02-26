// Package main implements a client for movieinfo service
package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/thargro5/cloudnativecourse/lab5/movieapi"
	"google.golang.org/grpc"
)

const (
	address      = "localhost:50050"
	defaultTitle = "Pulp Fiction"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := movieapi.NewMovieInfoClient(conn)

	// Contact the server and print out its response.
	title := defaultTitle
	if len(os.Args) > 1 {
		title = os.Args[1]
	}
	// Timeout if server doesn't respond
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Call SetMovieInfo to set movie data
	_, err = c.SetMovieInfo(ctx, &movieapi.MovieData{
		Title:    "The Shawshank Redemption",
		Year:     1994,
		Director: "Frank Darabont",
		Cast:     []string{"Tim Robbins", "Morgan Freeman", "Bob Gunton"},
	})
	if err != nil {
		log.Fatalf("could not set movie info: %v", err)
	}
	log.Println("Movie info set successfully.")

	// Call GetMovieInfo to retrieve movie data
	r, err := c.GetMovieInfo(ctx, &movieapi.MovieRequest{Title: title})
	if err != nil {
		log.Fatalf("could not get movie info: %v", err)
	}
	log.Printf("Movie Info for %s %d %s %v", title, r.GetYear(), r.GetDirector(), r.GetCast())
}
