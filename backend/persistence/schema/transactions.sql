-- postgresql
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    transaction_type INT NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    description TEXT,
    category_id INT NOT NULL,
    user_id INT NOT NULL,
    currency_id INT NOT NULL,
    date DATE NOT NULL,
    receipt_url TEXT,
    FOREIGN KEY (category_id) REFERENCES categories (id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (currency_id) REFERENCES currencies (id)
);