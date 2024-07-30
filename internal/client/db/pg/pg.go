package pg

import (
	"context"
	"fmt"
	"log"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mikhailsoldatkin/chat-server/internal/client/db"
	"github.com/mikhailsoldatkin/chat-server/internal/client/db/prettier"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type key string

// TxKey is used as a key to store and retrieve a database transaction from a context.
const TxKey key = "tx"

type pg struct {
	dbc *pgxpool.Pool
}

// ScanOneContext executes a query that returns a single row and scans the result into the destination.
func (p *pg) ScanOneContext(ctx context.Context, dest any, q db.Query, args ...any) error {
	logQuery(ctx, q, args...)

	row, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanOne(dest, row)
}

// ScanAllContext executes a query that returns multiple rows and scans all results into the destination.
func (p *pg) ScanAllContext(ctx context.Context, dest any, q db.Query, args ...any) error {
	logQuery(ctx, q, args...)

	rows, err := p.QueryContext(ctx, q, args...)
	if err != nil {
		return err
	}

	return pgxscan.ScanAll(dest, rows)
}

// RecordExists checks if a record with a given ID exists in a database table.
func (p *pg) RecordExists(ctx context.Context, id int64, table string) error {
	var exists bool
	notFoundMsg := fmt.Sprintf("record with ID %d not found in %s table", id, table)
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id=%d)", table, id)
	q := db.Query{
		Name:     "RecordExists",
		QueryRaw: query,
	}
	err := p.ScanOneContext(ctx, &exists, q)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to check existence: %v", err)
	}
	if !exists {
		return status.Errorf(codes.NotFound, notFoundMsg)
	}
	return nil
}

// ExecContext executes a query without returning any rows and returns the command tag and error if any.
func (p *pg) ExecContext(ctx context.Context, q db.Query, args ...any) (pgconn.CommandTag, error) {
	logQuery(ctx, q, args...)

	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Exec(ctx, q.QueryRaw, args...)
	}

	return p.dbc.Exec(ctx, q.QueryRaw, args...)
}

// QueryContext executes a query that returns multiple rows and returns the rows and error if any.
func (p *pg) QueryContext(ctx context.Context, q db.Query, args ...any) (pgx.Rows, error) {
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Query(ctx, q.QueryRaw, args...)
	}

	return p.dbc.Query(ctx, q.QueryRaw, args...)
}

// QueryRowContext executes a query that returns a single row and returns the row.
func (p *pg) QueryRowContext(ctx context.Context, q db.Query, args ...any) pgx.Row {
	logQuery(ctx, q, args...)

	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.QueryRow(ctx, q.QueryRaw, args...)
	}

	return p.dbc.QueryRow(ctx, q.QueryRaw, args...)
}

// BeginTx starts a new transaction with the given options and returns the transaction object and error if any.
func (p *pg) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return p.dbc.BeginTx(ctx, txOptions)
}

// Ping checks if the database connection is alive.
func (p *pg) Ping(ctx context.Context) error {
	return p.dbc.Ping(ctx)
}

// Close closes the database connection pool.
func (p *pg) Close() {
	p.dbc.Close()
}

// MakeContextTx returns a new context with the provided transaction object stored in it.
func MakeContextTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, TxKey, tx)
}

// logQuery logs the query and its arguments, including a pretty-printed version of the query.
func logQuery(_ context.Context, q db.Query, args ...any) {
	prettyQuery := prettier.Pretty(q.QueryRaw, prettier.PlaceholderDollar, args...)
	log.Printf("query: %s sql: %s\n", q.Name, prettyQuery)
}
