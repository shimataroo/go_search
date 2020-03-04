CREATE SCHEMA IF NOT EXISTS gosearch DEFAULT CHARACTER SET utf8mb4;

USE gosearch;

CREATE TABLE IF NOT EXISTS documents(
    PRIMARY KEY (document_id),
    document_id INT UNSIGNED AUTO_INCREMENT NOT NULL,
    document_title TEXT                     NOT NULL
)ENGINE=InnoDB DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_bin;