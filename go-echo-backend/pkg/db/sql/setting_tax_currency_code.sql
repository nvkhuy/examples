UPDATE setting_taxes
SET currency_code = (
    CASE
        WHEN country_code = 'RU' THEN 'RUB'
        WHEN country_code = 'SG' THEN 'SGD'
        WHEN country_code = 'US' THEN 'USD'
        WHEN country_code = 'VN' THEN 'VND'
        END
    )
WHERE country_code in ('RU', 'SG', 'US', 'VN');