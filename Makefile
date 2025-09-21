TARGET = pico2-w
RP2350 = pico2-w
RP2040 = pico-w
SOURCE = main.go
BINARY = main.uf2
#LDFLAGS = -size short -monitor -scheduler tasks -gc=conservative -size=full -stack-size=20kb
LDFLAGS = -size short -monitor

build:
	tinygo build -o $(BINARY) $(LDFLAGS) -target $(TARGET) $(SOURCE)

flash:
	tinygo flash $(LDFLAGS) -target $(TARGET) $(SOURCE)

rp2350:
	tinygo flash $(LDFLAGS) -target $(RP2350) -tags rp2350 main_rp2350.go	

rp2040:
	tinygo flash $(LDFLAGS) -target $(RP2040) -tags rp2040 main_rp2040.go
	
monitor: 
	tinygo monitor -target=$(TARGET)	

fmt:
	go fmt *.go

clean:
	-rm -f $(BINARY)