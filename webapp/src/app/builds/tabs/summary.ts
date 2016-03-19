import { Component } from 'angular2/core';
import { Input } from 'angular2/core'

@Component({
    selector: 'summary-tab',
    template: require('./summary.html'),
    styles: [require('./summary.css')]
})
export class SummaryTab {

    @Input('summary') summary;
    scoreLevel: string;

    constructor() {
    }


}
