-- API admin/data_analytics/product_trendings/group
-- CREATE
CREATE
MATERIALIZED VIEW product_trending_groups
AS
SELECT url,
       array_agg(id)                                       as product_trending_ids,
       max(images)                                         as images,
       max(name)                                           as name,
       max(domain)                                         as domain,
       max(category)                                       as category,
       max(sub_category)                                   as sub_category,
       max(price)                                          as price,
       COALESCE(false != ANY (array_agg(is_publish)), true) AS is_publish,
       max(created_at)                                     as created_at
FROM product_trendings pt
GROUP BY pt.url
ORDER BY created_at DESC
WITH NO DATA;

-- REFRESH MATERIALIZED VIEW product_trending_groups;
-- DROP MATERIALIZED VIEW product_trending_groups;