package pdf

import (
	"strings"
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
	FooterLogo  string
	PapperWidth float64
	// setting papper height
	Papperheight float64
	// setting page orientation by default "p" / "l"
	PageOrientation string
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

func (opt *PdfTableOptions) getHeaderTitle() string {
	if opt.HeaderTitle == "" {
		return "Bursa Efek Indonesia"
	}
	return opt.HeaderTitle
}
