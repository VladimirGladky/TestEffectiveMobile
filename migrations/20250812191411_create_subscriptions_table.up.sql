CREATE TABLE IF NOT EXISTS subscriptions (
    id VARCHAR(255) NOT NULL,
    service_name VARCHAR(255) NOT NULL,
    price INT NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL
)