package main

import (
	"bytes"
	"context"
	"html/template"
	"log"
	"os"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type ReportItem struct {
	Examination string
	Value       string
	Unit        *string
	Normal      *string
	Description string
}

type PageData struct {
	Title                 string
	Content               string
	DemoItems             []ReportItem
	DemoWithUnitAndNormal []ReportItem
}

func stringToPtr(s string) *string {
	return &s
}

func add(a, b int) int {
	return a + b
}

func main() {
	data := constructData()

	htmlContent, err := renderHTML(data)
	if err != nil {
		log.Fatal(err)
	}

	if err := generatePDF(htmlContent); err != nil {
		log.Fatal(err)
	}

	log.Println("PDF report generated successfully.")
}

func constructData() PageData {
	return PageData{
		Title:   "Sample Report",
		Content: "This is a long content...",
		DemoItems: []ReportItem{
			{
				Examination: "Blood Pressure",
				Value:       "120/80",
				Description: "Description naja eiei",
			},
			{
				Examination: "Heart Rate",
				Value:       "72",
				Description: "Description naja eiei Description naja eiei2 Description naja eiei3",
			},
			{
				Examination: "Temperature",
				Value:       "36.5",
				Description: "Description naja eiei",
			},
			{
				Examination: "Weight",
				Value:       "67",
				Description: "Description naja eiei",
			},
			{
				Examination: "Height",
				Value:       "171",
				Description: "Description naja eiei",
			},
			{
				Examination: "BMI",
				Value:       "22.9",
				Description: "Description naja eiei",
			},
			{
				Examination: "Blood Sugar",
				Value:       "90",
				Description: "Description naja eiei",
			},
			{
				Examination: "Cholesterol",
				Value:       "200",
				Description: "Description naja eiei",
			},
			{
				Examination: "Uric Acid",
				Value:       "5.5",
				Description: "Description naja eiei",
			},
		},
		DemoWithUnitAndNormal: []ReportItem{
			{
				Examination: "Hemoglobin",
				Value:       "15",
				Unit:        stringToPtr("g/dL"),
				Normal:      stringToPtr("12-16"),
				Description: "Description naja eiei",
			},
			{
				Examination: "Hematocrit",
				Value:       "45",
				Unit:        stringToPtr("%"),
				Normal:      stringToPtr("40-50"),
				Description: "Description naja eiei",
			},
			{
				Examination: "RBC",
				Value:       "5.5",
				Unit:        stringToPtr("M/uL"),
				Normal:      stringToPtr("4.5-6.0"),
				Description: "Description naja eiei Description naja eiei Description naja eiei Description naja eiei final",
			},
			{
				Examination: "WBC",
				Value:       "7.5",
				Unit:        stringToPtr("K/uL"),
				Normal:      stringToPtr("4.0-10.0"),
				Description: "Description naja eiei",
			},
			{
				Examination: "Platelet",
				Value:       "250",
				Unit:        stringToPtr("K/uL"),
				Normal:      stringToPtr("150-400"),
				Description: "Description naja eiei",
			},
			{
				Examination: "MCV",
				Value:       "85",
				Unit:        stringToPtr("fL"),
				Normal:      stringToPtr("80-100"),
				Description: "Description naja eiei",
			},
		},
	}
}

func renderHTML(data PageData) (string, error) {
	var tmplFile = "report_template.html"
	funcMap := template.FuncMap{
		"add": add,
	}
	tmpl, err := template.New(tmplFile).Funcs(funcMap).ParseFiles(tmplFile)

	if err != nil {
		return "", err
	}

	var resultHTML bytes.Buffer
	if err := tmpl.ExecuteTemplate(&resultHTML, tmplFile, data); err != nil {
		return "", err
	}
	return resultHTML.String(), nil
}

func generatePDF(htmlContent string) error {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	time.Sleep(2 * time.Second)

	var buf []byte
	if err := chromedp.Run(ctx, printToPDF(htmlContent, &buf)); err != nil {
		return err
	}

	if err := os.WriteFile("report.pdf", buf, 0644); err != nil {
		return err
	}

	return nil
}

func printToPDF(htmlContent string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}

			return page.SetDocumentContent(frameTree.Frame.ID, htmlContent).Do(ctx)
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().
				WithDisplayHeaderFooter(true).
				WithPrintBackground(true).
				WithHeaderTemplate(" ").
				WithFooterTemplate(`
				<div style="width: 100%; font-size: 10px; font-weight: 100; text-align: center;">
					<span class="pageNumber"></span>
				</div>
				`).
				WithPrintBackground(false).Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}
