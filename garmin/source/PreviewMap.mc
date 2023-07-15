import Toybox.WatchUi;
import Toybox.Graphics;
import Toybox.Position;
import Toybox.Lang;
using Toybox.Graphics;

class PreviewMap extends WatchUi.MapView {
    function initialize(coordinates as String, pos as Position.Location, radius as Float) {
        MapView.initialize();
        setMapMode(MAP_MODE_BROWSE);
        var area = Utils.calculateVisibleArea(pos, radius);
        var polyline = new WatchUi.MapPolyline();
        polyline.addLocation(Utils.DecodePolyline(coordinates));
        polyline.setColor(Graphics.COLOR_RED);
        polyline.setWidth(4);
        setPolyline(polyline);
        setMapVisibleArea(area[0], area[1]);
        setScreenVisibleArea(0, 0, System.getDeviceSettings().screenWidth, System.getDeviceSettings().screenHeight);

    }
}