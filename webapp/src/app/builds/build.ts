import { Component } from 'angular2/core'
import { RouteParams } from 'angular2/router';

import { SummaryTab } from './tabs/summary'
import { IssuesTab } from './tabs/issues'
import { FilesTab } from './tabs/files'
import { BuildsTab } from './tabs/builds'

import { BuildsService } from '../services/builds.service';

@Component({
    selector: 'build-resuts',
    template: require('./build.html'),
    styles: [require('./build.css')],
    directives: [SummaryTab, IssuesTab, FilesTab, BuildsTab],
    providers: [BuildsService]
})
export class Build {

    repoName: string;
    buildId: string = '';
    tab: string = "summary"

    summary: any = {};

    constructor(
        private _routeParams: RouteParams,
        private _buildsService: BuildsService
    ) {
        this.repoName = _routeParams.get('repoName');
        this.buildId = _routeParams.get('buildId');
    }

    ngOnInit() {
        console.log(`Loaded Build view for build ${this.repoName}, build #${this.buildId}`);
        let realRepoName = this.repoName.replace(/:/g, '/');

        this._buildsService.getBuildSummary(realRepoName, this.buildId)
            .map(res => res.json())
            .subscribe(
                summary => this.summary = summary,
                err => console.error('err:', err));
    }

    showSummary() {
        this.tab = 'summary';
    }
    showFiles() {
        this.tab = 'files';
    }
    showIssues() {
        this.tab = 'issues';
    }
    showBuilds() {
        this.tab = 'builds';
    }

}
