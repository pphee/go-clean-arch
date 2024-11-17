-- bmi.sql
CREATE TABLE bmi_records (
                             id BIGINT AUTO_INCREMENT PRIMARY KEY,
                             height DOUBLE NOT NULL,
                             weight DOUBLE NOT NULL,
                             value DOUBLE NOT NULL,
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);