import { Injectable } from 'angular2/core';
import { Http, URLSearchParams } from 'angular2/http';
import * as Rx from 'rxjs';

@Injectable()
export class UserService {

    host = 'http://docker:8080/';

    constructor(private http: Http){}

    getUserRepos(user : string) {
        return this.http.get(this.host + 'user_repos/?user=' + user);
    }

}
