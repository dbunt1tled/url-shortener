package hasher

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jxskiss/base62"
	"github.com/pkg/errors"
)

type JWTManager struct {
	algorithm     string
	publicKey     interface{}
	privateKey    interface{}
	signingMethod jwt.SigningMethod
}

type Hasher struct {
	jwtManager *JWTManager
}

func NewHasher(jwtAlgorithm string, jwtPublicKey string, jwtPrivateKey string) (*Hasher, error) {
	jm, err := getJWTManager(jwtAlgorithm, jwtPublicKey, jwtPrivateKey)
	if err != nil {
		return nil, err
	}
	return &Hasher{
		jwtManager: jm,
	}, nil
}

func (h *Hasher) RandomBytes(n uint) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (h *Hasher) RandomString(length uint) (string, error) {
	var realLength uint
	if length%2 == 0 {
		realLength = length / 2 //nolint:mnd // 1 symbol = 2 bytes
	} else {
		realLength = (length / 2) + 1 //nolint:mnd // 1 symbol = 2 bytes
	}
	bytes, err := h.RandomBytes(realLength)
	if err != nil {
		return "", err
	}
	return h.HexString(bytes)[:length], nil
}

func (h *Hasher) HexString(bytes []byte) string {
	return base62.EncodeToString(bytes)
}

func (h *Hasher) NewID() (string, error) {
	ms := time.Now().UnixMilli()
	tsBytes := make([]byte, 4)                               //nolint:mnd // reserved 4 bytes
	binary.BigEndian.PutUint32(tsBytes, uint32(ms&0xFFFFFF)) //nolint:mnd // timestamp

	entropy, err := h.RandomBytes(7) //nolint:mnd // generate 7 bytes
	if err != nil {
		return "", err
	}
	combined := append(tsBytes, entropy...)

	return base62.EncodeToString(combined), nil
}

// Encode time.Now().Add(time.Second * time.Duration(exp)).Unix().
func (h *Hasher) Encode(payload map[string]interface{}) (string, error) {
	claims := jwt.MapClaims(payload)
	token := jwt.NewWithClaims(h.jwtManager.signingMethod, claims)
	privateKey := h.jwtManager.privateKey
	return token.SignedString(privateKey)
}

func (h *Hasher) Decode(token string, checkExpire bool) (map[string]interface{}, error) {
	var claims jwt.MapClaims

	tokenData, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method != h.jwtManager.signingMethod {
			return nil, errors.New("invalid token algorithm")
		}
		return h.jwtManager.publicKey, nil
	})

	if !tokenData.Valid || err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) && !checkExpire {
			return claims, nil
		}
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func getJWTManager(jwtAlgorithm string, jwtPublicKey string, jwtPrivateKey string) (*JWTManager, error) {
	var err error
	jm := &JWTManager{
		algorithm: jwtAlgorithm,
	}
	switch jwtAlgorithm {
	case "RS512":
		jm.signingMethod = jwt.SigningMethodRS512
		jm.publicKey, jm.privateKey, err = loadRSAKeys(jwtPublicKey, jwtPrivateKey)
		if err != nil {
			return nil, err
		}
	case "RS256":
		jm.signingMethod = jwt.SigningMethodRS256
		jm.publicKey, jm.privateKey, err = loadRSAKeys(jwtPublicKey, jwtPrivateKey)
		if err != nil {
			return nil, err
		}
	case "ES512":
		jm.signingMethod = jwt.SigningMethodES512
		jm.publicKey, jm.privateKey, err = loadECKeys(jwtPublicKey, jwtPrivateKey)
		if err != nil {
			return nil, err
		}
	case "ES256":
		jm.signingMethod = jwt.SigningMethodES256
		jm.publicKey, jm.privateKey, err = loadECKeys(jwtPublicKey, jwtPrivateKey)
		if err != nil {
			return nil, err
		}

	default:
		return nil, errors.New("invalid algorithm")
	}
	return jm, nil
}

func loadRSAKeys(jwtPublicKey string, jwtPrivateKey string) (*rsa.PublicKey, *rsa.PrivateKey, error) {
	var (
		privateKeyBytes []byte
		publicKeyBytes  []byte
		publicKey       *rsa.PublicKey
		privateKey      *rsa.PrivateKey
		err             error
	)
	privateKeyBytes, err = base64.StdEncoding.DecodeString(jwtPrivateKey)
	if err != nil {
		return nil, nil, err
	}

	publicKeyBytes, err = base64.StdEncoding.DecodeString(jwtPublicKey)
	if err != nil {
		return nil, nil, err
	}

	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		return nil, nil, err
	}

	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return nil, nil, err
	}

	return publicKey, privateKey, nil
}

func loadECKeys(jwtPublicKey string, jwtPrivateKey string) (*ecdsa.PublicKey, *ecdsa.PrivateKey, error) {
	var (
		privateKeyBytes []byte
		publicKeyBytes  []byte
		publicKey       *ecdsa.PublicKey
		privateKey      *ecdsa.PrivateKey
		err             error
	)
	privateKeyBytes, err = base64.StdEncoding.DecodeString(jwtPrivateKey)
	if err != nil {
		return nil, nil, err
	}

	publicKeyBytes, err = base64.StdEncoding.DecodeString(jwtPublicKey)
	if err != nil {
		return nil, nil, err
	}

	publicKey, err = jwt.ParseECPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		return nil, nil, err
	}

	privateKey, err = jwt.ParseECPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return nil, nil, err
	}

	return publicKey, privateKey, nil
}
