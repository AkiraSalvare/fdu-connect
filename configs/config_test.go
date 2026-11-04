package configs

import "testing"

func TestDefault(t *testing.T) {
	cfg := Default()
	if cfg.Protocol != "easyconnect" {
		t.Fatalf("Protocol = %q, want easyconnect", cfg.Protocol)
	}
	if cfg.ServerAddress != "stuvpn.fudan.edu.cn" {
		t.Fatalf("ServerAddress = %q, want stuvpn.fudan.edu.cn", cfg.ServerAddress)
	}
	if cfg.ServerPort != 443 {
		t.Fatalf("ServerPort = %d, want 443", cfg.ServerPort)
	}
	if cfg.SocksBind != ":1080" || cfg.HTTPBind != ":1081" {
		t.Fatalf("proxy defaults = %q, %q", cfg.SocksBind, cfg.HTTPBind)
	}
}
