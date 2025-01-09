package main

import (
	"fmt"
	"time"

	"context"

	"github.com/gin-gonic/gin"
	"github.com/nat-prohmpiriya/goobserv/pkg/core"
	middleware "github.com/nat-prohmpiriya/goobserv/pkg/middleware/gin"
	"github.com/nat-prohmpiriya/goobserv/pkg/output"
)

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
}

type UserUseCase interface {
	CreateUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, id string) (*User, error)
}

type userRepository struct {
	obs *core.Observer
}

func NewUserRepository(obs *core.Observer) UserRepository {
	return &userRepository{obs: obs}
}

func (r *userRepository) Create(ctx context.Context, user *User) error {
	span, ctx := r.obs.StartSpan(ctx, "repository.CreateUser")
	defer r.obs.EndSpan(span)

	// Simulate database operation
	time.Sleep(100 * time.Millisecond)

	// Log repository operation
	r.obs.Info(ctx, "Creating user in database").
		WithField("user_id", user.ID).
		WithField("user_email", user.Email)

	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*User, error) {
	span, ctx := r.obs.StartSpan(ctx, "repository.GetUser")
	defer r.obs.EndSpan(span)

	// Simulate database operation
	time.Sleep(50 * time.Millisecond)

	// Log repository operation
	r.obs.Info(ctx, "Getting user from database").
		WithField("user_id", id)

	// Simulate not found error
	if id == "not_found" {
		err := fmt.Errorf("user not found")
		r.obs.Error(ctx, "Failed to get user").
			WithError(err).
			WithField("user_id", id)
		return nil, err
	}

	return &User{
		ID:    id,
		Email: "user@example.com",
		Name:  "Test User",
	}, nil
}

type userUseCase struct {
	repo UserRepository
	obs  *core.Observer
}

func NewUserUseCase(repo UserRepository, obs *core.Observer) UserUseCase {
	return &userUseCase{
		repo: repo,
		obs:  obs,
	}
}

func (u *userUseCase) CreateUser(ctx context.Context, user *User) error {
	span, ctx := u.obs.StartSpan(ctx, "usecase.CreateUser")
	defer u.obs.EndSpan(span)

	// Log business operation
	u.obs.Info(ctx, "Creating new user").
		WithField("user_email", user.Email)

	// Create user in repository
	if err := u.repo.Create(ctx, user); err != nil {
		u.obs.Error(ctx, "Failed to create user").
			WithError(err).
			WithField("user_email", user.Email)
		return err
	}

	return nil
}

func (u *userUseCase) GetUser(ctx context.Context, id string) (*User, error) {
	span, ctx := u.obs.StartSpan(ctx, "usecase.GetUser")
	defer u.obs.EndSpan(span)

	// Log business operation
	u.obs.Info(ctx, "Getting user details").
		WithField("user_id", id)

	// Get user from repository
	user, err := u.repo.GetByID(ctx, id)
	if err != nil {
		u.obs.Error(ctx, "Failed to get user").
			WithError(err).
			WithField("user_id", id)
		return nil, err
	}

	return user, nil
}

type UserHandler struct {
	useCase UserUseCase
	obs     *core.Observer
}

func NewUserHandler(useCase UserUseCase, obs *core.Observer) *UserHandler {
	return &UserHandler{
		useCase: useCase,
		obs:     obs,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	// Get observer context
	ctx := middleware.GetContext(c)

	// Start handler span
	span, ctx := h.obs.StartSpan(ctx, "handler.CreateUser")
	defer h.obs.EndSpan(span)

	// Parse request body
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		h.obs.Error(ctx, "Invalid request body").
			WithError(err)
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	// Add metadata
	h.obs.Info(ctx, "Creating user").
		WithField("user_email", user.Email)

	// Call use case
	if err := h.useCase.CreateUser(ctx, &user); err != nil {
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(201, user)
}

func (h *UserHandler) GetUser(c *gin.Context) {
	// Get observer context
	ctx := middleware.GetContext(c)

	// Start handler span
	span, ctx := h.obs.StartSpan(ctx, "handler.GetUser")
	defer h.obs.EndSpan(span)

	// Get user ID from path
	id := c.Param("id")

	// Add metadata
	h.obs.Info(ctx, "Getting user").
		WithField("user_id", id)

	// Call use case
	user, err := h.useCase.GetUser(ctx, id)
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	c.JSON(200, user)
}

func main() {
	// Create observer
	obs := core.NewObserver(core.Config{
		BufferSize:    1000,
		FlushInterval: 1 * time.Second,
	})
	defer obs.Close()

	// Add stdout output with color
	stdout := output.NewStdoutOutput(output.StdoutConfig{
		Colored: true,
	})
	obs.AddOutput(stdout)

	// Create repository
	repo := NewUserRepository(obs)

	// Create use case
	useCase := NewUserUseCase(repo, obs)

	// Create handler
	handler := NewUserHandler(useCase, obs)

	// Create gin engine
	r := gin.New()

	// Add middleware
	r.Use(gin.Recovery())
	r.Use(middleware.Middleware(middleware.Config{
		Observer: obs,
	}))

	// Add routes
	r.POST("/users", handler.CreateUser)
	r.GET("/users/:id", handler.GetUser)

	// Start server
	r.Run(":8080")
}
