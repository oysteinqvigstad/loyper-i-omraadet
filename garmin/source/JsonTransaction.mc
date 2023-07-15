import Toybox.System;
import Toybox.WatchUi;
import Toybox.Communications;
import Toybox.Lang;



class JsonTransaction extends WatchUi.BehaviorDelegate {
    typedef JsonResponse as Dictionary or String or Null;
    private var _notify as Method(args as JsonResponse) as Void;

    public function initialize(handler as Method(args as JsonResponse) as Void) {
        WatchUi.BehaviorDelegate.initialize();
        _notify = handler;
    }


    function onReceive(responseCode as Number, data as JsonResponse) as Void {
        if (responseCode == 200) {
            System.println("Request successful");
        } else {
            System.println("Response: " + responseCode as Number);
            System.println(data);
        }
        _notify.invoke(data);
    }

    function makeRequest(type as String, lat as Double, lon as Double) as Void {
        var url = "https://loyper-i-omraadet.onrender.com/near/";
        // var url = "http://localhost:8080/near/";
        var params = {                                              
            "c" => lat.toString() + "," + lon.toString(),
            "t" => type,
            "l" => "1-8",
            "r" => "50"
        };
        System.println(lat.toString() + "," + lon.toString());
        var options = {                                             
            :method => Communications.HTTP_REQUEST_METHOD_GET,      
        };
        Communications.makeWebRequest(url, params, options, method(:onReceive));
    }

}
