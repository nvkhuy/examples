-- API admin/data_analytics/products/group
-- CREATE
CREATE MATERIALIZED VIEW product_groups AS
SELECT
    url,
    max(id) as id,
    max(domain) as domain,
    max(country_code) as country_code,
    max(name) as name,
    max(overall_growth_rate) as overall_growth_rate,
    max(images) as images,
    max(category) as category,
    max(sub_category) as sub_category,
    max(price) as price,
    max(created_at) as created_at,
    max(updated_at) as updated_at
FROM
    products p
GROUP BY
    p.url
ORDER BY
    created_at DESC
WITH
    NO DATA;

-- REFRESH MATERIALIZED VIEW product_groups;
-- DROP MATERIALIZED VIEW product_groups;