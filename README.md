# VaultChain 🚀
A simple yet powerful blockchain implementation written in Go, with CLI support using Cobra and PostgreSQL validation.

---

## ✨ Features

- Build and manage a blockchain from scratch
- Add blocks with transaction data via CLI
- Validate blockchain integrity against a PostgreSQL database
- Persistent storage of block metadata
- Modular, clean Go architecture
- Easily extendable for future features

---

# How to Run the Project

This guide explains how to clone, build, and run the GoChain blockchain project with CLI support using Cobra and PostgreSQL validation.

---

## 1. Clone the Repository

```bash
git clone https://github.com/avirup-ghosal/VaultChain.git
cd VaultChain
```
## 2. Install Dependencies
```bash
go mod tidy
```
## 3. Set up your environment variables
Create a .env file in your root folder with the command
```bash
touch .env
```
Now add your Database url like this:
```text
DB_URL=your database url
```
## 4. Run the project
```bash
go run main.go
```
### Add a block
```bash
go run main.go addblock --data "Bought 5kg Basmati Rice for ₹500"
```
### Validate the blockchain
```bash
go run main.go validateblock
```
