package amneziawg

import (
	_ "embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/thebeyond-net/node-agent/assets"
	"github.com/thebeyond-net/node-agent/internal/core/domain"
)

var clientTmpl = template.Must(template.New("client").Parse(assets.AmneziaClientTmpl))

type clientConfigData struct {
	magicParams
	PrivateKey      string
	Address         string
	DNS             string
	Endpoint        string
	ServerPublicKey string
	ListenPort      string
}

func (a *Adapter) BuildClientConfig(privKey, addr string, cfg domain.ServerConfig) (string, error) {
	var buf strings.Builder
	if err := clientTmpl.Execute(&buf, clientConfigData{
		magicParams:     a.magicParams,
		PrivateKey:      privKey,
		Address:         addr,
		DNS:             a.clientDNS,
		Endpoint:        a.endpoint,
		ServerPublicKey: cfg.PublicKey,
		ListenPort:      cfg.ListenPort,
	}); err != nil {
		return "", fmt.Errorf("execute client template: %w", err)
	}
	return buf.String(), nil
}
