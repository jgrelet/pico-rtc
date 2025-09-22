TARGET = pico-w
RP2350 = pico2-w
RP2040 = pico-w
SOURCE = main.go
BINARY = main.uf2
#LDFLAGS = -size short -monitor -scheduler tasks -gc=conservative -size=full -stack-size=20kb
LDFLAGS = -size short -monitor #-tags $(TARGET)

build:
	tinygo build -o $(BINARY) $(LDFLAGS) -target $(TARGET) $(SOURCE)

flash:
	tinygo flash $(LDFLAGS) -target $(TARGET) $(SOURCE)
	
monitor: 
	tinygo monitor -target=$(TARGET)	

fmt:
	go fmt *.go

clean:
	-rm -f $(BINARY)