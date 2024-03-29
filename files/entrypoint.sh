#!/bin/bash

# Checkout Torque3D source code
git clone https://github.com/TorqueGameEngines/Torque3D.git /Torque3D
cd /Torque3D || exit
git checkout "${T3D_BRANCH}"

# Build
mkdir -p /Torque3D/My\ Projects/Stock/buildFiles/ubuntu
cd /Torque3D/My\ Projects/Stock/buildFiles/ubuntu || exit
cmake ../../../.. -DTORQUE_SCRIPT_EXTENSION=cs -DTORQUE_APP_NAME=Stock -DCMAKE_BUILD_TYPE=Release -DTORQUE_DEDICATED=ON -DVIDEO_WAYLAND=OFF
make

cd /Torque3D || exit
cp /Goxygen/Doxyfile ./Doxyfile
doxygen
cd /Torque3D/My\ Projects/Stock/game || exit
cat > ./main.cs <<EOF
dumpEngineDocs("consoledoc.h");
quit();
EOF
./Stock

cp /Goxygen/script.Doxyfile ./Doxyfile
doxygen

mkdir /DoxygenOutput
ls -la /Torque3D/My\ Projects/Stock/game
cp -r /Torque3D/My\ Projects/Stock/game/script-doxygen /DoxygenOutput/
cp -r /Torque3D/doxygen /DoxygenOutput/

cd /DoxygenOutput || exit
/Goxygen/DoxygenConverter

mkdir /Hugo
git clone https://github.com/lukaspj/T3DDocs.git /Hugo/t3ddocs
cp -r /DoxygenOutput/hugo/content /Hugo/t3ddocs/content/_generated

cd /Hugo/t3ddocs || exit
cp config.toml config.toml.bck
printf "uglyURLs = true\n" | cat - config.toml > /tmp/out && mv /tmp/out config.toml
printf "\nt3dversion = \"%s\"" "${T3D_VERSION}" >> config.toml
printf "\noffline = true\n" >> config.toml
export HUGO_ENV='production'
hugo -v --minify --enableGitInfo
mkdir static
zip -9 -r static/offline.zip public/
rm config.toml
mv config.toml.bck config.toml
printf "\nt3dversion = \"%s\"\n" "${T3D_VERSION}" >> config.toml
export HUGO_ENV='production'
hugo -v --minify --enableGitInfo
azcopy sync public/ "${AZURE_STORAGE_CONTAINER_URL}/${T3D_SLUG}${AZURE_STORAGE_SAS_TOKEN}" --delete-destination=true
