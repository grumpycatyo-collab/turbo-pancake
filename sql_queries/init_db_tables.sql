# Creating tables
CREATE TABLE sources (
                         id INT AUTO_INCREMENT PRIMARY KEY,
                         name VARCHAR(255)
);

CREATE TABLE campaigns (
                           id INT AUTO_INCREMENT PRIMARY KEY,
                           name VARCHAR(255)
);

CREATE TABLE source_campaign (
                                 source_id INT,
                                 campaign_id INT,
                                 FOREIGN KEY (source_id) REFERENCES sources(id),
                                 FOREIGN KEY (campaign_id) REFERENCES campaigns(id)
);

# ======================================================================================================================
# Random data seeding

-- For sources table
DELIMITER //
CREATE PROCEDURE insert_sources()
BEGIN
    DECLARE i INT DEFAULT 0;
    WHILE i < 100 DO
            INSERT INTO sources (name) VALUES (CONCAT('Source_', FLOOR(RAND() * 1000)));
            SET i = i + 1;
        END WHILE;
END //
DELIMITER ;
CALL insert_sources();

-- For campaigns table
DELIMITER //
CREATE PROCEDURE insert_campaigns()
BEGIN
    DECLARE i INT DEFAULT 0;
    WHILE i < 100 DO
            INSERT INTO campaigns (name) VALUES (CONCAT('Campaign_', FLOOR(RAND() * 1000)));
            SET i = i + 1;
        END WHILE;
END //
DELIMITER ;
CALL insert_campaigns();

# ======================================================================================================================
# Inserting connections between sources and campaign
INSERT INTO source_campaign (source_id, campaign_id)
SELECT source_id, campaign_id
FROM
    (SELECT sources.id AS source_id,
            campaigns.id AS campaign_id,
            ROW_NUMBER() OVER (PARTITION BY sources.id) AS rn
     FROM sources
              JOIN campaigns
     WHERE RAND() < 0.1) AS tmp
WHERE rn <= 10;