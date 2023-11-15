package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type InType uint8

const (
	InTypeDetails InType = iota
	InTypeCreateRoom
	InTypeJoinRoom
	InTypeGameReady
	InTypeStartRound
	InTypeRoundFinished
)

type InPkg struct {
	Type InType
	Data any
}

func ParseInPkg(data []byte) (*InPkg, error) {
	if len(data) < 1 {
		return nil, fmt.Errorf("data too short")
	}

	buf := bytes.NewBuffer(data)
	typeb, _ := buf.ReadByte()
	inType := InType(typeb)

	var inPkg InPkg
	inPkg.Type = inType

	if len(data) > 1 {
		dataBytes := buf.Next(len(data) - 1)
		dataRes, err := parseInPkg(dataBytes, inType)

		if err != nil {
			return nil, err
		}

		inPkg.Data = dataRes
	}

	return &inPkg, nil
}

func parseInPkg(data []byte, t InType) (any, error) {
	switch t {
	case InTypeDetails:
		return parseInDataDetails(data)

	case InTypeCreateRoom:
		return parseInDataRoom(data)

	case InTypeJoinRoom:
		return parseInDataRoom(data)

	case InTypeRoundFinished:
		return parseInDataRoundFinished(data)

	default:
		return nil, fmt.Errorf("unknown type")
	}
}

func parseInDataRoom(data []byte) (InDataRoom, error) {
	var inData InDataRoom
	err := json.Unmarshal(data, &inData)

	if err != nil {
		return InDataRoom{}, err
	}

	return inData, nil
}

func parseInDataDetails(data []byte) (InDataDetails, error) {
	var inData InDataDetails
	err := json.Unmarshal(data, &inData)

	if err != nil {
		return InDataDetails{}, err
	}

	return inData, nil
}

func parseInDataRoundFinished(data []byte) (InDataRoundFinished, error) {
	var inData InDataRoundFinished
	err := json.Unmarshal(data, &inData)

	if err != nil {
		return InDataRoundFinished{}, err
	}

	return inData, nil
}
