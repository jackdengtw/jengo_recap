
var request = require("request")
var mlog = require("mocha-logger");

var httpResp = require('./http_resp.js');

function UserService(host, port) {
    this.host = host;
    this.port = port;
    this.baseUri = "http://" + this.host + ":" + this.port;
}

UserService.prototype.createUser = function (userName, token, auth, scm, cb) {
    var options = {
        url: this.baseUri + "/v0.2/user?" +
        "login_name=" + userName +
        "&token=" + token +
        "&auth=" + auth +
        "&scm=" + scm,
        headers: {
            "User-Agent": "javascript-request"
        }
    };
    mlog.log("POST " + options.url);
    request.post(options, function (error, response, body) {
        var obj = {};
        httpResp(error, response, body, obj, response.statusCode >= 200 && response.statusCode < 300, cb);
    });
}

UserService.prototype.getUser = function (userId, cb) {
    var options = {
        url: this.baseUri + "/v0.2/internal_user/" + userId,
        headers: {}
    }
    mlog.log("GET " + options.url);
    request(options, function (error, response, body) {
        var user = {};
        httpResp(error, response, body, user, response.statusCode >= 200 && response.statusCode < 300, cb);
    });
}

module.exports = UserService;
