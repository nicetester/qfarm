import { Injectable } from 'angular2/core';
import { Http } from 'angular2/http';
import * as Rx from 'rxjs';

@Injectable()
export class IssuesService {

    host = 'http://docker:8080/';

    constructor(private http: Http){}

    getAllIssues(repoName: string, buildId: string, first: number, limit: number) {
        return this.http.get(this.host + 'issues/?repo=github.com/qfarm/bad-go-code&filter=error&skip=0&size=10');
    }


}
