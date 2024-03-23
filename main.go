package main

import (
	"bytes"
	"context"
	"html/template"
	"log"
	"os"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type ReportItem struct {
	Examination string
	Value       string
	Description string
}

type PageData struct {
	Title     string
	Content   string
	DemoItems []ReportItem
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
	}
}

func renderHTML(data PageData) (string, error) {
	tmpl, err := template.ParseFiles("report_template.html")
	if err != nil {
		return "", err
	}

	var resultHTML bytes.Buffer
	if err := tmpl.ExecuteTemplate(&resultHTML, "report_template.html", data); err != nil {
		return "", err
	}
	return resultHTML.String(), nil
}

func generatePDF(htmlContent string) error {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// time.Sleep(2 * time.Second)

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
