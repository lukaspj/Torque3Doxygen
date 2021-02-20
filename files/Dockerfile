FROM golang:1.15.4 AS t3ddocs_builder

WORKDIR /go/src/app

COPY . .

RUN CGO_ENABLED=0 go build -i -v -o ScriptExecServer

FROM ubuntu:20.04

ENV TZ=Europe/Copenhagen
ENV AZCOPY_VERSION=10

# First setup timezone to avoid prompt during install
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone \
# Then install dependencies
 && apt-get update \
 && apt-get upgrade -y \
 && apt-get install -y \
        # Binaries \
        git build-essential nasm xorg-dev \
        ninja-build gcc-multilib g++-multilib \
        cmake cmake-qt-gui \
        doxygen graphviz \
        # Libraries \
        libogg-dev libxft-dev libx11-dev libxxf86vm-dev \
        libopenal-dev libfreetype6-dev libxcursor-dev \
        libxinerama-dev libxi-dev libxrandr-dev \
        libxss-dev libglu1-mesa-dev libgtk-3-dev \
 && npm install -g postcss postcss-cli autoprefixer \
 && mkdir -p /home/hugo/ \
 && curl -L https://personalfrontend.blob.core.windows.net/misc/chugo_linux -o /home/hugo/hugo \
 && chmod +x /home/hugo/hugo \
 && mv /home/hugo/hugo /usr/bin/ \
 && mkdir -p /home/azcopy/ \
 && curl -L https://aka.ms/downloadazcopy-v${ AZCOPY_VERSION }-linux | tar -zxf - --directory /home/azcopy/ \
 && mv $(find /home/azcopy/ -type f -name azcopy) /usr/bin/

COPY Doxyfile /Torque3D/
COPY script.Doxyfile /Torque3D/

COPY --from=builder /go/src/app/ScriptExecServer /Goxygen/DoxygenConverter

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]