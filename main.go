package main

import (
	"fmt"
	"log/slog"
	"machine"
	"net/netip"
	"time"

	"github.com/soypat/cyw43439/examples/common"
	"github.com/soypat/seqs/eth/ntp"
	"github.com/soypat/seqs/stacks"
	"github.com/jgrelet/pico-rtc/rtcutil"
)

//const ntpServer = "pool.ntp.org:123"
const ntpServer = "pool.ntp.org"

//const ntpEpochOffset = 2208988800

type ntpConn struct {
	Hostname    string
	RequestedIP string
	UDPPorts    uint16
	// Number of TCP ports to open for the stack.
	TCPPorts uint16
	addrs    []netip.Addr
	routerhw [6]byte
	stack    *stacks.PortStack
}

func newNTPConn(hostname string, requestedIP string, udpPorts uint16) (*ntpConn, error) {
	logger := slog.New(slog.NewTextHandler(machine.Serial, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	time.Sleep(100 * time.Millisecond)
	dhcpc, stack, _, err := common.SetupWithDHCP(common.SetupConfig{
		Hostname:    hostname,
		Logger:      logger,
		RequestedIP: requestedIP,
		UDPPorts:    udpPorts,
	})
	if err != nil {
		return nil, fmt.Errorf("setup failed: %w", err.Error())
	}
	routerhw, err := common.ResolveHardwareAddr(stack, dhcpc.Router())
	if err != nil {
		return nil, fmt.Errorf("router hwaddr resolving: %w", err.Error())
	}

	resolver, err := common.NewResolver(stack, dhcpc)
	if err != nil {
		return nil, fmt.Errorf("resolver create: %w", err.Error())
	}
	// Résoudre l'adresse IP du serveur NTP
	addrs, err := resolver.LookupNetIP(ntpServer)
	if err != nil {
		return nil, fmt.Errorf("DNS lookup failed: %w", err.Error())
	}
	return &ntpConn{
		Hostname:    hostname,
		RequestedIP: requestedIP,
		UDPPorts:    udpPorts,
		addrs:       addrs,
		routerhw:    routerhw,
		stack:       stack,
	}, nil
}

func (c *ntpConn) String() string {
	return fmt.Sprintf("NTP conn to %s via %s", c.Hostname, c.stack.Addr())
}

func (c *ntpConn) getNTPTime() (time.Time, error) {
	ntpaddr := c.addrs[0]
	ntpc := stacks.NewNTPClient(c.stack, ntp.ClientPort)
	fmt.Println("NTP request to", ntpaddr.String())
	// Démarrer la requête NTP
	// Note: BeginDefaultRequest() est non-bloquant, il faut attendre avec IsDone()
	err := ntpc.BeginDefaultRequest(c.routerhw, ntpaddr)
	if err != nil {
		fmt.Errorf("NTP create: " + err.Error())
	}
	for !ntpc.IsDone() {
		time.Sleep(time.Second)
		//println("still ntping")
	}
	t := ntp.BaseTime().Add(ntpc.Offset())
	return t, nil
}

func main() {

	time.Sleep(2 * time.Second) // attendre que tout soit prêt
	fmt.Println("NTP dhcp started ...")

	// 1) Configurer le tick 1 Hz (à faire une fois au boot)
	//initRTC1Hz()
	// 1) Config 1 Hz (valeur standard Pico/Pico W)
	//rtcutil.Init1Hz(46874)

	// Lire la fréquence effective de clk_rtc
	freq := rtcutil.GetClkRTC()
	fmt.Println("clk_rtc =", freq, "Hz")

	/* // Supposons clk_rtc = 46875 Hz (cas standard Pico/Pico W : XOSC/12MHz ÷ 256)
	div := rtcutil.CalcDiv(freq) */
	rtcutil.InitRTC1Hz()

	// 2) Initialiser le Wi-Fi et la connexion NTP
	conn, err := newNTPConn("DHCP-pico-w", "192.168.1.150", 10)
	if err != nil {
		fmt.Println("Erreur connexion Wi-Fi :", err)
		return
	}
	println(conn.String())

	t, err := conn.getNTPTime()
	if err != nil {
		fmt.Println("NTP error:", err)
	} else {
		// format HH:MM:SS MM/DD/YYYY
		fmt.Printf("%02d:%02d:%02d %02d/%02d/%04d\n",
			t.Hour(), t.Minute(), t.Second(),
			t.Day(), t.Month(), t.Year())
	}

	// Convertir en Europe/Paris (CET/CEST)
	//tLocal := toEuropeParis(tUTC)

	// Régler le RTC
	//setRTC(t)
	rtcutil.Set(t)

	fmt.Println("RTC set.")

	for {
		now := rtcutil.Now()
		println(now.Format("15:04:05 02/01/2006"))
		time.Sleep(1 * time.Second)
	}
}
