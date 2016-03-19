import { Component, SimpleChange, Input } from 'angular2/core'

declare var Highcharts:any;

@Component({
    selector: 'summary-tab',
    template: require('./summary.html'),
    styles: [require('./summary.css')]
})
export class SummaryTab {

    @Input('summary') summary;
    scoreLevel: string;
    lintersList = [];

    constructor() {}

    ngOnInit() {
        if (this.summary.score) {
            this.showScoreChart();
            this.checkLintersList();
        }
    }

    ngOnChanges(changes: {[propName: string]: SimpleChange}) {
        if (changes['summary'].currentValue && changes['summary'].currentValue['score']) {
            this.showScoreChart();
            this.checkLintersList();
        }
    }

    showScoreChart() {
        new Highcharts.Chart(document.getElementById('score-chart'), {
            chart: {
                type: 'gauge',
                plotBackgroundColor: null,
                plotBackgroundImage: null,
                plotBorderWidth: 0,
                plotShadow: false
            },

            title: null,

            pane: {
                startAngle: -90,
                endAngle: 90,
                background: [{
                    backgroundColor: {
                        stops: [
                            [0, '#FFF'],
                            [1, '#333']
                        ]
                    },
                    borderWidth: 0,
                    outerRadius: '109%'
                }, {
                    backgroundColor: {
                        stops: [
                            [0, '#333'],
                            [1, '#FFF']
                        ]
                    },
                    borderWidth: 1,
                    outerRadius: '107%'
                }, {
                    backgroundColor: '#DDD',
                    borderWidth: 0,
                    outerRadius: '105%',
                    innerRadius: '103%'
                }]
            },

            // the value axis
            yAxis: {
                min: 0,
                max: 100,

                minorTickInterval: 'auto',
                minorTickWidth: 1,
                minorTickLength: 10,
                minorTickPosition: 'inside',
                minorTickColor: '#666',

                tickPixelInterval: 30,
                tickWidth: 2,
                tickPosition: 'inside',
                tickLength: 10,
                tickColor: '#666',
                labels: {
                    step: 2,
                    rotation: 'auto'
                },
                title: {
                    text: null
                },
                plotBands: [{
                    from: 0,
                    to: 60,
                    color: '#DF5353' // green
                }, {
                    from: 60,
                    to: 80,
                    color: '#DDDF0D' // yellow
                }, {
                    from: 80,
                    to: 100,
                    color: '#55BF3B' // red
                }]
            },

            series: [{
                name: 'Quality Score',
                data: [this.summary.score],
                tooltip: {
                    valueSuffix: null
                }
            }]
        });
    }

    checkLintersList() {
        var all = ['aligncheck', 'deadcode', 'dupl', 'errcheck', 'goconst', 'gocyclo', 'gofmt', 'goimports', 'golint', 'gotype', 'ineffassign', 'interfacer', 'lll', 'structcheck', 'test', 'testify', 'varcheck', 'vet', 'vetshadow', 'unconvert', 'coverage'];
        var ran = this.summary.config.linters;
        var res = [];

        for(var l of all) {
            res.push({
                name: l,
                ran: ran.indexOf(l) !== -1
            });
        }

        this.lintersList = res;
    }

}
