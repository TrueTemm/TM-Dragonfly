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

package cmd

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/internal/sliceutil"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"math/rand/v2"
	"reflect"
	"slices"
	"sort"
	"strconv"
	"strings"
)

type Line struct {
	args []string
	seen []string
	src  Source
	cmd  Command
}

func (line *Line) SyntaxError() error {
	if len(line.args) == 0 {
		return MessageSyntax.F(strings.Join(line.seen, " "), "", "")
	}
	next := strings.Join(line.args[1:], " ")
	if next != "" {
		next = " " + next
	}
	return MessageSyntax.F(strings.Join(line.seen, " ")+" ", line.args[0], next)
}

func (line *Line) UsageError() error {
	return MessageUsage.F(line.cmd.Usage())
}

func (line *Line) Next() (string, bool) {
	v, ok := line.NextN(1)
	if !ok {
		return "", false
	}
	return v[0], true
}

func (line *Line) NextN(n int) ([]string, bool) {
	if len(line.args) < n {
		return nil, false
	}
	v := line.args[:n]
	return v, true
}

func (line *Line) RemoveNext() {
	line.RemoveN(1)
}

func (line *Line) RemoveN(n int) {
	if len(line.args) < n {
		line.args = nil
		return
	}
	line.seen = append(line.seen, line.args[:n]...)
	line.args = line.args[n:]
}

func (line *Line) Leftover() []string {
	v := line.args
	line.args = nil
	return v
}

func (line *Line) Len() int {
	return len(line.args)
}

type parser struct {
	currentField string
}

func (p parser) parseArgument(line *Line, v reflect.Value, optional bool, name string, source Source, tx *world.Tx) (error, bool) {
	var err error
	i := v.Interface()
	if line.Len() == 0 && optional {

		v.Set(reflect.Zero(v.Type()))
		return nil, false
	}
	switch i.(type) {
	case int, int8, int16, int32, int64:
		err = p.int(line, v)
	case uint, uint8, uint16, uint32, uint64:
		err = p.uint(line, v)
	case float32, float64:
		err = p.float(line, v)
	case string:
		err = p.string(line, v)
	case bool:
		err = p.bool(line, v)
	case mgl64.Vec3:
		err = p.vec3(line, v)
	case Varargs:
		err = p.varargs(line, v)
	case []Target:
		err = p.targets(line, v, tx)
	case SubCommand:
		err = p.sub(line, name)
	default:
		if param, ok := i.(Parameter); ok {
			err = param.Parse(line, v)
			break
		}
		if enum, ok := i.(Enum); ok {
			err = p.enum(line, v, enum, source)
			break
		}
		panic(fmt.Sprintf("non-command parameter type %T in command structure", i))
	}
	if err == nil {

		line.RemoveNext()
	}
	return err, err == nil
}

func (p parser) int(line *Line, v reflect.Value) error {
	arg, ok := line.Next()
	if !ok {
		return line.UsageError()
	}
	value, err := strconv.ParseInt(arg, 10, v.Type().Bits())
	if err != nil {
		return MessageNumberInvalid.F(arg)
	}
	v.SetInt(value)
	return nil
}

func (p parser) uint(line *Line, v reflect.Value) error {
	arg, ok := line.Next()
	if !ok {
		return line.UsageError()
	}
	value, err := strconv.ParseUint(arg, 10, v.Type().Bits())
	if err != nil {
		return MessageNumberInvalid.F(arg)
	}
	v.SetUint(value)
	return nil
}

func (p parser) float(line *Line, v reflect.Value) error {
	arg, ok := line.Next()
	if !ok {
		return line.UsageError()
	}
	value, err := strconv.ParseFloat(arg, v.Type().Bits())
	if err != nil {
		return MessageNumberInvalid.F(arg)
	}
	v.SetFloat(value)
	return nil
}

func (p parser) string(line *Line, v reflect.Value) error {
	arg, ok := line.Next()
	if !ok {
		return line.UsageError()
	}
	v.SetString(arg)
	return nil
}

func (p parser) bool(line *Line, v reflect.Value) error {
	arg, ok := line.Next()
	if !ok {
		return line.UsageError()
	}
	value, err := strconv.ParseBool(arg)
	if err != nil {
		return MessageBooleanInvalid.F(arg)
	}
	v.SetBool(value)
	return nil
}

func (p parser) enum(line *Line, val reflect.Value, v Enum, source Source) error {
	arg, ok := line.Next()
	if !ok {
		return line.UsageError()
	}
	opts := v.Options(source)
	ind := slices.IndexFunc(opts, func(s string) bool {
		return strings.EqualFold(s, arg)
	})
	if ind < 0 {
		if _, loose := v.(Suggestion); loose {

			val.SetString(arg)
			return nil
		}
		return MessageParameterInvalid.F(arg)
	}
	val.SetString(opts[ind])
	return nil
}

func (p parser) sub(line *Line, name string) error {
	arg, ok := line.Next()
	if !ok {
		return line.UsageError()
	}
	if strings.EqualFold(name, arg) {
		return nil
	}
	return MessageParameterInvalid.F(arg)
}

func (p parser) vec3(line *Line, v reflect.Value) error {
	if err := p.float(line, v.Index(0)); err != nil {
		return err
	}
	line.RemoveNext()
	if err := p.float(line, v.Index(1)); err != nil {
		return err
	}
	line.RemoveNext()
	return p.float(line, v.Index(2))
}

func (p parser) varargs(line *Line, v reflect.Value) error {
	v.SetString(strings.Join(line.Leftover(), " "))
	return nil
}

func (p parser) targets(line *Line, v reflect.Value, tx *world.Tx) error {
	targets, err := p.parseTargets(line, tx)
	if err != nil {
		return err
	}
	if len(targets) == 0 {
		return MessageNoTargets.F()
	}
	v.Set(reflect.ValueOf(targets))
	return nil
}

func (p parser) parseTargets(line *Line, tx *world.Tx) ([]Target, error) {
	entities, players := targets(tx)
	first, ok := line.Next()
	if !ok {
		return nil, line.UsageError()
	}
	if strings.HasPrefix(first, "@") && tx == nil {
		return nil, MessageNoTargets.F()
	}
	switch first[:min(len(first), 2)] {
	case "@p":
		pos := line.src.Position()
		playerDistances := make([]float64, len(players))
		for i, p := range players {
			playerDistances[i] = p.Position().Sub(pos).Len()
		}
		sort.Slice(players, func(i, j int) bool {
			return playerDistances[i] < playerDistances[j]
		})
		if len(players) == 0 {
			return nil, nil
		}
		return sliceutil.Convert[Target](players[0:1]), nil
	case "@e":
		return entities, nil
	case "@a":
		return sliceutil.Convert[Target](players), nil
	case "@s":
		return []Target{line.src}, nil
	case "@r":
		if len(players) == 0 {
			return nil, nil
		}
		return []Target{players[rand.IntN(len(players))]}, nil
	default:
		target, err := p.parsePlayer(first, players)
		if err != nil {
			return nil, err
		}
		return []Target{target}, nil
	}
}

func (p parser) parsePlayer(name string, players []NamedTarget) (Target, error) {
	if ind := slices.IndexFunc(players, func(target NamedTarget) bool {
		return strings.EqualFold(target.Name(), name)
	}); ind != -1 {
		return players[ind], nil
	}
	return nil, MessagePlayerNotFound.F()
}
