#!/bin/bash

set -euo pipefail

if [ -z "$(ls -A /home/ubuntu/paper)" ]; then
  echo "eula=true" > /home/ubuntu/paper/eula.txt
fi

java -Xms4G -Xmx4G -jar /home/ubuntu/paper.jar nogui

