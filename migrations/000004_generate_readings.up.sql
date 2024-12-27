CREATE OR REPLACE PROCEDURE generate_readings()
LANGUAGE plpgsql
AS $function$
DECLARE
    sensNumber INT;
    sensId INT;
    n INT;
BEGIN
  SELECT COUNT(DISTINCT id)
  INTO sensNumber
  FROM sensors;

    FOR sensId IN 1..sensNumber LOOP
        FOR n IN 1..floor(random() * sensNumber + 1)::INT LOOP
            INSERT INTO readings(sensorId, temperature, createdAt) 
            VALUES(
              sensId, 
              random() * 90 - 30,
              NOW() - INTERVAL '1 year' * random()
            );
        END LOOP;
    END LOOP;
END;
$function$;

CALL generate_readings();
