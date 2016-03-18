(ns qfarm-websocket.core
  (:gen-class)
  (:require [clojure.tools.logging :refer [info]]
            [clojure.tools.namespace.repl :as tn]
            [mount.core :as mount]))

(defn go
  ([]
   (go {:port 8081
        :redis-host "192.168.99.100"
        :redis-port 6379
        :redis-topic "events"}))
  ([opts]
   (info "Staring with args..." opts)
   (mount.core/start-with-args opts)
   :ready))

(defn reset []
  (mount.core/stop)
  (tn/refresh :after 'qfarm-websocket.core/go))

(defn -main
  [& args]
  (mount.core/start-with-args {:port 8081
                               :redis-host "redis"
                               :redis-port 6379
                               :redis-topic "events"}))
