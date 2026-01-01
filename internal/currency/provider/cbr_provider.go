package provider

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"
)

type CBRProvider struct {
	client *http.Client
}

func NewCBRProvider() *CBRProvider {
	return &CBRProvider{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type ValCurs struct {
	Date    string   `xml:"Date,attr"`
	Valutes []Valute `xml:"Valute"`
}

type Valute struct {
	CharCode string `xml:"CharCode"`
	Nominal  int    `xml:"Nominal"`
	Value    string `xml:"Value"`
}

func (p *CBRProvider) GetRates(ctx context.Context) (map[string]float64, time.Time, error) {
	return p.loadRates(ctx)
	//req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.cbr.ru/scripts/XML_daily.asp", nil)
	//resp, err := p.client.Do(req)
	//if err != nil {
	//	return nil, err
	//}
	//defer resp.Body.Close()
	//
	//// Декодируем Windows-1251
	//decoder := xml.NewDecoder(resp.Body)
	//decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
	//	if strings.EqualFold(charset, "windows-1251") {
	//		return charmap.Windows1251.NewDecoder().Reader(input), nil
	//	}
	//	return nil, fmt.Errorf("unsupported charset: %s", charset)
	//}
	//
	//var valCurs ValCurs
	//if err := decoder.Decode(&valCurs); err != nil {
	//	return nil, err
	//}
	//
	//// валюты, которые нужны
	//codes := map[string]bool{"USD": true, "EUR": true, "AED": true}
	//result := make(map[string]float64)
	//
	//for _, v := range valCurs.Valutes {
	//	if codes[v.CharCode] {
	//		// заменить запятую на точку
	//		valStr := strings.Replace(v.Value, ",", ".", 1)
	//		val, _ := strconv.ParseFloat(valStr, 64)
	//		// делим на номинал
	//		result[v.CharCode] = val / float64(v.Nominal)
	//	}
	//}
	//
	//return result, nil
}

func (p *CBRProvider) ForceRefresh(ctx context.Context) (map[string]float64, time.Time, error) {
	//return p.GetRates(ctx)
	return p.loadRates(ctx)
}

func (p *CBRProvider) loadRates(ctx context.Context) (map[string]float64, time.Time, error) {
	req, _ := http.NewRequestWithContext(ctx,
		http.MethodGet,
		"https://www.cbr.ru/scripts/XML_daily.asp",
		nil)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, time.Time{}, err
	}
	defer resp.Body.Close()

	decoder := xml.NewDecoder(resp.Body)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if strings.EqualFold(charset, "windows-1251") {
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		}
		return nil, fmt.Errorf("unsupported charset: %s", charset)
	}

	var valCurs ValCurs

	if err := decoder.Decode(&valCurs); err != nil {
		return nil, time.Time{}, err
	}

	////
	rateDate, err := time.Parse("02.01.2006", valCurs.Date)
	if err != nil {
		return nil, time.Time{}, err
	}

	codes := map[string]bool{"USD": true, "EUR": true, "AED": true}
	result := make(map[string]float64)

	for _, v := range valCurs.Valutes {
		if codes[v.CharCode] {

			valStr := strings.TrimSpace(v.Value)
			valStr = strings.Replace(valStr, ",", ".", 1)

			val, err := strconv.ParseFloat(valStr, 64)
			if err != nil {
				return nil, time.Time{}, fmt.Errorf(
					"failed to parse rate %s for %s",
					v.Value,
					v.CharCode,
				)
			}

			if v.Nominal == 0 {
				return nil, time.Time{}, fmt.Errorf("nominal is zero for %s", v.CharCode)
			}

			rate := val / float64(v.Nominal)
			rate = round2(rate)

			result[v.CharCode] = rate
		}
	}
	return result, rateDate, nil
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
