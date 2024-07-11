CREATE OR REPLACE FUNCTION insert_casbin_rules() RETURNS VOID AS $$
BEGIN
    -- Truncate
    TRUNCATE casbin_rule RESTART IDENTITY;

    -- Group 1, super_group
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/*', 'super_group');
    INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'super_admin', 'super_group', '/*');
    INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'client', 'super_group', '/*');
    INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'seller', 'super_group', '/*');
    INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'leader:dev', 'super_group', '/*');
    INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'staff:dev', 'super_group', '/*');
    INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'leader:operator', 'super_group', '/*');

    -- Group 2, lead_sales_group
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/sales/*', 'lead_sales_group');
    INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'leader:sales', 'lead_sales_group', '/*');
    INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'staff:sales', 'lead_sales_group', 'GET');

    -- Group 3, staff_sales_write_group
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/sales/write-but-for-staff', 'staff_sales_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'staff:sales', 'staff_sales_write_group', 'PUT');
    INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'staff:sales', 'staff_sales_write_group', 'PATCH');

    -- Group 4, marketing_write_group
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/purchase_orders', 'marketing_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/purchase_orders/*', 'marketing_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/categories', 'marketing_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/categories/*', 'marketing_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/shops', 'marketing_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/shops/*', 'marketing_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/pages', 'marketing_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/pages/*', 'marketing_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/posts', 'marketing_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/posts/*', 'marketing_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/inquiries', 'marketing_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/inquiries/*', 'marketing_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/products', 'marketing_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/products/*', 'marketing_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/collections', 'marketing_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/collections/*', 'marketing_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/blog', 'marketing_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/blog/*', 'marketing_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/settings', 'marketing_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/settings/*', 'marketing_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/bulk_purchase_orders', 'marketing_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/bulk_purchase_orders/*', 'marketing_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/users/search', 'marketing_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/users/search/*', 'marketing_write_group');

    INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'leader:marketing', 'marketing_write_group', '/*');
    INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'staff:marketing', 'marketing_write_group', '/*');

    -- Group 5, sales_write_group
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/inquiries', 'sales_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/inquiries/*', 'sales_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/purchase_orders', 'sales_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/purchase_orders/*', 'sales_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/settings', 'sales_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/settings/*', 'sales_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/users/search', 'sales_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/users/search/*', 'sales_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/bulk_purchase_orders', 'sales_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/bulk_purchase_orders/*', 'sales_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/payment_transactions', 'sales_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/payment_transactions/*', 'sales_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'leader:sales', 'sales_write_group', '/*');
    INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'staff:sales', 'sales_write_group', '/*');

    -- Group 6, leader_write_group
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/users', 'leader_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/users/*', 'leader_write_group');

    INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'leader:sales', 'leader_write_group', '/*');
    INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'leader:finance', 'leader_write_group', '/*');
    INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'leader:marketing', 'leader_write_group', '/*');
    INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'leader:designer', 'leader_write_group', '/*');
    INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'leader:customer_service', 'leader_write_group', '/*');
    INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'leader:operator', 'leader_write_group', '/*');

    -- Group 7, finance_write_group
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/inquiries', 'finance_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/inquiries/*', 'finance_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/purchase_orders', 'finance_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/purchase_orders/*', 'finance_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/settings', 'finance_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/settings/*', 'finance_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/payment_transactions', 'finance_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/payment_transactions/*', 'finance_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/bulk_purchase_orders', 'finance_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/bulk_purchase_orders/*', 'finance_write_group');

    INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'leader:finance', 'finance_write_group', '/*');
    INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'staff:finance', 'finance_write_group', '/*');

    -- Group 8, lead_operator_write_group
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/inquiries', 'lead_operator_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/inquiries/*', 'lead_operator_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/purchase_orders', 'lead_operator_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/purchase_orders/*', 'lead_operator_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/payment_transactions', 'lead_operator_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/payment_transactions/*', 'lead_operator_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/categories/get_category_tree', 'lead_operator_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/categories/get_category_tree/*', 'lead_operator_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/shops', 'lead_operator_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/shops/*', 'lead_operator_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/pages', 'lead_operator_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/pages/*', 'lead_operator_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/posts', 'lead_operator_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/posts/*', 'lead_operator_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/inquiries', 'lead_operator_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/inquiries/*', 'lead_operator_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/products', 'lead_operator_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/products/*', 'lead_operator_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/collections', 'lead_operator_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/collections/*', 'lead_operator_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/settings/sizes', 'lead_operator_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/settings/sizes/*', 'lead_operator_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/settings/taxes', 'lead_operator_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/settings/taxes/*', 'lead_operator_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/settings/banks', 'lead_operator_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/settings/banks/*', 'lead_operator_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/notifications', 'lead_operator_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/notifications/*', 'lead_operator_write_group');

    INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'leader:operator', 'lead_operator_write_group', '/*');
    INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', 'staff:operator', 'lead_operator_write_group', '/*');

    -- Group 9, all_read_group
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/notifications', 'all_read_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/notifications/*', 'all_read_group');

    INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', '*', 'all_read_group', 'GET');
    INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', '*', 'all_read_group', 'OPTIONS');

    -- Group 10, all_write_group
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/common', 'all_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/common/*', 'all_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/notifications', 'all_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/notifications/*', 'all_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/me', 'all_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/me/*', 'all_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/resources', 'all_write_group');
    INSERT INTO casbin_rule (ptype, v0, v1) VALUES ('g2', '/api/v1/admin/resources/*', 'all_write_group');

    INSERT INTO casbin_rule (ptype, v0, v1, v2) VALUES ('p', '*', 'all_write_group', '/*');
END;
$$ LANGUAGE plpgsql;
