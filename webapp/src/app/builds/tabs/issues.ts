import { Component, SimpleChange } from 'angular2/core';
import { Input } from 'angular2/core';

import { IssuesService } from '../../services/issues.service';

@Component({
    selector: 'issues-tab',
    template: require('./issues.html'),
    styles: [require('./issues.css')],
    providers: [IssuesService]
})
export class IssuesTab {

    issues: any;
    filter: string;
    warns = true;
    errors = true;

    @Input('summary') summary;

    constructor(private _issuesService: IssuesService) {}

    ngOnInit() {
      if (this.summary.repo && this.summary.no) {
        this.getIssues();
      }
    }

    ngOnChanges(changes: {[propName: string]: SimpleChange}) {
      if (changes['summary'].currentValue && changes['summary'].currentValue['repo']) {
        this.getIssues();
      }
    }

    getIssues() {
        if (this.warns && this.errors) {
            this.filter = "";
        }
        if (this.warns && !this.errors) {
            this.filter = "warning";
        }
        if (!this.warns && this.errors) {
            this.filter = "error";
        }
        if (!this.warns && !this.errors) {
            this.filter = "none";
        }

        console.log("warn: ", this.warns, "errors: ", this.errors, "filter: ", this.filter);

        this._issuesService.getAllIssues(this.summary.repo, this.summary.no, 0, 50, this.filter)
            .map(res => res.json())
            .subscribe(
                (issues) => this.issues = issues,
                (err) => console.error('err', err));
    }
}
