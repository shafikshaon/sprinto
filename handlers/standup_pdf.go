package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/go-pdf/fpdf"

	"sprinto/models"
	"sprinto/repository"
)

func (h *StandupHandler) PDF(c *gin.Context) {
	projectID := activeProjectIDFromCtx(c)
	f := repository.StandupFilter{
		DateFrom: c.Query("date_from"),
		DateTo:   c.Query("date_to"),
		Search:   c.Query("search"),
	}

	entries, _, _ := h.svc.All(projectID, f, 1, 10000)

	pdf := buildStandupPDF(entries, f)

	filename := "standups"
	if f.DateFrom != "" {
		filename += "_from_" + f.DateFrom
	}
	if f.DateTo != "" {
		filename += "_to_" + f.DateTo
	}
	filename += ".pdf"

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	if err := pdf.Output(c.Writer); err != nil {
		c.String(http.StatusInternalServerError, "PDF error: %v", err)
	}
}

// safe converts a string to latin-1 safe for fpdf (replaces non-latin chars).
func safe(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r <= 0xFF {
			b.WriteRune(r)
		} else {
			b.WriteRune('?')
		}
	}
	return b.String()
}

// truncate shortens a string to maxRunes runes.
func truncate(s string, maxRunes int) string {
	if utf8.RuneCountInString(s) <= maxRunes {
		return s
	}
	runes := []rune(s)
	return string(runes[:maxRunes]) + "..."
}

func buildStandupPDF(entries []models.StandupEntry, f repository.StandupFilter) *fpdf.Fpdf {
	const (
		marginL    = 15.0
		marginR    = 15.0
		marginT    = 15.0
		marginB    = 18.0
		pageW      = 210.0
		pageH      = 297.0
		contentW   = pageW - marginL - marginR
		headerH    = 28.0
		footerH    = 10.0
	)

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(marginL, marginT, marginR)
	pdf.SetAutoPageBreak(true, marginB+footerH)

	// ── Footer (page number) ────────────────────────────────────────────────
	pdf.SetFooterFunc(func() {
		pdf.SetY(-(marginB + footerH - 4))
		pdf.SetDrawColor(220, 220, 230)
		pdf.Line(marginL, pageH-marginB-footerH+2, pageW-marginR, pageH-marginB-footerH+2)
		pdf.SetFont("Helvetica", "", 8)
		pdf.SetTextColor(160, 160, 175)
		pdf.CellFormat(contentW/2, 5, safe("Sprinto · Daily Standups"), "", 0, "L", false, 0, "")
		pdf.CellFormat(contentW/2, 5, fmt.Sprintf("Page %d", pdf.PageNo()), "", 0, "R", false, 0, "")
	})

	pdf.AddPage()

	// ── Document header ─────────────────────────────────────────────────────
	// Accent bar
	pdf.SetFillColor(79, 70, 229) // indigo-600
	pdf.Rect(marginL, marginT, contentW, 1.2, "F")
	pdf.Ln(3)

	// Title
	pdf.SetFont("Helvetica", "B", 20)
	pdf.SetTextColor(30, 27, 75)
	pdf.CellFormat(contentW, 9, "Daily Standups", "", 1, "L", false, 0, "")

	// Subtitle / filter info
	subtitle := "Exported on " + time.Now().Format("2 Jan 2006")
	if f.DateFrom != "" || f.DateTo != "" {
		parts := []string{}
		if f.DateFrom != "" {
			parts = append(parts, "from "+f.DateFrom)
		}
		if f.DateTo != "" {
			parts = append(parts, "to "+f.DateTo)
		}
		subtitle += "  ·  " + strings.Join(parts, " ")
	}
	if f.Search != "" {
		subtitle += `  ·  search: "` + truncate(f.Search, 40) + `"`
	}
	subtitle += fmt.Sprintf("  ·  %d record(s)", len(entries))

	pdf.SetFont("Helvetica", "", 9)
	pdf.SetTextColor(120, 115, 160)
	pdf.CellFormat(contentW, 5, safe(subtitle), "", 1, "L", false, 0, "")
	pdf.Ln(5)

	if len(entries) == 0 {
		pdf.SetFont("Helvetica", "I", 11)
		pdf.SetTextColor(160, 160, 175)
		pdf.CellFormat(contentW, 10, "No standup entries found for the selected filters.", "", 1, "C", false, 0, "")
		return pdf
	}

	// ── Entries ─────────────────────────────────────────────────────────────
	for i, e := range entries {
		// Estimate height needed; add a page break if < 40mm left.
		if i > 0 {
			remaining := pageH - marginB - footerH - pdf.GetY()
			if remaining < 38 {
				pdf.AddPage()
				pdf.Ln(2)
			} else {
				pdf.Ln(4)
			}
		}

		drawStandupCard(pdf, e, contentW, marginL)
	}

	return pdf
}

