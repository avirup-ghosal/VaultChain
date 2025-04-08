package main

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

// Function to calculate the hash based on `SetHash()`
func calculateHash(timestamp int64, prevHash, data string) string {
	timeBytes := []byte(strconv.FormatInt(timestamp, 10))
	dataBytes := []byte(data)

	//  Convert prevHash to Hex before hashing (Match `SetHash()`)
	prevHashHex := hex.EncodeToString([]byte(prevHash))
	prevHashBytes := []byte(prevHashHex)

	headers := bytes.Join([][]byte{timeBytes, prevHashBytes, dataBytes}, []byte{})
	hash := sha256.Sum256(headers)
	return hex.EncodeToString(hash[:])
}

// Function to read and compute hashes from logs.txt
func readAndComputeHashes(filename string) ([]map[string]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := stat.Size()
	buffer := make([]byte, fileSize)

	_, err = file.Read(buffer)
	if err != nil {
		return nil, err
	}

	// Convert buffer to string and split by lines
	lines := strings.Split(strings.TrimSpace(string(buffer)), "\n")
	if len(lines)%4 != 0 { // 4 lines per block: data, timestamp, prev_hash, current_hash
		return nil, fmt.Errorf("invalid log file format")
	}

	var blocks []map[string]string

	for i := 0; i < len(lines); i += 4 {
		data := lines[i]
		timestamp := lines[i+1]
		prevHash := lines[i+2]
		expectedHash := lines[i+3]

		// Convert timestamp string to int64
		timestampInt, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid timestamp format in file")
		}

		// Compute hash using updated logic
		computedHash := calculateHash(timestampInt, prevHash, data)

		blocks = append(blocks, map[string]string{
			"data":          data,
			"timestamp":     timestamp,
			"prev_hash":     prevHash,
			"expected_hash": expectedHash,
			"computed_hash": computedHash,
		})
	}

	return blocks, nil
}

func fetchHashesFromDB(db *sql.DB) (map[string]map[string]string, error) {
	rows, err := db.Query("SELECT prev_hash, hash, EXTRACT(EPOCH FROM timestamp AT TIME ZONE 'UTC') FROM blocks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dbHashes := make(map[string]map[string]string)

	fmt.Println(" Fetching Hashes from PostgreSQL:")
	for rows.Next() {
		var prevHash, hash string
		var timestamp float64 // Store as float for precise conversion

		if err := rows.Scan(&prevHash, &hash, &timestamp); err != nil {
			return nil, err
		}

		// Convert to int64 and format correctly
		unixTimestamp := strconv.FormatInt(int64(timestamp-19800), 10)

		dbHashes[hash] = map[string]string{
			"prev_hash": prevHash,
			"timestamp": unixTimestamp,
		}

		fmt.Println("   DB Hash:", hash)
		fmt.Println("   DB Prev Hash:", prevHash)
		fmt.Println("   DB Timestamp (Unix):", unixTimestamp)
	}

	return dbHashes, nil
}

// Function to validate blockchain integrity
func validateBlockchain(filename string, db *sql.DB) error {
	// Get hashes from logs.txt
	fileHashes, err := readAndComputeHashes(filename)
	if err != nil {
		return fmt.Errorf(" Error reading file: %v", err)
	}

	// Get hashes from PostgreSQL
	dbHashes, err := fetchHashesFromDB(db)
	if err != nil {
		return fmt.Errorf(" Error fetching database hashes: %v", err)
	}

	// Compare computed hashes with stored database hashes
	for _, block := range fileHashes {
		expectedHash := block["expected_hash"]
		computedHash := block["computed_hash"]
		prevHash := block["prev_hash"]
		fileTimestamp := block["timestamp"] // Unix format from logs.txt

		// Fetch corresponding timestamp from PostgreSQL
		dbData, exists := dbHashes[computedHash]
		if !exists {
			fmt.Println(" Block not found in PostgreSQL!")
			fmt.Println("   Expected Hash:", computedHash)
			fmt.Println("   Stored Hashes in DB:", dbHashes) // Print all stored hashes
			return errors.New(fmt.Sprintf(" Blockchain validation failed! Block missing in PostgreSQL: %s", computedHash))
		}

		// Ensure timestamp consistency
		dbTimestamp := dbData["timestamp"] // Already converted to Unix format
		if fileTimestamp != dbTimestamp {
			fmt.Println(" Timestamp mismatch!")
			fmt.Printf("   File Timestamp (Unix): %s\n", fileTimestamp)
			fmt.Printf("   PostgreSQL Timestamp (Unix): %s\n", dbTimestamp)
			return errors.New(fmt.Sprintf(" Blockchain validation failed! Timestamp mismatch: File: %s, DB: %s", fileTimestamp, dbTimestamp))
		}

		// Ensure computed hash matches expected hash from logs.txt
		if computedHash != expectedHash {
			fmt.Println(" Hash mismatch!")
			fmt.Printf("   Expected: %s\n", expectedHash)
			fmt.Printf("   Computed: %s\n", computedHash)
			return errors.New(fmt.Sprintf(" Blockchain validation failed! Hash mismatch: Expected %s, Got %s", expectedHash, computedHash))
		}

		// Ensure previous hash matches the one stored in PostgreSQL
		if dbData["prev_hash"] != prevHash {
			fmt.Println(" Previous hash mismatch!")
			fmt.Printf("   Expected: %s\n", prevHash)
			fmt.Printf("   Found in DB: %s\n", dbData["prev_hash"])
			return errors.New(fmt.Sprintf(" Blockchain validation failed! Previous hash mismatch: Expected %s, Found %s", prevHash, dbData["prev_hash"]))
		}
	}

	fmt.Println(" Blockchain integrity verified! No tampering detected.")
	return nil
}
