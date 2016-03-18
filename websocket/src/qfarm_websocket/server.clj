(ns qfarm-websocket.server
  (:require [aleph.http :as http]
            [clojure.tools.logging :refer [info]]
            [manifold
             [deferred :as d]
             [stream :as s]]
            [mount.core :refer [args defstate]]
            [qfarm-websocket.redis :as redis]))

(defn handler [req]
  (info "Websocket connecting...")
  (-> (http/websocket-connection req)
      (d/chain
       (fn [ws]
         (info "Websocket connected")
         (s/on-closed ws #(info "Websocket closed"))
         (s/connect redis/events ws)))
      (d/catch
          (fn [e]
            {:status 400
             :headers {"content-type" "application/text"}
             :body "Expected a websocket request"}))))

(defstate server
  :start (http/start-server #'handler {:port (:port (args))})
  :stop (.close server))
