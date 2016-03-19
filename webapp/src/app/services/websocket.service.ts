import { Injectable } from 'angular2/core';
import * as Rx from 'rxjs';

@Injectable()
export class WebSocketService {

    socket: Rx.Observable<any>;

    init() {
        this.socket = Rx.Observable.create(function (obs) {
            let host = '192.168.99.100';
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
}
