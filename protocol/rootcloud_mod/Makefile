
OBJS = rootcloud.go
LIB = rootcloud.so
LIB_PATH = ../../lib/protocol

copy: $(LIB)
	cp $(LIB) $(LIB_PATH)

mac: $(LIB)
	GOOS=darwin GOARCH=amd64 go build -buildmode=plugin -o $(LIB) $(OBJS)

windows: $(LIB)
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -buildmode=plugin -o $(LIB) $(OBJS)

linux: $(LIB)
	CGO_ENABLED=0 GOOS=linux GOARCH=386  go build -buildmode=plugin -o $(LIB) $(OBJS)

arm: $(LIB)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm  GOARM=7 go build -buildmode=plugin -o $(LIB) $(OBJS)

all:
	go build -buildmode=plugin -o $(LIB) $(OBJS)
	cp $(LIB) $(LIB_PATH)
clean:
	rm $(LIB)
	rm $(LIB_PATH)/$(LIB)
