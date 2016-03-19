import { Component, SimpleChange } from 'angular2/core';
import { Input } from 'angular2/core'

import { FilesService } from '../../services/files.service';

@Component({
    selector: 'files-tab',
    template: require('./files.html'),
	providers: [FilesService]
})
export class FilesTab {

	files:any;

	@Input('summary') summary;

    constructor(private _filesService: FilesService){}

	ngOnInit() {
      if (this.summary.repo && this.summary.no) {
        this.getFiles();
      }
    }

    ngOnChanges(changes: {[propName: string]: SimpleChange}) {
      if (changes['summary'].currentValue && changes['summary'].currentValue['repo']) {
        this.getFiles();
      }
    }

	getFiles() {
      this._filesService.getAllFiles(this.summary.repo, this.summary.no)
        .map(res => res.json())
        .subscribe(
          (files) => {this.files = files},
          (err) => console.error('err', err));
    }
}
