(defproject qfarm-websocket "0.1.0-SNAPSHOT"
  :description "quality.farm websocket service"
  :url "https://github.com/qfarm/qfarm"
  :license {:name "The MIT License"
            :url "https://opensource.org/licenses/MIT"}
  :dependencies [[org.clojure/clojure "1.8.0"]
                 [aleph "0.4.1-beta5"]
                 [mount "0.1.10"]
                 [org.clojure/tools.logging "0.3.1"]
                 [com.taoensso/carmine "2.12.2"]
                 [log4j/log4j "1.2.17"]]
  :main ^:skip-aot qfarm-websocket.core
  :target-path "target/%s"
  :profiles {:dev {:dependencies [[org.clojure/tools.namespace "0.2.11"]]}
             :uberjar {:aot :all}})
