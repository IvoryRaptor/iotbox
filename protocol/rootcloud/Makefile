
OBJS = rootcloud.go
LIB = rootcloud.so

all:
	go build -buildmode=plugin -o $(LIB) $(OBJS)
	$(MAKE) copy

clean:
	-rm $(LIB)
	-rm $(LIBDIR)/protocol/$(LIB)

arm:
	CGO_ENABLED=1 CC=arm-linux-gnueabihf-gcc GOOS=linux GOARCH=arm  GOARM=7 go build -buildmode=plugin -o $(LIB) $(OBJS)
	$(MAKE) copy

copy: $(LIB)
	cp $(LIB) $(LIBDIR)/protocol