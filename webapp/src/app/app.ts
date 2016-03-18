import {Component} from 'angular2/core';
import {RouteConfig, Router} from 'angular2/router';

import { Entry } from './entry/entry';
import { Build } from './builds/build';

@Component({
    selector: 'app',
    pipes: [ ],
    providers: [ ],
    directives: [ ],
    styles: [`
      nav ul {
        display: inline;
        list-style-type: none;
        margin: 0;
        padding: 0;
        width: 60px;
      }
      nav li {
        display: inline;
      }
      nav li.active {
        background-color: lightgray;
      }
    `],
    template: require('./app.html')
})
@RouteConfig([
    { path: '/',      name: 'Entry', component: Entry, useAsDefault: true },
    { path: '/build/:repoName/', component: Build, name: 'Last Build' },
    { path: '/build/:repoName/:buildId/', component: Build, name: 'Build' }
])
export class App {
    name = 'Quality Farm';

    constructor() {}

}
