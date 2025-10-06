package parsers

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

type ValCurs struct {
	Date    string   `xml:"Date,attr"`
	Name    string   `xml:"name,attr"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	ID       string `xml:"ID,attr"`
	NumCode  int    `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Nominal  string `xml:"Nominal"`
	Name     string `xml:"Name"`
	Value    string `xml:"Value"`
}

type Currency struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func ParseConfigFile(path string) (*Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("file %s does not exist: %w", path, err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %w", err)
	}

	cfg := &Config{InputFile: "", OutputFile: ""}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal file: %w", err)
	}

	return cfg, nil
}

func ParseXMLFile(path string) (*ValCurs, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("file %s does not exist: %w", path, err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %w", err)
	}

	reader := bytes.NewReader(data)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = createCharsetReader

	valCurs := &ValCurs{Date: "", Name: "", Valutes: nil}
	err = decoder.Decode(valCurs)

	if err != nil {
		return nil, fmt.Errorf("cannot decode file: %w", err)
	}

	return valCurs, nil
}

func createCharsetReader(charset string, input io.Reader) (io.Reader, error) {
	if charset == "windows-1251" {
		return charmap.Windows1251.NewDecoder().Reader(input), nil
	}

	return input, nil
}

func convertFloat(float string) (float64, error) {
	partsOfFloat := strings.Split(float, ",")
	result, err := strconv.ParseFloat(partsOfFloat[0]+"."+partsOfFloat[1], 64)

	if err != nil {
		return 0, fmt.Errorf("cannot convert to float: %w", err)
	}

	return result, nil
}

func ConvertToJSON(valutes []Valute) ([]Currency, error) {
	currencies := make([]Currency, len(valutes))

	for idx, valute := range valutes {
		value, err := convertFloat(valute.Value)
		if err != nil {
			return nil, fmt.Errorf("convertation failed: %w", err)
		}

		currencies[idx] = Currency{
			NumCode:  valute.NumCode,
			CharCode: valute.CharCode,
			Value:    value,
		}
	}

	return currencies, nil
}

func SaveToJSON(currencies []Currency, outputPath string) error {
	jsonData, err := json.MarshalIndent(currencies, "", "  ")
	if err != nil {
		return fmt.Errorf("cannot marshal to JSON file: %w", err)
	}

	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0); err != nil {
		return fmt.Errorf("cannot make directory: %w", err)
	}

	if err = os.WriteFile(outputPath, jsonData, 0); err != nil {
		return fmt.Errorf("cannot write file: %w", err)
	}

	return nil
}
