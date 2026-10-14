package main

import "testing"

func TestDatabaseURLIsRequired(t *testing.T) {
	t.Setenv("DATABASE_URL", " ")
	if _, err := databaseURLFromEnv(); err == nil {
		t.Fatal("expected an error when DATABASE_URL is empty")
	}
}

func TestDatabaseURLIsTrimmed(t *testing.T) {
	t.Setenv("DATABASE_URL", " postgres://localhost/pos ")
	value, err := databaseURLFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if value != "postgres://localhost/pos" {
		t.Fatalf("expected trimmed URL, got %q", value)
	}
}
