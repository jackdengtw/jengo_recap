
var request = require("request")
var mlog = require("mocha-logger");

const GITHUB_API_URL = "https://api.github.com"

function GithubDriver() {
}

GithubDriver.prototype.getJengoWebHook = function (user, project, cb) {
    mlog.log(JSON.stringify(user));
    var options = {
        url: GITHUB_API_URL + "/repos/" + user.auth.login_name + "/" + project.meta.name + "/hooks",
        headers: {
            "Authorization": "token " + user.auth.token,
            "User-Agent": "javascript-request"
        }
    };
    request(options, function (error, response, body) {
        var hook = {};
        mlog.log(options.url);
        if (!error && response.statusCode == 200) {
            mlog.log(body);
            var hooks = JSON.parse(body);
            for (var i = 0; i < hooks.length; i++) {
                if (hooks[i].name == "web") {
                    break;
                }
            }
            if (i < hooks.length) {
                hook = hooks[i];
            }
        } else {
            if (!error) {
                mlog.error(response.statusCode + " " + response.statusMessage);
            }
        }
        cb(error, hook, response);
    });
}

module.exports = GithubDriver;
