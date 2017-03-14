
OUTFILE=${CIRCLE_ARTIFACTS}/pictures
FOLDERS=template/ static/

all:
	GOOS=windows go build .
	go build .
	zip ${OUTFILE}.zip pictures.exe -r ${FOLDERS}
	tar -czf ${OUTFILE}.tar.gz pictures ${FOLDERS}

