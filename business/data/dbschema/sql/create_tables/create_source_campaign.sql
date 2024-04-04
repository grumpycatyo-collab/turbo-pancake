CREATE TABLE IF NOT EXISTS source_campaign (
                                               source_id INT,
                                               campaign_id INT,
                                               FOREIGN KEY (source_id) REFERENCES sources(id),
                                               FOREIGN KEY (campaign_id) REFERENCES campaigns(id)
);