FROM eclipse-temurin:21-jre
WORKDIR /setup
RUN useradd -ms /bin/bash steve
USER steve
WORKDIR /home/steve
RUN wget https://api.papermc.io/v2/projects/paper/versions/1.21.3/builds/81/downloads/paper-1.21.3-81.jar -O paper.jar
RUN echo "eula=true" > eula.txt
CMD ["java", "-Xms4G", "-Xmx4G", "-jar", "paper.jar", "nogui"]

