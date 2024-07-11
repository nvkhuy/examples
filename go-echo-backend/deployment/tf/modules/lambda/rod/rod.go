package main

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/rotisserie/eris"
	"github.com/ysmood/gson"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/devices"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/launcher/flags"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/stealth"
)

var defaultPaperWidth = 8.5
var defaultPaperHeight = 11.3

type Option struct {
	Name   flags.Flag
	Values []string
}

func getChromeLauncher(isLandscapeViewport bool) (*launcher.Launcher, error) {
	path, found := launcher.LookPath()
	if !found {
		return nil, fmt.Errorf("launcher is not found")
	}

	var l = launcher.New().
		// where lambda runtime stores chromium
		Bin(path).

		// recommended flags to run in serverless environments
		// see https://github.com/alixaxel/chrome-aws-lambda/blob/master/source/index.ts
		Set("allow-running-insecure-content").
		Set("autoplay-policy", "user-gesture-required").
		Set("disable-component-update").
		Set("disable-domain-reliability").
		Set("disable-features", "AudioServiceOutOfProcess", "IsolateOrigins", "site-per-process").
		Set("disable-print-preview").
		Set("disable-setuid-sandbox").
		Set("disable-site-isolation-trials").
		Set("disable-speech-api").
		Set("disable-web-security").
		Set("disk-cache-size", "33554432").
		Set("enable-features", "SharedArrayBuffer").
		Set("hide-scrollbars").
		Set("ignore-gpu-blocklist").
		Set("in-process-gpu").
		Set("mute-audio").
		Set("no-default-browser-check").
		Set("no-pings").
		Set("no-sandbox").
		Set("no-zygote").
		Set("single-process").
		Set("use-gl", "swiftshader").
		Set("window-size", "1920", "1080")

	return l, nil
}

type GetPDFParams struct {
	URL      string `json:"url" mapstructure:"url" validate:"required"`
	Selector string `json:"selector" mapstructure:"selector" validate:"required"`

	JWTToken  string `json:"jwt_token" mapstructure:"jwt_token"`
	Landscape string `json:"landscape" mapstructure:"landscape"`

	PrintBackground   string `json:"print_background" mapstructure:"print_background"`
	PreferCSSPageSize string `json:"prefer_css_page_size" mapstructure:"prefer_css_page_size"`

	PaperWidth   string `json:"paper_width" mapstructure:"paper_width"`
	PaperHeight  string `json:"paper_height" mapstructure:"paper_height"`
	MarginTop    string `json:"margin_top" mapstructure:"margin_top"`
	MarginBottom string `json:"margin_bottom" mapstructure:"margin_bottom"`
	MarginLeft   string `json:"margin_left" mapstructure:"margin_left"`
	MarginRight  string `json:"margin_right" mapstructure:"margin_right"`
}

func getPDF(params GetPDFParams) ([]byte, error) {
	landscape, _ := strconv.ParseBool(params.Landscape)
	preferCSSPageSize, _ := strconv.ParseBool(params.PreferCSSPageSize)
	printBackground, _ := strconv.ParseBool(params.PrintBackground)

	const (
		navigateTimeout    = 10 * time.Second
		navigationTimeout  = 10 * time.Second
		requestIdleTimeout = 10 * time.Second
		htmlTimeout        = 20 * time.Second
	)

	l, err := getChromeLauncher(landscape)
	if err != nil {
		return nil, eris.Wrap(err, "Get launcher error")
	}

	// https://github.com/go-rod/go-rod.github.io/blob/main/error-handling.md
	urlLaunch, err := l.Launch()
	if err != nil {
		return nil, eris.Wrap(err, "Launch error")
	}

	var browser = rod.New().ControlURL(urlLaunch).DefaultDevice(devices.IPadMini)
	err = browser.Connect()
	if err != nil {
		return nil, eris.Wrap(err, "Browser connect error")
	}

	defer browser.MustClose()

	page, err := stealth.Page(browser)
	if err != nil {
		return nil, eris.Wrap(err, "Stealth page error")
	}
	defer page.Close()

	err = page.Emulate(devices.Clear)
	if err != nil {
		return nil, eris.Wrap(err, "Emulate device error")
	}

	err = page.Timeout(navigateTimeout).Navigate(params.URL)
	if err != nil {
		return nil, eris.Wrapf(err, "Navigate url %s error", params.URL)
	}

	if params.JWTToken != "" {
		page.MustEval(`jwt => window.localStorage.setItem('TOKEN_KEY', jwt)`, params.JWTToken)
		fmt.Println("Get token", page.MustEval(`jwt => window.localStorage.getItem('TOKEN_KEY')`))
	}

	var wait = page.Timeout(navigationTimeout).MustWaitNavigation()
	wait()

	waitRequestIdle := page.Timeout(requestIdleTimeout).MustWaitRequestIdle()
	waitRequestIdle()

	if params.Selector != "" {
		err = page.Timeout(htmlTimeout).WaitElementsMoreThan(params.Selector, 0)
		if err != nil {
			return nil, eris.Wrapf(err, "Wait selector %s error", params.Selector)
		}
	}

	var printParams = &proto.PagePrintToPDF{
		Landscape:         landscape,
		PrintBackground:   printBackground,
		PreferCSSPageSize: preferCSSPageSize,
		PaperWidth:        gson.Num(defaultPaperWidth),
		PaperHeight:       gson.Num(defaultPaperHeight),
		MarginTop:         gson.Num(0),
		MarginBottom:      gson.Num(0),
		MarginLeft:        gson.Num(0),
		MarginRight:       gson.Num(0),
	}

	if f, err := strconv.ParseFloat(params.PaperWidth, 64); err == nil {
		printParams.PaperWidth = gson.Num(f)
	}

	if f, err := strconv.ParseFloat(params.PaperHeight, 64); err == nil {
		printParams.PaperHeight = gson.Num(f)
	}

	if f, err := strconv.ParseFloat(params.MarginTop, 64); err == nil {
		printParams.MarginTop = gson.Num(f)
	}

	if f, err := strconv.ParseFloat(params.MarginBottom, 64); err == nil {
		printParams.MarginBottom = gson.Num(f)
	}

	if f, err := strconv.ParseFloat(params.MarginLeft, 64); err == nil {
		printParams.MarginLeft = gson.Num(f)
	}

	if f, err := strconv.ParseFloat(params.MarginRight, 64); err == nil {
		printParams.MarginRight = gson.Num(f)
	}

	reader, err := page.PDF(printParams)
	if err != nil {
		return nil, eris.Wrapf(err, "Get pdf error")
	}

	return io.ReadAll(reader)
}
