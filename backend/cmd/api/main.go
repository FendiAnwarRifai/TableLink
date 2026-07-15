package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"backend/pb"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Model Data untuk Bahan Makanan
type Ingredient struct {
	UUID        uuid.UUID  `json:"uuid"`
	Name        string     `json:"name"`
	CauseAlergy bool       `json:"cause_alergy"`
	Type        int        `json:"type"`
	Status      int        `json:"status"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

// Model Data untuk Menu Item
type Item struct {
	UUID        uuid.UUID     `json:"uuid"`
	Name        string        `json:"name"`
	Price       float64       `json:"price"`
	Status      int           `json:"status"`
	CreatedAt   *time.Time    `json:"created_at"`
	UpdatedAt   *time.Time    `json:"updated_at"`
	DeletedAt   *time.Time    `json:"deleted_at,omitempty"`
	Ingredients []*Ingredient `json:"ingredients,omitempty"`
}

// Format data yang dikirim dari Frontend untuk Item
type CreateOrUpdateItemRequest struct {
	Name            string      `json:"name"`
	Price           float64     `json:"price"`
	Status          int         `json:"status"`
	IngredientUUIDs []uuid.UUID `json:"ingredients"`
}

// Global database pool and gRPC client connection variables
var db *pgxpool.Pool
var relationClient pb.RelationServiceClient

func main() {
	// Load environment variables from .env
	if err := godotenv.Load(); err != nil {
		log.Println("Using system environment variables...")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	ctx := context.Background()
	var err error
	db, err = pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}
	log.Println("Connected to PostgreSQL database!")

	// Connect to gRPC server (for table relationships)
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}
	grpcAddr := fmt.Sprintf("localhost:%s", grpcPort)
	grpcConn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer grpcConn.Close()
	log.Printf("Connected to gRPC server on %s", grpcAddr)

	relationClient = pb.NewRelationServiceClient(grpcConn)

	// Initialize Fiber Web Application
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		},
	})

	// Enable CORS for frontend API calls
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	api := app.Group("/api")

	// Ingredients (Bahan) API Routes
	api.Get("/ingredients", getIngredients)
	api.Get("/ingredients/:uuid", getIngredient)
	api.Post("/ingredients", createIngredient)
	api.Put("/ingredients/:uuid", updateIngredient)
	api.Delete("/ingredients/:uuid", deleteIngredient)

	// Items (Menu) API Routes
	api.Get("/items", getItems)
	api.Get("/items/:uuid", getItem)
	api.Post("/items", createItem)
	api.Put("/items/:uuid", updateItem)
	api.Delete("/items/:uuid", deleteItem)

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "3000"
	}
	log.Printf("Fiber HTTP server running on port %s", httpPort)
	if err := app.Listen(fmt.Sprintf(":%s", httpPort)); err != nil {
		log.Fatalf("HTTP server failed to start: %v", err)
	}
}

// ==========================================
// Ingredients (Bahan) Route Handlers
// ==========================================

func getIngredients(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	var total int
	err := db.QueryRow(c.Context(), `SELECT count(*) FROM tm_ingredient WHERE deleted_at IS NULL`).Scan(&total)
	if err != nil {
		return err
	}

	rows, err := db.Query(c.Context(), `SELECT uuid, name, cause_alergy, type, status, created_at, updated_at, deleted_at FROM tm_ingredient WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return err
	}
	defer rows.Close()

	var ingredients []*Ingredient
	for rows.Next() {
		var ing Ingredient
		err := rows.Scan(&ing.UUID, &ing.Name, &ing.CauseAlergy, &ing.Type, &ing.Status, &ing.CreatedAt, &ing.UpdatedAt, &ing.DeletedAt)
		if err != nil {
			return err
		}
		ingredients = append(ingredients, &ing)
	}

	return c.JSON(fiber.Map{
		"data":  ingredients,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func getIngredient(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("uuid"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ingredient UUID"})
	}

	var ing Ingredient
	query := `SELECT uuid, name, cause_alergy, type, status, created_at, updated_at, deleted_at FROM tm_ingredient WHERE uuid = $1 AND deleted_at IS NULL`
	err = db.QueryRow(c.Context(), query, id).Scan(&ing.UUID, &ing.Name, &ing.CauseAlergy, &ing.Type, &ing.Status, &ing.CreatedAt, &ing.UpdatedAt, &ing.DeletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Ingredient not found"})
		}
		return err
	}
	return c.JSON(ing)
}

func createIngredient(c *fiber.Ctx) error {
	var ing Ingredient
	if err := c.BodyParser(&ing); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if ing.Name == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Name is required"})
	}

	// Validate name uniqueness
	var count int
	err := db.QueryRow(c.Context(), `SELECT count(*) FROM tm_ingredient WHERE LOWER(name) = LOWER($1) AND deleted_at IS NULL`, ing.Name).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Ingredient name already exists"})
	}

	ing.UUID = uuid.New()
	now := time.Now()
	ing.CreatedAt = &now
	ing.UpdatedAt = &now

	query := `INSERT INTO tm_ingredient (uuid, name, cause_alergy, type, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = db.Exec(c.Context(), query, ing.UUID, ing.Name, ing.CauseAlergy, ing.Type, ing.Status, ing.CreatedAt, ing.UpdatedAt)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(ing)
}

func updateIngredient(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("uuid"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ingredient UUID"})
	}

	var ing Ingredient
	if err := c.BodyParser(&ing); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if ing.Name == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Name is required"})
	}

	// Validate name uniqueness (excluding current UUID)
	var count int
	err = db.QueryRow(c.Context(), `SELECT count(*) FROM tm_ingredient WHERE LOWER(name) = LOWER($1) AND uuid != $2 AND deleted_at IS NULL`, ing.Name, id).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Ingredient name already exists"})
	}

	now := time.Now()
	ing.UUID = id
	ing.UpdatedAt = &now

	query := `UPDATE tm_ingredient SET name = $1, cause_alergy = $2, type = $3, status = $4, updated_at = $5 WHERE uuid = $6 AND deleted_at IS NULL`
	_, err = db.Exec(c.Context(), query, ing.Name, ing.CauseAlergy, ing.Type, ing.Status, ing.UpdatedAt, ing.UUID)
	if err != nil {
		return err
	}

	return c.JSON(ing)
}

func deleteIngredient(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("uuid"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ingredient UUID"})
	}

	// Soft delete (setting deleted_at timestamp)
	query := `UPDATE tm_ingredient SET deleted_at = $1 WHERE uuid = $2 AND deleted_at IS NULL`
	_, err = db.Exec(c.Context(), query, time.Now(), id)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"message": "Ingredient successfully deleted"})
}

// ==========================================
// Items (Menu) Route Handlers
// ==========================================

func getItems(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	var total int
	err := db.QueryRow(c.Context(), `SELECT count(*) FROM tm_item WHERE deleted_at IS NULL`).Scan(&total)
	if err != nil {
		return err
	}

	rows, err := db.Query(c.Context(), `SELECT uuid, name, price, status, created_at, updated_at, deleted_at FROM tm_item WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return err
	}
	defer rows.Close()

	var items []*Item
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.UUID, &item.Name, &item.Price, &item.Status, &item.CreatedAt, &item.UpdatedAt, &item.DeletedAt)
		if err != nil {
			return err
		}

		// Fetch ingredient relations from gRPC Server
		resp, err := relationClient.GetIngredientsByItem(c.Context(), &pb.GetIngredientsRequest{ItemUuid: item.UUID.String()})
		if err == nil {
			var ingredients []*Ingredient
			for _, ingUUIDStr := range resp.IngredientUuids {
				ingUUID, parseErr := uuid.Parse(ingUUIDStr)
				if parseErr == nil {
					var ing Ingredient
					q := `SELECT uuid, name, cause_alergy, type, status, created_at, updated_at, deleted_at FROM tm_ingredient WHERE uuid = $1 AND deleted_at IS NULL`
					err = db.QueryRow(c.Context(), q, ingUUID).Scan(&ing.UUID, &ing.Name, &ing.CauseAlergy, &ing.Type, &ing.Status, &ing.CreatedAt, &ing.UpdatedAt, &ing.DeletedAt)
					if err == nil {
						ingredients = append(ingredients, &ing)
					}
				}
			}
			item.Ingredients = ingredients
		}

		items = append(items, &item)
	}

	return c.JSON(fiber.Map{
		"data":  items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func getItem(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("uuid"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid item UUID"})
	}

	var item Item
	query := `SELECT uuid, name, price, status, created_at, updated_at, deleted_at FROM tm_item WHERE uuid = $1 AND deleted_at IS NULL`
	err = db.QueryRow(c.Context(), query, id).Scan(&item.UUID, &item.Name, &item.Price, &item.Status, &item.CreatedAt, &item.UpdatedAt, &item.DeletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Item not found"})
		}
		return err
	}

	// Fetch ingredients from gRPC Server
	resp, err := relationClient.GetIngredientsByItem(c.Context(), &pb.GetIngredientsRequest{ItemUuid: item.UUID.String()})
	if err == nil {
		var ingredients []*Ingredient
		for _, ingUUIDStr := range resp.IngredientUuids {
			ingUUID, parseErr := uuid.Parse(ingUUIDStr)
			if parseErr == nil {
				var ing Ingredient
				q := `SELECT uuid, name, cause_alergy, type, status, created_at, updated_at, deleted_at FROM tm_ingredient WHERE uuid = $1 AND deleted_at IS NULL`
				err = db.QueryRow(c.Context(), q, ingUUID).Scan(&ing.UUID, &ing.Name, &ing.CauseAlergy, &ing.Type, &ing.Status, &ing.CreatedAt, &ing.UpdatedAt, &ing.DeletedAt)
				if err == nil {
					ingredients = append(ingredients, &ing)
				}
			}
		}
		item.Ingredients = ingredients
	}

	return c.JSON(item)
}

func createItem(c *fiber.Ctx) error {
	var req CreateOrUpdateItemRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if req.Name == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Name is required"})
	}
	if req.Price <= 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Price must be greater than 0"})
	}
	if len(req.IngredientUUIDs) == 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "At least one ingredient is required"})
	}

	// Validate name uniqueness
	var count int
	err := db.QueryRow(c.Context(), `SELECT count(*) FROM tm_item WHERE LOWER(name) = LOWER($1) AND deleted_at IS NULL`, req.Name).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Item name already exists"})
	}

	var item Item
	item.UUID = uuid.New()
	item.Name = req.Name
	item.Price = req.Price
	item.Status = req.Status
	now := time.Now()
	item.CreatedAt = &now
	item.UpdatedAt = &now

	query := `INSERT INTO tm_item (uuid, name, price, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = db.Exec(c.Context(), query, item.UUID, item.Name, item.Price, item.Status, item.CreatedAt, item.UpdatedAt)
	if err != nil {
		return err
	}

	// Associate ingredients using gRPC Server
	var ingUUIDsStr []string
	for _, id := range req.IngredientUUIDs {
		ingUUIDsStr = append(ingUUIDsStr, id.String())
	}

	_, err = relationClient.AssociateIngredients(c.Context(), &pb.AssociateRequest{
		ItemUuid:        item.UUID.String(),
		IngredientUuids: ingUUIDsStr,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(item)
}

func updateItem(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("uuid"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid item UUID"})
	}

	var req CreateOrUpdateItemRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if req.Name == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Name is required"})
	}
	if req.Price <= 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Price must be greater than 0"})
	}
	if len(req.IngredientUUIDs) == 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "At least one ingredient is required"})
	}

	// Validate name uniqueness (excluding current UUID)
	var count int
	err = db.QueryRow(c.Context(), `SELECT count(*) FROM tm_item WHERE LOWER(name) = LOWER($1) AND uuid != $2 AND deleted_at IS NULL`, req.Name, id).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Item name already exists"})
	}

	now := time.Now()
	var item Item
	item.UUID = id
	item.Name = req.Name
	item.Price = req.Price
	item.Status = req.Status
	item.UpdatedAt = &now

	query := `UPDATE tm_item SET name = $1, price = $2, status = $3, updated_at = $4 WHERE uuid = $5 AND deleted_at IS NULL`
	_, err = db.Exec(c.Context(), query, item.Name, item.Price, item.Status, item.UpdatedAt, item.UUID)
	if err != nil {
		return err
	}

	// Update ingredients using gRPC Server
	var ingUUIDsStr []string
	for _, id := range req.IngredientUUIDs {
		ingUUIDsStr = append(ingUUIDsStr, id.String())
	}

	_, err = relationClient.AssociateIngredients(c.Context(), &pb.AssociateRequest{
		ItemUuid:        item.UUID.String(),
		IngredientUuids: ingUUIDsStr,
	})
	if err != nil {
		return err
	}

	return c.JSON(item)
}

func deleteItem(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("uuid"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid item UUID"})
	}

	query := `UPDATE tm_item SET deleted_at = $1 WHERE uuid = $2 AND deleted_at IS NULL`
	_, err = db.Exec(c.Context(), query, time.Now(), id)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"message": "Item successfully deleted"})
}
