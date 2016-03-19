import {Component} from 'angular2/core';
import {RouteConfig, Router} from 'angular2/router';
import {WebSocketService} from './services/websocket.service';

import { Entry } from './entry/entry';
import { Build } from './builds/build';

@Component({
    selector: 'app',
    pipes: [ ],
    providers: [ WebSocketService ],
    directives: [ ],
    styles: [require('./app.css')],
    template: require('./app.html')
})
@RouteConfig([
    { path: '/',      name: 'Entry', component: Entry, useAsDefault: true },
    { path: '/build/:repoName/', component: Build, name: 'Last Build' },
    { path: '/build/:repoName/:buildId/', component: Build, name: 'Build' }
])
export class App {
    name = 'Quality Farm';

    constructor(private _websocketService : WebSocketService) {
    }

}
