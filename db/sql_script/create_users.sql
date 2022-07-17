CREATE TABLE Users (
	userid     VARCHAR(32)  NOT NULL PRIMARY KEY,
	username   VARCHAR(20)  NOT NULL UNIQUE,
	password   BINARY(64)   NOT NULL /* 32 bytes SHA256, 32 bytes salt*/
)
