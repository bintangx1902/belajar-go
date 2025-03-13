package Utils

import (
	"TestApp/Models"
	"fmt"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/list"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
	"strconv"
)

type InvoiceItem Models.PrintInvoiceItems

func ExportPDF(m core.Maroto, title, content string, contents [][]string) {
	addItemHeader(m)
	addItemContent(m, title, content)
	addItemsList(m, contents)
}

func ToPDF(m core.Maroto, title, content string) {
	addPDFHeader(m)
	addPDFContent(m, title, content)
	addItemList(m)
}

func addPDFHeader(m core.Maroto) {
	m.AddRow(50, image.NewFromFileCol(12, "assets/logo.jpg", props.Rect{
		Center:  true,
		Percent: 75,
	}))
}

func addItemHeader(m core.Maroto) {
	m.AddRow(50, image.NewFromFileCol(12, "assets/logo.jpg", props.Rect{
		Center:  true,
		Percent: 75,
	}))
}

func addPDFContent(m core.Maroto, title, content string) {
	m.AddRow(20,
		text.NewCol(12, title, props.Text{
			Top:   5,
			Style: fontstyle.Bold,
			Align: align.Center,
			Size:  16,
		}),
	)

	m.AddRow(30,
		text.NewCol(12, content, props.Text{
			Top:   5,
			Align: align.Center,
			Size:  12,
		}),
	)
}

func addItemContent(m core.Maroto, title, content string) {
	m.AddRow(20,
		text.NewCol(12, title, props.Text{
			Top:   5,
			Style: fontstyle.Bold,
			Align: align.Center,
			Size:  16,
		}),
	)

	m.AddRow(30,
		text.NewCol(12, content, props.Text{
			Top:   5,
			Align: align.Center,
			Size:  12,
		}),
	)
}

func addItemList(m core.Maroto) {
	content := [][]string{{}}
	rows, err := list.Build[InvoiceItem](getObjects(content)) // Use InvoiceItem here
	if err != nil {
		return
	}
	m.AddRows(rows...)
}

func addItemsList(m core.Maroto, content [][]string) {
	rows, err := list.Build[InvoiceItem](getObjects(content))
	if err != nil {
		return
	}
	m.AddRows(rows...)
}

func (i InvoiceItem) GetHeader() core.Row {
	return row.New(10).Add(
		text.NewCol(1, "No.", props.Text{Style: fontstyle.Bold}),
		text.NewCol(2, "Item", props.Text{Style: fontstyle.Bold}),
		text.NewCol(3, "Description", props.Text{Style: fontstyle.Bold}),
		text.NewCol(1, "Quantity", props.Text{Style: fontstyle.Bold}),
		text.NewCol(2, "Price", props.Text{Style: fontstyle.Bold}),
		text.NewCol(2, "Discounted Price", props.Text{Style: fontstyle.Bold}),
		text.NewCol(2, "Total", props.Text{Style: fontstyle.Bold}),
	)
}

func (o InvoiceItem) GetContent(i int) core.Row {
	r := row.New(6).Add(
		text.NewCol(1, o.No),
		text.NewCol(2, o.Item),
		text.NewCol(3, o.Description),
		text.NewCol(1, o.Quantity),
		text.NewCol(2, o.Price),
		text.NewCol(2, o.DiscountedPrice),
		text.NewCol(2, o.Total),
	)

	if i%2 == 0 {
		r.WithStyle(&props.Cell{
			BackgroundColor: &props.Color{Red: 240, Green: 240, Blue: 240},
		})
	}

	return r
}

func getObjects(content [][]string) []InvoiceItem {
	if len(content) == 0 || len(content[0]) == 0 {
		content = [][]string{
			{"Laptop", "16 inch", "1", "Rp. 1.000.000", "Rp. 990.000", "Rp. 990.000"},
			{"Mouse", "Logitech", "1", "Rp. 100.000", "Rp. 90.000", "Rp. 90.000"},
			{"Keyboard", "Logitech", "1", "Rp. 50.000", "Rp. 45.000", "Rp. 45.000"},
		}
	}

	var items []InvoiceItem
	for i := 0; i < len(content); i++ {
		if len(content[i]) >= 6 {
			items = append(items, InvoiceItem{
				No:              strconv.Itoa(i + 1),
				Item:            content[i][0],
				Description:     content[i][1],
				Quantity:        content[i][2],
				Price:           content[i][3],
				DiscountedPrice: content[i][4],
				Total:           content[i][5],
			})
		}
	}
	return items
}

func ConvertInvoiceItemsToContent(items []Models.Items) [][]string {
	var content [][]string

	for _, item := range items {
		discountedPrice := item.Price * (1 - item.Discount)
		totalPrice := discountedPrice * float64(item.Quantity)

		// Format the data as strings
		content = append(content, []string{
			item.Item,                       // Item name
			item.Description,                // Item description
			strconv.Itoa(item.Quantity),     // Quantity
			formatCurrency(item.Price),      // Price
			formatCurrency(discountedPrice), // Discounted Price
			formatCurrency(totalPrice),      // Total Price
		})
	}

	return content
}

func formatCurrency(value float64) string {
	return fmt.Sprintf("$ %.2f", value)
}
