FROM eclipse-temurin:21-jre

WORKDIR /home/ubuntu
RUN wget https://api.papermc.io/v2/projects/paper/versions/1.21.3/builds/81/downloads/paper-1.21.3-81.jar -O paper.jar
COPY ./start.sh start.sh

WORKDIR /home/ubuntu/paper
RUN echo "eula=true" > eula.txt
CMD ["/home/ubuntu/start.sh"]

