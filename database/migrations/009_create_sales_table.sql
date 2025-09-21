-- Migration: Create sales table
-- Description: Table for managing animal sales (RF10)
-- Created: 2024-01-21

CREATE TABLE IF NOT EXISTS sales (
    id INT PRIMARY KEY AUTO_INCREMENT,
    animal_id INT NOT NULL,
    buyer_name VARCHAR(255) NOT NULL,
    buyer_document VARCHAR(20) COMMENT 'CPF or CNPJ of buyer',
    buyer_contact VARCHAR(255) COMMENT 'Phone or email',
    sale_date DATE NOT NULL,
    sale_price DECIMAL(10,2) NOT NULL COMMENT 'Sale price in local currency',
    payment_method ENUM('cash', 'bank_transfer', 'check', 'installments', 'other') DEFAULT 'cash',
    payment_status ENUM('pending', 'partial', 'completed') DEFAULT 'pending',
    contract_pdf_path VARCHAR(500) COMMENT 'Path to generated PDF contract',
    animal_history_pdf_path VARCHAR(500) COMMENT 'Path to animal history PDF',
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    
    FOREIGN KEY (animal_id) REFERENCES animals(id) ON DELETE CASCADE,
    
    INDEX idx_animal_id (animal_id),
    INDEX idx_sale_date (sale_date),
    INDEX idx_buyer_name (buyer_name),
    INDEX idx_payment_status (payment_status),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;