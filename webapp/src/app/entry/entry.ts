import { Component } from 'angular2/core';

import { AddRepo } from './directives/add-repo';
import { AllBuilds } from './directives/all-builds';

@Component({
    selector: 'entry',
    template: require('./entry.html'),
    styles: [require('./entry.css')],
    directives: [AddRepo, AllBuilds]
})
export class Entry {

    constructor(){}

}
