import { Component, SimpleChange, Input } from 'angular2/core';

import { BuildsService } from '../../services/builds.service';

@Component({
  selector: 'builds-tab',
  template: require('./builds.html'),
  styles: [require('./builds.css')],
  providers: [ BuildsService ]
})
export class BuildsTab {

  @Input('summary') summary;
  buildsList = [];

  constructor(private _buildsService: BuildsService) {}

  ngOnInit() {
    if (this.summary.repo && this.summary.no) {
      this.getLastBuilds();
    }
  }

  ngOnChanges(changes: {[propName: string]: SimpleChange}) {
    if (changes['summary'].currentValue && changes['summary'].currentValue['repo']) {
      this.getLastBuilds();
    }
  }

  getLastBuilds() {
    this._buildsService.getRepoBuilds(this.summary.repo)
      .map(res => res.json())
      .subscribe(
        (buildsList) => {
            console.log(buildsList);
            for(var key in buildsList) {
                buildsList[key].link = '#/build/' + buildsList[key].repo.replace(/\//g, ':') + '/' + buildsList[key].no;
            }
            this.buildsList = buildsList;
        },
        (err) => console.error('err', err));
  }
}
