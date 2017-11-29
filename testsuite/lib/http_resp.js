
var mlog = require("mocha-logger");

var httpResp = function (error, response, body, out, success, cb) {
    if (!error && success) {
        mlog.log(body);
        try {
            out = JSON.parse(body);
        } catch (e) {
            mlog.error(e);
        } finally {
            cb(null, out, response);
        }
    } else {
        if (!error) {
            mlog.error(response.statusCode + " " + response.statusMessage);
        }
        cb(error, out, response);
    }
}

module.exports = httpResp;

