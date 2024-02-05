-- init_db.sql

-- Create a movies table
CREATE TABLE IF NOT EXISTS movies (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title VARCHAR(255) NOT NULL,
    release_year INTEGER NOT NULL,
    director VARCHAR(255),
    genre VARCHAR(100)
);

-- Example of how to insert initial data (optional)
-- INSERT INTO movies (title, release_year, director, genre) VALUES ('The Shawshank Redemption', 1994, 'Frank Darabont', 'Drama');
-- INSERT INTO movies (title, release_year, director, genre) VALUES ('The Godfather', 1972, 'Francis Ford Coppola', 'Crime, Drama');
