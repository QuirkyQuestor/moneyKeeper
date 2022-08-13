CREATE TABLE account (
    accountId SERIAL,
    typeId BIGINT UNSIGNED NOT NULL,
    name VARCHAR(40) CHARACTER SET utf8 UNIQUE,
    description VARCHAR(200) CHARACTER SET utf8,
    active BOOLEAN
);

--  SERIAL is an alias for BIGINT UNSIGNED NOT NULL AUTO_INCREMENT UNIQUE.

CREATE TABLE accountType (
    typeId SERIAL,
    name VARCHAR(40) CHARACTER SET utf8 UNIQUE,
    description VARCHAR(200) CHARACTER SET utf8
);

CREATE TABLE category (
    categoryId SERIAL,
    parentId BIGINT UNSIGNED NOT NULL,
    name VARCHAR(40) CHARACTER SET utf8,
    description VARCHAR(200) CHARACTER SET utf8,
    expence BOOLEAN
);
