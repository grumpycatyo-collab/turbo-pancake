SELECT
    s.id,
    s.name,
    COUNT(sc.campaign_id) AS campaign_count
FROM
    sources s
        LEFT JOIN
    source_campaign sc
    ON
            s.id = sc.source_id
GROUP BY
    s.id
ORDER BY
    campaign_count DESC
LIMIT 5;