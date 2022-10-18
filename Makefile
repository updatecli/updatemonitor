

.PHONY: app.build
app.build:
	go build -o bin/updatefactory

app.start: app.build
	./bin/updatefactory

.PHONY: db
db.reset: db.delete db.start

.PHONY: db.connect
db.connect:
	docker exec -i -t updatefactory-mongodb-1 mongosh mongodb://admin:password@localhost:27017

.PHONY: db.start
db.start:
	docker compose up -d mongodb

.PHONY: db.delete
db.delete:
	docker compose down mongodb -v
	
