package core

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

func Dbinit() (*sql.DB, error) {
	connStr := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	// Ping database to check connection
	if err = db.Ping(); err != nil {
		db.Close()
		log.Fatalf("Error pinging database: %v", err)
	}

	fmt.Println("Connected to database successfully")

	// Call the createTable function
	err = createTable(db)
	if err != nil {
		return db, fmt.Errorf("error creating table: %v", err)
	}
	return db, nil
}
func createTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS blocks (
		id SERIAL PRIMARY KEY,
		data TEXT NOT NULL,
		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		prev_hash TEXT NOT NULL,
		hash TEXT NOT NULL UNIQUE
	);`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating table: %v", err)
	} else {
		fmt.Println("Table created (or already exists)")
	}
	return nil
}
func getData(db *sql.DB) error {
	query := `SELECT * FROM blocks`
	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("error fetching data: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var data, prevHash, hash, timestamp string

		if err := rows.Scan(&id, &data, &timestamp, &prevHash, &hash); err != nil {
			return fmt.Errorf("error scanning row: %w", err)
		}
		fmt.Printf("ID: %d, Data: %s, Timestamp: %s, Prev Hash: %s, Hash: %s\n",
			id, data, timestamp, prevHash, hash)
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("error iterating over rows: %w", err)
	}

	return nil
}

func saveBlockToDB(db *sql.DB, data, timestamp, prevHash, newHash string) error {
	unixTimestamp, err := strconv.ParseInt(timestamp, 10, 64)

	formattedTimestamp := time.Unix(unixTimestamp, 0).Format("2006-01-02 15:04:05")

	fmt.Println("🛢️ Saving to Database:")
	fmt.Println("   Data:", data)
	fmt.Println("   Timestamp:", formattedTimestamp)
	fmt.Println("   Previous Hash (hex):", prevHash)
	fmt.Println("   Current Block Hash (hex):", newHash)

	_, err = db.Exec("INSERT INTO blocks(data,timestamp,prev_hash,hash) VALUES ($1,$2,$3,$4)", data, formattedTimestamp, prevHash, newHash)

	if err != nil {
		fmt.Println(" SQL Insert Failed:", err)
		return fmt.Errorf(" Failed to save block: %w", err)
	}
	fmt.Println("Data added successfully to database.")
	return nil
}

func alterTable(db *sql.DB) error {
	_, err := db.Exec("ALTER TABLE blocks ALTER COLUMN timestamp TYPE BIGINT USING EXTRACT(EPOCH FROM timestamp)::BIGINT;")
	if err != nil {
		return fmt.Errorf("error altering table: %w", err)
	}
	fmt.Println("Table altered successfully!")
	return nil
}
