import {Component} from 'angular2/core';
import {RouteConfig, Router} from 'angular2/router';

import {Home} from './home/home';

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
  { path: '/',      name: 'Index', component: Home, useAsDefault: true },
])
export class App {
    name = 'Quality Farm';

    constructor() {}

}
