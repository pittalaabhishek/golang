package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type BlogPost struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Title     string `gorm:"size:255" json:"title"`
	Content   string `json:"content"`
	Author    string `gorm:"size:100" json:"author"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt int64  `gorm:"autoUpdateTime" json:"updated_at"`
}

var db *gorm.DB

func main() {

	dsn := "root:Mo73cu73!@tcp(127.0.0.1:3306)/myDB?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate
	if err := db.AutoMigrate(&BlogPost{}); err != nil {
		log.Fatal("Auto migration failed:", err)
	}

	// Set up routes
	http.HandleFunc("/posts", postsHandler)
	http.HandleFunc("/posts/", postHandler)
	http.HandleFunc("/posts/search", searchHandler)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func postsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getPosts(w, r)
	case http.MethodPost:
		createPost(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/posts/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPut:
		updatePost(w, r, uint(id))
	case http.MethodDelete:
		deletePost(w, r, uint(id))
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func createPost(w http.ResponseWriter, r *http.Request) {
	var post BlogPost
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if post.Title == "" || post.Content == "" || post.Author == "" {
		http.Error(w, "Title, content and author are required", http.StatusBadRequest)
		return
	}

	if err := db.Create(&post).Error; err != nil {
		log.Println("Create error:", err)
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	log.Println("Created post ID:", post.ID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func getPosts(w http.ResponseWriter, r *http.Request) {
	var posts []BlogPost
	if err := db.Find(&posts).Error; err != nil {
		log.Println("Get posts error:", err)
		http.Error(w, "Failed to get posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func updatePost(w http.ResponseWriter, r *http.Request, id uint) {
	var post BlogPost
	if err := db.First(&post, id).Error; err != nil {
		log.Println("Update find error:", err)
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	var updateData struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if updateData.Title == "" || updateData.Content == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	post.Title = updateData.Title
	post.Content = updateData.Content

	if err := db.Save(&post).Error; err != nil {
		log.Println("Update save error:", err)
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		return
	}

	log.Println("Updated post ID:", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func deletePost(w http.ResponseWriter, r *http.Request, id uint) {
	if err := db.Delete(&BlogPost{}, id).Error; err != nil {
		log.Println("Delete error:", err)
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}
	log.Println("Deleted post ID:", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Post deleted"})
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	title := r.URL.Query().Get("title")
	author := r.URL.Query().Get("author")

	var posts []BlogPost
	query := db

	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}
	if author != "" {
		query = query.Where("author LIKE ?", "%"+author+"%")
	}

	if err := query.Find(&posts).Error; err != nil {
		log.Println("Search error:", err)
		http.Error(w, "Failed to search posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}