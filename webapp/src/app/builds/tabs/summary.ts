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
        new Highcharts.Chart('score-chart', {

            chart: {
              type: 'solidgauge',
              marginTop: 20
            },

            title: null,

            tooltip: {
              borderWidth: 0,
              backgroundColor: 'none',
              shadow: false,
              style: {
                fontSize: '16px'
              },
              pointFormat: '{series.name}<br><span style="font-size:2em; color: {point.color}; font-weight: bold">{point.y}%</span>',
              positioner: function (labelWidth, labelHeight) {
                return {
                  x: 190 - labelWidth / 2,
                  y: 150
                };
              }
            },

            pane: {
              startAngle: 0,
              endAngle: 360,
              background: [{ // Track for Score
                outerRadius: '112%',
                innerRadius: '88%',
                backgroundColor: Highcharts.Color(Highcharts.getOptions().colors[0]).setOpacity(0.3).get(),
                borderWidth: 0
              }, { // Track for Coverage
                outerRadius: '87%',
                innerRadius: '63%',
                backgroundColor: Highcharts.Color(Highcharts.getOptions().colors[1]).setOpacity(0.3).get(),
                borderWidth: 0
              }]
            },

            yAxis: {
              min: 0,
              max: 100,
              lineWidth: 0,
              tickPositions: []
            },

            plotOptions: {
              solidgauge: {
                borderWidth: '15px',
                dataLabels: {
                  enabled: false
                },
                linecap: 'round',
                stickyTracking: false
              }
            },

            series: [{
              name: 'Score',
              borderColor: Highcharts.getOptions().colors[0],
              data: [{
                color: Highcharts.getOptions().colors[0],
                radius: '100%',
                innerRadius: '100%',
                y: this.summary.score
              }]
            }, {
              name: 'Coverage',
              borderColor: Highcharts.getOptions().colors[1],
              data: [{
                color: Highcharts.getOptions().colors[1],
                radius: '75%',
                innerRadius: '75%',
                y: Number((this.summary.coverage).toFixed(2))
              }]
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
