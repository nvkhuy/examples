create or replace function calculate_avg_growth_rate(series double precision[]) returns double precision
    language plpgsql
as
$$
DECLARE
total       float   := 0;
    i           integer := 0;
    growth_rate float[];
BEGIN
    IF array_length(series, 1) IS NOT NULL THEN
        growth_rate := ARRAY [1.0];
FOR i IN 2..array_length(series, 1)
            LOOP
                IF series[i - 1] <> 0 THEN
                    growth_rate := array_append(growth_rate, (series[i] - series[i - 1]) / series[i - 1]);
                    total := total + growth_rate[i];
END IF;
END LOOP;
        IF array_length(growth_rate, 1) - 1 <> 0 THEN
            RETURN total / (array_length(growth_rate, 1) - 1);
ELSE
            RETURN 0;
END IF;
ELSE
        RETURN 0;
END IF;
END;
$$;

alter function calculate_avg_growth_rate(double precision[]) owner to inflow;