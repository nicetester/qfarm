import { Injectable } from 'angular2/core';
import { Http } from 'angular2/http';
import * as Rx from 'rxjs';

@Injectable()
export class IssuesService {

    getAllIssues(repoName: string, buildId: stirng, first: number, limit: number) {
        return Rx.Observable.fromPromise(Promise.resolve([
            {
                file: 'path/to/file',
                line: 612,
                message: 'error: error not found',
                level: 'major'
            },
            {
                file: 'path/to/file',
                line: 612,
                message: 'error: error not found',
                level: 'minor'
            },
            {
                file: 'path/to/file',
                line: 612,
                message: 'error: error not found',
                level: 'major'
            },
            {
                file: 'path/to/file',
                line: 612,
                message: 'error: error not found',
                level: 'major'
            },
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
