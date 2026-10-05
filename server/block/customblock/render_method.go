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

type Method struct {
	renderMethod
}

func OpaqueRenderMethod() Method {
	return Method{0}
}

func AlphaTestRenderMethod() Method {
	return Method{1}
}

func BlendRenderMethod() Method {
	return Method{2}
}

func DoubleSidedRenderMethod() Method {
	return Method{3}
}

type renderMethod uint8

func (m renderMethod) Uint8() uint8 {
	return uint8(m)
}

func (m renderMethod) String() string {
	switch m {
	case 0:
		return "opaque"
	case 1:
		return "alpha_test"
	case 2:
		return "blend"
	case 3:
		return "double_sided"
	}
	panic("should never happen")
}

func (m renderMethod) AmbientOcclusion() bool {
	return m != 1 && m != 2
}
