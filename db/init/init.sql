CREATE USER 'myuser' IDENTIFIED BY 'dev000000';

CREATE DATABASE user_api;
USE user_api;

CREATE TABLE Users (
	username   VARCHAR(20)  NOT NULL PRIMARY KEY,
	password   BINARY(64)   NOT NULL /* 32 bytes SHA256, 32 bytes salt*/
);

CREATE TABLE Subdomains (
	subdomain   VARCHAR(20) NOT NULL PRIMARY KEY,
	username    VARCHAR(20) NOT NULL,
	FOREIGN KEY (username) REFERENCES Users(username) ON DELETE CASCADE
);

GRANT ALL PRIVILEGES ON *.* TO 'myuser';
