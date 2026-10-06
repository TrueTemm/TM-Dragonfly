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

package packsecure

import (
	"archive/zip"
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/sandertv/gophertunnel/minecraft/resource"
)

const contentsMagic uint32 = 0x9BCFB9FC
const headerSize = 0x100
const keyAlphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var exempt = map[string]bool{"manifest.json": true, "pack_icon.png": true, "bug_pack_icon.png": true, "contents.json": true}

func Load(dir string) ([]*resource.Pack, error) {
	_ = os.MkdirAll(dir, 0o777)
	out := filepath.Join(dir, ".secured") // encrypted copies live here
	_ = os.MkdirAll(out, 0o777)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read dir: %w", err)
	}
	packs := make([]*resource.Pack, 0, len(entries))
	for _, e := range entries {
		name := e.Name()
		if name == ".secured" {
			continue
		}
		if !e.IsDir() {
			switch strings.ToLower(filepath.Ext(name)) {
			case ".mcpack", ".zip":
			default:
				continue
			}
		}
		src := filepath.Join(dir, name)
		id := strings.TrimSuffix(name, filepath.Ext(name))
		encPath := filepath.Join(out, id+".mcpack")
		keyPath := filepath.Join(out, id+".key")

		key, kerr := os.ReadFile(keyPath)
		if kerr != nil || len(key) != 32 || stale(src, encPath) {
			k := string(key)
			if len(k) != 32 {
				k = genKey()
			}
			enc, err := encryptPack(src, k)
			if err != nil {
				return nil, fmt.Errorf("secure %v: %w", name, err)
			}
			if err := os.WriteFile(encPath, enc, 0o644); err != nil {
				return nil, err
			}
			if err := os.WriteFile(keyPath, []byte(k), 0o600); err != nil {
				return nil, err
			}
			key = []byte(k)
		}
		pack, err := resource.ReadPath(encPath)
		if err != nil {
			return nil, fmt.Errorf("read secured %v: %w", name, err)
		}
		packs = append(packs, pack.WithContentKey(string(key)))
	}
	return packs, nil
}

func stale(src, enc string) bool {
	ei, err := os.Stat(enc)
	if err != nil {
		return true
	}
	return newest(src).After(ei.ModTime())
}

func newest(p string) time.Time {
	info, err := os.Stat(p)
	if err != nil {
		return time.Now()
	}
	if !info.IsDir() {
		return info.ModTime()
	}
	t := info.ModTime()
	_ = filepath.WalkDir(p, func(_ string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			if fi, e := d.Info(); e == nil && fi.ModTime().After(t) {
				t = fi.ModTime()
			}
		}
		return nil
	})
	return t
}

type file struct {
	path string
	data []byte
}

func encryptPack(in, contentKey string) ([]byte, error) {
	files, err := read(in)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	type entry struct {
		Path string `json:"path"`
		Key  string `json:"key,omitempty"`
	}
	var contents struct {
		Content []entry `json:"content"`
	}
	for _, f := range files {
		w, err := zw.Create(f.path)
		if err != nil {
			return nil, err
		}
		if exempt[strings.ToLower(path.Base(f.path))] || exempt[strings.ToLower(f.path)] {
			if _, err := w.Write(f.data); err != nil {
				return nil, err
			}
			contents.Content = append(contents.Content, entry{Path: f.path})
			continue
		}
		fk := genKey()
		enc, err := cfb8(([]byte(fk)), f.data)
		if err != nil {
			return nil, err
		}
		if _, err := w.Write(enc); err != nil {
			return nil, err
		}
		contents.Content = append(contents.Content, entry{Path: f.path, Key: fk})
	}
	payload, err := json.Marshal(contents)
	if err != nil {
		return nil, err
	}
	encPayload, err := cfb8([]byte(contentKey), payload)
	if err != nil {
		return nil, err
	}
	header := make([]byte, headerSize)
	binary.LittleEndian.PutUint32(header[4:], contentsMagic)
	id := packID(files)
	header[16] = byte(len(id))
	copy(header[17:], id)
	cw, err := zw.Create("contents.json")
	if err != nil {
		return nil, err
	}
	if _, err := cw.Write(append(header, encPayload...)); err != nil {
		return nil, err
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func read(in string) ([]file, error) {
	info, err := os.Stat(in)
	if err != nil {
		return nil, err
	}
	var files []file
	if info.IsDir() {
		err = filepath.WalkDir(in, func(p string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}
			data, err := os.ReadFile(p)
			if err != nil {
				return err
			}
			rel, _ := filepath.Rel(in, p)
			files = append(files, file{path: filepath.ToSlash(rel), data: data})
			return nil
		})
	} else {
		files, err = readZip(in)
	}
	if err != nil {
		return nil, err
	}
	return rootAtManifest(files)
}

func readZip(in string) ([]file, error) {
	zr, err := zip.OpenReader(in)
	if err != nil {
		return nil, err
	}
	defer zr.Close()
	var files []file
	for _, zf := range zr.File {
		if zf.FileInfo().IsDir() {
			continue
		}
		rc, err := zf.Open()
		if err != nil {
			return nil, err
		}
		data, err := io.ReadAll(rc)
		_ = rc.Close()
		if err != nil {
			return nil, err
		}
		files = append(files, file{path: filepath.ToSlash(zf.Name), data: data})
	}
	return files, nil
}

func rootAtManifest(files []file) ([]file, error) {
	prefix := ""
	for _, f := range files {
		if strings.ToLower(path.Base(f.path)) == "manifest.json" {
			prefix = strings.TrimSuffix(f.path, path.Base(f.path))
			break
		}
	}
	out := make([]file, 0, len(files))
	for _, f := range files {
		if prefix != "" && !strings.HasPrefix(f.path, prefix) {
			continue
		}
		f.path = strings.TrimPrefix(f.path, prefix)
		if f.path == "" {
			continue
		}
		out = append(out, f)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no manifest.json in pack")
	}
	return out, nil
}

func packID(files []file) string {
	for _, f := range files {
		if strings.ToLower(f.path) != "manifest.json" {
			continue
		}
		var m struct {
			Header struct {
				UUID string `json:"uuid"`
			} `json:"header"`
		}
		if json.Unmarshal(f.data, &m) == nil && m.Header.UUID != "" {
			return m.Header.UUID
		}
	}
	return "00000000-0000-0000-0000-000000000000"
}

func cfb8(key, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	iv := make([]byte, aes.BlockSize)
	copy(iv, key[:aes.BlockSize])
	out := make([]byte, len(data))
	ks := make([]byte, aes.BlockSize)
	for i := range data {
		block.Encrypt(ks, iv)
		c := data[i] ^ ks[0]
		out[i] = c
		copy(iv, iv[1:])
		iv[aes.BlockSize-1] = c // feed one cipher byte back
	}
	return out, nil
}

func genKey() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	for i := range b {
		b[i] = keyAlphabet[int(b[i])%len(keyAlphabet)]
	}
	return string(b)
}
