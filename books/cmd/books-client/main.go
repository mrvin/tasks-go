package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/mrvin/tasks-go/books/internal/booksapi"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
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
		fmt.Printf("0 - Exit\n1 - Save book\n2 - Search by title\n3 - Search by author\n4 - List of books\n")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		switch []byte(input)[0] {
		case '0':
			fmt.Printf("Exit\n")
			break exit
		case '1':
			fmt.Printf("Title:")
			title, _ := reader.ReadString('\n')
			title = strings.TrimSuffix(title, "\n")
			fmt.Printf("Authors:")
			authors, _ := reader.ReadString('\n')

			req := &booksapi.Book{
				Title:   title,
				Authors: strings.Split(strings.TrimSuffix(authors, "\n"), ", "),
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if _, err := client.CreateBook(ctx, req); err != nil {
				log.Fatalf("Save book: %v", err)
			}
		case '2':
			fmt.Printf("Title:")
			title, _ := reader.ReadString('\n')
			title = strings.TrimSuffix(title, "\n")

			req := &booksapi.Title{
				Title: title,
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			res, err := client.GetBookByTitle(ctx, req)
			if err != nil {
				log.Fatalf("Search by title: %v", err)
			}

			printListBooks("Book with title "+title, []*booksapi.Book{res})
		case '3':
			fmt.Printf("Author:")
			author, _ := reader.ReadString('\n')
			author = strings.TrimSuffix(author, "\n")

			req := &booksapi.Author{
				Author: author,
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			res, err := client.ListBooksByAuthor(ctx, req)
			if err != nil {
				log.Fatalf("Search by author: %v", err)
			}

			printListBooks("List of books by author "+author, res.Books)
		case '4':
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			res, err := client.ListAllBooks(ctx, &emptypb.Empty{})
			if err != nil {
				log.Fatalf("List of books: %v", err)
			}

			printListBooks("List of books", res.Books)
		default:
		}
	}
}

func printListBooks(header string, slBooks []*booksapi.Book) {
	fmt.Println(header)
	const formatHeader = "%s\t%s\t%s\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, formatHeader, "Number", "Title", "Authors")
	fmt.Fprintf(tw, formatHeader, "------", "-----", "------")

	for i, book := range slBooks {
		fmt.Fprintf(tw, "%d\t%s\t%s\n", i+1, book.Title, book.Authors)
	}

	tw.Flush()
}
