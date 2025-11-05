TARGET = pico2-w
RP2350 = pico2-w
RP2040 = pico-w
#SOURCE = main.go
SOURCE = .
BINARY = main.uf2
#LDFLAGS = -size short -monitor -scheduler tasks -gc=conservative -size=full -stack-size=20kb
LDFLAGS = -opt=1 -size short -monitor 

build:
	tinygo build -o $(BINARY) -target $(TARGET) $(LDFLAGS) $(SOURCE)

flash:
	tinygo flash -target $(TARGET) $(LDFLAGS) $(SOURCE)

info:
	tinygo flash -target $(TARGET) -tags=info $(LDFLAGS) $(SOURCE) 

debug:
	tinygo flash -target $(TARGET) -tags=debug $(LDFLAGS) $(SOURCE) 
	
monitor: 
	tinygo monitor -target=$(TARGET)	

gotests:
	gotests -all -ai -w *.go

gotestsall:
	gotests -all -ai -w ./...

test:
	go test -v ./...

fmt:
	go fmt *.go
	go fmt logger/*.go

clean:
	-rm -f $(BINARY)