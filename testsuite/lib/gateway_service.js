
var request = require("request")
var mlog = require("mocha-logger");

var httpResp = require('./http_resp.js');

function GatewayService(host, port) {
    this.host = host;
    this.port = port;
    this.baseUri = "http://" + this.host + ":" + this.port;
    this.pathMap = {
        "github": "github_webhook"
    };
}

GatewayService.prototype.sendEvent = function (scm, event, cb) {
    var options = {
        url: this.baseUri + "/" + this.pathMap[scm],
        headers: {
            "User-Agent": "javascript-request",
            "X-GitHub-Event": "push",
            "X-GitHub-Delivery": createUUID(),
            "Content-Type": "application/json"
        },
        body: JSON.stringify(event)
    };
    mlog.log("POST " + options.url);
    request.post(options, function (error, response, body) {
        var obj = {};
        //  Note: http client request is sent with a Expect: 100 header when receiving a big payload
        //        server side will first respond with a 100 Continue
        httpResp(error, response, body, obj, response.statusCode >= 100 && response.statusCode < 300, cb);
    });
}

function createUUID() {
    // http://www.ietf.org/rfc/rfc4122.txt
    var s = [];
    var hexDigits = "0123456789abcdef";
    for (var i = 0; i < 36; i++) {
        s[i] = hexDigits.substr(Math.floor(Math.random() * 0x10), 1);
    }
    s[14] = "4";  // bits 12-15 of the time_hi_and_version field to 0010
    s[19] = hexDigits.substr((s[19] & 0x3) | 0x8, 1);  // bits 6-7 of the clock_seq_hi_and_reserved to 01
    s[8] = s[13] = s[18] = s[23] = "-";

    var uuid = s.join("");
    return uuid;
}

module.exports = GatewayService;
