CREATE OR REPLACE PROCEDURE generate_sensors()
LANGUAGE plpgsql
AS $function$
DECLARE
    regionsNumber INT;
    regId INT;
    n INT;
BEGIN
  SELECT COUNT(DISTINCT id)
  INTO regionsNumber
  FROM regions;

    FOR regId IN 1..regionsNumber LOOP
        FOR n IN 1..floor(random() * regionsNumber + 1)::INT LOOP
            INSERT INTO sensors (Region_Id, Name, Longitude, Latitude) 
            VALUES (
                regId, 
                substr(md5(random()::text), 1, 12), 
                random() * 180 - 90, 
                random() * 360 - 180
            );
        END LOOP;
    END LOOP;
END;
$function$;

CALL generate_sensors();
