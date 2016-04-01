import { Injectable } from 'angular2/core';
import { Http } from 'angular2/http';
import * as Rx from 'rxjs';

@Injectable()
export class IssuesService {

    host = 'http://docker:8080/';

    constructor(private http: Http){}

    getAllIssues(repoName: string, buildId: string, skip: number, size: number, filter: string) {
        return this.http.get(this.host + 'issues/?repo=' + repoName + '&filter=' + filter + '&skip=' + skip + '&size=' + size+ '&no=' + buildId);
    }


}
