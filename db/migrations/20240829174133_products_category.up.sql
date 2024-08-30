CREATE TABLE IF NOT EXISTS products_category (
    product_id UUID NOT NULL,
    category_id UUID NOT NULL,

    FOREIGN KEY (product_id) REFERENCES products(id),
    FOREIGN KEY (category_id) REFERENCES category(id)
);