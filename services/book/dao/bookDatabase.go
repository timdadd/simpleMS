package dao

import (
	pb "book/pb/pb_book_v1"
	"context"
)

// BookDatabase provides thread-safe access to a database of books.
type BookDatabase interface {
	// ListBooks returns a list of books, ordered by title.
	ListBooks(context.Context) ([]*pb.Book, error)

	// GetBook retrieves a book by its ID.
	GetBook(ctx context.Context, id string) (*pb.Book, error)

	// AddBook saves a given book, assigning it a new ID.
	AddBook(ctx context.Context, book *pb.Book) (id string, err error)

	// DeleteBook removes a given book by its ID.
	DeleteBook(ctx context.Context, id string) error

	// UpdateBook updates the entry for a given book.
	UpdateBook(ctx context.Context, book *pb.Book) error
}
