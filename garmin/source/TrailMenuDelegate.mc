import Toybox.Graphics;
import Toybox.Lang;
import Toybox.WatchUi;
import Toybox.System;

class TrailMenuDelegate extends WatchUi.Menu2InputDelegate {
    private var _data as JsonTransaction.JsonResponse;

    public function initialize(data as JsonTransaction.JsonResponse) {
        Menu2InputDelegate.initialize();
        self._data = data;
    }

    public function onSelect(item as MenuItem) as Void {
        var index = 0;
        for (var i = 0; i < _data.size(); i++) {
            if (item.getId() == _data[i]["id"]) {
                index = i;
                break;
            }
        }
        var coords = new Position.Location({
            :latitude => _data[index]["location"]["latitude"].toDouble(),
            :longitude => _data[index]["location"]["longitude"].toDouble(),
            :format => :degrees
        });
        var view = new PreviewMap(_data[index]["polyline_detail_medium"], coords, 2.0);
        WatchUi.pushView(view, null, WatchUi.SLIDE_UP);

    }

    public function onBack() as Void {
        WatchUi.popView(WatchUi.SLIDE_DOWN);
    }
}