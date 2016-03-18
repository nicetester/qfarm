import { Injectable } from 'angular2/core';
import { Http } from 'angular2/http';
import * as Rx from 'rxjs';

import { Build } from './Build'

@Injectable()
export class BuildsService {

    getAllBuilds() {
        return Rx.Observable.fromPromise(Promise.resolve([
            new Build(3, 'userA/repoX'),
            new Build(1, 'userB/repoY'),
            new Build(2, 'userA/repoX'),
            new Build(1, 'userA/repoX')
        ]));
    }

    startNewBuild(repoName: string) {
        return Rx.Observable.fromPromise(
            Promise.resolve(new Build(1, repoName))
        );
    }

    getRepoBuilds(repoName: string) {
        return Rx.Observable.fromPromise(Promise.resolve([
            new Build(3, repoName),
            new Build(2, repoName),
            new Build(1, repoName)
        ]));
    }

}
