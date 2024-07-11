create or replace function array_unique (a text[]) returns text[] as $$
  select array (
    select distinct v from unnest(a) as b(v)
  )
$$ language sql;