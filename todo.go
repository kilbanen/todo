package main

import (
  "database/sql"
  "fmt"
  "os"
  "gopkg.in/yaml.v2"
  _ "github.com/lib/pq"
)

type Config struct {
  Host string `yaml:"host"`
  Port int `yaml:"port"`
  User string `yaml:"user"`
  Dbname string `yaml:"dbname"`
}

func list(db *sql.DB) {
  rows, err := db.Query("SELECT * FROM task")
  if err != nil {
    panic(err)
  }
  defer rows.Close()
  for rows.Next() {
    var id int
    var name string
    var description string
    var urgency string
    err = rows.Scan(&id, &name, &description, &urgency)
    if err != nil {
      panic(err)
    }
    fmt.Println(id, name, description, urgency)
  }
  err = rows.Err()
  if err != nil {
    panic(err)
  }
}

func main() {
  f, err := os.Open("config.yml")
  if err != nil {
      panic(err)
  }

  defer f.Close()

  var cfg Config
  decoder := yaml.NewDecoder(f)
  err = decoder.Decode(&cfg)
  if err != nil {
      panic(err)
  }
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }
  defer db.Close()

  err = db.Ping()
  if err != nil {
    panic(err)
  }
  fmt.Println("Successfully connected!")

  command := os.Args[1]
  if command == "list" {
    list(db)
  } else if command == "add" {
    fmt.Println("Adding new task")
  } else {
    fmt.Println("Unknown command")
  }
}
