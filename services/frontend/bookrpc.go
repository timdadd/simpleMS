package main

import (
	"context"

	pb "frontend/pb/pb_book_v1"
)

// Lists books. The order is unspecified but deterministic. Newly created
// books will not necessarily appear at the end of this list.
func (fe *frontendServer) ListBooks(ctx context.Context) ([]*pb.Book, error) {
	req := pb.ListBooksRequest{}
	resp, err := pb.NewBookServiceClient(fe.bookSvcConn).ListBooks(ctx, &req)
	return resp.Books, err
}

// Creates a book, and returns the new Book.
func (fe *frontendServer) AddBook(ctx context.Context, b *pb.Book) (id string, err error) {
	req := pb.CreateBookRequest{Book: b}
	resp, err := pb.NewBookServiceClient(fe.bookSvcConn).CreateBook(ctx, &req)
	return resp.Id, err
}

func (fe *frontendServer) GetBook(ctx context.Context, id string) (*pb.Book, error) {
	req := pb.GetBookRequest{Id: id}
	resp, err := pb.NewBookServiceClient(fe.bookSvcConn).GetBook(ctx, &req)
	return resp, err
}

// DeleteBook removes a given book by its ID.
func (fe *frontendServer) DeleteBook(ctx context.Context, id string) error {
	req := pb.DeleteBookRequest{Id: id}
	_, err := pb.NewBookServiceClient(fe.bookSvcConn).DeleteBook(ctx, &req)
	return err
}

// UpdateBook updates the entry for a given book.
func (fe *frontendServer) UpdateBook(ctx context.Context, b *pb.Book) (*pb.Book, error) {
	req := pb.UpdateBookRequest{Book: b}
	resp, err := pb.NewBookServiceClient(fe.bookSvcConn).UpdateBook(ctx, &req)
	return resp, err
}
