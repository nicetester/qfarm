import { Component } from 'angular2/core';

import { AddRepo } from './directives/add-repo';
import { AllBuilds } from './directives/all-builds';

@Component({
    template: require('./entry.html'),
    directives: [AddRepo, AllBuilds]
})
export class Entry {

    constructor(){}

}
