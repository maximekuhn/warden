#!/bin/bash

set -euo pipefail

if [ -z "$(ls -A /home/steve/paper)" ]; then
  echo "eula=true" > /home/steve/paper/eula.txt
fi

java -Xms4G -Xmx4G -jar /home/steve/paper.jar nogui

