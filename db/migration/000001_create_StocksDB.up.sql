CREATE TABLE "stocks" (
    "id" BIGSERIAL PRIMARY KEY,
    "name" TEXT NOT NULL,
    "price" BIGINT NOT NULL,
    "company" TEXT NOT NULL
);