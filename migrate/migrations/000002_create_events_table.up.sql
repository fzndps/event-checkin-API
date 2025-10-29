CREATE TABLE IF NOT EXISTS events (
    id VARCHAR(36) PRIMARY KEY,
    organizer_id BIGINT UNSIGNED NOT NULL,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    date DATETIME NOT NULL,
    venue VARCHAR(500) NOT NULL,
    participant_count INT NOT NULL DEFAULT 0,
    total_price INT NOT NULL DEFAULT 0,
    payment_status ENUM('pending', 'verified', 'active') NOT NULL DEFAULT 'pending',
    payment_proof_url VARCHAR(500) NULL,
    scanner_pin CHAR(4) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (organizer_id) REFERENCES organizers(id) ON DELETE CASCADE
);
