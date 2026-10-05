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

package sliceutil

import "slices"

func Convert[A, B any, S ~[]B](v S) []A {
	a := make([]A, len(v))
	for i, b := range v {
		a[i] = (any)(b).(A)
	}
	return a
}

func SearchValue[A any, S ~[]A](v S, f func(a A) bool) (a A, ok bool) {
	for _, val := range v {
		if f(val) {
			return val, true
		}
	}
	return
}

func Filter[E any](s []E, c func(E) bool) []E {
	a := make([]E, 0, len(s))
	for _, e := range s {
		if c(e) {
			a = append(a, e)
		}
	}
	return a
}

func DeleteVal[E any](s []E, v E) []E {
	for i, vs := range s {
		if (any)(v) == (any)(vs) {
			return slices.Delete(s, i, i+1)
		}
	}
	return s
}
