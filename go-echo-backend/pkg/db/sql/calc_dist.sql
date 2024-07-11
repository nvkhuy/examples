create or replace function calc_distance_in_miles(lat1 decimal,lng1 decimal,lat2 decimal,lng2 decimal)
   returns bigint 
   language plpgsql
  as
$$
declare 
  distance decimal;

begin
 -- logic
  distance = 3958.8 * ACOS(((COS(((PI() / 2) - RADIANS((90.0 - lat1)))) *
                COS(PI() / 2 - RADIANS(90.0 - lat2)) *
                COS((RADIANS(lng1) - RADIANS(lng2))))
        + (SIN(((PI() / 2) - RADIANS((90.0 - lat1)))) *
          SIN(((PI() / 2) - RADIANS(90.0 - lat2))))))
  
  return distance;
end;
$$