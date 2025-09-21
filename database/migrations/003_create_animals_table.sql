-- Migration: Create animals table
-- Description: Core table for animal management (RF02, RF03)
-- Created: 2024-01-21

CREATE TABLE IF NOT EXISTS animals (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    lot_id INT NULL,
    name VARCHAR(255) NOT NULL,
    ear_tag VARCHAR(100) UNIQUE NOT NULL COMMENT 'Número do brinco - unique identifier',
    birth_date DATE NOT NULL,
    race VARCHAR(255),
    gender ENUM('male', 'female') NOT NULL,
    mother_id INT NULL COMMENT 'Reference to genitora (mother)',
    father_id INT NULL COMMENT 'Reference to genitor (father)',
    status ENUM('active', 'sold', 'deceased', 'transferred') DEFAULT 'active',
    current_weight DECIMAL(8,2) NULL COMMENT 'Current weight in kg',
    daily_milk_production DECIMAL(8,2) DEFAULT 0.00 COMMENT 'Daily milk production in liters',
    notes TEXT COMMENT 'Additional health and general notes',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (lot_id) REFERENCES lots(id) ON DELETE SET NULL,
    FOREIGN KEY (mother_id) REFERENCES animals(id) ON DELETE SET NULL,
    FOREIGN KEY (father_id) REFERENCES animals(id) ON DELETE SET NULL,
    
    INDEX idx_user_id (user_id),
    INDEX idx_lot_id (lot_id),
    INDEX idx_ear_tag (ear_tag),
    INDEX idx_status (status),
    INDEX idx_gender (gender),
    INDEX idx_mother_id (mother_id),
    INDEX idx_father_id (father_id),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;