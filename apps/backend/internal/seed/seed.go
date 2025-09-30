package seed

import (
	"log"
	"os"

	"github.com/hariomop12/clearrouter/apps/backend/internal/models"
	"gorm.io/gorm"
)

// SeedDefaultUser creates an initial user if it doesn't already exist.
// Config via env (with sensible defaults):
//   SEED_DEFAULT_USER_EMAIL    (default: admin@clearrouter.local)
//   SEED_DEFAULT_USER_NAME     (default: Admin)
//   SEED_DEFAULT_USER_PASSWORD (default: admin123)
//   SEED_ENABLE                (optional: if set to "false", seeding is skipped)
func SeedDefaultUser(db *gorm.DB) {
	if os.Getenv("SEED_ENABLE") == "false" {
		log.Println("Seed: skipped (SEED_ENABLE=false)")
		return
	}

	email := getenvDefault("SEED_DEFAULT_USER_EMAIL", "admin@clearrouter.local")
	name := getenvDefault("SEED_DEFAULT_USER_NAME", "Admin")
	password := getenvDefault("SEED_DEFAULT_USER_PASSWORD", "admin123")

	var existing models.User
	if err := db.Where("email = ?", email).First(&existing).Error; err == nil {
		log.Printf("Seed: user already exists (%s), skipping\n", email)
		return
	}

	u := models.User{
		Name:          name,
		Email:         email,
		EmailVerified: true,
	}
	if err := u.HashPassword(password); err != nil {
		log.Printf("Seed: failed to hash password: %v\n", err)
		return
	}

	if err := db.Create(&u).Error; err != nil {
		log.Printf("Seed: failed to create user: %v\n", err)
		return
	}
	log.Printf("Seed: created default user %s (email verified)\n", email)
}

func getenvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
