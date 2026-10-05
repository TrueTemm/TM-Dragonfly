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

package customblock

type Material struct {
	texture string

	renderMethod Method

	faceDimming bool

	ambientOcclusion float32
}

func NewMaterial(texture string, method Method) Material {
	m := Material{
		texture:          texture,
		renderMethod:     method,
		faceDimming:      true,
		ambientOcclusion: 1,
	}
	if !method.AmbientOcclusion() {
		m.ambientOcclusion = 0
	}
	return m
}

func (m Material) WithFaceDimming() Material {
	m.faceDimming = true
	return m
}

func (m Material) WithoutFaceDimming() Material {
	m.faceDimming = false
	return m
}

func (m Material) WithAmbientOcclusion() Material {
	m.ambientOcclusion = 1
	return m
}

func (m Material) WithoutAmbientOcclusion() Material {
	m.ambientOcclusion = 0
	return m
}

func (m Material) Encode() map[string]any {
	return map[string]any{
		"texture":           m.texture,
		"render_method":     m.renderMethod.String(),
		"face_dimming":      m.faceDimming,
		"ambient_occlusion": m.ambientOcclusion,
	}
}
