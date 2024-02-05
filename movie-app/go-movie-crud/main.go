package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Movie struct {
	ID          int
	Title       string
	ReleaseYear int
	Director    string
	Genre       string
}

func AddMovie(db *sql.DB, movie Movie) error {
	query := `INSERT INTO movies (title, release_year, director, genre) VALUES (?, ?, ?, ?)`
	statement, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(movie.Title, movie.ReleaseYear, movie.Director, movie.Genre)
	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Println("Added new movie:", movie.Title)
	return nil
}

func DeleteMovie(db *sql.DB, movieTitle string) error {
	// Prepare the delete SQL statement to ignore case
	stmt, err := db.Prepare("DELETE FROM movies WHERE LOWER(title) = LOWER(?)")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer stmt.Close()

	// Execute the delete statement with the title
	result, err := stmt.Exec(movieTitle)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// Check how many rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
		return err
	}
	if rowsAffected == 0 {
		fmt.Println("No movie found with the given title.")
		return nil
	}

	fmt.Printf("Movie titled \"%s\" was deleted successfully.\n", movieTitle)
	return nil
}

func main() {

	// Open the database connection
	db, err := sql.Open("sqlite3", "movies.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Do you want to add, delete, or edit a movie? (add/delete/edit)")
	var choice string
	fmt.Scanln(&choice)

	switch choice {
	case "add":
		Title, ReleaseYear, Director, Genre := getUserInput()
		newMovie := Movie{
			Title:       Title,
			ReleaseYear: int(ReleaseYear),
			Director:    Director,
			Genre:       Genre,
		}
		if err := AddMovie(db, newMovie); err != nil {
			log.Fatalf("AddMovie() error: %v", err)
		}
	case "delete":
		fmt.Println("Enter the title of the movie you want to delete:")
		var titleToDelete string
		fmt.Scanln(&titleToDelete)
		if err := DeleteMovie(db, titleToDelete); err != nil {
			log.Fatalf("DeleteMovie() error: %v", err)
		}
	case "edit":
		// Placeholder for edit functionality
		fmt.Println("Edit functionality not implemented yet.")
	default:
		fmt.Println("Invalid choice.")
	}
}

func getUserInput() (string, int, string, string) {
	var Title string
	var ReleaseYear int
	var Director string
	var Genre string

	fmt.Println("Enter the title of movie: ")
	fmt.Scan(&Title)

	fmt.Println("Enter the release year of movie: ")
	fmt.Scan(&ReleaseYear)

	fmt.Println("Enter the director name of movie: ")
	fmt.Scan(&Director)

	fmt.Println("Enter the genre of movie: ")
	fmt.Scan(&Genre)

	return Title, ReleaseYear, Director, Genre
}
