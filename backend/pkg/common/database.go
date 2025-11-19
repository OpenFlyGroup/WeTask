package common

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/wetask/backend/pkg/models"
)

var DB *gorm.DB

// ? InitPostgreSQL initializes PostgreSQL connection
func InitPostgreSQL() error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "kanban"),
		getEnv("DB_PASSWORD", "kanban123"),
		getEnv("DB_NAME", "kanban"),
		getEnv("DB_PORT", "5432"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("PostgreSQL connected successfully")
	return nil
}

// ? MigrateAuthModels migrates models for auth service
func MigrateAuthModels() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}
	return DB.AutoMigrate(
		&models.AuthUser{},
		&models.RefreshToken{},
	)
}

// ? MigrateUsersModels migrates models for users service
func MigrateUsersModels() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}
	// ? Users service has its own copy of user data
	return DB.AutoMigrate(
		&models.User{},
	)
}

// ? MigrateTeamsModels migrates models for teams service
func MigrateTeamsModels() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}
	return DB.AutoMigrate(
		&models.Team{},
		&models.TeamMember{},
	)
}

// ? MigrateBoardsModels migrates models for boards service
func MigrateBoardsModels() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}
	return DB.AutoMigrate(
		&models.Board{},
		&models.Column{},
	)
}

// ? MigrateTasksModels migrates models for tasks service
func MigrateTasksModels() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}
	// Ensure task table exists
	if err := DB.AutoMigrate(&models.Task{}); err != nil {
		return err
	}

	// Remove any leftover DB-level foreign key constraint to the Users
	// table (this service uses the Users service and must not enforce
	// cross-database foreign keys). Drop constraint if it exists.
	migrator := DB.Migrator()
	if migrator.HasConstraint(&models.Task{}, "User") {
		_ = migrator.DropConstraint(&models.Task{}, "User")
	}
	// Some DBs name the FK like fk_tasks_user; also attempt to drop that.
	if migrator.HasConstraint(&models.Task{}, "fk_tasks_user") {
		_ = migrator.DropConstraint(&models.Task{}, "fk_tasks_user")
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
