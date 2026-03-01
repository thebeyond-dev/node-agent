package amneziawg

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/thebeyond-net/node-agent/internal/core/domain"
	"github.com/thebeyond-net/node-agent/pkg/keypair"
)

type configLines []string

type parsedConfig struct {
	publicKey   string
	listenPort  string
	magicParams magicParams
}

func (a *Adapter) loadAndParse() (parsedConfig, error) {
	f, err := os.Open(a.configPath)
	if err != nil {
		return parsedConfig{}, fmt.Errorf("open config: %w", err)
	}
	defer f.Close()

	var cfg parsedConfig
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.Contains(line, "=") || strings.HasPrefix(line, "[") {
			continue
		}

		key, val, _ := strings.Cut(line, "=")
		key, val = strings.TrimSpace(key), strings.TrimSpace(val)

		switch key {
		case "PrivateKey":
			if pub, err := keypair.DerivePublicKey(val); err == nil {
				cfg.publicKey = pub
			}
		case "ListenPort":
			cfg.listenPort = val
		default:
			a.fillMagicParam(&cfg.magicParams, key, val)
		}
	}

	if err := scanner.Err(); err != nil {
		return cfg, fmt.Errorf("read config: %w", err)
	}

	return cfg, nil
}

func (a *Adapter) GetServerConfig() (domain.ServerConfig, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.serverPubKey != "" {
		return domain.ServerConfig{
			PublicKey:  a.serverPubKey,
			ListenPort: a.listenPort,
		}, nil
	}

	cfg, err := a.loadAndParse()
	if err != nil {
		return domain.ServerConfig{}, err
	}

	a.magicParams = cfg.magicParams
	a.serverPubKey = cfg.publicKey
	a.listenPort = cfg.listenPort

	return domain.ServerConfig{
		PublicKey:  a.serverPubKey,
		ListenPort: a.listenPort,
	}, nil
}

func (a *Adapter) fillMagicParam(m *magicParams, key, val string) {
	switch key {
	case "Jc":
		m.Jc = val
	case "Jmin":
		m.Jmin = val
	case "Jmax":
		m.Jmax = val
	case "S1":
		m.S1 = val
	case "S2":
		m.S2 = val
	case "S3":
		m.S3 = val
	case "S4":
		m.S4 = val
	case "H1":
		m.H1 = val
	case "H2":
		m.H2 = val
	case "H3":
		m.H3 = val
	case "H4":
		m.H4 = val
	case "I1":
		m.I1 = val
	case "I2":
		m.I2 = val
	case "I3":
		m.I3 = val
	case "I4":
		m.I4 = val
	case "I5":
		m.I5 = val
	}
}

func (a *Adapter) readConfigLines() (configLines, error) {
	f, err := os.Open(a.configPath)
	if err != nil {
		return nil, fmt.Errorf("open config: %w", err)
	}
	defer f.Close()

	var lines configLines
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	return lines, nil
}

func (a *Adapter) writeAtomic(lines configLines) error {
	dir := filepath.Dir(a.configPath)
	tmp, err := os.CreateTemp(dir, "awg.conf.*")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	defer os.Remove(tmp.Name())

	writer := bufio.NewWriter(tmp)
	for _, line := range lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("write line: %w", err)
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("flush writer: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temp file: %w", err)
	}
	if err := os.Rename(tmp.Name(), a.configPath); err != nil {
		return fmt.Errorf("rename temp file: %w", err)
	}
	return nil
}
