CREATE OR REPLACE FUNCTION clean_categories_hierarchical() RETURNS VOID AS $$
BEGIN
    update categories
    set parent_category_id = ''
    where id in (select b.id
                 from categories a
                          join categories b ON a.parent_category_id = b.id
                 where b.parent_category_id != ''
    group by b.id);
END;
$$ LANGUAGE plpgsql;