package pdf

import (
	"log"
	"strings"
	"time"

	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

type PdfTableOptions struct {
	// default is "Bursa Effek Indonesia"
	HeaderTitle string
	// specify each collumn name
	HeaderRows []string
	// "A3", "A4", "Legal", "Letter", "A5" default is "A4"
	PageSize string
	// path to logo default is globe idx
	HeaderLogo string
	// path to logo
	FooterLogo string
	// line header and footer color default is maroon
	LineColor *color.Color
	// even indexed rows bg color in table, default is gray
	TableBgCol *color.Color
	// setting paper width
	PapperWidth float64
	// setting papper height
	Papperheight float64
	// setting page orientation by default "p" / "l"
	PageOrientation string
}

func createHeader(page pdf.Maroto, config PdfTableOptions, columnTotal uint) {
	title := config.HeaderTitle
	if title == "" {
		title = "Bursa Effek Indonesia"
	}

	logoPath := config.FooterLogo
	if logoPath == "" {
		logoPath = "./icon-globe-idx.png"
	}

	page.RegisterHeader(func() {
		logoWidth := 1
		if columnTotal >= 4 {
			logoWidth = int(columnTotal) / 4
		}

		page.Row(24, func() {
			page.Col(uint(logoWidth), func() {
				_ = page.FileImage(logoPath, props.Rect{
					// left space to make gap with line below
					Percent: 80,
				})
			})

			page.Col(columnTotal-uint(logoWidth), func() {
				page.Text(title, props.Text{
					Style: consts.Bold,
					Size:  18,
					Align: consts.Middle,
				})
			})
		})
		drawLine(page, config)
	})
}

func createFooter(page pdf.Maroto, config PdfTableOptions, columnTotal uint) {
	footerHeight := float64(8)
	footerPaddingTop := footerHeight / 4

	logoPath := config.FooterLogo
	if logoPath == "" {
		logoPath = "./idx-logo-2.png"
	}

	dateAndNameWidth := 1
	if columnTotal > 8 {
		dateAndNameWidth = 2
	}

	page.RegisterFooter(func() {
		page.Row(footerHeight, func() {
			drawLine(page, config)
		})
		page.Row(footerHeight, func() {
			page.Col(uint(dateAndNameWidth), func() {
				page.Text(time.Now().Format("02 January 2006"), props.Text{
					Align: consts.Left,
					Size:  8,
					Top:   footerPaddingTop,
				})
			})

			page.Col(1, func() {
				errImg := page.FileImage(logoPath, props.Rect{
					Top: footerPaddingTop,
				})
				if errImg != nil {
					log.Println(errImg)
				}
			})
			page.Col(uint(dateAndNameWidth), func() {
				page.Text("Sistem Portal Anggota Bursa", props.Text{
					Align: consts.Left,
					Size:  8,
					Top:   footerPaddingTop,
				})
			})
		})
	})
}

func drawLine(page pdf.Maroto, config PdfTableOptions) {
	pageWidth, _ := page.GetPageSize()
	lineColor := config.LineColor
	if lineColor == nil {
		lineColor = &color.Color{
			Red: 159, Green: 14, Blue: 15,
		}
	}

	page.Line(.5, props.Line{
		Width: pageWidth,
		Style: consts.Solid,
		// red color
		Color: *lineColor,
	})
}

func (opt *PdfTableOptions) getHeaderTitle() string {
	if opt.HeaderTitle == "" {
		return "Bursa Effek Indonesia"
	}
	return opt.HeaderTitle
}

func addUserPermissions(types string) []string {
	var result []string
	if strings.EqualFold(types, "internal") {
		result = append(result, "Maker", "Checker", "Approval")
		return result
	}

	if strings.EqualFold("external", types) {
		result = append(result, "Maker", "Checker")
		return result
	}

	return nil
}

func addUserGroup(types string) []ManagementFormUserGroup {
	var result []ManagementFormUserGroup
	if strings.EqualFold(types, "internal") {
		result = append(result,
			ManagementFormUserGroup{
				Name:             "Internal I",
				Desc:             "EAB",
				User_Permissions: addUserPermissions(types),
			},
			ManagementFormUserGroup{
				Name:             "Internal I",
				Desc:             "PLAB",
				User_Permissions: addUserPermissions(types)},
			ManagementFormUserGroup{
				Name:             "Internal III",
				Desc:             "Partisipant",
				User_Permissions: addUserPermissions(types),
			})
		return result
	}
	if strings.EqualFold(types, "external") {
		result = append(result, ManagementFormUserGroup{
			Name:             "External I",
			Desc:             "AB",
			User_Permissions: addUserPermissions(types),
		}, ManagementFormUserGroup{
			Name:             "External II",
			Desc:             "Partisipant",
			User_Permissions: addUserPermissions(types),
		}, ManagementFormUserGroup{
			Name:             "External III",
			Desc:             "PJSPPA",
			User_Permissions: addUserPermissions(types),
		}, ManagementFormUserGroup{
			Name:             "External IV",
			Desc:             "DU",
			User_Permissions: addUserPermissions(types),
		})
		return result
	}

	return nil

}

func CreateManagementFormTableHeader(userType []string) []ManagementFormTableHeader {
	var result []ManagementFormTableHeader

	for _, usrType := range userType {
		item := ManagementFormTableHeader{
			User_Type:  usrType,
			User_Group: addUserGroup(usrType),
		}

		result = append(result, item)
	}

	return result
}
