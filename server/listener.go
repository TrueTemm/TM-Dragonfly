/*
 _____               _____
|_   _| __ _   _  __|_   _|__ _ __ ___  _ __ ___
  | || '__| | | |/ _ \| |/ _ \ '_ ` _ \| '_ ` _ \
  | || |  | |_| |  __/| |  __/ | | | | | | | | | |
  |_||_|   \__,_|\___||_|\___|_| |_| |_|_| |_| |_|

 _____ __  __       ____                               __ _
|_   _|  \/  |     |  _ \ _ __ __ _  __ _  ___  _ __  / _| |_   _
  | | | |\/| |_____| | | | '__/ _` |/ _` |/ _ \| '_ \| |_| | | | |
  | | | |  | |_____| |_| | | | (_| | (_| | (_) | | | |  _| | |_| |
  |_| |_|  |_|     |____/|_|  \__,_|\__, |\___/|_| |_|_| |_|\__, |
                                    |___/                   |___/

@author TrueTemm
@link   https://github.com/TrueTemm
TM-Dragonfly Project
*/

package server

import (
	"fmt"
	"io"
	"log/slog"

	"github.com/df-mc/dragonfly/multiversion/native"
	"github.com/df-mc/dragonfly/multiversion/v1001"
	"github.com/df-mc/dragonfly/multiversion/v2168"
	"github.com/df-mc/dragonfly/multiversion/v2169"
	"github.com/df-mc/dragonfly/multiversion/v671"
	"github.com/df-mc/dragonfly/multiversion/v685"
	"github.com/df-mc/dragonfly/multiversion/v686"
	"github.com/df-mc/dragonfly/multiversion/v712"
	"github.com/df-mc/dragonfly/multiversion/v729"
	"github.com/df-mc/dragonfly/multiversion/v748"
	"github.com/df-mc/dragonfly/multiversion/v766"
	"github.com/df-mc/dragonfly/multiversion/v776"
	"github.com/df-mc/dragonfly/multiversion/v786"
	"github.com/df-mc/dragonfly/multiversion/v800"
	"github.com/df-mc/dragonfly/multiversion/v818"
	"github.com/df-mc/dragonfly/multiversion/v819"
	"github.com/df-mc/dragonfly/multiversion/v827"
	"github.com/df-mc/dragonfly/multiversion/v844"
	"github.com/df-mc/dragonfly/multiversion/v859"
	"github.com/df-mc/dragonfly/multiversion/v898"
	"github.com/df-mc/dragonfly/multiversion/v924"
	"github.com/df-mc/dragonfly/multiversion/v944"
	"github.com/df-mc/dragonfly/multiversion/v975"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/sandertv/gophertunnel/minecraft"
)

type Listener interface {
	Accept() (session.Conn, error)

	Disconnect(conn session.Conn, reason string) error
	io.Closer
}

func (uc UserConfig) listenerFunc(conf Config) (Listener, error) {
	cfg := minecraft.ListenConfig{
		MaximumPlayers:         conf.MaxPlayers,
		StatusProvider:         conf.StatusProvider,
		AuthenticationDisabled: conf.AuthDisabled,
		ResourcePacks:          conf.Resources,
		TexturePacksRequired:   conf.ResourcesRequired,
		Compression:            conf.Compression,
		Allow:                  conf.Allower.Allow,
		HTTPClient:             authClient,

		AcceptedProtocols: []minecraft.Protocol{v671.Protocol{}, v685.Protocol{}, v686.Protocol{}, v712.Protocol{}, v729.Protocol{}, v748.Protocol{}, v766.Protocol{}, v776.Protocol{}, v786.Protocol{}, v800.Protocol{}, v818.Protocol{}, v819.Protocol{}, v827.Protocol{}, v844.Protocol{}, v859.Protocol{}, v898.Protocol{}, v924.Protocol{}, v944.Protocol{}, v975.Protocol{}, v1001.Protocol{}, v2168.Protocol{}, v2169.Protocol{}, native.Protocol{}},
	}

	cfg.ErrorLog = conf.Log.With("net origin", "gophertunnel")
	l, err := cfg.Listen("raknet", uc.Network.Address)
	if err != nil {
		return nil, fmt.Errorf("create minecraft listener: %w", err)
	}
	conf.Log.Info("Listener running.", "addr", l.Addr())
	return listener{Listener: l, log: conf.Log}, nil
}

type listener struct {
	*minecraft.Listener
	log *slog.Logger
}

func (l listener) Accept() (session.Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}

	if mc, ok := conn.(*minecraft.Conn); ok && l.log != nil {
		id, cd := mc.IdentityData(), mc.ClientData()
		l.log.Info("Player connecting",
			"name", id.DisplayName,
			"ip", mc.RemoteAddr().String(),
			"version", cd.GameVersion,
			"device", cd.DeviceModel,
		)
	}
	return conn.(session.Conn), err
}

func (l listener) Disconnect(conn session.Conn, reason string) error {
	return l.Listener.Disconnect(conn.(*minecraft.Conn), reason)
}
