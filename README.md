# web-app-analyser-service

To run this application it is suggested to have docker installed in the system.
To start Server using following command
- ##docker-compose up

this will run the docker container and expose the port 8080
To test the end-point '/page-analytics', run the following curl request

- ### curl --location --request GET 'http://localhost:8080/page-analytics?url=http://www.data.gov.lk/'
- t http://www.data.gov.lk/ as the resource website

To have a more readable response can pipe through jq if this is hosted in unix based env.

- #### curl --location --request GET 'http://localhost:8080/page-analytics?url=http://www.data.gov.lk/' | jq
