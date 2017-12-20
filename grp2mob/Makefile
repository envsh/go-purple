VERSION         :=      $(shell cat ./VERSION)
GOVVV=`govvv -flags -version ${VERSION}|sed 's/=/=GOVVV-/g'`

all:

bd:
	go build -v -ldflags "${GOVVV}" .
