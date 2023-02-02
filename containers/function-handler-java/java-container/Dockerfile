#
# Build stage
#
FROM maven:3.8.7-eclipse-temurin-17-alpine AS build
COPY src /home/app/src
COPY pom.xml /home/app
RUN mvn -f /home/app/pom.xml clean package

#
# Package stage
#
FROM eclipse-temurin:17-jre-alpine
COPY --from=build /home/app/target/serverless-scaleway-containers-1.0-jar-with-dependencies.jar /usr/local/lib/serverless-scaleway-containers.jar
ENTRYPOINT ["java","-jar","/usr/local/lib/serverless-scaleway-containers.jar"]