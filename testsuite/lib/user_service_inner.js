
var MongoClient = require('mongodb').MongoClient;
var mlog = require('mocha-logger');

function UserServiceInner(host, mongoPort) {
    this.host = host;
    this.mongoPort = mongoPort;
    this.db = "users";
    this.collection = "user02";
}

UserServiceInner.prototype.removeUser = function (userName, auth, cb) {
    var url = "mongodb://" + this.host + ":" + this.mongoPort + "/" + this.db;
    var self = this; // don't use this in callbacks
    mlog.log("mongo: " + url);
    MongoClient.connect(url, function (err, db) {
        // Note: MongoClient Bug
        // Callback function WONT be executed when connection failure.
        // The error handling code below could be useless due to the bug.
        if (err) {
            mlog.error("connecting: " + err);
            cb(err, 0);
        } else {
            var query = {
                "auths.login_name": userName,
                "auths.auth_source_id": auth,
                "auths.primary": true
            };
            db.collection(self.collection).deleteMany(query, function (err, ret) {
                if (err) {
                    mlog.error("querying: " + err);
                    cb(err, 0);
                } else {
                    cb(err, ret.result.n);
                }
            });
        }
    });
}

module.exports = UserServiceInner;
