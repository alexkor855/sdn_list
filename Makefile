BINDIR=${CURDIR}/bin
PACKAGE=${CURDIR}

build: bindir
	GOOS=linux GOARCH=amd64 go build -o ${BINDIR}/app ${PACKAGE}

bindir:
	mkdir -p ${BINDIR}

run: 
	docker-compose up --force-recreate --build

run-build: build
	docker-compose up --force-recreate --build
