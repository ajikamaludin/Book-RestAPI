package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
)

var (
	baseURL string
	reader  = bufio.NewReader(os.Stdin)
)

type Book struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	PublishedAt time.Time `json:"published_at"`
	Edition     string    `json:"edition"`
	Description string    `json:"description"`
	Genre       string    `json:"genre"`
}

func main() {
	flag.StringVar(&baseURL, "url", "http://localhost:3000", "Base URL of the REST API")
	flag.Parse()

	for {
		showMenu()
		choice := readInput("Enter your choice: ")

		switch choice {
		case "1":
			filterBooks()
		case "2":
			addBook()
		case "3":
			updateBook()
		case "4":
			deleteBook()
		case "5":
			fmt.Println("Exiting...")
			os.Exit(0)
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

func showMenu() {
	fmt.Println("----- Book Menu -----")
	fmt.Println("1. List Books")
	fmt.Println("2. Add a Book")
	fmt.Println("3. Update a Book")
	fmt.Println("4. Delete a Book")
	fmt.Println("5. Exit")
	fmt.Println("----------------------")
}

func filterBooks() {
	var queryParams []string

	filterByAuthor := readYesNoInput("Filter by author? (y/N): ")
	if filterByAuthor {
		author := readInput("Enter the author: ")
		queryParams = append(queryParams, "author="+author)
	}

	filterByGenre := readYesNoInput("Filter by genre? (y/N): ")
	if filterByGenre {
		genre := readInput("Enter the genre: ")
		queryParams = append(queryParams, "genre="+genre)
	}

	filterByDate := readYesNoInput("Filter by date range? (y/N): ")
	if filterByDate {
		start := readInput("Enter the start date (YYYY-MM-DD): ")
		end := readInput("Enter the end date (YYYY-MM-DD): ")
		queryParams = append(queryParams, "start="+start)
		queryParams = append(queryParams, "end="+end)
	}

	queryParamsStr := strings.Join(queryParams, "&")
	getBooks(queryParamsStr)
}

func getBooks(queryParams string) {
	url := fmt.Sprintf("%s/books", baseURL)
	if queryParams != "" {
		url = fmt.Sprintf("%s?%s", url, queryParams)
	}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Failed to get books:", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		var books []Book
		err = json.NewDecoder(resp.Body).Decode(&books)
		if err != nil {
			fmt.Println("Failed to decode response body:", err.Error())
			return
		}

		printBookTable(books)
	} else {
		fmt.Println("Failed to get books:", resp.Status)
	}
}

func printBookTable(books []Book) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"ID", "Title", "Author", "Published At", "Edition", "Description", "Genre"})

	for _, book := range books {
		t.AppendRow([]interface{}{
			book.ID,
			book.Title,
			book.Author,
			book.PublishedAt.Format("2006-01-02"),
			book.Edition,
			book.Description,
			book.Genre,
		})
	}

	t.Render()
}

func addBook() {
	title := readInput("Enter the title: ")
	author := readInput("Enter the author: ")
	publishedAt := readInput("Enter the published date (YYYY-MM-DD): ")
	edition := readInput("Enter the edition: ")
	description := readInput("Enter the description: ")
	genre := readInput("Enter the genre: ")

	parsedPublishedAt, err := time.Parse("2006-01-02", publishedAt)
	if err != nil {
		fmt.Println("Invalid published date")
		return
	}

	book := Book{
		Title:       title,
		Author:      author,
		PublishedAt: parsedPublishedAt,
		Edition:     edition,
		Description: description,
		Genre:       genre,
	}

	body, err := json.Marshal(book)
	if err != nil {
		fmt.Println("Failed to serialize book object")
		return
	}

	resp, err := http.Post(fmt.Sprintf("%s/books", baseURL), "application/json", strings.NewReader(string(body)))
	if err != nil {
		fmt.Println("Failed to add book:", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		fmt.Println("Book added successfully")
	} else {
		fmt.Println("Failed to add book:", resp.Status)
	}
}

func updateBook() {
	idStr := readInput("Enter the ID of the book to update: ")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Invalid book ID")
		return
	}

	title := readInput("Enter the updated title: ")
	author := readInput("Enter the updated author: ")
	publishedAt := readInput("Enter the updated published date (YYYY-MM-DD): ")
	edition := readInput("Enter the updated edition: ")
	description := readInput("Enter the updated description: ")
	genre := readInput("Enter the updated genre: ")

	parsedPublishedAt, err := time.Parse("2006-01-02", publishedAt)
	if err != nil {
		fmt.Println("Invalid published date")
		return
	}

	book := Book{
		ID:          uint(id),
		Title:       title,
		Author:      author,
		PublishedAt: parsedPublishedAt,
		Edition:     edition,
		Description: description,
		Genre:       genre,
	}

	body, err := json.Marshal(book)
	if err != nil {
		fmt.Println("Failed to serialize book object")
		return
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/books/%d", baseURL, id), strings.NewReader(string(body)))
	if err != nil {
		fmt.Println("Failed to create request:", err.Error())
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Failed to update book:", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		fmt.Println("Book updated successfully")
	} else {
		fmt.Println("Failed to update book:", resp.Status)
	}
}

func deleteBook() {
	idStr := readInput("Enter the ID of the book to delete: ")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Invalid book ID")
		return
	}

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/books/%d", baseURL, id), nil)
	if err != nil {
		fmt.Println("Failed to create request:", err.Error())
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Failed to delete book:", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		fmt.Println("Book deleted successfully")
	} else {
		fmt.Println("Failed to delete book:", resp.Status)
	}
}

func readInput(prompt string) string {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func readYesNoInput(prompt string) bool {
	input := strings.ToLower(readInput(prompt))
	return input == "y" || input == "yes"
}
