import { Component } from 'angular2/core';

import { BuildsService } from '../../services/builds.service';

@Component({
    selector: 'all-builds',
    template: require('./all-builds.html'),
    providers: [BuildsService]
})
export class AllBuilds {

    builds: any[];

    constructor(private _buildsService: BuildsService) {}

    ngOnInit() {
        this._buildsService.getAllBuilds()
            .subscribe(
                builds => this.builds = builds,
                err => console.error('err:', err));
    }

}
