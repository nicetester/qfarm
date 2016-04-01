import {Component} from 'angular2/core';
import {Router, RouteConfig, ROUTER_DIRECTIVES} from 'angular2/router';

import {BuildsService} from './services/builds.service';
import {WebSocketService} from './services/websocket.service';
import {UserService} from './services/user.service';

import {Home} from './home/home';
import {Build} from './builds/build';
import {Config} from './config/config';

import '../style/app.scss';

@Component({
    selector: 'app',
    directives: [...ROUTER_DIRECTIVES],
    providers: [ BuildsService, WebSocketService, UserService ],
    styles: [require('./app.scss')],
    template: require('./app.html')
})
@RouteConfig([
    { path: '/', name: 'Home', component: Home, useAsDefault: true},
    { path: '/build/:repoName/', component: Build, name: 'Last Build' },
    { path: '/build/:repoName/:buildId/', component: Build, name: 'Build' },
    { path: '/build/:repoName/:buildId/:file/', component: Build, name: 'Build - File View' },
    { path: '/config', component: Config, name: 'Config - Get QFarm Configuration' }
])
export class App {

        name = 'Quality Farm';
    buildsList: any;
    userRepos: any;

    user : string;

    constructor(private _buildsService : BuildsService,
                private _websocketService : WebSocketService,
                private _userService : UserService,
                private _router : Router) {
        let buildPath = /^build\//;

        _router.subscribe(path => {
            if (buildPath.test(path)) {
                let repo = path.split(/\//)[1];
                this.user = repo.split(/:/)[1];
                this.getUserRepos();
            } else {
                this.user = null;
                this.getLastBuilds();
            }
        });
    }

    getLastBuilds() {
        this._buildsService.getLastBuilds()
            .map(res => res.json())
            .subscribe(
                (buildsList) => {
                    this.buildsList = buildsList.map((b) => {
                        return {
                            repo: b.repo,
                            no: b.no,
                            link: '#/build/' + b.repo.replace(/\//g, ':') + '/' + b.no
                        };
                    }).sort((a, b) => b.no - a.no);
                },
                (err) => console.error('err', err));
    }

    getUserRepos() {
        this._userService.getUserRepos(this.user)
            .map(res => res.json())
            .subscribe(
                (repos) => {
                    this.userRepos = repos.map((r) => {
                        return {
                            repo: r,
                            link: '#/build/' + r.replace(/\//g, ':') + '/'
                        };
                    });
                },
                (err) => console.error('err', err));
    }


}
