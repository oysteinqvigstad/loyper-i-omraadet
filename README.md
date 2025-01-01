# Løyper i Området

## About

A Garmin Connect IQ app that helps you discover nearby trails across Norway, directly from your Garmin device.

### Features
Find nearby trails based on your current GPS location:
- Hiking trails (8,657 routes)
- Biking trails (353 routes)
- Ski trails (1,089 routes)
- Miscellaneous trails (628 routes)

### Technical Implementation
Due to Garmin devices' 1MB storage capacity limit, the app uses a RESTful API to fetch nearby trails based on GPS coordinates. For non-cellular Garmin devices, requests are proxied through the paired Garmin Connect mobile app via Bluetooth.

Below is a demonstration of the app running in low-resolution mode. While this mode reduces trail accuracy, it provides enhanced performance and faster response times:

https://github.com/user-attachments/assets/3255087e-540f-40b6-967d-5c8a718937e5

The project served as a valuable learning experience with:
- Garmin's Connect IQ platform and its Monkey C programming language
- Geographic data processing and trail mapping algorithms
- Working within hardware and API constraints

## Project Status
While the core functionality works, development was discontinued after 2-3 weeks due to several practical limitations:

- MapView with polylines support is limited to newer Garmin models
- Server infrastructure would require ongoing maintenance costs
- Quality control for generated trails would be too time-intensive
- Privacy concerns regarding GPS coordinate transmission


## Trails generations

The raw data is sourced from [Turrutebasen](https://kartkatalog.geonorge.no/metadata/turrutebasen/d1422d17-6d95-4ef1-96ab-8af31744dd63) a national trail database made available by Geonorge. Two main challenges were addressed during data processing:

#### Challenge 1
The source data is organized in unordered, undirectional trail segments to avoid redundancy from overlapping trails. Each trail is defined by references to these shared segments rather than complete routes.
#### Challenge 2
Garmin's MapView only supports single polylines (sequential GPS coordinates), while real trails often have:
- Multiple entry/exit points
- Alternative pathways and shortcuts
- Non-linear route structures

#### Solution Implementation
1. Segment Assembly:
- Starts with a random segment
- Greedily splices the nearest segment by evaluating all possible endpoint combinations
- "Rotates" segments to achieve optimal splicing
2. Path Processing:
- Implements backtracking when needed to maintain a single sequential line
- Reduces coordinate density for bandwidth optimization
- Uses multiple iterations of reduction/interpolation to preserve path accuracy

The transformation process is handled by /backend/cmd/xml_parser.go, with the processed trails stored in /backend/resources/trails.json.
