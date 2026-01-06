ALTER TABLE participants 
ADD COLUMN qr_sent BOOLEAN DEFAULT FALSE AFTER qr_token,
ADD COLUMN qr_sent_at TIMESTAMP NULL AFTER qr_sent;

CREATE INDEX idx_qr_sent ON participants(qr_sent);