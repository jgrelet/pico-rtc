//go:build pico || pico_w || pico2 || pico2_w

package ntputil

import (
	"log/slog"
	"fmt"
	"net/netip"
	"time"
	"github.com/jgrelet/pico-rtc/logger"

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

// NewNTPConn initializes a new NTP connection by setting up network interfaces,
// configuring DHCP, DNS, and resolving the NTP server address. It returns an
// ntpConn instance on success, or an error if any step fails.
//
// Parameters:
//   - hostname: the desired hostname for the device.
//   - requestedIP: the preferred IP address to request via DHCP.
//   - udpPorts: the number of UDP ports to allocate.
//   - logs: a logger instance for logging network events.
//
// Returns:
//   - *ntpConn: a pointer to the initialized NTP connection structure.
//   - error: an error if setup fails at any stage.
func NewNTPConn(hostname string, requestedIP string, udpPorts uint16, logs *slog.Logger) (*ntpConn, error) {
	time.Sleep(100 * time.Millisecond)
	// Configurer le Wi-Fi, DHCP, DNS, etc.
	dhcpc, stack, _, err := common.SetupWithDHCP(common.SetupConfig{
		Hostname:    hostname,
		Logger:      logs,
		RequestedIP: requestedIP,
		UDPPorts:    udpPorts,
	})

	if err != nil {
		return nil, fmt.Errorf("setup failed: %w", err.Error())
	}
	// Obtenir l'adresse MAC de la passerelle (routeur)
	routerhw, err := common.ResolveHardwareAddr(stack, dhcpc.Router())
	if err != nil {
		return nil, fmt.Errorf("router hwaddr resolving: %v", err.Error())
	}
	// Créer un résolveur DNS
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

// String returns a string representation of the ntpConn, including the hostname and IP address.
func (c *ntpConn) String() string {
	return fmt.Sprintf("NTP conn to %s, IP: %s", c.Hostname, c.stack.Addr())
}

// GetNTPTime sends an NTP request to the first configured NTP server address,
// waits for the response, and returns the synchronized time.
// The function initiates a non-blocking NTP request and waits until the request is complete.
// It returns the calculated time based on the NTP offset, or an error if the request fails.
func (c *ntpConn) GetNTPTime() (time.Time, error) {
	ntpaddr := c.addrs[0]
	ntpc := stacks.NewNTPClient(c.stack, ntp.ClientPort)
	logger.Logger.Info("NTP request to", ntpaddr.String())
	// Démarrer la requête NTP
	// Note: BeginDefaultRequest() est non-bloquant, il faut attendre avec IsDone()
	err := ntpc.BeginDefaultRequest(c.routerhw, ntpaddr)
	if err != nil {
		fmt.Errorf("NTP create: " + err.Error())
	}
	for !ntpc.IsDone() {
		time.Sleep(time.Second)
		logger.Logger.Info("still ntping")
	}
	t := ntp.BaseTime().Add(ntpc.Offset())
	return t, nil
}
