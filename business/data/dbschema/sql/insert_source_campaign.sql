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