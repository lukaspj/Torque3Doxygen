FROM ubuntu:20.04 AS builder

ENV TZ=Europe/Copenhagen

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
        # Libraries \
        libogg-dev libxft-dev libx11-dev libxxf86vm-dev \
        libopenal-dev libfreetype6-dev libxcursor-dev \
        libxinerama-dev libxi-dev libxrandr-dev \
        libxss-dev libglu1-mesa-dev libgtk-3-dev

# Checkout Torque3D source code
RUN git clone https://github.com/TorqueGameEngines/Torque3D.git /Torque3D \
 && cd /Torque3D \
 && git checkout Preview4_0

# Build
RUN mkdir -p /Torque3D/My\ Projects/Stock/buildFiles/ubuntu \
 && cd /Torque3D/My\ Projects/Stock/buildFiles/ubuntu \
 && cmake ../../../.. -DTORQUE_APP_NAME=Stock -DCMAKE_BUILD_TYPE=Release -DTORQUE_DEDICATED=ON -DVIDEO_WAYLAND=OFF \
 && make

FROM ubuntu:20.04 AS runtime

RUN apt-get update \
 && apt-get upgrade -y

RUN apt-get install -y libx11-6 libxft2 libgtk-3-0 libglib2.0-0

COPY --from=builder [ \
    "/Torque3D/My Projects/Stock/game/Stock", \
    "/Torque3D/My Projects/Stock/game/libSDL2-2.0d.so.0.12.0", \
    "/Torque3D/" \
]

COPY main.cs /Torque3D/

# HACK: SDL linking error fix
ENV LD_LIBRARY_PATH=/Torque3D

RUN ln -s /Torque3D/libSDL2-2.0d.so.0.12.0 /Torque3D/libSDL2-2.0d.so.0 \
 && ln -s /Torque3D/libSDL2-2.0d.so.0 /Torque3D/libSDL2-2.0d.so

WORKDIR /Torque3D

ENTRYPOINT ["./Stock"]