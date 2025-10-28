package currency

import (
    "encoding/xml"
    "fmt"
    "strconv"
    "strings"
)

type ValCurs struct {
    XMLName xml.Name `xml:"ValCurs"`
    Date    string   `xml:"Date,attr"`
    Name    string   `xml:"name,attr"`
    Valutes []Valute `xml:"Valute"`
}

type Valute struct {
    ID       string  `xml:"ID,attr"`
    NumCode  int     `json:"num_code" xml:"NumCode"`
    CharCode string  `json:"char_code" xml:"CharCode"`
    Nominal  int     `xml:"Nominal"`
    Name     string  `xml:"Name"`
    Value    float64 `json:"value" xml:"Value"`
}

type currencyValue float64

func (v *currencyValue) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var str string
    if err := d.DecodeElement(&str, &start); err != nil {
        return fmt.Errorf("failed to decode value: %w", err)
    }

    str = strings.ReplaceAll(str, ",", ".")

    value, err := strconv.ParseFloat(str, 64)
    if err != nil {
        return fmt.Errorf("failed to parse value '%s': %w", str, err)
    }

    *v = currencyValue(value)
    return nil
}

type valuteXML struct {
    ID       string        `xml:"ID,attr"`
    NumCode  int           `xml:"NumCode"`
    CharCode string        `xml:"CharCode"`
    Nominal  int           `xml:"Nominal"`
    Name     string        `xml:"Name"`
    Value    currencyValue `xml:"Value"`
}

func (v *Valute) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    var temp valuteXML
    if err := d.DecodeElement(&temp, &start); err != nil {
        return err
    }

    v.ID = temp.ID
    v.NumCode = temp.NumCode
    v.CharCode = temp.CharCode
    v.Nominal = temp.Nominal
    v.Name = temp.Name
    v.Value = float64(temp.Value)

    return nil
}
