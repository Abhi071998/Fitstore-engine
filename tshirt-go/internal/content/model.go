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

// HeroContent stores the admin-editable homepage hero banner (tag line,
// two-tone heading, description, image, and the two call-to-action links).
// Only one row is expected to ever exist.
type HeroContent struct {
	ID                      uint      `gorm:"primaryKey" json:"id"`
	HeroTag                 string    `gorm:"type:varchar(255)" json:"hero_tag"`
	HeroHeadingLine1        string    `gorm:"type:varchar(255)" json:"hero_heading_line1"`
	HeroHeadingHighlight    string    `gorm:"type:varchar(255)" json:"hero_heading_highlight"`
	HeroHeadingLine2        string    `gorm:"type:varchar(255)" json:"hero_heading_line2"`
	HeroDescription         string    `gorm:"type:text" json:"hero_description"`
	HeroImage               string    `gorm:"type:text" json:"hero_image"`
	HeroPrimaryButtonText   string    `gorm:"type:varchar(255)" json:"hero_primary_button_text"`
	HeroPrimaryButtonLink   string    `gorm:"type:text" json:"hero_primary_button_link"`
	HeroSecondaryButtonText string    `gorm:"type:varchar(255)" json:"hero_secondary_button_text"`
	HeroSecondaryButtonLink string    `gorm:"type:text" json:"hero_secondary_button_link"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

func (HeroContent) TableName() string { return "hero_content" }

// HeroDTO carries the editable fields for create/update requests.
type HeroDTO struct {
	HeroTag                 string `json:"hero_tag"`
	HeroHeadingLine1        string `json:"hero_heading_line1"`
	HeroHeadingHighlight    string `json:"hero_heading_highlight"`
	HeroHeadingLine2        string `json:"hero_heading_line2"`
	HeroDescription         string `json:"hero_description"`
	HeroImage               string `json:"hero_image"`
	HeroPrimaryButtonText   string `json:"hero_primary_button_text"`
	HeroPrimaryButtonLink   string `json:"hero_primary_button_link"`
	HeroSecondaryButtonText string `json:"hero_secondary_button_text"`
	HeroSecondaryButtonLink string `json:"hero_secondary_button_link"`
}
