package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	json "github.com/json-iterator/go"
	"github.com/transerver/commons/redis"
	"strings"
	"time"
)

type RsaKeyNotExist struct {
	requestId string
}

func (e RsaKeyNotExist) Error() string {
	return fmt.Sprintf("the requestId[%q] is not exist", e.requestId)
}

type RsaObj struct {
	requestId  string
	PrivateKey []byte
	PublicKey  []byte
}

func (o *RsaObj) MarshalBinary() (data []byte, err error) {
	return json.Marshal(o)
}

func (o *RsaObj) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, o)
}

// Encrypt returns to encode data and if an error
func (o *RsaObj) Encrypt(data []byte) ([]byte, error) {
	block, _ := pem.Decode(o.PublicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.EncryptPKCS1v15(rand.Reader, pubKey, data)
}

// Decrypt returns to decode data and if an error
func (o *RsaObj) Decrypt(ciphertext []byte) ([]byte, error) {
	block, _ := pem.Decode(o.PrivateKey)
	if block == nil {
		return nil, errors.New("private key error")
	}
	private, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, private, ciphertext)
}

func (o *RsaObj) Release() error {
	return redis.Client().Del(o.requestId).Err()
}

type globalOption struct {
	prefix     string
	bits       int
	expiration time.Duration
}

type generator struct {
	requestId  string
	bits       int
	expiration time.Duration
	renew      bool
}

type Option func(*generator)

func WithBits(bits int) Option {
	return func(g *generator) {
		g.bits = bits
	}
}

// WithExpiration settings the global rsa expiration
// Zero expiration means the key has no expiration time.
func WithExpiration(expiration time.Duration) Option {
	return func(g *generator) {
		g.expiration = expiration
	}
}

// WithNoGen when the requestId is not exist don't create
func WithNoGen(g *generator) {
	g.renew = false
}

var global = &globalOption{bits: 1024, expiration: time.Minute * 10}

func SetRsaBits(bits int) {
	global.bits = bits
}

func SetExpiration(expiration time.Duration) {
	global.expiration = expiration
}

// SetRsaKeyPrefix settings the redis cached rsa key
func SetRsaKeyPrefix(prefix string) {
	if !strings.HasSuffix(prefix, ":") {
		prefix += ":"
	}
	global.prefix = prefix
}

func RsaKeyPrefix() string {
	return global.prefix
}

// FetchRsaKey called globalOption.fetch
func FetchRsaKey(requestId string, opts ...Option) (*RsaObj, error) {
	if len(requestId) == 0 {
		return nil, errors.New("empty requestId for fetch rsa key")
	}

	g := &generator{requestId: requestId, renew: true, expiration: time.Duration(-1)}
	for _, opt := range opts {
		opt(g)
	}
	g.init()
	return g.fetch()
}

func (g *generator) init() {
	if len(global.prefix) > 0 && strings.HasPrefix(g.requestId, ":") {
		g.requestId = g.requestId[1:]
	}
	g.requestId = global.prefix + g.requestId
	if g.bits <= 0 {
		g.bits = global.bits
	}
	if g.expiration < 0 {
		g.expiration = global.expiration
	}
}

// fetch first pull from redis, if is not exist, then create and saved to redis
func (g *generator) fetch() (*RsaObj, error) {
	cmd := redis.Client().Get(g.requestId)
	var rsaObj RsaObj
	if cmd.Err() == redis.Nil {
		if !g.renew {
			return nil, RsaKeyNotExist{g.requestId}
		}

		key, err := g.createRsaKey()
		if err != nil {
			return nil, err
		}
		rsaObj = *key
	} else if cmd.Err() != nil {
		return nil, cmd.Err()
	} else {
		err := cmd.Scan(&rsaObj)
		if err != nil {
			return nil, err
		}
	}
	rsaObj.requestId = g.requestId
	return &rsaObj, nil
}

// createRsaKey only generate rsa key and saved to redis
func (g *generator) createRsaKey() (*RsaObj, error) {
	key, err := g.genRsaKey()
	if err != nil {
		return nil, err
	}
	status, err := redis.Client().Set(g.requestId, key, g.expiration).Result()
	if err != nil {
		return nil, err
	}
	if "OK" != status {
		return nil, errors.New("fetch rsa key error")
	}
	return key, nil
}

// genRsaKey generated rsa private and public key
func (g *generator) genRsaKey() (*RsaObj, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, g.bits)
	if err != nil {
		return nil, err
	}

	// private key
	privateData := x509.MarshalPKCS1PrivateKey(privateKey)
	privateBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateData,
	}
	var privateBuf bytes.Buffer
	err = pem.Encode(&privateBuf, privateBlock)
	if err != nil {
		return nil, err
	}

	// public key
	publicKey := &privateKey.PublicKey
	publicData := x509.MarshalPKCS1PublicKey(publicKey)
	publicBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicData,
	}
	var publicBuf bytes.Buffer
	err = pem.Encode(&publicBuf, publicBlock)
	if err != nil {
		return nil, err
	}

	return &RsaObj{
		PrivateKey: privateBuf.Bytes(),
		PublicKey:  publicBuf.Bytes(),
	}, nil
}
