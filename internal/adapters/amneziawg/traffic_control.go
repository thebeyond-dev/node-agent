package amneziawg

import (
	"fmt"
	"net/netip"
	"strings"

	"github.com/thebeyond-net/node-agent/pkg/exec"
)

func (a *Adapter) SetPeerBandwidth(ip string, bandwidth int) error {
	if bandwidth <= 0 {
		return a.removeBandwidthLimit(ip)
	}
	return a.applyBandwidthLimit(ip, bandwidth)
}

func (a *Adapter) applyBandwidthLimit(allowedIP string, bandwidth int) error {
	ipAddr := strings.TrimSuffix(allowedIP, "/32")
	classID, err := generateClassID(ipAddr)
	if err != nil {
		return err
	}

	rate := fmt.Sprintf("%dmbit", bandwidth)

	if err := exec.Run("tc", "class", "replace", "dev", a.intf,
		"parent", "1:1",
		"classid", classID,
		"htb",
		"rate", rate,
		"ceil", rate,
		"burst", "128k",
		"cburst", "128k",
	); err != nil {
		return fmt.Errorf("tc class: %w", err)
	}

	if err := exec.Run("tc", "qdisc", "replace", "dev", a.intf,
		"parent", classID,
		"fq_codel",
		"limit", "10240",
		"target", "5ms",
		"interval", "100ms",
	); err != nil {
		return fmt.Errorf("fq_codel: %w", err)
	}

	if err := exec.Run("tc", "filter", "replace", "dev", a.intf,
		"protocol", "ip",
		"parent", "1:0",
		"prio", "1", "u32",
		"match", "ip", "dst", ipAddr,
		"flowid", classID,
	); err != nil {
		return fmt.Errorf("tc filter: %w", err)
	}

	return nil
}

func (a *Adapter) removeBandwidthLimit(allowedIP string) error {
	ipAddr := strings.TrimSuffix(allowedIP, "/32")
	classID, err := generateClassID(ipAddr)
	if err != nil {
		return err
	}

	if err := exec.Run("tc", "qdisc", "del", "dev", a.intf, "parent", classID); err != nil {
		return fmt.Errorf("tc qdisc del: %w", err)
	}

	return exec.Run("tc", "class", "del", "dev", a.intf, "classid", classID)
}

func generateClassID(ipAddr string) (string, error) {
	ip, err := netip.ParseAddr(ipAddr)
	if err != nil {
		return "", fmt.Errorf("invalid ip %s: %w", ipAddr, err)
	}

	if !ip.Is4() {
		return "", fmt.Errorf("ipv6 tc limit not implemented")
	}

	b := ip.As4()
	return fmt.Sprintf("1:%x", uint16(b[2])<<8|uint16(b[3])), nil
}
