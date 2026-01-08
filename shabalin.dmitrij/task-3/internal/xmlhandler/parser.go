package xmlhandler

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/dmitei/task-3/internal/models"
	"golang.org/x/text/encoding/charmap"
)

var ErrCharsetNotSupported = errors.New("charset is not supported")

func ParseXMLFile(xmlContent []byte, currencyContainer *models.CurrencyContainer) error {
	contentDecoder := xml.NewDecoder(bytes.NewReader(xmlContent))

	contentDecoder.CharsetReader = func(charsetName string, inputReader io.Reader) (io.Reader, error) {
		if strings.ToLower(charsetName) == "windows-1251" {
			decodedReader := charmap.Windows1251.NewDecoder().Reader(inputReader)

			return decodedReader, nil
		}

		return nil, fmt.Errorf("parsing charset %w: %q", ErrCharsetNotSupported, charsetName)
	}

	if decodeError := contentDecoder.Decode(currencyContainer); decodeError != nil {
		return fmt.Errorf("cannot decode XML content: %w", decodeError)
	}

	return nil
}
