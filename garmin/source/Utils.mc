import Toybox.Position;
import Toybox.Lang;
import Toybox.Math;
import Toybox.System;
import Toybox.Graphics;
import Toybox.WatchUi;

module Utils {

    function calculateVisibleArea(pos as Position.Location, radiusInKm as Float) as Array<Position.Location> {
            var deltaSquareDegrees = radiusInKm / 111.0;
            var upperLeft = new Position.Location({
                :latitude => (pos.toDegrees()[0] + deltaSquareDegrees),
                :longitude => (pos.toDegrees()[1] - deltaSquareDegrees),
                :format => :degrees
            });
            var lowerRight = new Position.Location({
                :latitude => pos.toDegrees()[0] - deltaSquareDegrees,
                :longitude => pos.toDegrees()[1] + deltaSquareDegrees,
                :format => :degrees
            });
            return [upperLeft, lowerRight];
    }

    const MAX_POLYLINE_OBJECT_COUNT = 233;
    function DecodePolyline(polyline) {
            polyline = polyline.toCharArray();
            var len = polyline.size();
            var indexJump = Math.ceil(len / MAX_POLYLINE_OBJECT_COUNT);
            var poly = [];
            var index = 0;
            var lat = 0;
            var lng = 0;
            var skipIndex = 0;

            while (index < len) {

                var byte = 0;
                var shift = 0;
                var result = 0;

                do {
                    byte = polyline[index].toNumber() - 63;
                    result = result | ((byte & 31) << shift);
                    shift += 5;
                    index++;
                } while (byte >= 32 && index < len);

                var dlat = ((result & 1) ? ~(result >> 1) : (result >> 1));

                shift = 0;
                result = 0;

                do {
                    byte = polyline[index].toNumber() - 63;
                    result = result | ((byte & 31) << shift);
                    shift += 5;
                    index++;
                } while (byte >= 32 && index < len);

                var dlng = ((result & 1) ? ~(result >> 1) : (result >> 1));

                lat += dlat;
                lng += dlng;

                if (indexJump == 0 || skipIndex % indexJump == 0) {
                    var p = new Position.Location({:latitude=>lat/1e5, :longitude=>lng/1e5, :format=>:degrees});
                    poly.add(p);
                }

                skipIndex++;
            }
            return poly;
        }

}