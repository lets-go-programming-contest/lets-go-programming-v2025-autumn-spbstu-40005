package converter

func ConvertXMLToJSON(inputFile, outputFile string) error {
	valCurs, err := ParseXML(inputFile)
	if err != nil {
		return err
	}

	currencies := ConvertAndSortCurrencies(valCurs)

	if err := SaveJSON(outputFile, currencies); err != nil {
		return err
	}

	return nil
}