func drawStandupCard(pdf *fpdf.Fpdf, e models.StandupEntry, contentW, marginL float64) {
	startY := pdf.GetY()

	// ── Card header background ───────────────────────────────────────────────
	pdf.SetFillColor(245, 245, 255) // near-white indigo tint
	pdf.SetDrawColor(220, 218, 240)

	// We'll draw the header rect after computing its height
	headerY := startY

	// Date + project chips
	pdf.SetXY(marginL+3, headerY+2.5)
	pdf.SetFont("Helvetica", "B", 10)
	pdf.SetTextColor(30, 27, 75)
	dateStr := safe(e.Date)
	pdf.CellFormat(0, 5, dateStr, "", 0, "L", false, 0, "")

	// Project pill
	if e.Project.Name != "" {
		dateW := pdf.GetStringWidth(dateStr)
		pdf.SetXY(marginL+3+dateW+4, headerY+3)
		pdf.SetFillColor(237, 233, 254) // violet-100
		pdf.SetDrawColor(221, 214, 254)
		pdf.SetTextColor(109, 40, 217) // violet-700
		pdf.SetFont("Helvetica", "", 7.5)
		pName := safe(truncate(e.Project.Name, 30))
		pW := pdf.GetStringWidth(pName) + 5
		pdf.RoundedRect(pdf.GetX(), headerY+2, pW, 5, 1.2, "1234", "FD")
		pdf.CellFormat(pW, 5, pName, "", 0, "C", false, 0, "")
	}

	// Timestamps on the right
	pdf.SetFont("Helvetica", "", 7.5)
	pdf.SetTextColor(150, 145, 170)
	created := "Created " + timeAgoStr(e.CreatedAt)
	edited := ""
	if e.UpdatedAt.Unix()-e.CreatedAt.Unix() > 1 {
		edited = "  ·  Updated " + timeAgoStr(e.UpdatedAt)
	}
	tsStr := safe(created + edited)
	pdf.SetXY(marginL, headerY+2.5)
	pdf.CellFormat(contentW-3, 5, tsStr, "", 0, "R", false, 0, "")

	headerBoxH := 9.5
	// Draw header bg (behind text) — done here because fpdf draws in order
	// so we re-draw the background using a clipping trick: just draw rect first
	// Actually fpdf doesn't have layers, so we draw bg then re-draw text above.
	// Instead we'll draw the full rounded rect as background before text.
	// Restructure: draw bg rect, then text cells.

	// Draw card header bg rect
	pdf.SetFillColor(245, 245, 255)
	pdf.SetDrawColor(220, 218, 240)
	pdf.Rect(marginL, headerY, contentW, headerBoxH, "FD")

	// Re-draw text over bg
	pdf.SetFont("Helvetica", "B", 10)
	pdf.SetTextColor(30, 27, 75)
	pdf.SetXY(marginL+3, headerY+2.5)
	pdf.CellFormat(0, 5, dateStr, "", 0, "L", false, 0, "")

	if e.Project.Name != "" {
		dateW := pdf.GetStringWidth(dateStr)
		pdf.SetXY(marginL+3+dateW+4, headerY+2.5)
		pdf.SetFillColor(237, 233, 254)
		pdf.SetDrawColor(221, 214, 254)
		pdf.SetTextColor(109, 40, 217)
		pdf.SetFont("Helvetica", "", 7.5)
		pName := safe(truncate(e.Project.Name, 30))
		pW := pdf.GetStringWidth(pName) + 5
		pdf.RoundedRect(pdf.GetX(), headerY+2, pW, 5.5, 1.2, "1234", "FD")
		pdf.CellFormat(pW, 5.5, pName, "", 0, "C", false, 0, "")
	}

	pdf.SetFont("Helvetica", "", 7.5)
	pdf.SetTextColor(150, 145, 170)
	pdf.SetXY(marginL, headerY+2.5)
	pdf.CellFormat(contentW-3, 5, tsStr, "", 0, "R", false, 0, "")

	pdf.SetY(headerY + headerBoxH)

	// ── Card body ────────────────────────────────────────────────────────────
	bodyStartY := pdf.GetY()
	pdf.SetDrawColor(220, 218, 240)
	pdf.SetFillColor(255, 255, 255)
	// We'll draw the border after computing body height

	bodyX := marginL + 3
	bodyW := contentW - 6

	currentY := bodyStartY + 3

	if e.Summary != "" {
		currentY = drawSection(pdf, bodyX, currentY, bodyW,
			"DISCUSSION SUMMARY", e.Summary,
			100, 95, 120, // label color: slate
			40, 40, 60,   // text color: near-black
			248, 248, 252, // bg: very light
			220, 218, 240) // border
	}

	if e.Dependencies != "" {
		if currentY > bodyStartY+3 {
			currentY += 2
		}
		currentY = drawSection(pdf, bodyX, currentY, bodyW,
			"DEPENDENCIES", e.Dependencies,
			146, 110, 0,   // amber label
			120, 80, 0,    // amber-dark text
			255, 251, 235, // amber-50 bg
			253, 230, 138) // amber-200 border
	}

	if e.Issues != "" {
		if currentY > bodyStartY+3 {
			currentY += 2
		}
		currentY = drawSection(pdf, bodyX, currentY, bodyW,
			"ISSUES / BLOCKERS", e.Issues,
			185, 28, 28,   // red label
			153, 20, 20,   // red text
			254, 242, 242, // red-50 bg
			254, 202, 202) // red-200 border
	}

	if e.ActionItems != "" {
		if currentY > bodyStartY+3 {
			currentY += 2
		}
		currentY = drawSection(pdf, bodyX, currentY, bodyW,
			"ACTION ITEMS", e.ActionItems,
			21, 128, 61,   // green label
			20, 90, 50,    // green text
			240, 253, 244, // green-50 bg
			187, 247, 208) // green-200 border
	}

	if e.Summary == "" && e.Dependencies == "" && e.Issues == "" && e.ActionItems == "" {
		pdf.SetFont("Helvetica", "I", 9)
		pdf.SetTextColor(180, 175, 195)
		pdf.SetXY(bodyX, currentY)
		pdf.CellFormat(bodyW, 6, "No details recorded.", "", 1, "L", false, 0, "")
		currentY += 6
	}

	bodyEndY := currentY + 3

	// Draw card outline
	pdf.SetDrawColor(220, 218, 240)
	pdf.SetFillColor(255, 255, 255)
	pdf.Rect(marginL, bodyStartY, contentW, bodyEndY-bodyStartY, "D")

	pdf.SetY(bodyEndY)
}

