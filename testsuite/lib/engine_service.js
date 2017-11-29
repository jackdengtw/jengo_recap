
var request = require("request")
var mlog = require("mocha-logger");

var httpResp = require('./http_resp.js');

function EngineService(host, port) {
    this.host = host;
    this.port = port;
    this.baseUri = "http://" + this.host + ":" + this.port;
}

EngineService.prototype.getRuns = function (userId, projectId, runId, eventId, cb) {
    var options = {
        url: this.baseUri + "/v0.1/runs?" + "userId=" +
        userId + "&projectId=" + projectId + "&runId=" + runId + "&event_id=" + eventId,
        headers: {
            "User-Agent": "javascript-request",
        },
    };
    mlog.log("GET " + options.url);
    request(options, function (error, response, body) {
        var resp = [];
        httpResp(error, response, body, resp, response.statusCode >= 200 && response.statusCode < 300, cb);
    });
}

EngineService.prototype.getRealtimeLog = function (runId, offset, limit, cb) {
    var options = {
        url: this.baseUri + "/v0.1/run/" + runId + "/run_log?offset=" + offset + "&limit=" + limit,
        headers: {
            "User-Agent": "javascript-request",
        },
        gzip: true
    };
    mlog.log("GET " + options.url);
    request(options, function (error, response, body) {
        if (!error && response.statusCode >= 200 && response.statusCode < 300) {
            mlog.log("Headers:" + JSON.stringify(response.headers));
            var encoding = response.headers['content-type'].split("charset=")[1];
            body = body.toString(encoding);
            mlog.log(body.slice(0, 100));
            mlog.log("......");
            mlog.log(body.slice(-20));
            cb(null, body, response);
        } else {
            mlog.error(response.statusCode + " " + response.statusMessage);
            cb(error, body, response);
        }
    });
}

module.exports = EngineService;
