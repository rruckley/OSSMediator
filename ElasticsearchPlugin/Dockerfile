# Staging: for building go app (ElasticsearchPlugin)
FROM golang:1.16 as builder

# set environment variable
ENV JAVA_BIN="/usr/local/jre1.8.0_171/bin/"
ENV PATH=$PATH:$GOBIN:$JAVA_BIN

# Install java for running elasticsearch for go test
RUN mkdir -p /tmp/java && \
        wget -O jre-8u171-linux-x64.tar.gz http://javadl.oracle.com/webapps/download/AutoDL?BundleId=233162_512cd62ec5174c3487ac17c61aaa89e8 && \
        tar -C /usr/local -xzf jre-8u171-linux-x64.tar.gz

# Install elasticsearch for unit test
RUN mkdir /elasticsearch/ && cd /elasticsearch && \
        wget https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-7.0.1-linux-x86_64.tar.gz && \
        tar -xzf elasticsearch-7.0.1-linux-x86_64.tar.gz

# create required directory for ElasticsearchPlugin
RUN mkdir -p /ElasticsearchPlugin/ /elasticsearch/elasticsearch-7.0.1/data/

# copy project directory to be built
COPY . /ElasticsearchPlugin/.
COPY ./VERSION ../.

# set the working directory
WORKDIR /ElasticsearchPlugin/

#Creating user to run elastic search
RUN groupadd -f -g 1000 elasticsearch && useradd -m -u 1000 -g 1000 -s /bin/bash elasticsearch
RUN chown -R elasticsearch /elasticsearch/elasticsearch-7.0.1/bin/elasticsearch \
        /elasticsearch/elasticsearch-7.0.1/data/ \
        /elasticsearch/elasticsearch-7.0.1/config/ \
        /elasticsearch/elasticsearch-7.0.1/logs/

# build the go app (ElasticsearchPlugin)
RUN make build test
