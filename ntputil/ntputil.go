//go:build pico || pico_w || pico2 || pico2_w

package ntputil

import (
	"fmt"
	"log/slog"
	"machine"
	"net/netip"
	"time"

	"github.com/soypat/cyw43439/examples/common"
	"github.com/soypat/seqs/eth/ntp"
	"github.com/soypat/seqs/stacks"
)

// const ntpServer = "pool.ntp.org:123"
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

