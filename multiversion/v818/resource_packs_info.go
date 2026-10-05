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

package v818

import (
	"github.com/df-mc/dragonfly/multiversion/v786"
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

type ResourcePacksInfo struct {
	TexturePackRequired bool

	HasAddons bool

	HasScripts bool

	ForceDisableVibrantVisuals bool

	WorldTemplateUUID uuid.UUID

	WorldTemplateVersion string

	TexturePacks []protocol.TexturePackInfo
}

func (*ResourcePacksInfo) ID() uint32 {
	return IDResourcePacksInfo
}

func (pk *ResourcePacksInfo) Marshal(io protocol.IO) {
	io.Bool(&pk.TexturePackRequired)
	io.Bool(&pk.HasAddons)
	io.Bool(&pk.HasScripts)
	io.Bool(&pk.ForceDisableVibrantVisuals)
	io.UUID(&pk.WorldTemplateUUID)
	io.String(&pk.WorldTemplateVersion)
	v786.SliceUint16Length786(io, &pk.TexturePacks)
}
