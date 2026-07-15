package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"backend/pb"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type relationServer struct {
	pb.UnimplementedRelationServiceServer
	db *pgxpool.Pool
}

func (s *relationServer) AssociateIngredients(ctx context.Context, req *pb.AssociateRequest) (*pb.AssociateResponse, error) {
	itemUUID, err := uuid.Parse(req.ItemUuid)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid item uuid: %v", err)
	}

	var ingredientUUIDs []uuid.UUID
	for _, idStr := range req.IngredientUuids {
		u, err := uuid.Parse(idStr)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid ingredient uuid %s: %v", idStr, err)
		}
		ingredientUUIDs = append(ingredientUUIDs, u)
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to start transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `DELETE FROM tm_item_ingredient WHERE uuid_item = $1`, itemUUID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete old relations: %v", err)
	}

	for _, ingUUID := range ingredientUUIDs {
		_, err = tx.Exec(ctx, `INSERT INTO tm_item_ingredient (uuid_item, uuid_ingredient) VALUES ($1, $2)`, itemUUID, ingUUID)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to insert relations: %v", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to commit transaction: %v", err)
	}

	return &pb.AssociateResponse{
		Success: true,
		Message: "associated successfully",
	}, nil
}

func (s *relationServer) GetIngredientsByItem(ctx context.Context, req *pb.GetIngredientsRequest) (*pb.GetIngredientsResponse, error) {
	itemUUID, err := uuid.Parse(req.ItemUuid)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid item uuid: %v", err)
	}

	rows, err := s.db.Query(ctx, `SELECT uuid_ingredient FROM tm_item_ingredient WHERE uuid_item = $1`, itemUUID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "query failed: %v", err)
	}
	defer rows.Close()

	var idStrings []string
	for rows.Next() {
		var u uuid.UUID
		if err := rows.Scan(&u); err != nil {
			return nil, status.Errorf(codes.Internal, "scan failed: %v", err)
		}
		idStrings = append(idStrings, u.String())
	}

	return &pb.GetIngredientsResponse{
		IngredientUuids: idStrings,
	}, nil
}

func (s *relationServer) GetItemsByIngredient(ctx context.Context, req *pb.GetItemsRequest) (*pb.GetItemsResponse, error) {
	ingredientUUID, err := uuid.Parse(req.IngredientUuid)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid ingredient uuid: %v", err)
	}

	rows, err := s.db.Query(ctx, `SELECT uuid_item FROM tm_item_ingredient WHERE uuid_ingredient = $1`, ingredientUUID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "query failed: %v", err)
	}
	defer rows.Close()

	var idStrings []string
	for rows.Next() {
		var u uuid.UUID
		if err := rows.Scan(&u); err != nil {
			return nil, status.Errorf(codes.Internal, "scan failed: %v", err)
		}
		idStrings = append(idStrings, u.String())
	}

	return &pb.GetItemsResponse{
		ItemUuids: idStrings,
	}, nil
}

func (s *relationServer) RemoveAllAssociationsForItem(ctx context.Context, req *pb.RemoveAllRequest) (*pb.RemoveAllResponse, error) {
	itemUUID, err := uuid.Parse(req.ItemUuid)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid item uuid: %v", err)
	}

	_, err = s.db.Exec(ctx, `DELETE FROM tm_item_ingredient WHERE uuid_item = $1`, itemUUID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "delete failed: %v", err)
	}

	return &pb.RemoveAllResponse{
		Success: true,
	}, nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
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
	dbPool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(ctx); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}
	log.Println("gRPC Server connected to PostgreSQL")

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", grpcPort, err)
	}

	s := grpc.NewServer()
	pb.RegisterRelationServiceServer(s, &relationServer{db: dbPool})

	log.Printf("gRPC server listening on port %s", grpcPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
