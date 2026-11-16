package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/deblinux123/semantic-scholar-go/semanticscholar"
)

func main() {
	// apiKey := os.Getenv("SEMANTIC_SCHOLAR_API_KEY")

	client := semanticscholar.NewClient(
		// semanticscholar.WithAPIKey(apiKey),
		semanticscholar.WithTimeout(40 * time.Second),
	)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		45*time.Second,
	)

	defer cancel()

	query := "retrival augmented generation"

	result, err := client.SearchPapers(
		ctx,
		query,
		5,
		0,
	)

	if err != nil {
		log.Fatalf("paper search failed: %v", err)
	}

	fmt.Printf("Total results: %d\n", result.Total)
	fmt.Printf("Returned papers: %d\n\n", len(result.Data))

	for i, paper := range result.Data {
		fmt.Printf("Paper #%d\n", i+1)
		fmt.Printf("ID: %s\n", paper.PaperID)
		fmt.Printf("Title: %s\n", paper.Title)
		fmt.Printf("Year: %d\n", paper.Year)
		fmt.Printf("Citations: %d\n", paper.CitationCount)
		fmt.Printf("URL: %s\n", paper.URL)

		if len(paper.Authors) > 0 {
			fmt.Print("Authors: ")

			for j, author := range paper.Authors {
				if j > 0 {
					fmt.Print(", ")
				}
				fmt.Print(author.Name)
			}

			fmt.Println()
		}

		fmt.Println("----------------------------------------")
	}
}
