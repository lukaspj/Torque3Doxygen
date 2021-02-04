FROM ubuntu:20.04 AS base

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
        doxygen graphviz \
        # Libraries \
        libogg-dev libxft-dev libx11-dev libxxf86vm-dev \
        libopenal-dev libfreetype6-dev libxcursor-dev \
        libxinerama-dev libxi-dev libxrandr-dev \
        libxss-dev libglu1-mesa-dev libgtk-3-dev

FROM base AS docgen


# Checkout Torque3D source code
RUN git clone https://github.com/TorqueGameEngines/Torque3D.git /Torque3D \
 && cd /Torque3D \
 && git checkout Preview4_0

# Build
RUN mkdir -p /Torque3D/My\ Projects/Stock/buildFiles/ubuntu \
 && cd /Torque3D/My\ Projects/Stock/buildFiles/ubuntu \
 && cmake ../../../.. -DTORQUE_APP_NAME=Stock -DCMAKE_BUILD_TYPE=Release -DTORQUE_DEDICATED=ON -DVIDEO_WAYLAND=OFF \
 && make

COPY Doxyfile /Torque3D/

COPY test /Test