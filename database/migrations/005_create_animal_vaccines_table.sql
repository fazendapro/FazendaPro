-- Migration: Create animal_vaccines table
-- Description: Junction table for animal vaccination records (RF06)
-- Created: 2024-01-21

CREATE TABLE IF NOT EXISTS animal_vaccines (
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