package main

import (
  "database/sql"
  "fmt"
  "os"
  "bufio"
  "strings"
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

func add(db *sql.DB) {
  fmt.Println("Adding new task")
  in := bufio.NewReader(os.Stdin)
  fmt.Println("Name:")
  var line string
  var err error
  line, err = in.ReadString('\n') 
  if err != nil {
    panic(err)
  }
  name := strings.TrimSpace(line)
  fmt.Println("Description:")
  line, err = in.ReadString('\n') 
  if err != nil {
    panic(err)
  }
  description := strings.TrimSpace(line)
  fmt.Println("Urgency:")
  if err != nil {
    panic(err)
  }
  line, err = in.ReadString('\n') 
  urgency := strings.TrimSpace(line)
  fmt.Printf("Your task's name is: %s\n", name)
  fmt.Printf("Your task's description is: %s\n", description)
  fmt.Printf("Your task's urgency is: %s\n", urgency)
  query := "INSERT INTO task (name, description, urgency) VALUES ('" + name + 
                                                  "', '" + description + "', '" + urgency + "');"
  _, err = db.Exec(query)
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
      add(db)
  } else {
    fmt.Println("Unknown command")
  }
}
