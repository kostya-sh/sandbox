#!/usr/bin/env python
import sys
import mapnik

# pathes
DATA_DIR="../../data"

# Some commonly used SRS definitions
SRS_900913 = "+proj=merc +a=6378137 +b=6378137 +lat_ts=0.0 +lon_0=0.0 +x_0=0.0 +y_0=0 +k=1.0 +units=m +nadgrids=@null +no_defs +over"
SRS_4326 = "+proj=longlat +datum=WGS84 +ellps=WGS84 +no_defs"

def add_layer(m, name, datasource, symbols, srs = SRS_900913):
    # create style
    s = mapnik.Style()
    r = mapnik.Rule()
    [r.symbols.append(symbol) for symbol in symbols]
    s.rules.append(r)
    m.append_style(name, s)

    # create layer
    layer = mapnik.Layer(name, srs)
    layer.datasource = datasource

    layer.styles.append(name)
    m.layers.append(layer)

def add_pg_layer(m, name, table, symbols):
    add_layer(m, name, mapnik.PostGIS(dbname='iceland', table=table), symbols, SRS_900913) 

# Create Map object
# Use 900913 SRC because Iceland looks square in it
m = mapnik.Map(2400, 2400, SRS_900913)
m.background = mapnik.Color('blue')
# Set extent 
# Coordinates in 900913 were extracted from postgres DB with Iceland data using SQL
#    select ST_Extent(way) from osm_planet_roads
m.maximum_extent = mapnik.Box2d(-2784508.0, 9123829.0, -1449420.0, 10257366.0)

# Add Layers
# 1. Coastline
add_layer(m, "coastline",
    mapnik.Shapefile(file='%s/world_boundries/10m_coastline.shp' % DATA_DIR),
    [mapnik.PolygonSymbolizer(mapnik.Color('white')),
     mapnik.LineSymbolizer(mapnik.Color('black'), 1)],
    srs = SRS_4326)

# 2. Roads
add_pg_layer(m, "roads",
    "planet_osm_roads",
    [mapnik.LineSymbolizer(mapnik.Color('red'), 2)])


# Render map
m.zoom_all()

mapnik.render_to_file(m, 'map.png', 'png')
mapnik.save_map(m,"map.xml")

print "Done"
