OSM2PGSQL=../../tools/osm2pgsql
MAP_FILE=../../data/iceland/iceland.osm.bz2

populatedb:
	echo ${OSM2PGSQL} -s -d iceland -S osm2pgsql.style ${MAP_FILE}

map.png map.xml: mkmap.py
	python mkmap.py


