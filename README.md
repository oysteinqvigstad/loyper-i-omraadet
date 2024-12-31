# Løyper i Området

## About

A Garmin app for finding nearby trails based on proximity (limited to Norway)

The number of trails in the database is:
- Foot trails: 8657
- Bike trails: 353
- Ski trails: 1089
- Misc trails: 628

The majority of the project time was spent on data processing for generating the trails.

To overcome compute and data capacity (1MB) limitations on Garmin devices, the app relies on RESTful API calls to a backend to retrieve nearby trails by GPS coordinates. For non-cellular Garmin devices, these runtime request are automatically Bluetooth proxied through the paired Garmin Connect app on your mobile phone.



# Trails generations

The raw data comes from [Turrutebasen](https://kartkatalog.geonorge.no/metadata/turrutebasen/d1422d17-6d95-4ef1-96ab-8af31744dd63) made avilable by Geonorge. There are two main problems with how the data is encoded:
1. To avoid redundent data due to overlapping trails, the data is organized in trail segments, and then linked to trails that use them. Moreover, the segments are unordered and undirectional.
2. Garmin's MapView only accepts a single polyline (sequential GPS coordinates), but most trails aren't linear in reality. They may have alternative pathways (shortcuts, etc), with multiple starts/exits along the way.

To overcome the first problem, an algorithm will start with a random segment, and then greedily splice the nearest segment, considering either endpoints of both segments to be sliced (meaning "rotating" to find the best splices).

To overcome the second problem, when there's a nearer coordinate than the ends of the current polyline/splice, the path must backtrack through the same coordinates to maintain a single sequential line. Moreover, to conserve bandwidth when transferring trails, the number of coordinates are reduced, but this creates some problems with backtracking if not handled correctly, so several iterations of reducing/interpolating coordinates are needed.

These transformations are performed only once, and handled by (entrypoint `/backend/cmd/xml_parser.go`). The results are then stored in JSON, at location `/backend/resources/trails.jon`.
