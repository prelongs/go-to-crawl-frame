package browserutil

import (
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/tebeka/selenium"
)

func GetXpathCondition(xpathPattern string) selenium.Condition {
	return func(wd selenium.WebDriver) (bool, error) {
		ele := GetByXpath(wd, xpathPattern)
		if ele != nil {
			return true, nil
		}
		return false, nil
	}
}

func GetTextByXpath(wd selenium.WebDriver, xpathPattern string) string {
	element := GetByXpath(wd, xpathPattern)
	if element == nil {
		return ""
	}
	text, _ := element.Text()
	return text
}

func GetTextArrayByXpath(wd selenium.WebDriver, xpathPattern string) []string {
	elements, _ := wd.FindElements(selenium.ByXPATH, xpathPattern)
	var arr []string
	for _, ele := range elements {
		text, _ := ele.Text()
		arr = append(arr, text)
	}
	return arr
}

func GetTextByRegex(wd selenium.WebDriver, regexPattern string) []string {
	html, err := wd.PageSource()
	if err != nil {
		return nil
	}

	rets, err := gregex.MatchString(regexPattern, html)
	return rets
}

func GetByXpath(wd selenium.WebDriver, xpathPattern string) selenium.WebElement {
	element, _ := wd.FindElement(selenium.ByXPATH, xpathPattern)
	return element
}
