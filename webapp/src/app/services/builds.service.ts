import { Injectable } from 'angular2/core';
import { Http } from 'angular2/http';
import * as Rx from 'rxjs';

import { Build } from './Build'

@Injectable()
export class BuildsService {

    host = 'http://localhost:8080/';

    constructor(private http: Http){}

    getAllBuilds() {
        return Rx.Observable.fromPromise(Promise.resolve([
            new Build(3, 'userA/repoX'),
            new Build(1, 'userB/repoY'),
            new Build(2, 'userA/repoX'),
            new Build(1, 'userA/repoX')
        ]));

    }

    startNewBuild(repoName: string) {
        return this.http.post(
            this.host + 'build/',
            JSON.stringify({repo: repoName}));
    }

    getRepoBuilds(repoName: string) {
        return this.http.get(this.host + 'last_repo_builds/?repo=' + repoName);
    }

    getBuildSummary(repoName: string, buildId: string) {
        return Rx.Observable.fromPromise(Promise.resolve(
            {
                path: 'github.com/qfarm',
                repoName: 'bad-go-code',
                no: 123,
                score: Math.floor(Math.random()*100),
                time: Date.now()
            }
        ));
    }

}
