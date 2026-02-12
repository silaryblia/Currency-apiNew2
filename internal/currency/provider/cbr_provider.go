package provider

import (
	"Currency-apiNew2/internal/config"
	"Currency-apiNew2/internal/metrics"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"
)

type CBRProvider struct {
	client *http.Client
	config *config.CBRConfig
}

func NewCBRProvider(cfg *config.CBRConfig) *CBRProvider {
	if cfg == nil {
		cfg = config.DefaultCBRConfig()
	}

	if cfg.URL == "" {
		cfg.URL = "https://www.cbr.ru/scripts/XML_daily.asp"
	}

	return &CBRProvider{
		config: cfg,
		client: &http.Client{
			Timeout: cfg.HTTPConfig.Timeout,
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
}

func (p *CBRProvider) ForceRefresh(
	ctx context.Context,
) (map[string]float64, time.Time, error) {

	start := time.Now()
	metrics.ProviderRequests.WithLabelValues("cbr").Inc()

	defer func() {
		metrics.ProviderLatency.
			WithLabelValues("cbr").
			Observe(time.Since(start).Seconds())
	}()

	rates, date, err := p.fetch(ctx)
	if err != nil {
		metrics.ProviderErrors.
			WithLabelValues("cbr", classifyProviderError(err)).
			Inc()
		return nil, time.Time{}, err
	}

	return rates, date, nil
}

func (p *CBRProvider) loadRates(ctx context.Context) (map[string]float64, time.Time, error) {

	if len(p.config.RequiredCodes) == 0 {
		return nil, time.Time{}, fmt.Errorf("no required currencies configured")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.config.URL, nil)
	if err != nil {
		return nil, time.Time{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set(
		"User-Agent",
		"Mozilla/5.0 (compatible; CurrencyService/1.0; +https://example.com)",
	)
	req.Header.Set(
		"Accept",
		"application/xml,text/xml;q=0.9,*/*;q=0.8",
	)
	req.Header.Set(
		"Accept-Language",
		"ru-RU,ru;q=0.9,en;q=0.8",
	)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, time.Time{}, err
	}
	defer func() { _ = resp.Body.Close() }()
	//defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, time.Time{}, fmt.Errorf(
			"CBR bad status: %d, body: %s",
			resp.StatusCode,
			string(body),
		)
	}

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

	rateDate, err := time.Parse("02.01.2006", valCurs.Date)
	if err != nil {
		return nil, time.Time{}, err
	}

	codes := make(map[string]bool)
	for _, code := range p.config.RequiredCodes {
		codes[code] = true
	}

	result := make(map[string]float64)
	foundCurrencies := make(map[string]bool)

	for _, v := range valCurs.Valutes {
		if !codes[v.CharCode] {
			continue
		}

		foundCurrencies[v.CharCode] = true

		valStr := strings.TrimSpace(v.Value)
		valStr = strings.Replace(valStr, ",", ".", 1)
		valStr = strings.ReplaceAll(valStr, " ", "")

		val, err := strconv.ParseFloat(valStr, 64)
		if err != nil {
			return nil, time.Time{}, fmt.Errorf("failed to parse value for %s: %w", v.CharCode, err)
		}

		if v.Nominal == 0 {
			return nil, time.Time{}, fmt.Errorf("nominal is zero for %s", v.CharCode)
		}

		rate := val / float64(v.Nominal)
		result[v.CharCode] = rate
	}

	if len(foundCurrencies) != len(codes) {
		missing := []string{}
		for code := range codes {
			if !foundCurrencies[code] {
				missing = append(missing, code)
			}
		}
		return nil, time.Time{}, fmt.Errorf("missing currencies: %v", missing)
	}

	// RUB всегда 1.0
	//result["RUB"] = 1.0

	return result, rateDate, nil
}

func (p *CBRProvider) fetch(
	ctx context.Context,
) (map[string]float64, time.Time, error) {
	return p.loadRates(ctx)
}

func classifyProviderError(err error) string {
	if err == nil {
		return "unknown"
	}

	// context
	if err == context.DeadlineExceeded || err == context.Canceled {
		return "timeout"
	}

	// HTTP errors
	var httpErr *url.Error
	if errors.As(err, &httpErr) {
		return "network"
	}

	msg := strings.ToLower(err.Error())

	switch {
	case strings.Contains(msg, "timeout"):
		return "timeout"
	case strings.Contains(msg, "bad status"):
		return "http_status"
	case strings.Contains(msg, "decode"),
		strings.Contains(msg, "xml"),
		strings.Contains(msg, "parse"):
		return "decode"
	default:
		return "unknown"
	}
}

//func round2(v float64) float64 {
//	return math.Round(v*100) / 100
//}
