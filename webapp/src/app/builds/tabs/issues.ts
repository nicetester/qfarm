import { Component } from 'angular2/core';

import { IssuesService } from '../../services/issues.service';

@Component({
    selector: 'issues-tab',
    template: require('./issues.html'),
    styles: [require('./issues.css')],
    providers: [IssuesService]
})
export class IssuesTab {

    issues:any;

    constructor(private _issuesService: IssuesService){
        this._issuesService.getAllIssues('', '', 1, 1)
            .map(res => res.json())
            .subscribe(
                (issues) => {this.issues = issues},
                (err) => console.error('err', err));
    }


}
