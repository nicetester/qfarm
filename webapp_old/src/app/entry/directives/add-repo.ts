import { Component } from 'angular2/core';
import { Router } from 'angular2/router';

import { BuildsService } from '../../services/builds.service';
import { WebSocketService } from '../../services/websocket.service';

@Component({
    selector: 'add-repo',
    template: require('./add-repo.html'),
    styles: [require('./add-repo.css')],
    providers: [BuildsService, WebSocketService]
})
export class AddRepo {

    repoName: string = 'github.com/qfarm/bad-go-code';
    events: any = {};

    constructor(private _router: Router,
                private _buildsService: BuildsService,
                private _websocketService: WebSocketService){}

    ngOnInit() {
        let linterEvent = /-done|-error$/;

        this._websocketService.socket().subscribe(
            event => {
                if (this.events['in-progress'] && this.repoName === event.repo) {
                    if (event.type === 'all-done') {
                        this.goToBuild(event.repo, event.payload);
                    }

                    if (event.type === 'error') {
                        this.events['error'] = true;
                        console.log(event.description);
                        delete this.events['in-progress'];
                    }
                    
                    if (linterEvent.test(event.type)) {
                        this.events[event.type] = true;
                    }
                }
            }
        )
    }

    startBuild() {
        this.events = { 'in-progress' : true };
        this._buildsService.startNewBuild(this.repoName).subscribe(
            build => this.listenToResults(this.repoName),
            err => console.error('Err', err)
        )
    }

    private listenToResults(build : string) {
        //
    }

    private goToBuild(repoName: string, buildId : string) {
        var safeRepoName = repoName.replace(/\//g,':');

        this.events = {};
        console.log("Build started on repo:", safeRepoName);
        this._router.navigate(['Build', { repoName: safeRepoName, buildId: buildId }]);
    }

}
