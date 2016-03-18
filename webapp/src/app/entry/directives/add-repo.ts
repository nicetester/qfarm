import { Component } from 'angular2/core';
import { Router } from 'angular2/router';

import { BuildsService } from '../../services/builds.service';

@Component({
    selector: 'add-repo',
    template: require('./add-repo.html'),
    providers: [BuildsService]
})
export class AddRepo {

    repoName: string = 'github.com/qfarm/bad-go-code'

    constructor(private _router: Router,
                private _buildsService: BuildsService){}

    startBuild() {
        this._buildsService.startNewBuild(this.repoName).subscribe(
            build => this.goToBuild(build.repo),
            err => console.error('Err', err)
        )
    }

    private goToBuild(repoName: string) {
        var safeRepoName = repoName.replace(/\//g,':');

        console.log("Build started on repo:", safeRepoName);
        this._router.navigate(['Last Build', { repoName: safeRepoName }]);
    }

}
