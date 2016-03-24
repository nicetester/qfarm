import { Injectable } from 'angular2/core';
import { Http } from 'angular2/http';

@Injectable()
export class FilesService {

    host = 'http://docker:8080/';

    constructor(private http: Http) {}

    getAllFiles(repoName: string, buildId: string) {
        return this.http.get(this.host + 'files/?repo=' + repoName + '&no=' + buildId);
    }


}
