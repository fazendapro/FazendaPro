-- FazendaPro Database Complete Schema
-- Description: Complete database schema for the farm management system
-- Created: 2024-01-21
-- 
-- This file contains all tables for the FazendaPro system based on functional requirements.
-- Execute migrations in order: 001-009

-- First drop all tables in reverse order to avoid foreign key constraints
DROP TABLE IF EXISTS sales;
DROP TABLE IF EXISTS suppliers;
DROP TABLE IF EXISTS pregnancy_records;
DROP TABLE IF EXISTS weight_records;
DROP TABLE IF EXISTS animal_vaccines;
DROP TABLE IF EXISTS vaccines;
DROP TABLE IF EXISTS animals;
DROP TABLE IF EXISTS lots;
DROP TABLE IF EXISTS users;

-- Create all tables in the correct order

-- Users table (RF01 - Authentication)
CREATE TABLE users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    
    INDEX idx_email (email),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Lots table (RF08 - Lot Management)
CREATE TABLE lots (
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

-- Animals table (RF02, RF03 - Core animal management)
CREATE TABLE animals (
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

-- Vaccines table (RF11 - Vaccine management)
CREATE TABLE vaccines (
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

-- Animal vaccines junction table (RF06 - Animal vaccination records)
CREATE TABLE animal_vaccines (
    id INT PRIMARY KEY AUTO_INCREMENT,
    animal_id INT NOT NULL,
    vaccine_id INT NOT NULL,
    application_date DATE NOT NULL,
    next_application_date DATE NULL COMMENT 'Next scheduled vaccination date',
    batch_number VARCHAR(100) COMMENT 'Vaccine batch number for tracking',
    administered_by VARCHAR(255) COMMENT 'Person who administered the vaccine',
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    
    FOREIGN KEY (animal_id) REFERENCES animals(id) ON DELETE CASCADE,
    FOREIGN KEY (vaccine_id) REFERENCES vaccines(id) ON DELETE CASCADE,
    
    INDEX idx_animal_id (animal_id),
    INDEX idx_vaccine_id (vaccine_id),
    INDEX idx_application_date (application_date),
    INDEX idx_next_application_date (next_application_date),
    INDEX idx_deleted_at (deleted_at),
    
    UNIQUE KEY unique_animal_vaccine_date (animal_id, vaccine_id, application_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Weight records table (RF07 - Weight tracking)
CREATE TABLE weight_records (
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

-- Pregnancy records table (RF09 - Pregnancy tracking)
CREATE TABLE pregnancy_records (
    id INT PRIMARY KEY AUTO_INCREMENT,
    animal_id INT NOT NULL COMMENT 'The pregnant female animal',
    father_id INT NULL COMMENT 'The male animal (bull) - if known',
    pregnancy_date DATE NOT NULL COMMENT 'Date when pregnancy was confirmed',
    expected_birth_date DATE NOT NULL COMMENT 'Calculated expected birth date',
    actual_birth_date DATE NULL COMMENT 'Actual birth date when it occurs',
    pregnancy_status ENUM('confirmed', 'completed', 'aborted', 'false_positive') DEFAULT 'confirmed',
    notification_sent BOOLEAN DEFAULT FALSE COMMENT 'WhatsApp notification sent flag',
    notification_date TIMESTAMP NULL COMMENT 'When notification was sent',
    notes TEXT COMMENT 'Additional pregnancy notes',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    
    FOREIGN KEY (animal_id) REFERENCES animals(id) ON DELETE CASCADE,
    FOREIGN KEY (father_id) REFERENCES animals(id) ON DELETE SET NULL,
    
    INDEX idx_animal_id (animal_id),
    INDEX idx_father_id (father_id),
    INDEX idx_pregnancy_date (pregnancy_date),
    INDEX idx_expected_birth_date (expected_birth_date),
    INDEX idx_pregnancy_status (pregnancy_status),
    INDEX idx_notification_sent (notification_sent),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Suppliers table (RF12 - Supplier management)
CREATE TABLE suppliers (
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

-- Sales table (RF10 - Sales management)
CREATE TABLE sales (
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