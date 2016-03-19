import { Component } from 'angular2/core';
import { Input } from 'angular2/core'

declare var Highcharts:any;

@Component({
    selector: 'summary-tab',
    template: require('./summary.html'),
    styles: [require('./summary.css')]
})
export class SummaryTab {

    @Input('summary') summary;
    scoreLevel: string;

    constructor() {
    }

    ngOnInit() {
        console.log('oninit', this.summary)

        new Highcharts.Chart(document.getElementById('build-score-widget'), {
            chart: {
                type: 'gauge',
                plotBackgroundColor: null,
                plotBackgroundImage: null,
                plotBorderWidth: 0,
                plotShadow: false
            },

            title: {
                text: 'Quality Score'
            },

            pane: {
                startAngle: -150,
                endAngle: 150,
                background: [{
                    backgroundColor: {
                        linearGradient: { x1: 0, y1: 0, x2: 0, y2: 1 },
                        stops: [
                            [0, '#FFF'],
                            [1, '#333']
                        ]
                    },
                    borderWidth: 0,
                    outerRadius: '109%'
                }, {
                    backgroundColor: {
                        linearGradient: { x1: 0, y1: 0, x2: 0, y2: 1 },
                        stops: [
                            [0, '#333'],
                            [1, '#FFF']
                        ]
                    },
                    borderWidth: 1,
                    outerRadius: '107%'
                }, {
                    // default background
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
                data: [80],
                tooltip: {
                    valueSuffix: null
                }
            }]
        });

    }

}
