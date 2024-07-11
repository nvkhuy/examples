create or replace function get_avg_growth_rate(purl character varying, field_name character varying) returns double precision
    language plpgsql
as
$$
DECLARE
values          float[];
    avg_growth_rate float;
BEGIN
SELECT ARRAY_AGG(field_value)
INTO values
FROM (SELECT avg(CASE
    WHEN field_name = 'price' THEN price
    WHEN field_name = 'sold' THEN sold
    WHEN field_name = 'stock' THEN stock
    ELSE NULL
    END) as field_value
    FROM product_changes ps
    WHERE url = purl
    GROUP BY url, scrape_date
    ORDER BY scrape_date desc) sub;

avg_growth_rate := calculate_avg_growth_rate(values);
RETURN avg_growth_rate;
END;
$$;

alter function get_avg_growth_rate(varchar, varchar) owner to inflow;

