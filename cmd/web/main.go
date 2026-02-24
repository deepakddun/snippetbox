package main

import (
	"context"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"github.com/deepakddun/snippetbox/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type application struct {
	logger        *slog.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	//

	// Default level is log
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	db, dbError := openDB()

	if dbError != nil {
		logger.Error(dbError.Error())

		os.Exit(1)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()

	if err != nil {
		logger.Error(dbError.Error())

		os.Exit(1)
	}

	app := &application{
		logger:        logger,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	addr := flag.String("addr", ":4000", " HTTP server network address")

	flag.Parse()
	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.

	logger.Info(" Starting the server at", "addr", *addr)
	err = http.ListenAndServe(*addr, app.routes())
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

}

func openDB() (*pgxpool.Pool, error) {

	dbPool, err := pgxpool.New(context.Background(), "postgres://snippet:snippet_dev_password@localhost:5432/snippet?sslmode=disable")

	if err != nil {
		return nil, err
	}

	err = dbPool.Ping(context.Background())

	if err != nil {
		dbPool.Close()
		return nil, err
	}

	return dbPool, nil

}
