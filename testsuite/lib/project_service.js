var request = require("request")
var mlog = require("mocha-logger");

var httpResp = require('./http_resp.js');

function ProjectService(host, port) {
    this.host = host;
    this.port = port;
    this.baseUri = "http://" + this.host + ":" + this.port;
}

ProjectService.prototype.refresh = function (userId, cb) {
    var options = {
        url: this.baseUri + "/v0.1/user/" + userId + "/projects/action?" +
        "method=update",
        headers: {
            "User-Agent": "javascript-request"
        }
    };
    mlog.log("GET " + options.url);
    request(options, function (error, response, body) {
        var obj = {};
        httpResp(error, response, body, obj, response.statusCode >= 200 && response.statusCode < 300, cb);
    });
}

ProjectService.prototype.listProjects = function (userId, cb) {
    var options = {
        url: this.baseUri + "/v0.1/user/" + userId + "/projects",
        headers: {
            "User-Agent": "javascript-request"
        }
    };
    mlog.log("GET " + options.url);
    request(options, function (error, response, body) {
        var obj = [];
        httpResp(error, response, body, obj, response.statusCode >= 200 && response.statusCode < 300, cb);
    });
}

ProjectService.prototype.getProject = function (userId, projectId, cb) {
    var options = {
        // TODO: project service should provide get project by Id API
        url: this.baseUri + "/v0.1/user/" + userId + "/projects",
        headers: {
            "User-Agent": "javascript-request"
        }
    };
    mlog.log("GET " + options.url);
    request(options, function (error, response, body) {
        var obj = [];
        httpResp(error, response, body, obj, response.statusCode >= 200 && response.statusCode < 300,
            function (error, tmp, response) {
                for (var i = 0; i < tmp.length; i++) {
                    if (tmp[i].id == projectId) {
                        mlog.log(JSON.stringify(tmp[i]));
                        cb(error, tmp[i], response);
                        return;
                    }
                }
                mlog.log("project id " + projectId + " not found!");
                cb(error, {}, response);
            });
    });
}

ProjectService.prototype.actionOnProject = function (userId, projectId, action, cb) {
    var options = {
        url: this.baseUri + "/v0.1/user/" + userId + "/project/" + projectId + "/action?" +
        "method=" + action,
        headers: {
            "User-Agent": "javascript-request"
        }
    };
    mlog.log("GET " + options.url);
    request(options, function (error, response, body) {
        var obj = {};
        httpResp(error, response, body, obj, response.statusCode >= 200 && response.statusCode < 300, cb);
    });
}

ProjectService.prototype.getRun = function (userId, projectId, runId, cb) {
    var options = {
        url: this.baseUri + "/v0.1/user/" + userId + "/project/" + projectId + "/builds",
        headers: {
            "User-Agent": "javascript-request"
        }
    };
    mlog.log("GET " + options.url);
    request(options, function (error, response, body) {
        var obj = [];
        httpResp(error, response, body, obj, response.statusCode >= 200 && response.statusCode < 300,
            function (error, builds, response) {
                for (var i = 0; i < builds.length; i++) {
                    for (var j = 0; j < builds[i].runs.length; j++) {
                        var tmp = builds[i].runs[j]
                        if (tmp.id == runId) {
                            mlog.log(JSON.stringify(tmp));
                            cb(error, tmp, response);
                            return;
                        }
                    }
                }
                mlog.log("run id " + runId + " not found!");
                cb(error, {}, response);
            });
    });
}

ProjectService.prototype.getBuilds = function (userId, projectId, max_count, cb) {
    var options = {
        url: this.baseUri + "/v0.1/user/" + userId + "/project/" + projectId + "/builds?max_count=" + max_count,
        headers: {
            "User-Agent": "javascript-request"
        }
    };
    mlog.log("GET " + options.url);
    request(options, function (error, response, body) {
        var obj = [];
        httpResp(error, response, body, obj, response.statusCode >= 200 && response.statusCode < 300, cb);
    });
}

ProjectService.prototype.getLog = function (logId, cb) {
    var options = {
        url: this.baseUri + "/v0.1/log/" + logId,
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

module.exports = ProjectService;
