package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

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

func promptForUpdate(prompt string, currentValue string) string {
	fmt.Printf("%s [%s]: ", prompt, currentValue)
	var input string
	fmt.Scanln(&input)
	if input == "" {
		return currentValue
	}
	return input
}

func EditMovie(db *sql.DB) {
	fmt.Println("Enter the title of the movie you want to edit:")
	var searchTitle string
	fmt.Scanln(&searchTitle)

	// First, find the movie by title to get its current details
	movie, err := FindMovieByTitle(db, searchTitle)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No movie found with the given title.")
		} else {
			log.Fatal("Error querying movie:", err)
		}
		return
	}

	// Prompt for new values, allowing Enter to keep existing values
	newTitle := promptForUpdate("New title (press Enter to keep current)", movie.Title)
	newDirector := promptForUpdate("New director (press Enter to keep current)", movie.Director)
	newGenre := promptForUpdate("New genre (press Enter to keep current)", movie.Genre)
	fmt.Println("Enter new release year (press Enter to keep current):")
	var newReleaseYearStr string
	fmt.Scanln(&newReleaseYearStr)
	newReleaseYear := movie.ReleaseYear // Initialize with current year in case of no input
	if newReleaseYearStr != "" {
		newReleaseYear, err = strconv.Atoi(newReleaseYearStr)
		if err != nil {
			fmt.Println("Invalid input for release year. Keeping the current value.")
		}
	}

	// Update the movie in the database
	if err := UpdateMovie(db, movie.ID, newTitle, newReleaseYear, newDirector, newGenre); err != nil {
		log.Fatal("Error updating movie:", err)
	}
	fmt.Println("Movie updated successfully.")
}

func FindMovieByTitle(db *sql.DB, title string) (Movie, error) {
	var movie Movie
	query := `SELECT id, title, release_year, director, genre FROM movies WHERE LOWER(title) = LOWER(?)`
	err := db.QueryRow(query, title).Scan(&movie.ID, &movie.Title, &movie.ReleaseYear, &movie.Director, &movie.Genre)
	if err != nil {
		return Movie{}, err
	}
	return movie, nil
}

func UpdateMovie(db *sql.DB, id int, title string, releaseYear int, director string, genre string) error {
	query := `UPDATE movies SET title = ?, release_year = ?, director = ?, genre = ? WHERE id = ?`
	_, err := db.Exec(query, title, releaseYear, director, genre, id)
	return err
}

func main() {
	db, err := sql.Open("sqlite3", "./movies.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for {
		fmt.Println("Do you want to add, delete, edit a movie, or exit? (add/delete/edit/exit)")
		var choice string
		fmt.Scanln(&choice)

		// Convert the choice to lowercase to make the command case-insensitive
		choice = strings.ToLower(choice)

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
			EditMovie(db)
			fmt.Println("Update successfull.")
		case "exit":
			fmt.Println("Exiting application...")
			return
		default:
			fmt.Println("Invalid choice. Please type add, delete, edit, or exit.")
		}

		fmt.Println() // Print an empty line for better readability between operations
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
