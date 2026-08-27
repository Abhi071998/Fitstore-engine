// Package content lets admins manage marketing copy/images for static
// storefront pages (About Us first) without redeploying the frontend.
package content

import "time"

// AdminContent stores admin-editable text/images shown on public pages.
// Only one row is expected to ever exist for the About Us section.
type AdminContent struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	AboutUsImg         string    `gorm:"type:text" json:"about_us_img"`
	AboutUsTitle       string    `gorm:"type:varchar(255)" json:"about_us_title"`
	AboutUsDescription string    `gorm:"type:text" json:"about_us_description"`
	AboutUsTagline1    string    `gorm:"type:varchar(255)" json:"about_us_tagline1"`
	AboutUsTagline2    string    `gorm:"type:varchar(255)" json:"about_us_tagline2"`
	AboutUsTagline3    string    `gorm:"type:varchar(255)" json:"about_us_tagline3"`
	AboutUsTagline4    string    `gorm:"type:varchar(255)" json:"about_us_tagline4"`
	AboutUsEstbYear    string    `gorm:"type:varchar(10)" json:"about_us_estb_year"`
	AboutUsVisitUs     string    `gorm:"type:text" json:"about_us_visit_us"`
	AboutUsEmail       string    `gorm:"type:varchar(255)" json:"about_us_email"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (AdminContent) TableName() string { return "admin_content" }

// AboutUsDTO carries the editable fields for create/update requests.
type AboutUsDTO struct {
	AboutUsImg         string `json:"about_us_img"`
	AboutUsTitle       string `json:"about_us_title"`
	AboutUsDescription string `json:"about_us_description"`
	AboutUsTagline1    string `json:"about_us_tagline1"`
	AboutUsTagline2    string `json:"about_us_tagline2"`
	AboutUsTagline3    string `json:"about_us_tagline3"`
	AboutUsTagline4    string `json:"about_us_tagline4"`
	AboutUsEstbYear    string `json:"about_us_estb_year"`
	AboutUsVisitUs     string `json:"about_us_visit_us"`
	AboutUsEmail       string `json:"about_us_email"`
}
