-- Migration: Create weight_records table
-- Description: Table for tracking animal weight over time (RF07)
-- Created: 2024-01-21

CREATE TABLE IF NOT EXISTS weight_records (
    id INT PRIMARY KEY AUTO_INCREMENT,
    animal_id INT NOT NULL,
    weight DECIMAL(8,2) NOT NULL COMMENT 'Weight in kilograms',
    measurement_date DATE NOT NULL,
    measurement_type ENUM('weekly', 'monthly', 'special') DEFAULT 'monthly',
    notes TEXT COMMENT 'Additional notes about the measurement',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    
    FOREIGN KEY (animal_id) REFERENCES animals(id) ON DELETE CASCADE,
    
    INDEX idx_animal_id (animal_id),
    INDEX idx_measurement_date (measurement_date),
    INDEX idx_measurement_type (measurement_type),
    INDEX idx_deleted_at (deleted_at),
    
    UNIQUE KEY unique_animal_measurement_date (animal_id, measurement_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;