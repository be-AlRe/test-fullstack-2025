package main

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

var (
	rdb = newRedisClient()
	ctx = context.Background()
)

type userRecord struct {
	RealName string `json:"realname"`
	Email    string `json:"email"`
	Password string `json:"password"` // SHA-1 hex
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Message  string `json:"message"`
	Username string `json:"username"`
	RealName string `json:"realname"`
	Email    string `json:"email"`
}

func newRedisClient() *redis.Client {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	password := os.Getenv("REDIS_PASSWORD") // kosongkan jika tidak pakai password
	db := 0

	return redis.NewClient(&redis.Options{
		Addr:        addr,
		Password:    password,
		DB:          db,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	})
}

func sha1Hex(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func main() {
	app := fiber.New(fiber.Config{
		Prefork:      false,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})

	app.Post("/login", handleLogin)

	log.Println("listening on :3000")
	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}

func handleLogin(c *fiber.Ctx) error {
	var req loginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if req.Username == "" || req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "username and password are required")
	}

	key := "login_" + req.Username
	val, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		// user tidak ditemukan
		return fiber.NewError(fiber.StatusUnauthorized, "invalid username or password")
	} else if err != nil {
		// error koneksi/redis
		return fiber.NewError(fiber.StatusInternalServerError, "redis error")
	}

	var rec userRecord
	if err := json.Unmarshal([]byte(val), &rec); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "invalid user record format")
	}

	// cek password: SHA-1 hex
	inputHash := sha1Hex(req.Password)
	if inputHash != rec.Password {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid username or password")
	}

	// sukses
	return c.Status(fiber.StatusOK).JSON(loginResponse{
		Message:  "login success",
		Username: req.Username,
		RealName: rec.RealName,
		Email:    rec.Email,
	})
}
