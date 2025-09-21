-- Migration: Create lots table
-- Description: Table for managing production lots (RF08)
-- Created: 2024-01-21

CREATE TABLE IF NOT EXISTS lots (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    min_milk_production DECIMAL(8,2) DEFAULT 0.00 COMMENT 'Minimum milk production in liters',
    max_milk_production DECIMAL(8,2) DEFAULT 999999.99 COMMENT 'Maximum milk production in liters',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    
    INDEX idx_production_range (min_milk_production, max_milk_production),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;