#!/usr/bin/env python
import sys
import mapnik

# pathes
DATA_DIR="../../data"

# Some commonly used SRS definitions
SRS_900913 = "+proj=merc +a=6378137 +b=6378137 +lat_ts=0.0 +lon_0=0.0 +x_0=0.0 +y_0=0 +k=1.0 +units=m +nadgrids=@null +no_defs +over"
SRS_4326 = "+proj=longlat +datum=WGS84 +ellps=WGS84 +no_defs"

PALETTE = []
for i in range(0, 16):
    c = i * 17
    PALETTE.append(mapnik.Color(c, c, c, 255))
BLACK = PALETTE[0]
OCEAN_COLOR = PALETTE[12]
LAND_COLOR = PALETTE[15]

def add_layer(m, name, datasource, rules_dict, srs=SRS_900913):
    # create style
    s = mapnik.Style()
    for (expr, symbolizers) in rules_dict.items():
        r = mapnik.Rule()
        [r.symbols.append(symbolizer) for symbolizer in symbolizers]
        if expr is not None:
            r.filter = mapnik.Expression(expr)
        s.rules.append(r)

    m.append_style(name, s)

    # create layer
    layer = mapnik.Layer(name, srs)
    layer.datasource = datasource

    layer.styles.append(name)
    m.layers.append(layer)

def add_pg_layer(m, name, table, rules):
    add_layer(m, name, mapnik.PostGIS(dbname='iceland', table=table), rules, SRS_900913) 

# Create Map object
# Use 900913 SRC because Iceland looks square in it
m = mapnik.Map(1560, 1776, SRS_900913)
m.background = OCEAN_COLOR
# Set extent 
# Coordinates in 900913 were extracted from postgres DB with Iceland data using SQL
#    select ST_Extent(way) from osm_planet_roads
m.maximum_extent = mapnik.Box2d(-2784508.0, 9123829.0, -1449420.0, 10257366.0)

# Add Layers
# Coastline
add_layer(m, "coastline",
    mapnik.Shapefile(file='%s/world_boundries/10m_coastline.shp' % DATA_DIR),
    { None: [mapnik.PolygonSymbolizer(LAND_COLOR), 
             mapnik.LineSymbolizer(BLACK, 0.1)] 
    },
    srs = SRS_4326)

# Places
# TODO: use ShieldSymbolizer
add_pg_layer(m, "places", "planet_osm_point",
    { "[place]='city'": [mapnik.TextSymbolizer(mapnik.Expression("[name]"), 'DejaVu Sans Book', 12, BLACK),
                         mapnik.PointSymbolizer()],
      "[place]='town'": [mapnik.TextSymbolizer(mapnik.Expression("[name]"), 'DejaVu Sans Book', 10, BLACK),
                         mapnik.PointSymbolizer()],
      "[place]='village'": [mapnik.TextSymbolizer(mapnik.Expression("[name]"), 'DejaVu Sans Book', 8, BLACK),
                            mapnik.PointSymbolizer()],
    })

# 2. Roads
#add_pg_layer(m, "roads",
#    "planet_osm_roads",
#    [mapnik.LineSymbolizer(mapnik.Color('red'), 2)])


# Render map
m.zoom_all()

mapnik.render_to_file(m, 'map.png', 'png')
mapnik.save_map(m, "map.xml")

print "Done"
