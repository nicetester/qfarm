import { Component } from 'angular2/core';

import { AddRepo } from './directives/add-repo';
import { AllBuilds } from './directives/all-builds';

@Component({
    selector: 'home',
    template: require('./home.html'),
    styles: [require('./home.css')],
    directives: [AddRepo, AllBuilds]
})
export class Home {}
