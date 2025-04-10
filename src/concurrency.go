package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

const (
	baseURL       = "https://reqres.in/api"
	maxConcurrent = 5
)

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"first_name"` // Note: reqres.in uses first_name/last_name
}

type UsersResponse struct {
	Data []User `json:"data"`
}

func listUsers() ([]int, error) {
	resp, err := http.Get(baseURL + "/users")
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response UsersResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Extract just the user IDs
	var userIds []int
	for _, user := range response.Data {
		userIds = append(userIds, user.ID)
	}

	return userIds, nil
}

func getUser(userId int, wg *sync.WaitGroup, sem chan struct{}, results chan<- User) {
	defer wg.Done()

	// Acquire a semaphore slot
	sem <- struct{}{}
	defer func() { <-sem }()

	resp, err := http.Get(fmt.Sprintf("%s/users/%d", baseURL, userId))
	if err != nil {
		fmt.Printf("Error fetching user %d: %v\n", userId, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Unexpected status code for user %d: %d\n", userId, resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response for user %d: %v\n", userId, err)
		return
	}

	var userResponse struct {
		Data User `json:"data"`
	}
	if err := json.Unmarshal(body, &userResponse); err != nil {
		fmt.Printf("Error parsing user %d: %v\n", userId, err)
		return
	}

	results <- userResponse.Data
}

func main() {
	// Step 1: Get list of user IDs
	userIds, err := listUsers()
	if err != nil {
		fmt.Printf("Error listing users: %v\n", err)
		return
	}

	if len(userIds) == 0 {
		fmt.Println("No users found")
		return
	}

	// Step 2: Fetch each user concurrently with max 5 requests at a time
	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrent)
	results := make(chan User, len(userIds))

	for _, userId := range userIds {
		wg.Add(1)
		go getUser(userId, &wg, sem, results)
	}

	// Close the results channel when all goroutines are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	var users []User
	for user := range results {
		users = append(users, user)
	}

	// Print the results
	fmt.Printf("\nFetched %d users:\n", len(users))
	for _, user := range users {
		fmt.Printf("ID: %d, Name: %s, Email: %s\n", user.ID, user.Name, user.Email)
	}
}