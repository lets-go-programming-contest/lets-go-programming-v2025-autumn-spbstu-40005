package currency

import (
    "encoding/json"
    "encoding/xml"
    "fmt"
    "os"
    "path/filepath"
    "sort"
    "strconv"
    "strings"

    "golang.org/x/text/encoding/charmap"
)

type ValCurs struct {
    Date    string   `xml:"Date,attr"`
    Name    string   `xml:"name,attr"`
    Valutes []Valute `xml:"Valute"`
}

type Valute struct {
    ID       string `xml:"ID,attr"`
    NumCode  string `xml:"NumCode"`
    CharCode string `xml:"CharCode"`
    Nominal  int    `xml:"Nominal"`
    Name     string `xml:"Name"`
    Value    string `xml:"Value"`
}

type CurrencyResult struct {
    NumCode  int     `json:"num_code"`
    CharCode string  `json:"char_code"`
    Value    float64 `json:"value"`
}

func decodeWindows1251(data []byte) ([]byte, error) {
    decoder := charmap.Windows1251.NewDecoder()
    return decoder.Bytes(data)
}

func ParseXMLData(filePath string) (*ValCurs, error) {
    if filePath == "" {
        return nil, fmt.Errorf("input file path is empty")

    data, err := os.ReadFile(filePath)
    if err != nil {
        if os.IsNotExit(err) {
            return nil, fmt.Errorf("file does not exist")
        }
        return nil, fmt.Errorf("cannot read input file")
    }

    decodedData, err := decodeWindows1251(data)
    if err != nil {
        decodedData = data
    }

    var valCurs ValCurs
    err = xml.Unmarshal(decodedData, &valCurs)
    if err != nil {
        return nil, fmt.Errorf("decoding to XML failed")
    }

    return &valCurs, nil
}

func convertValue(valueStr string) (float64, error) {
    if valueStr == "" {
        return 0, fmt.Errorf("empty value string")
    }

    normalizStr := strings.Replace(valueStr, ",", ".", -1)
    value, err := strconv.ParseFloat(normalizStr, 64)
    if err != nil {
        return 0, fmt.Errorf("conversion of the number to float failed")
    }
    return value, nil
}

func convertNumCode(numCodeStr string) (int, error) {
    if numCodeStr == "" {
        return 0, fmt.Errorf("empty num code string")
    }

    numCode, err := strconv.Atoi(numCodeStr)
    if err != nil {
        return 0, fmt.Errorf("conversion of the number to int failed")
    }
    return numCode, nil
}

func ProcessCurrencies(valCurs *ValCurs) ([]CurrencyResult, error) {
    if valCurs == nil {
        return nil, fmt.Errorf("nil valCurs")
    }

    var results []CurrencyResult
    for _, valute := range valCurs.Valutes {
        value, err := convertValue(valute.Value)
        if err != nil {
            return nil, err
        }
        numCode, err := convertNumCode(valute.NumCode)
        if err != nil {
            return nil, err
        }
        result := CurrencyResult{
            NumCode:  numCode,
            CharCode: valute.CharCode,
            Value:    value,
        }
        results = append(results, result)
    }

    sort.Slice(results, func(i, j int) bool {
        return results[i].Value > results[j].Value
    })

    return results, nil
}

func SaveResults(results []CurrencyResult, outputPath string) error {
    if outputPath == "" {
        return fmt.Errorf("output file path is empty")
    }

    dir := filepath.Dir(outputPath)
    err := os.MkdirAll(dir, 0755)
    if err != nil {
        return fmt.Errorf("creation of the directory failed")
    }

    file, err := os.Create(outputPath)
    if err != nil {
        return fmt.Errorf("creation of the file failed")
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    err = encoder.Encode(results)
    if err != nil {
        return fmt.Errorf("encoding in JSON failed")
    }

    return nil
}

