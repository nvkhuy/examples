create or replace function count_elements(elements text[], to_find text[])
   returns bigint 
   language plpgsql
  as
$$
declare 
  element_count integer;

begin
 -- logic
  select count(*)
  into element_count
  from unnest(elements) element 
  where element = ANY(to_find);
  
  return element_count;
end;
$$