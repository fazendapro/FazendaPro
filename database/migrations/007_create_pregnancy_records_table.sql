-- Migration: Create pregnancy_records table
-- Description: Table for tracking animal pregnancy (RF09)
-- Created: 2024-01-21

CREATE TABLE IF NOT EXISTS pregnancy_records (
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