CREATE TABLE IF NOT EXISTS
medical_records (
    id VARCHAR(16) PRIMARY KEY,
    user_id VARCHAR(16) NOT NULL,
    patient_id VARCHAR(16) NOT NULL,
    -- identity_number CHAR(16) UNIQUE,
    symptoms VARCHAR(2000) NOT NULL,
    medications VARCHAR(2000) NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp
);

ALTER TABLE medical_records
	ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE medical_records
	ADD CONSTRAINT fk_patient_id FOREIGN KEY (patient_id) REFERENCES medical_patients(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS medical_records_user_id
	ON medical_records USING HASH(user_id);
CREATE INDEX IF NOT EXISTS medical_records_patient_id
	ON medical_records USING HASH(patient_id);
-- CREATE INDEX IF NOT EXISTS medical_records_identity_number
-- 	ON medical_records(identity_number);
CREATE INDEX IF NOT EXISTS medical_records_created_at_desc
	ON medical_records(created_at DESC);
CREATE INDEX IF NOT EXISTS medical_records_created_at_asc
	ON medical_records(created_at ASC);