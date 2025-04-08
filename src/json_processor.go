package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"` // Using placeholder photo URL
}

func main() {
	resp, err := http.Get("https://jsonplaceholder.typicode.com/users")
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var users []User
	if err := json.Unmarshal(body, &users); err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	// Create output directory
	if err := os.Mkdir("user_data", 0755); err != nil && !os.IsExist(err) {
		fmt.Println("Error creating directory:", err)
		return
	}

	// Process first 10 users (instead of 100 as JSONPlaceholder only has 10 users)
	for i, user := range users {
		if i >= 10 {
			break
		}

		// Download avatar
		if err := downloadFile(user.Avatar, fmt.Sprintf("user_data/avatar_%d.jpg", user.ID)); err != nil {
			fmt.Printf("Error downloading avatar for user %d: %v\n", user.ID, err)
			continue
		}
	}

	// Save user data to JSON
	userJSON, err := json.MarshalIndent(users[:10], "", "  ")
	if err != nil {
		fmt.Println("Error marshaling user data:", err)
		return
	}

	if err := os.WriteFile("user_data/users.json", userJSON, 0644); err != nil {
		fmt.Println("Error writing user data file:", err)
		return
	}

	fmt.Println("Successfully processed users. Data saved in 'user_data' directory.")
}

func downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}