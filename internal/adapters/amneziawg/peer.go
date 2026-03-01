package amneziawg

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func (a *Adapter) RegisterPeer(pubKey, allowedIP string) error {
	if err := a.addToConfig(pubKey, allowedIP); err != nil {
		return fmt.Errorf("add to config: %w", err)
	}
	if err := a.addPeerLive(pubKey, allowedIP); err != nil {
		_ = a.removeFromConfig(pubKey)
		return fmt.Errorf("add peer live: %w", err)
	}
	return nil
}

func (a *Adapter) RemovePeer(pubKey string) error {
	_ = a.removePeerLive(pubKey)
	return a.removeFromConfig(pubKey)
}

func (a *Adapter) addToConfig(pubKey, allowedIP string) error {
	lines, err := a.readConfigLines()
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}
	if len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) != "" {
		lines = append(lines, "")
	}
	lines = append(lines,
		"[Peer]",
		"PublicKey = "+pubKey,
		"AllowedIPs = "+allowedIP+"/32",
	)
	return a.writeAtomic(lines)
}

func (a *Adapter) removeFromConfig(pubKey string) error {
	lines, err := a.readConfigLines()
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	var result configLines
	skip := false

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		if strings.EqualFold(line, "[Peer]") {
			if isTargetPeer(lines, i, pubKey) {
				skip = true
				continue
			}
		}

		if skip && strings.HasPrefix(line, "[") && !strings.EqualFold(line, "[Peer]") {
			skip = false
		}
		if !skip {
			result = append(result, lines[i])
		}
	}

	return a.writeAtomic(cleanEmptyLines(result))
}

func isTargetPeer(lines []string, start int, pubKey string) bool {
	for j := start + 1; j < len(lines); j++ {
		next := strings.TrimSpace(lines[j])
		if strings.HasPrefix(next, "[") {
			break
		}
		if strings.Contains(next, "PublicKey") && strings.Contains(next, pubKey) {
			return true
		}
	}
	return false
}

func (a *Adapter) addPeerLive(pubKey, allowedIP string) error {
	return a.execAwg("set", a.intf, "peer", pubKey, "allowed-ips", allowedIP+"/32")
}

func (a *Adapter) removePeerLive(pubKey string) error {
	return a.execAwg("set", a.intf, "peer", pubKey, "remove")
}

func (a *Adapter) execAwg(args ...string) error {
	cmd := exec.Command("awg", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("awg %s failed: %s (%w)", args[0], strings.TrimSpace(stderr.String()), err)
	}
	return nil
}

func cleanEmptyLines(lines []string) []string {
	var res []string
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			if i > 0 && strings.TrimSpace(lines[i-1]) == "" {
				continue
			}
		}
		res = append(res, line)
	}
	for len(res) > 0 && strings.TrimSpace(res[len(res)-1]) == "" {
		res = res[:len(res)-1]
	}
	return res
}
