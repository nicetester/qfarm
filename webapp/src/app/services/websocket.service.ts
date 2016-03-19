import { Injectable } from 'angular2/core';
import * as Rx from 'rxjs';

@Injectable()
export class WebSocketService {

    socketObservable: Rx.Observable<any>;

    socket() : Rx.Observable<any> {
        if (!this.socketObservable) {
            this.socketObservable = Rx.Observable.create(function (obs) {
                // let host = location.hostname.split(/:/)[0];
                let host = 'docker';
                let connect = () => {
                    let ws = new WebSocket(`ws://${host}:8081/`);
                    console.log('Websocket: Connecting...')
                    ws.onopen = (s) => { console.log("Websocket: connected."); }
                    ws.onmessage = (e) => {
                        try {
                            let msg = JSON.parse(e.data);
                            obs.next(msg);
                        } catch (e) {
                        }
                    }
                    ws.onclose = (e) => {
                        try {
                            ws.close();
                        } catch (e) {}
                        setTimeout(() => {
                            connect();
                        }, 5000);
                    };
                };

                connect();
            });
        }
        return this.socketObservable;
    }
}
