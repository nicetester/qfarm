import { Component, SimpleChange } from 'angular2/core';
import { EventEmitter, Input, Output } from 'angular2/core'
import { Router } from 'angular2/router';

import { FilesService } from '../../services/files.service';

@Component({
    selector: 'files-tab',
    template: require('./files.html'),
    styles: [require('./files.css')],
    providers: [FilesService]
})
export class FilesTab {

	  files: any;
    file: any;

  @Input('summary') summary;
  @Input('file') filePath;
  @Output() exitFileView = new EventEmitter();


    constructor(private _router: Router,
                private _filesService: FilesService){}

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
            (files) => {
                this.files = files;
                if (this.filePath) {
                    this.getFile();
                }
            },
          (err) => console.error('err', err));
  }

    getFile() {
        for (var f of this.files) {
            if (f.path === this.filePath) {
                this.file = f;
                this.file.decodedContent = atob(this.file.content);
            }
        }
    }

    showFile(file) {
        if(!file.dir) {
          this.filePath = file.path.slice(0);
            this.getFile();
        }
    }

    backToFiles() {
        this.file = null;
        this.exitFileView.emit();
    }

}
