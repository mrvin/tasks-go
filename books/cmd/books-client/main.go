package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mrvin/tasks-go/books/internal/booksapi"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:55555", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("books-client: %v", err)
	}
	defer conn.Close()
	client := booksapi.NewBookServiceClient(conn)
exit:
	for {
		fmt.Printf("0 - Exit\n1 - Save book\n2 - Search by author\n3 - Search by title\n")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		switch []byte(input)[0] {
		case '0':
			fmt.Printf("Exit\n")
			break exit
		case '1':
			fmt.Printf("Title:")
			title, _ := reader.ReadString('\n')
			fmt.Printf("Authors:")
			authors, _ := reader.ReadString('\n')

			req := &booksapi.CreateBookRequest{
				Title:   strings.TrimSuffix(title, "\n"),
				Authors: strings.Split(strings.TrimSuffix(authors, "\n"), ", "),
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if _, err := client.CreateBook(ctx, req); err != nil {
				log.Fatalf("Save book: %v", err)
			}
		case '2':
			fmt.Printf("Author:")
			author, _ := reader.ReadString('\n')

			req := &booksapi.GetBooksByAuthorRequest{
				Author: strings.TrimSuffix(author, "\n"),
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			res, err := client.GetBooksByAuthor(ctx, req)
			if err != nil {
				log.Fatalf("Search by author: %v", err)
			}

			fmt.Printf("Titles: %v\n", res.Titles)
		case '3':
			fmt.Printf("Title:")
			title, _ := reader.ReadString('\n')

			req := &booksapi.GetAuthorsByTitleRequest{
				Title: strings.TrimSuffix(title, "\n"),
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			res, err := client.GetAuthorsByTitle(ctx, req)
			if err != nil {
				log.Fatalf("Search by title: %v", err)
			}

			fmt.Printf("Authors: %v\n", res.Authors)
		default:
		}
	}
}
