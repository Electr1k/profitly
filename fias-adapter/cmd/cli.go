package main

import (
	"fias-adapter/tmp/parser"
	"fias-adapter/tmp/postgres"
	"flag"
	"log"
	"path/filepath"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}

	archivePath := flag.String("archive", "", "Path to FIAS ZIP archive")
	flag.Parse()

	if *archivePath == "" {
		log.Fatal("Please provide -archive flag with path to the FIAS ZIP file")
	}

	absPath, err := filepath.Abs(*archivePath)
	if err != nil {
		log.Fatalf("Invalid archive path: %v", err)
	}

	db, err := postgres.Connect()
	if err != nil {
		log.Fatalf("DB connect: %v", err)
	}
	defer db.Close()

	if err := postgres.Migrate(db); err != nil {
		log.Fatalf("Migration: %v", err)
	}

	log.Printf("Starting parse of %s", absPath)
	if err := parser.ParseArchive(db, absPath); err != nil {
		log.Fatalf("Parse error: %v", err)
	}

	log.Println("Parse completed successfully")
}
