package spider

import (
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// extractHiddenInputs 从 HTML 中提取所有 <input type="hidden"> 的 name -> value。
// 用于 ASP.NET 页面中的 __VIEWSTATE、__EVENTVALIDATION 等。
func extractHiddenInputs(htmlStr string) map[string]string {
	result := map[string]string{}
	re := regexp.MustCompile(`(?is)<input[^>]*>`)
	for _, tag := range re.FindAllString(htmlStr, -1) {
		if !strings.Contains(strings.ToLower(tag), `type="hidden"`) &&
			!strings.Contains(strings.ToLower(tag), `type='hidden'`) {
			continue
		}
		name := attrValue(tag, "name")
		value := attrValue(tag, "value")
		if name != "" {
			result[name] = value
		}
	}
	return result
}

func attrValue(tag, attr string) string {
	// 同时兼容单双引号
	pattern := `(?i)` + attr + `\s*=\s*"([^"]*)"`
	if m := regexp.MustCompile(pattern).FindStringSubmatch(tag); len(m) == 2 {
		return m[1]
	}
	pattern2 := `(?i)` + attr + `\s*=\s*'([^']*)'`
	if m := regexp.MustCompile(pattern2).FindStringSubmatch(tag); len(m) == 2 {
		return m[1]
	}
	return ""
}

// ParseHTMLTableRows 从 HTML 中解析指定 matcher 匹配的 <table> 内所有数据行的单元格文本。
// tableFilter 可以是类名、id 等特征子串，空则匹配全部 <table>。
// 返回二维数组：每行一个 []string，代表一行的 td 文本。
func ParseHTMLTableRows(htmlStr, tableFilter string) [][]string {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return nil
	}
	var rows [][]string
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "table" && tableMatches(n, tableFilter) {
			collectRows(n, &rows)
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)
	return rows
}

func tableMatches(n *html.Node, filter string) bool {
	if filter == "" {
		return true
	}
	for _, a := range n.Attr {
		if strings.Contains(a.Val, filter) || strings.Contains(a.Key+"="+a.Val, filter) {
			return true
		}
	}
	return false
}

func collectRows(table *html.Node, out *[][]string) {
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "tr" {
			var row []string
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && (c.Data == "td" || c.Data == "th") {
					row = append(row, strings.TrimSpace(nodeText(c)))
				}
			}
			if len(row) > 0 {
				*out = append(*out, row)
			}
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(table)
}

func nodeText(n *html.Node) string {
	var sb strings.Builder
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.TextNode {
			sb.WriteString(n.Data)
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(n)
	return strings.Join(strings.Fields(sb.String()), " ")
}
