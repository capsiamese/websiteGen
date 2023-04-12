package scanner

import (
	"context"
	"log"
	"os"
	"testing"
)

func TestNewIssueScanner(t *testing.T) {
	ctx := context.Background()
	token := os.Getenv("GITHUB_TOKEN")
	scanner := NewIssueScanner(ctx, token, "capsiamese", "blog")
	res, err := scanner.Scan(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	for _, v := range res {
		t.Log(v)
	}
}
