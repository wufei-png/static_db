package utils

import (
	// #nosec
	"crypto/rc4"
	"encoding/json"
	"errors"
	"io/ioutil"
)

type NormalizeParam struct {
	SrcPoints []float32
	DstPoints []float32
}

type NormalizeParamFromSDK struct {
	SrcPoints []float32 `json:"src_points"`
	DstPoints []float32 `json:"dst_points"`
}

func (n *NormalizeParam) Normalize(score float32) float32 {
	size := len(n.SrcPoints)
	if score <= n.SrcPoints[0] {
		return n.DstPoints[0]
	}
	if score >= n.SrcPoints[size-1] {
		return n.DstPoints[size-1]
	}
	for i := 1; i < size; i++ {
		if score < n.SrcPoints[i] {
			return (score-n.SrcPoints[i-1])/(n.SrcPoints[i]-n.SrcPoints[i-1])*
				(n.DstPoints[i]-n.DstPoints[i-1]) + n.DstPoints[i-1]
		}
	}
	// XXX should be unreachable
	return -1.0
}

func LoadNormalizeParam(fn string) (*NormalizeParam, error) {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	// #nosec
	c, err := rc4.NewCipher([]byte("Sensetime"))
	if err != nil {
		return nil, err
	}
	c.XORKeyStream(data, data)
	var p NormalizeParam
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	if len(p.DstPoints) < 2 || len(p.DstPoints) != len(p.SrcPoints) {
		return nil, errors.New("invalid normalize param")
	}
	return &p, nil
}

func LoadNormalizeJSONFile(fn string) (*NormalizeParam, error) {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	n := NormalizeParam{}
	if err := json.Unmarshal(data, &n); err != nil {
		return nil, err
	}
	if len(n.DstPoints) < 2 || len(n.DstPoints) != len(n.SrcPoints) {
		return nil, errors.New("invalid normalize param")
	}
	return &n, nil
}
