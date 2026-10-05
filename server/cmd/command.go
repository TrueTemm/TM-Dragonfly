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
	"encoding/csv"
	"fmt"
	"go/ast"
	"reflect"
	"slices"
	"strings"

	"github.com/df-mc/dragonfly/server/world"
)

type Runnable interface {
	Run(src Source, o *Output, tx *world.Tx)
}

type Allower interface {
	Allow(src Source) bool
}

type Command struct {
	v           []reflect.Value
	name        string
	description string
	usage       string
	aliases     []string
}

func New(name, description string, aliases []string, r ...Runnable) Command {
	name = strings.ToLower(name)
	for i, alias := range aliases {
		aliases[i] = strings.ToLower(alias)
	}

	usages := make([]string, len(r))
	runnableValues := make([]reflect.Value, len(r))

	if len(aliases) > 0 && slices.Index(aliases, name) == -1 {
		aliases = append(aliases, name)
	}

	for i, runnable := range r {
		t := reflect.TypeOf(runnable)
		if t.Kind() != reflect.Struct && (t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct) {
			panic(fmt.Sprintf("Runnable r must be struct or pointer to struct, but got %v", t.Kind()))
		}
		original := reflect.ValueOf(runnable)
		if t.Kind() == reflect.Ptr {
			original = original.Elem()
		}

		cp := reflect.New(original.Type()).Elem()
		if err := verifySignature(cp); err != nil {
			panic(err.Error())
		}
		runnableValues[i], usages[i] = original, parseUsage(name, cp)
	}

	return Command{name: name, description: description, aliases: aliases, v: runnableValues, usage: strings.Join(usages, "\n")}
}

func (cmd Command) Name() string {
	return cmd.name
}

func (cmd Command) Description() string {
	return cmd.description
}

func (cmd Command) Usage() string {
	return cmd.usage
}

func (cmd Command) Aliases() []string {
	return cmd.aliases
}

func (cmd Command) Execute(args string, source Source, tx *world.Tx) {
	if source == nil {
		panic("execute: invalid command source: source must not be nil")
	}
	output := &Output{}
	defer source.SendCommandOutput(output)

	var leastErroneous error
	var leastArgsLeft *Line

	for _, v := range cmd.v {
		cp := reflect.New(v.Type())
		cp.Elem().Set(v)
		line, err := cmd.executeRunnable(cp, args, source, output, tx)
		if err == nil {

			return
		}
		if line == nil {

			if leastErroneous == nil {
				leastErroneous = err
			}
			continue
		}
		if leastArgsLeft == nil || line.Len() <= leastArgsLeft.Len() {

			leastErroneous = err
			leastArgsLeft = line
		}
	}

	if leastArgsLeft != nil {
		output.Error(leastArgsLeft.SyntaxError())
	}
	output.Error(leastErroneous)
}

type ParamInfo struct {
	Name     string
	Value    any
	Optional bool
	Suffix   string
}

func (cmd Command) Params(src Source) [][]ParamInfo {
	params := make([][]ParamInfo, 0, len(cmd.v))
	for _, runnable := range cmd.v {
		if allower, ok := runnable.Interface().(Allower); ok && !allower.Allow(src) {

			continue
		}

		if d, ok := runnable.Interface().(ParamDescriber); ok {
			params = append(params, d.DescribeParams(src))
			continue
		}

		elem := reflect.New(runnable.Type()).Elem()
		elem.Set(runnable)

		var fields []ParamInfo
		for _, t := range exportedFields(elem) {
			field := elem.FieldByName(t.Name)
			fields = append(fields, ParamInfo{
				Name:     name(t),
				Value:    unwrap(field).Interface(),
				Optional: optional(field),
				Suffix:   suffix(t),
			})
		}
		params = append(params, fields)
	}
	return params
}

func (cmd Command) Runnables(src Source) map[int]Runnable {
	m := make(map[int]Runnable, len(cmd.v))
	for i, runnable := range cmd.v {
		v := runnable.Interface().(Runnable)
		if allower, ok := v.(Allower); !ok || allower.Allow(src) {
			m[i] = v
		}
	}
	return m
}

func (cmd Command) String() string {
	return cmd.usage
}

func (cmd Command) executeRunnable(v reflect.Value, args string, source Source, output *Output, tx *world.Tx) (*Line, error) {
	if a, ok := v.Interface().(Allower); ok && !a.Allow(source) {
		return nil, MessageUnknown.F(cmd.name)
	}

	var argFrags []string
	if args != "" {
		r := csv.NewReader(strings.NewReader(args))
		r.Comma, r.LazyQuotes = ' ', true
		record, err := r.Read()
		if err != nil {

			return nil, MessageUsage.F(cmd.Usage())
		}
		argFrags = record
	}
	parser := parser{}
	arguments := &Line{args: argFrags, src: source, seen: []string{"/" + cmd.name}, cmd: cmd}

	signature := v.Elem()
	for _, t := range exportedFields(signature) {
		field := signature.FieldByName(t.Name)
		parser.currentField = t.Name
		opt := optional(field)

		val := field
		if opt {
			val = reflect.New(field.Field(0).Type()).Elem()
		}

		err, success := parser.parseArgument(arguments, val, opt, name(t), source, tx)
		if err != nil {

			return arguments, err
		}
		if success && opt {
			field.Set(reflect.ValueOf(field.Interface().(optionalT).with(val.Interface())))
		}
	}
	if arguments.Len() != 0 {
		return arguments, arguments.UsageError()
	}

	v.Interface().(Runnable).Run(source, output, tx)
	return arguments, nil
}

func parseUsage(commandName string, command reflect.Value) string {
	parts := make([]string, 0, command.NumField()+1)
	parts = append(parts, "/"+commandName)

	for _, t := range exportedFields(command) {
		field := command.FieldByName(t.Name)

		typeName := typeNameOf(field.Interface(), name(t))
		if _, ok := field.Interface().(optionalT); ok {
			typeName = typeNameOf(reflect.New(field.Field(0).Type()).Elem().Interface(), name(t))
		}
		if _, ok := field.Interface().(SubCommand); ok {
			parts = append(parts, typeName)
			continue
		}
		if optional(field) {
			parts = append(parts, "["+name(t)+": "+typeName+"]"+suffix(t))
			continue
		}
		parts = append(parts, "<"+name(t)+": "+typeName+">"+suffix(t))
	}
	return strings.Join(parts, " ")
}

func verifySignature(command reflect.Value) error {
	optionalField := false
	for _, t := range exportedFields(command) {
		field := command.FieldByName(t.Name)

		opt := optional(field)
		if !opt && optionalField {
			return fmt.Errorf("command must only have optional parameters at the end")
		}
		val := field
		if opt {
			val = reflect.New(field.Field(0).Type()).Elem()
		}
		if _, ok := val.Interface().(Enum); ok && val.Kind() != reflect.String {
			return fmt.Errorf("parameters implementing Enum must be of the type string")
		}
		optionalField = opt
	}
	return nil
}

func exportedFields(command reflect.Value) []reflect.StructField {
	visible := reflect.VisibleFields(command.Type())
	fields := make([]reflect.StructField, 0, len(visible))

	for _, t := range visible {
		if !ast.IsExported(t.Name) || name(t) == "-" || t.Anonymous {
			continue
		}
		field := command.FieldByName(t.Name)
		if !field.CanSet() {
			continue
		}
		fields = append(fields, t)
	}
	return fields
}
