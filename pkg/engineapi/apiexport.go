package engineapi

import (
	"ScriptExecServer/pkg/xmlhelper"
	"bytes"
	"encoding/xml"
)

func ReadEngineApiExportXml(path string) (*EngineExportScope, error) {
	data, err := xmlhelper.ReadFileWithBadUTF8(path)
	if err != nil {
		return nil, err
	}

	result := &EngineExportScope{}

	illegalUtf8SequencesToPurge := []string{"&#x03;", "&#x05;", "&#x10;", "&#x1E;", "&#x0F;"}

	for _, s := range illegalUtf8SequencesToPurge {
		data = bytes.ReplaceAll(data, []byte(s), []byte(""))
	}

	decoder := xml.NewDecoder(bytes.NewReader(data))
	decoder.Entity = map[string]string{
		"x03": "",
	}
	err = decoder.Decode(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
