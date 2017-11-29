
require('./utils/chai.js');
var mlog = require('mocha-logger');

var config = require("./config.js");

var UserService = require('../lib/user_service');
var UserServiceInner = require('../lib/user_service_inner');
var ProjectService = require('../lib/project_service');
var ProjectServiceInner = require('../lib/project_service_inner');
var EngineService = require('../lib/engine_service');
var GatewayService = require('../lib/gateway_service');
var GithubDriver = require('../lib/github_driver');

var projectId1 = "p_github_" + config.projectId1;
var userId1 = "u_github_" + config.userId1;

describe('When receiving a push event from github, ', function () {
    // suite vars
    var eventId1 = "";
    var run1 = {};
    var logId1 = "";
    var log1 = "";

    var userService = new UserService(config.server1, 8088);
    var userServiceInner = new UserServiceInner(config.server1, 27017);
    var projectService = new ProjectService(config.server1, 8080);
    var projectServiceInner = new ProjectServiceInner(config.server1, 27017);
    var engineService = new EngineService(config.server1, 8081);
    var gatewayService = new GatewayService(config.server1, 8082);
    var githubDriver = new GithubDriver();

    // this.timeout(120 * 1000);

    before(function (done) {
        // TODO:
        // Use Prepared user/projects for building
        // or
        // Create user, enable projects during runtime
        done();
    });

    it('gateway should process the event and give successfully resp', function (done) {
        gatewayService.sendEvent("github", testEvent, function (error, obj, response) {
            expect(error).equal(null);
            expect(response.statusCode).equal(200);
            expect(obj.event_id).to.include("e_github");
            eventId1 = obj.event_id;
            done();
        });
    });

    describe('then a run should be ongoing, ', function () {
        before(function (done) {
            expect(eventId1).is.not.equal("");
            this.timeout(3000);
            // sleep 2 seconds. Assuming an SLA of build starting in 2s
            setTimeout(function () {
                engineService.getRuns("", "", "", eventId1, function (error, runs, response) {
                    expect(error).equal(null);
                    expect(response.statusCode).equal(200);
                    expect(Array.isArray(runs)).equal(true);
                    expect(runs.length).equal(1);
                    expect(runs[0].event_id).equal(eventId1);
                    run1 = runs[0];

                    done()
                })
            }, 2000);
        });

        it('should get a running status of the project from engine service', function (done) {
            engineService.getRuns(userId1, projectId1, run1.id, eventId1, function (error, runs, response) {
                expect(error).equal(null);
                expect(response.statusCode).equal(200);
                expect(Array.isArray(runs)).equal(true);
                expect(runs.length).equal(1);
                expect(runs[0].status == "running" || runs[0].status == "prestarted").equal(true);
                done();
            });
        });

        it('should get a running status of the project from project service', function (done) {
            projectService.getRun(userId1, projectId1, run1.id, function (error, run, response) {
                expect(error).equal(null);
                expect(response.statusCode).equal(200);

                expect(run.status == "running" || run.status == "prestarted").equal(true);
                done();
            });
        });

        it('should get a living build log from engine service', function (done) {
            engineService.getRealtimeLog(run1.id, 0, 50, function (error, log, response) {
                expect(error).equal(null);
                expect(log.length).equal(50);
                log1 = log;
                done();
            });
        });
        it('should get an increasing build log from engine service', function (done) {
            engineService.getRealtimeLog(run1.id, 0, 999999, function (error, log, response) {
                expect(error).equal(null);
                var len = parseInt(response.headers["x-jengo-length"], 10);
                expect(log.length).equal(len);
                expect(log.slice(0, 50)).equal(log1);
                log1 = log;
                done();
            });
        });

        describe('then a run should success in 1 min, ', function () {
            this.timeout(40 * 1000);
            before(function (done) {
                // sleep. Assuming an simple build should end in 40 seconds
                setTimeout(function () { done() }, 40 * 999);
            });

            it('should get a finished status of the project from project service', function (done) {
                projectService.getRun(userId1, projectId1, run1.id, function (error, run, response) {
                    expect(error).equal(null);
                    expect(response.statusCode).equal(200);

                    expect(run.status).equal("finished");
                    done();
                });
            });
            it('should NOT be able to get this run from engine service', function (done) {
                // TODO: Should be able to trigger Engine cleanup
                mlog.log("TODO");
                /*
                engineService.getRuns(userId1, projectId1, run1.id, eventId1, function (error, runs, response) {
                    expect(error).equal(null);
                    expect(response.statusCode).equal(404);
                    done();
                });
                */
                done();
            });
            it('should NOT be able to get a living log console from engine service', function (done) {
                // TODO: Should be able to trigger Engine cleanup
                mlog.log("TODO");
                done();
            });
            it('should get a history run from project service', function (done) {
                projectService.getRun(userId1, projectId1, run1.id, function (error, obj, response) {
                    expect(error).equal(null);
                    expect(response.statusCode).equal(200);
                    expect(obj.event_id).equal(eventId1);
                    logId1 = obj.log_id;
                    done();
                });
            });
            it('should get a latest history log for the project from project service', function (done) {
                // expect a same log content from living log
                projectService.getLog(logId1, function (error, log, response) {
                    expect(error).equal(null);
                    var headers = response.headers;
                    var slen = headers["x-jengo-length"];
                    var len = parseInt(slen, 10);
                    expect(log.length).equal(len);
                    expect(log).equal(log1);
                    done();
                });
            });
        });
    });
    after(function () {
    });
    var testEvent = {
        "ref": "refs/heads/master",
        "commits": [
            {
                "message": "update travis yml",
                "author": {
                    "name": "Peng Chang",
                    "email": "pengchang@jd.com"
                },
                "url": "https://github.com/psmooth/helloworld/commit/1e3d6beba0b61cb19d929251c4cb8660d90dc54e",
                "distinct": true,
                "id": "1e3d6beba0b61cb19d929251c4cb8660d90dc54e",
                "tree_id": "104de2e3219784fe2274b7736a11b7b344d36c55",
                "timestamp": "2017-08-04T21:43:22+08:00",
                "committer": {
                    "name": "Peng Chang",
                    "email": "pengchang@jd.com"
                },
                "modified": [
                    ".travis.yml"
                ]
            }
        ],
        "before": "59663cf46382722800e738124ef6f6aa932f4fe7",
        "after": "1e3d6beba0b61cb19d929251c4cb8660d90dc54e",
        "created": false,
        "deleted": false,
        "forced": false,
        "compare": "https://github.com/psmooth/helloworld/compare/59663cf46382...1e3d6beba0b6",
        "repository": {
            "id": config.projectId1,
            "name": "helloworld",
            "full_name": "psmooth/helloworld",
            "owner": {
                "name": "psmooth",
                "email": "psmooth2@hotmail.com",
                "id": config.userId1
            },
            "private": false,
            "fork": false,
            "created_at": "2017-06-30T22:06:45+08:00",
            "pushed_at": "2017-08-04T21:43:37+08:00",
            "updated_at": "2017-06-30T14:14:53Z",
            "size": 0,
            "stargazers_count": 0,
            "watchers_count": 0,
            "language": "Go",
            "has_issues": true,
            "has_downloads": true,
            "has_wiki": true,
            "has_pages": false,
            "forks_count": 0,
            "open_issues_count": 0,
            "default_branch": "master",
            "master_branch": "master",
            "url": "https://github.com/psmooth/helloworld",
            "archive_url": "https://api.github.com/repos/psmooth/helloworld/{archive_format}{/ref}",
            "html_url": "https://github.com/psmooth/helloworld",
            "statuses_url": "https://api.github.com/repos/psmooth/helloworld/statuses/{sha}",
            "git_url": "git://github.com/psmooth/helloworld.git",
            "ssh_url": "git@github.com:psmooth/helloworld.git",
            "clone_url": "https://github.com/psmooth/helloworld.git",
            "svn_url": "https://github.com/psmooth/helloworld"
        },
        "head_commit": {
            "message": "update travis yml",
            "author": {
                "name": "Peng Chang",
                "email": "pengchang@jd.com"
            },
            "url": "https://github.com/psmooth/helloworld/commit/1e3d6beba0b61cb19d929251c4cb8660d90dc54e",
            "distinct": true,
            "id": "1e3d6beba0b61cb19d929251c4cb8660d90dc54e",
            "tree_id": "104de2e3219784fe2274b7736a11b7b344d36c55",
            "timestamp": "2017-08-04T21:43:22+08:00",
            "committer": {
                "name": "Peng Chang",
                "email": "pengchang@jd.com"
            },
            "modified": [
                ".travis.yml"
            ]
        },
        "pusher": {
            "name": "psmooth",
            "email": "psmooth2@hotmail.com"
        },
        "sender": {
            "login": "psmooth",
            "id": 123456,
            "avatar_url": "https://avatars3.githubusercontent.com/u/14289282?v=4",
            "html_url": "https://github.com/psmooth",
            "gravatar_id": "",
            "type": "User",
            "site_admin": false,
            "url": "https://api.github.com/users/psmooth",
            "events_url": "https://api.github.com/users/psmooth/events{/privacy}",
            "following_url": "https://api.github.com/users/psmooth/following{/other_user}",
            "followers_url": "https://api.github.com/users/psmooth/followers",
            "gists_url": "https://api.github.com/users/psmooth/gists{/gist_id}",
            "organizations_url": "https://api.github.com/users/psmooth/orgs",
            "received_events_url": "https://api.github.com/users/psmooth/received_events",
            "repos_url": "https://api.github.com/users/psmooth/repos",
            "starred_url": "https://api.github.com/users/psmooth/starred{/owner}{/repo}",
            "subscriptions_url": "https://api.github.com/users/psmooth/subscriptions"
        }
    };
});

