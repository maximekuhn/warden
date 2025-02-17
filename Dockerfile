FROM eclipse-temurin:21-jre

ARG GID=1000
ARG UID=1000
RUN groupadd -g $GID steve
RUN useradd -m -u $UID -g $GID -s /bin/bash steve

USER steve

WORKDIR /home/steve
RUN wget https://api.papermc.io/v2/projects/paper/versions/1.21.3/builds/81/downloads/paper-1.21.3-81.jar -O paper.jar
COPY ./start.sh start.sh

WORKDIR /home/steve/paper
RUN echo "eula=true" > eula.txt
CMD ["/home/steve/start.sh"]

