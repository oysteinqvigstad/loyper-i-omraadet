import Toybox.Lang;
import Toybox.WatchUi;
import Toybox.System;
import Toybox.Graphics;
import Toybox.Position;
import Toybox.Weather;

class ActivityMenuDelegate extends WatchUi.Menu2InputDelegate {
    private var _type = "";
    private var _locationInfo as Position.Info?;

    function initialize() {
        Menu2InputDelegate.initialize();
    }

    public function onSelect(item as MenuItem) as Void {
        _type = item.getId();
        WatchUi.pushView(new WatchUi.ProgressBar("Venter på GPS signal...", null), new WatchUi.BehaviorDelegate(), WatchUi.SLIDE_IMMEDIATE);
        if (isLastKnownLocationValid()) {
            onLocationRetrieved(_locationInfo);
        } else {
            Position.enableLocationEvents(Position.LOCATION_ONE_SHOT, method(:onLocationRetrieved));
        }

    }

    public function onLocationRetrieved(info as Position.Info) as Void {
        if (info.accuracy >= Position.QUALITY_POOR) {
            _locationInfo = info;
        }

        WatchUi.switchToView(new WatchUi.ProgressBar("Henter løyper...", null), new WatchUi.BehaviorDelegate(), WatchUi.SLIDE_IMMEDIATE);
        var location = info.position.toDegrees();
        var json = new JsonTransaction(method(:onHttpResponse));
        json.makeRequest(_type, location[0], location[1]);
    }

    public function onHttpResponse(data as JsonTransaction.JsonResponse) as Void {
        switch (data) {
            case instanceof Array:
                if (data.size() > 0) {
                    openSubmenu(data);
                } else {
                    WatchUi.switchToView(new WatchUi.ProgressBar("Ingen resultater", 100.0), new WatchUi.BehaviorDelegate(), WatchUi.SLIDE_IMMEDIATE); 
                }
                break;
            default:
                System.println("oopsie");
        } 
    }

    public function openSubmenu(data as JsonTransaction.JsonResponse) as Void {
        var menu = new WatchUi.Menu2({:title=>new $.DrawableMenuTitle()});
        for (var i = 0; i < data.size(); i++) {
            menu.addItem(new WatchUi.MenuItem(
                data[i]["name"] as String, 
                data[i]["DistanceFromUser"].format("%.2f") + " km unna", 
                data[i]["id"] as String, 
                null));
        }
        WatchUi.switchToView(menu, new TrailMenuDelegate(data), WatchUi.SLIDE_UP);
        
    }

    function isLastKnownLocationValid() as Boolean {
        var timeNow = Time.getCurrentTime({:currentTimeType => Time.CURRENT_TIME_RTC});
        var secondsInEightMinutes = 60 * 8;
        if (_locationInfo != null && _locationInfo.when != null) {
            System.println(timeNow.compare(_locationInfo.when) < secondsInEightMinutes);
        }

        return (_locationInfo != null && _locationInfo has :when && 
                timeNow.compare(_locationInfo.when) < secondsInEightMinutes);

    }

    public function onDone() as Void {
        WatchUi.popView(WatchUi.SLIDE_DOWN);
    }
}