SELECT
    c.id,
    c.name
FROM
    campaigns c
        LEFT JOIN
    source_campaign sc
    ON
            c.id = sc.campaign_id
WHERE
    sc.source_id IS NULL;