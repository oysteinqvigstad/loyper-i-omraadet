import Toybox.Application;
import Toybox.Lang;
import Toybox.WatchUi;
import Toybox.Position;
import Toybox.System;
import Toybox.Graphics;

class MapTestApp extends Application.AppBase {

    function initialize() {
        AppBase.initialize();
    }

    function onStart(state as Dictionary?) as Void {
    }
    
    function onStop(state as Dictionary?) as Void {
    }

    function getInitialView() as Array<Views or InputDelegates>? {
        var names = ["Fottur", "Ski", "Sykkel", "Andre"];
        var ids = ["f", "s", "b", "o"];
        var resources = [Rez.Drawables.HikingIcon, Rez.Drawables.SkiingIcon, Rez.Drawables.BikingIcon, Rez.Drawables.OtherIcon];
        var menu = new WatchUi.Menu2({:title=>new $.DrawableMenuTitle()});
        for (var i = 0; i < names.size(); i++) {
            var bitmap = new WatchUi.Bitmap({:rezId=>resources[i], :locX=>WatchUi.LAYOUT_HALIGN_CENTER, :locY=>WatchUi.LAYOUT_VALIGN_CENTER});
            menu.addItem(new WatchUi.IconMenuItem(names[i], null, ids[i], bitmap, null));
        }
        return [menu, new ActivityMenuDelegate()];
    }


}

function getApp() as MapTestApp {
    return Application.getApp() as MapTestApp;
}


class DrawableMenuTitle extends WatchUi.Drawable {

    public function initialize() {
        Drawable.initialize({});
    }

    public function draw(dc as Dc) as Void {
        var spacing = 2;
        var appIcon = WatchUi.loadResource($.Rez.Drawables.LauncherIcon) as BitmapResource;
        var bitmapWidth = appIcon.getWidth();
        var labelWidth = dc.getTextWidthInPixels("Finn løype", Graphics.FONT_TINY);

        var bitmapX = (dc.getWidth() - (bitmapWidth + spacing + labelWidth)) / 2;
        var bitmapY = (dc.getHeight() - appIcon.getHeight()) / 2;
        var labelX = bitmapX + bitmapWidth + spacing;
        var labelY = dc.getHeight() / 2;

        dc.setColor(Graphics.COLOR_BLACK, Graphics.COLOR_BLACK);
        dc.clear();

        dc.drawBitmap(bitmapX, bitmapY, appIcon);
        dc.setColor(Graphics.COLOR_WHITE, Graphics.COLOR_TRANSPARENT);
        dc.drawText(labelX, labelY, Graphics.FONT_TINY, "Finn løype", Graphics.TEXT_JUSTIFY_LEFT | Graphics.TEXT_JUSTIFY_VCENTER);
    }
}
