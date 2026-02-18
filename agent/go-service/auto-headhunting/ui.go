package autoheadhunting

import (
	"fmt"
	"html"
	"strings"

	"github.com/MaaXYZ/MaaEnd/agent/go-service/pkg/maafocus"
	maa "github.com/MaaXYZ/maa-framework-go/v4"
)

var starColors = map[string]string{
	"4": "#CF1DCC",
	"5": "#E0DD19",
	"6": "#F54927",
}

func LogMXUHTML(ctx *maa.Context, htmlText string) {
	htmlText = strings.TrimLeft(htmlText, " \t\r\n")
	maafocus.NodeActionStarting(ctx, htmlText)
}

// LogMXUSimpleHTMLWithColor logs a simple styled span, allowing a custom color.
func LogMXUSimpleHTMLWithColor(ctx *maa.Context, text string, color string) {
	HTMLTemplate := fmt.Sprintf(`<span style="color: %s; font-weight: 500;">%%s</span>`, color)
	LogMXUHTML(ctx, fmt.Sprintf(HTMLTemplate, text))
}

// LogMXUSimpleHTML logs a simple styled span with a default color.
func LogMXUSimpleHTML(ctx *maa.Context, text string) {
	// Call the more specific function with the default color "#00bfff".
	LogMXUSimpleHTMLWithColor(ctx, text, "#00bfff")
}

// getColorForStars 根据星级返回对应的颜色
func getColorForStars(stars string) string {
	if color, exists := starColors[stars]; exists {
		return color
	}
	return "#00bfff" // 默认颜色
}

// escapeHTML 简单封装 html.EscapeString，便于后续统一替换/扩展
func escapeHTML(s string) string {
	return html.EscapeString(s)
}

// formatOperatorNameColoredHTML 根据干员星级为干员名着色并返回 HTML 片段
func formatOperatorNameColoredHTML(name string, stars string) string {
	color := getColorForStars(stars)
	return fmt.Sprintf(`<span style="color: %s; font-weight: 600;">%s</span>`, color, escapeHTML(name))
}

// logTaskParamsHTML 输出任务参数配置摘要的 HTML
func logTaskParamsHTML(ctx *maa.Context, targetPulls int, targetLabel string, targetOperatorNum int, preferMode int) {
	var b strings.Builder
	b.WriteString(`<div style="font-weight: 900; color: #00bfff; margin-bottom: 2px;">🎰 Auto Headhunting</div>`)
	b.WriteString(`<table style="border-collapse: collapse; font-size: 12px;">`)

	rows := []struct {
		label string
		value string
	}{
		{t("target_pulls"), fmt.Sprintf("%d", targetPulls)},
		{t("target_operator"), targetLabel},
		{t("target_num"), fmt.Sprintf("%d", targetOperatorNum)},
		{t("prefer_mode"), fmt.Sprintf("%d", preferMode)},
	}

	for _, row := range rows {
		b.WriteString(fmt.Sprintf(
			`<tr><td style="padding: 1px 6px 1px 0; color: #888;">%s</td><td style="padding: 1px 0; color: #e0e0e0; font-weight: 500;">%s</td></tr>`,
			escapeHTML(row.label), escapeHTML(row.value),
		))
	}

	b.WriteString(`</table>`)
	LogMXUHTML(ctx, b.String())
}

// logPullResultsHTML 输出单轮抽卡结果的 HTML
func logPullResultsHTML(ctx *maa.Context, usedPulls int, targetPulls int, ansMp map[string]int) {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(
		`<div style="color: #00bfff; font-weight: 500;">%s %d/%d</div>`,
		escapeHTML(t("used_pulls")), usedPulls, targetPulls,
	))
	for name, count := range ansMp {
		_, stars := o(t(name))
		coloredName := formatOperatorNameColoredHTML(name, stars)
		b.WriteString(fmt.Sprintf(
			`<div>%s: %d</div>`, coloredName, count,
		))
	}
	LogMXUHTML(ctx, b.String())
}

// logFinalSummaryHTML 输出最终抽卡结果摘要的 HTML
func logFinalSummaryHTML(ctx *maa.Context, usedPulls int, targetCount int, targetLabel string, mp map[string]int) {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(
		`<div style="color: #00bfff; font-weight: 900; margin-top: 4px;">%s</div>`,
		escapeHTML(fmt.Sprintf(t("done"), usedPulls, targetCount, targetLabel)),
	))
	b.WriteString(`<table style="width: 100%; border-collapse: collapse; font-size: 12px;">`)
	b.WriteString(fmt.Sprintf(
		`<tr><th style="text-align:left; padding: 2px 4px;">%s</th><th style="text-align:right; padding: 2px 4px;">%s</th></tr>`,
		escapeHTML(t("target_operator")), escapeHTML(t("target_num")),
	))
	for name, count := range mp {
		// 跳过星级统计条目（key 为纯数字星级如 "4", "5", "6"）
		if _, exists := starColors[name]; exists {
			continue
		}
		_, stars := o(t(name))
		coloredName := formatOperatorNameColoredHTML(name, stars)
		b.WriteString("<tr>")
		b.WriteString(fmt.Sprintf(`<td style="padding: 2px 4px;">%s</td>`, coloredName))
		b.WriteString(fmt.Sprintf(`<td style="padding: 2px 4px; text-align: right;">%d</td>`, count))
		b.WriteString("</tr>")
	}
	b.WriteString(`</table>`)
	LogMXUHTML(ctx, b.String())
}
