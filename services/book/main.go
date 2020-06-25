package main

import (
	"book/dao"
	pb "book/pb/pb_book_v1"
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
	"lib/common"
	"net"
	"net/http"
)

const (
	serviceName string = "book" // Make sure same cfg name in Dockerfile
	version     string = "1.0.0"
)

type bookServer struct {
	pb.UnimplementedBookServiceServer
	DB    dao.BookDatabase
	log   *logrus.Logger
	empty empty.Empty
}

// getBook retrieves a book from the database given a book ID
var ErrNoIdForBook = errors.New("All books have an ID")

// Gets a book. Returns NOT_FOUND if the book does not exist.
func (b *bookServer) GetBook(ctx context.Context, req *pb.GetBookRequest) (*pb.Book, error) {
	if req.Id == "" {
		b.log.Error(ErrNoIdForBook)
		return nil, ErrNoIdForBook
	}
	book, err := b.DB.GetBook(ctx, req.Id)
	if err != nil {
		b.log.Errorf("could not find book: %v", err)
		return nil, fmt.Errorf("could not find book: %w", err)
	}
	return book, nil
}

// Lists books. The order is unspecified but deterministic. Newly created
// books will not necessarily appear at the end of this list.
// ListBooks lists all books
func (b *bookServer) ListBooks(ctx context.Context, req *pb.ListBooksRequest) (*pb.ListBooksResponse, error) {
	books, err := b.DB.ListBooks(ctx)
	if err != nil {
		b.log.Errorf("could not list books: %v:%v", req, err)
		return nil, fmt.Errorf("could not list books: %v:%w", req, err)
	}
	return &pb.ListBooksResponse{Books: books}, nil
}

// Creates a book, and returns the new Book.
func (b *bookServer) CreateBook(ctx context.Context, req *pb.CreateBookRequest) (*pb.Book, error) {
	id, err := b.DB.AddBook(ctx, req.GetBook())
	if err != nil {
		b.log.Errorf("could not save book: %v : %v", req.GetBook(), err)
		return nil, fmt.Errorf("could not save book: %v : %w", req.GetBook(), err)
	}
	return b.GetBook(ctx, &pb.GetBookRequest{Id: id})
}

// Deletes a book. Returns NOT_FOUND if the book does not exist.
func (b *bookServer) DeleteBook(ctx context.Context, req *pb.DeleteBookRequest) (*empty.Empty, error) {
	if err := b.DB.DeleteBook(ctx, req.Id); err != nil {
		b.log.Errorf("could not delete book: %s : %v", req.Id, err)
		return nil, fmt.Errorf("could not deete book: %s : %w", req.Id, err)
	}
	return &b.empty, nil
}

// Updates a book. Returns INVALID_ARGUMENT if the name of the book
// is non-empty and does not equal the existing name.
func (b *bookServer) UpdateBook(ctx context.Context, req *pb.UpdateBookRequest) (*pb.Book, error) {
	if err := b.DB.UpdateBook(ctx, req.GetBook()); err != nil {
		b.log.Errorf("could not update book: %v : %v", req.GetBook(), err)
		return nil, fmt.Errorf("could not update book: %v : %w", req.GetBook(), err)
	}
	return req.GetBook(), nil
}

func serialize(book *pb.Book) string {
	return fmt.Sprintf("%s %s", book.Id, book.Title)
}

func newServer(db dao.BookDatabase, log *logrus.Logger) *bookServer {
	b := &bookServer{DB: db, log: log}
	return b
}

func main() {
	// Command line stuff
	showversion := flag.Bool("version", false, "display version")
	flag.Parse()
	if *showversion {
		fmt.Printf("Version %s\n", version)
		return
	}

	// Service details
	c, err := common.LoadConfig(serviceName, "")
	if err != nil {
		fmt.Printf("Cannot load the configuration: %s", err)
	}

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s\n", version)
	})

	c.KeyPrefix("book")
	c.Log.Info(common.GetLocalIP(), ":", c.Port())

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", c.Port()))
	if err != nil {
		c.Log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	if c.TLS() {
		certFile := c.CertFile()
		if certFile == "" {
			certFile = testdata.Path("server1.pem")
		}
		keyFile := c.KeyFile()
		if keyFile == "" {
			keyFile = testdata.Path("server1.key")
		}
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			c.Log.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)
	db, err := dao.NewMemoryDB()
	if err != nil {
		c.Log.Fatalf("newFirestoreDB: %v", err)
	}
	pb.RegisterBookServiceServer(grpcServer, newServer(db, c.Log))
	grpcServer.Serve(lis)
}
