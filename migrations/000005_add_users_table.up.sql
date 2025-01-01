CREATE TABLE
    Users (
        Id BIGSERIAL PRIMARY KEY,
        Name VARCHAR(255) NOT NULL UNIQUE,
        Email VARCHAR(255) NOT NULL UNIQUE,
        Password_hash VARCHAR(255) NOT NULL
    ); 