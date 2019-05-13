#!/bin/bash

# Start Kafka Zookeeper
kafka_2.12-2.2.0/bin/zookeeper-server-start.sh kafka_2.12-2.2.0/config/zookeeper.properties

# Start Kafka server
kafka_2.12-2.2.0/bin/kafka-server-start.sh kafka_2.12-2.2.0/config/server.properties

# Create a new topic
kafka_2.12-2.2.0/bin/kafka-topics.sh --bootstrap-server localhost:9092 --topic users_topic --create --partitions 1 --replication-factor 1

# Describe topic
kafka_2.12-2.2.0/bin/kafka-topics.sh --bootstrap-server localhost:9092 --topic users_topic --describe

# List all topics
kafka_2.12-2.2.0/bin/kafka-topics.sh --bootstrap-server localhost:9092 --list

# Delete a topic
kafka_2.12-2.2.0/bin/kafka-topics.sh --bootstrap-server localhost:9092 --delete --topic userss_topic

# Start Kafka console consumer on specific topic, consume only new messages
kafka_2.12-2.2.0/bin/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic instruments_topic

# Start Kafka console consumer on specific topic, consume all messages from the beginning
kafka_2.12-2.2.0/bin/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic instruments_topic --from-beginning
