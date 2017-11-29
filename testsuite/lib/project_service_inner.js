
var MongoClient = require('mongodb').MongoClient;
var mlog = require('mocha-logger');

function ProjectServiceInner(host, mongoPort) {
    this.host = host;
    this.mongoPort = mongoPort;
    this.db = "projects";
    this.collection = "project";
}

ProjectServiceInner.prototype.removeProjects = function (userId, cb) {
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
            var query = { "$pull": { "project.users": userId } }
            db.collection(self.collection).updateMany({}, query, function (err, ret) {
                if (err) {
                    mlog.error("querying: " + err);
                    cb(err, 0);
                } else {
                    cb(err, ret.result.nModified);
                }
            });
        }
    });
}

module.exports = ProjectServiceInner;
