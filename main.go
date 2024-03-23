package main

import (
	"bytes"
	"context"
	"html/template"
	"io/ioutil"
	"log"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type PageData struct {
	Title      string
	Content    string
	PageNumber int
}

func main() {
	longContent := `
       This is a long content...
			 มีภาษาไทยนะจ๊ะ สวัสดีครับท่านผู้ชม ทั่วไป
    `

	data := PageData{
		Title:      "Sample Report",
		Content:    longContent,
		PageNumber: 1,
	}

	htmlContent, err := renderHTML(data)
	if err != nil {
		log.Fatal(err)
	}

	if err := generatePDF(htmlContent); err != nil {
		log.Fatal(err)
	}

	log.Println("PDF report generated successfully.")
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

	time.Sleep(2 * time.Second)

	var buf []byte
	if err := chromedp.Run(ctx, printToPDF(htmlContent, &buf)); err != nil {
		return err
	}

	if err := ioutil.WriteFile("report.pdf", buf, 0644); err != nil {
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
