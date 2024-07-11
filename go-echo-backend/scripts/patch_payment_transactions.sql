-- Patch payment transactions with Paid status --
DROP TABLE IF EXISTS temp_table;
CREATE
TEMPORARY TABLE temp_table AS
SELECT b.id              as bulk_id,
       p.id              as pm_id,
       b.tracking_status as bulk_tracking_status,
       p.status          as pm_status,
       p.milestone       as pm_milestone,
       p.total_amount,
       p.paid_amount,
       b.created_at,
       b.user_id
FROM bulk_purchase_orders b
         INNER JOIN payment_transactions p
                    ON b.id = p.bulk_purchase_order_id
                        AND (b.tracking_status IN
                             ('final_payment_confirmed', 'delivering', 'delivery_confirmed', 'delivered'))
                        AND p.status != 'paid'
ORDER BY created_at DESC;

update payment_transactions p
set status='paid'
where p.id in (select pm_id from temp_table);