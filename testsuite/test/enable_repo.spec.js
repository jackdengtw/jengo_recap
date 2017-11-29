
require('./utils/chai.js')
var mlog = require('mocha-logger');

var config = require('./config.js');

var UserService = require('../lib/user_service');
var UserServiceInner = require('../lib/user_service_inner');
var ProjectService = require('../lib/project_service');
var ProjectServiceInner = require('../lib/project_service_inner');
var GithubDriver = require('../lib/github_driver');

describe('When enabling/disabling projects for signup-ed user, ', function () {
    // suite variables
    var user1 = {}, projects = [], project1 = {};
    var userId1 = "u_github_" + config.userId1;
    var projectId1 = "p_github_" + config.projectId1;

    var userService = new UserService(config.server1, 8088);
    var userServiceInner = new UserServiceInner(config.server1, 27017);
    var projectService = new ProjectService(config.server1, 8080);
    var projectServiceInner = new ProjectServiceInner(config.server1, 27017);
    var githubDriver = new GithubDriver();

    this.timeout(15000);

    before(function (done) { // Note: Mocha Don't support multiple done() invocation
        // setTimeout(done, 10000);
        projectServiceInner.removeProjects(userId1, function (err, n) {
            mlog.log(n + " projects removed"); // ignore error

            userServiceInner.removeUser(config.userName1, config.auth1, function (err, n) {
                mlog.log(n + " users removed"); // ignore error

                userService.createUser(config.userName1, config.token1, config.auth1, config.scm1, function (err, resp, response) {
                    expect(err).equal(null);
                    expect(response.statusCode).equal(201);
                    expect(resp.id).is.not.equal(undefined);
                    mlog.log("user " + config.userName1 + " created!");

                    userService.getUser(resp.id, function (err, user, response) {
                        expect(err).equal(null);
                        expect(response.statusCode).equal(200);
                        user1 = user;
                        expect(user1).is.not.equal(null);

                        projectService.refresh(user1.user_id, function (error, obj, response) {
                            expect(err).equal(null);

                            projectService.listProjects(user1.user_id, function (error, projs, response) {
                                expect(err).equal(null);
                                projects = projs;

                                expect(projects).is.not.equal(null);
                                expect(Array.isArray(projects)).equal(true);
                                expect(projects.length).above(0);
                                project1 = projects[0];
                                done();
                            });
                        });
                    });
                });
            });
        });
    });

    describe('then enabling the first project, ', function () {
        var ret = {};
        before(function (done) {
            projectService.actionOnProject(user1.user_id, project1.meta.id, "enable", function (error, obj, response) {
                expect(error).equal(null);
                // expect(response.statusCode).equal(200);
                ret = obj;
                done();
            });
        });
        it('should return success', function () {
            expect(ret).is.not.equal(null);
            // expect(ret.status).equal("enabled");
        });
        it('should set github web hook correctly', function (done) {
            githubDriver.getJengoWebHook(user1, project1, function (error, hook, response) {
                expect(error).equal(null);
                expect(response.statusCode).equal(200);
                mlog.log(JSON.stringify(hook));
                expect(hook.events).include("push");
                // expect(hook.config.url).include("http://jengo.ci" ); TODO: check domain name after env finalized
                expect(hook.config.url).include("http://47.93.38.8");
                done();
            });
        });
        it('should be able to redo', function (done) {
            projectService.actionOnProject(user1.user_id, project1.meta.id, "enable", function (error, obj, response) {
                expect(error).equal(null);
                expect(obj).is.not.equal(null);
                expect(response.statusCode).equal(200);
                done();
            });
        });
    });
    describe('then disabling the first project, ', function () {
        var ret = {};
        before(function (done) {
            projectService.actionOnProject(user1.user_id, project1.meta.id, "disable", function (error, obj, response) {
                expect(error).equal(null);
                expect(response.statusCode).equal(200);
                ret = obj;
                done();
            });
        });
        it('should return success', function () {
            expect(ret).is.not.equal(null);
        });
        it('should remove github web hook correctly', function (done) {
            githubDriver.getJengoWebHook(user1, project1, function (error, hook, response) {
                expect(error).equal(null);
                expect(response.statusCode).equal(200);
                expect(hook).deep.equal({});
                done();
            });
        });
        it('should be able to redo', function (done) {
            projectService.actionOnProject(user1.user_id, project1.meta.id, "disable", function (error, obj, response) {
                expect(error).equal(null);
                expect(obj).is.not.equal(null);
                expect(response.statusCode).equal(200);
                done();
            });
        });
    });
    /*
    after(function (done) {
        projectServiceInner.removeProjects(userId1, function (err, n) {
            mlog.log(n + " projects removed"); // ignore error
            userServiceInner.removeUser(config.userName1, config.auth1, function (err, n) {
                mlog.log(n + " users removed"); // ignore error
                done()
            });
        });
    });
    */
});

