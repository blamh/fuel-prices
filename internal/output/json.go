package output

import (
	"encoding/json"
	"io"

	"fuel-prices/internal/model"
)

func WriteProducts(w io.Writer, prices []model.Price) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc.Encode(prices)
}
