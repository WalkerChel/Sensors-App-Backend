CREATE OR REPLACE PROCEDURE create_region(
  Name text
)
LANGUAGE plpgsql
AS $function$
BEGIN
INSERT INTO regions(name) VALUES(name);
END;
$function$;

-- europe: western europe, eastern europe, northern europe, southern europe, central europe.
-- asia: east asia, southeast asia, south asia, central asia, western asia.
-- africa: north africa, west africa, east africa, central africa, southern africa.
-- americas: north america, latin america (central and south america), the caribbean region.
-- australia and oceania: australia, melanesia, polynesia, micronesia.

CALL create_region(
  Name := 'western europe');

CALL create_region(
  Name := 'eastern europe');

CALL create_region(
  Name := 'northern europe');

CALL create_region(
  Name := 'southern europe');

CALL create_region(
  Name := 'central europe');

CALL create_region(
  Name := 'east asia');

CALL create_region(
  Name := 'southeast asia');

CALL create_region(
  Name := 'south asia');

CALL create_region(
  Name := 'central asia');

CALL create_region(
  Name := 'western asia');

CALL create_region(
  Name := 'north africa');

CALL create_region(
  Name := 'west africa');

CALL create_region(
  Name := 'east africa');

CALL create_region(
  Name := 'central africa');

CALL create_region(
  Name := 'southern africa');

CALL create_region(
  Name := 'north america');

CALL create_region(
  Name := 'central america');

CALL create_region(
  Name := 'south america');

CALL create_region(
  Name := 'australia');

CALL create_region(
  Name := 'melanesia');

CALL create_region(
  Name := 'polynesia');

CALL create_region(
  Name := 'micronesia');
  