FROM ubuntu:20.04

RUN apt-get update \
 && apt-get upgrade -y

RUN apt-get install -y libx11-6 libxft2 libgtk-3-0 libglib2.0-0
