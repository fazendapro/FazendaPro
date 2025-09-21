-- Migration: Create vaccines table
-- Description: Table for vaccine catalog management (RF11)
-- Created: 2024-01-21

CREATE TABLE IF NOT EXISTS vaccines (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    manufacturer VARCHAR(255),
    description TEXT,
    dosage VARCHAR(100) COMMENT 'Dosage information (e.g., "5ml", "1 dose")',
    application_method VARCHAR(100) COMMENT 'Application method (e.g., "Intramuscular", "Subcutaneous")',
    validity_period_days INT DEFAULT 365 COMMENT 'Validity period in days',
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    
    INDEX idx_user_id (user_id),
    INDEX idx_name (name),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;