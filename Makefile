
OUTFILE=${CIRCLE_ARTIFACTS}/pictures
FOLDERS=template/ static/

all:
	cd cmd/app
	GOOS=windows go build .
	go build .
	zip ${OUTFILE}.nw -r app app.exe css/ fonts/ index.html js/ package.json

