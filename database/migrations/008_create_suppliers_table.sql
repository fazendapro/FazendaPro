-- Migration: Create suppliers table
-- Description: Table for supplier management (RF12)
-- Created: 2024-01-21

CREATE TABLE IF NOT EXISTS suppliers (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    company_name VARCHAR(255),
    cnpj VARCHAR(18) COMMENT 'Brazilian company registry number',
    cpf VARCHAR(14) COMMENT 'Brazilian individual registry number',
    email VARCHAR(255),
    phone VARCHAR(20),
    address TEXT,
    city VARCHAR(100),
    state VARCHAR(2),
    zip_code VARCHAR(10),
    supplier_type ENUM('feed', 'medicine', 'equipment', 'services', 'other') DEFAULT 'other',
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    
    INDEX idx_user_id (user_id),
    INDEX idx_name (name),
    INDEX idx_cnpj (cnpj),
    INDEX idx_cpf (cpf),
    INDEX idx_supplier_type (supplier_type),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;