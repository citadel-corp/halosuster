DROP TYPE IF EXISTS gender;
CREATE TYPE gender AS ENUM('male', 'female');

CREATE TABLE IF NOT EXISTS
medical_patients (
    id CHAR(16) PRIMARY KEY,
    identity_number BIGINT UNIQUE,
    phone_number VARCHAR(16) NOT NULL,
    name VARCHAR(30) NOT NULL,
    birth_date DATE NOT NULL,
		gender gender NOT NULL,
		identity_card_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp
);

CREATE INDEX IF NOT EXISTS medical_patients_name
	ON medical_patients(lower(name));
CREATE INDEX IF NOT EXISTS medical_patients_phone_number
	ON medical_patients(phone_number);
CREATE INDEX IF NOT EXISTS medical_patients_identity_number
	ON medical_patients(identity_number);
CREATE INDEX IF NOT EXISTS medical_patients_created_at_desc
	ON medical_patients(created_at DESC);
CREATE INDEX IF NOT EXISTS products_created_at_asc
	ON medical_patients(created_at ASC);