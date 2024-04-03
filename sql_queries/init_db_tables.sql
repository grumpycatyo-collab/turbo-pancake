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
INSERT INTO sources (name)
SELECT CONCAT('Source_', FLOOR(RAND() * 1000))
FROM information_schema.tables
LIMIT 100;

INSERT INTO campaigns (name)
SELECT CONCAT('Campaign_', FLOOR(RAND() * 1000))
FROM information_schema.tables
LIMIT 100;

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
