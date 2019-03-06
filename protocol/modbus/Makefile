
OBJS = modbus.go net.go rtu.go crc.go
LIB = modbus.so

all:
	go build -buildmode=plugin -o $(LIB) $(OBJS)
	$(MAKE) copy

clean:
	-rm $(LIB)
	-rm $(LIBDIR)/protocol/$(LIB)

mac:
	GOOS=darwin GOARCH=amd64 go build -buildmode=plugin -o $(LIB) $(OBJS)
	$(MAKE) copy
linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=386  go build -buildmode=plugin -o $(LIB) $(OBJS)
	$(MAKE) copy
arm:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm  GOARM=7 go build -buildmode=plugin -o $(LIB) $(OBJS)
	$(MAKE) copy

copy: $(LIB)
	cp $(LIB) $(LIBDIR)/protocol