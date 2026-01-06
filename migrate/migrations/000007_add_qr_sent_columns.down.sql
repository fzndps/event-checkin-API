DROP INDEX IF EXISTS idx_qr_sent ON participants;

ALTER TABLE participants
DROP COLUMN IF EXISTS qr_sent_at,
DROP COLUMN IF EXISTS qr_sent;