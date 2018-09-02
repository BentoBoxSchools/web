CREATE TABLE IF NOT EXISTS school (
    id INT(9) UNSIGNED NOT NULL AUTO_INCREMENT, # TODO: Make sure this starts at a non-zero value
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    link VARCHAR(255) NOT NULL,
    CONSTRAINT school_pk PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS donation_detail (
    id INT(9) UNSIGNED NOT NULL AUTO_INCREMENT, # TODO: Make sure this starts at a non-zero value
    school_id INT(9) NOT NULL,
    grade VARCHAR(255) NOT NULL,
    account_name VARCHAR(255) NOT NULL,
    balance DECIMAL(13,2) NOT NULL,
    CONSTRAINT donation_detail_pk PRIMARY KEY(id)
);