// drawSection renders a labelled section block and returns the new Y position.
func drawSection(pdf *fpdf.Fpdf, x, y, w float64,
	label, text string,
	labelR, labelG, labelB int,
	textR, textG, textB int,
	bgR, bgG, bgB int,
	borderR, borderG, borderB int,
) float64 {
	const labelH = 5.0
	const pad = 3.0

	// Measure text height
	pdf.SetFont("Helvetica", "", 9)
	lines := pdf.SplitLines([]byte(safe(text)), w-pad*2)
	textH := float64(len(lines)) * 4.5

	totalH := pad + labelH + 1.5 + textH + pad

	// Background + border
	pdf.SetFillColor(bgR, bgG, bgB)
	pdf.SetDrawColor(borderR, borderG, borderB)
	pdf.Rect(x, y, w, totalH, "FD")

	// Label
	pdf.SetFont("Helvetica", "B", 7)
	pdf.SetTextColor(labelR, labelG, labelB)
	pdf.SetXY(x+pad, y+pad)
	pdf.CellFormat(w-pad*2, labelH, safe(label), "", 1, "L", false, 0, "")

	// Body text
	pdf.SetFont("Helvetica", "", 9)
	pdf.SetTextColor(textR, textG, textB)
	pdf.SetXY(x+pad, y+pad+labelH+1.5)
	pdf.MultiCell(w-pad*2, 4.5, safe(text), "", "L", false)

	return y + totalH
}

// timeAgoStr returns a human-readable duration without using the template funcMap.
func timeAgoStr(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		m := int(d.Minutes())
		if m == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", m)
	case d < 24*time.Hour:
		h := int(d.Hours())
		if h == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", h)
	case d < 7*24*time.Hour:
		days := int(d.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	default:
		return t.Format("2 Jan 2006")
	}
}
