FROM ubuntu:20.04

RUN apt-get update \
 && apt-get upgrade -y \
 && apt-get install -y libx11-6 libxft2 libgtk-3-0 libglib2.0-0 iptables sudo

RUN adduser --system --group  --home /workspace app

