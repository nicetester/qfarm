import { Injectable } from 'angular2/core';
import { Http } from 'angular2/http';

@Injectable()
export class UserService {

    host = 'http://docker:8080/';

    constructor(private http: Http) {}

    getUserRepos(user : string) {
        return this.http.get(this.host + 'user_repos/?user=' + user);
    }

}
