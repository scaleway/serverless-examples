deploy: node_modules
	serverless deploy

node_modules:
	npm i

clean: remove

remove:
	- serverless remove --service=apiGateway
        - serverless remove --service=getToken
	- serverless remove --service=myApp

dist-clean: remove
	- rm -rf node_modules package-lock.json
	- find . -name .serverless -exec rm -rf "{}" +
