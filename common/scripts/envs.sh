namespace="azazon"

baseapp="${namespace}-base"
apps=(auth notification order payment stock)
modules=(common auth notification order payment stock)

kafka_cluster_name="main-kafka-cluster"
kafka_nodepool_name="dual-role"
topics=(order_created order_confirmed order_canceled)
