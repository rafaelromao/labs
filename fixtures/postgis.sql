-- All booleans must be true. Geometry is planar coordinates; geography is explicit.
WITH shape AS (
 SELECT ST_GeomFromText('POLYGON((0 0,10 0,10 10,0 10,0 0),(3 3,3 7,7 7,7 3,3 3))') AS p
)
SELECT ST_IsValid(p)
 AND ST_Covers(p,ST_Point(0,5))
 AND NOT ST_Contains(p,ST_Point(0,5))
 AND NOT ST_Covers(p,ST_Point(5,5))
 AND ST_Covers(p,ST_Point(3,5))
 AND ST_DWithin(ST_SetSRID(ST_Point(179.9,0),4326)::geography,
                ST_SetSRID(ST_Point(-179.9,0),4326)::geography,23000)
FROM shape;
