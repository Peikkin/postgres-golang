CREATE TABLE "stocks" (
    "stock_id" SERIAL PRIMARY KEY,
    "name" TEXT NOT NULL,
    "price" INT NOT NULL,
    "company" TEXT NOT NULL
);