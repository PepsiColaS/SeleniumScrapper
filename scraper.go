package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

type Product struct {
	name, price, link string
}

var products []Product

// Урлы для проверки
// https://www.okeydostavka.ru/msk/miaso-ptitsa-kolbasy/ptitsa-20
// https://www.okeydostavka.ru/msk/ovoshchi-i-frukty/ovoshchi
// https://www.okeydostavka.ru/msk/miaso-ptitsa-kolbasy/miaso-20

const siteURL string = "https://www.okeydostavka.ru/msk/ovoshchi-i-frukty/ovoshchi"
const addres string = "Москва, Олимпийский проспект, 18/1"

func main() {
	startTime := time.Now()

	service, err := selenium.NewChromeDriverService("./chromedriver", 4444)
	if err != nil {
		log.Fatal("Error:", err)
	}
	defer service.Stop()

	caps := selenium.Capabilities{}
	chromeCaps := chrome.Capabilities{
		Args: []string{
			"--disable-gpu",
			"--no-sandbox",
			"--disable-dev-shm-usage",
			"--window-size=1920,1080",
			"user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		},
	}
	caps.AddChrome(chromeCaps)

	driver, err := selenium.NewRemote(caps, "")
	if err != nil {
		log.Fatal("Error:", err)
	}
	defer driver.Quit()

	err = driver.Get(siteURL)
	if err != nil {
		log.Fatal("Error:", err)
	}
	waitForElement(driver, "availableReceiptTimeslot")
	selectPlace, _ := driver.FindElement(selenium.ByID, "availableReceiptTimeslot")
	selectPlace.Click()

	waitForElement(driver, ".dijitReset .dijitInputInner")
	inputPlace, _ := driver.FindElement(selenium.ByCSSSelector, ".dijitReset .dijitInputInner")
	err = inputPlace.SendKeys(addres)
	if err != nil {
		log.Fatal("Error:", err)

	}

	waitForElement(driver, "addressSelectionButton")
	saveButton, _ := driver.FindElement(selenium.ByID, "addressSelectionButton")
	saveButton.Click()

	for {
		waitForElement(driver, ".product")
		productElements, _ := driver.FindElements(selenium.ByCSSSelector, ".product")

		for _, productElement := range productElements {
			nameElement, err := productElement.FindElement(selenium.ByCSSSelector, ".product-info .product-name a")
			if err != nil {
				log.Fatal("Error:", err)
			}

			priceElement, err := productElement.FindElement(selenium.ByCSSSelector, ".shopper-actions .price_and_cart .product-price__container .product-price .price")
			if err != nil {
				log.Fatal("Error:", err)
			}

			name, err := nameElement.Text()
			if err != nil {
				log.Fatal("Error:", err)
			}
			price, err := priceElement.Text()
			if err != nil {
				log.Fatal("Error:", err)
			}
			link, err := nameElement.GetAttribute("href")
			if err != nil {
				log.Fatal("Error:", err)
			}

			product := Product{name: name, price: price, link: link}
			products = append(products, product)

		}

		element, err := driver.FindElement(selenium.ByCSSSelector, ".paging_controls .right_arrow ")
		if err != nil {
			break
		}
		IsDisplayed, _ := element.IsDisplayed()
		if IsDisplayed {
			element.Click()
		} else {
			break
		}

	}
	createFile()
	fmt.Println(time.Since(startTime), " - work time")
}

func waitForElement(driver selenium.WebDriver, selector string) {
	timeout := time.After(5 * time.Second)
	tick := time.Tick(500 * time.Millisecond)

	var element selenium.WebElement
	var err error
	for {
		select {
		case <-timeout:
			log.Fatal("Timeout waiting for element to be visible")
			return
		case <-tick:
			if rune(selector[0]) == '.' {
				element, err = driver.FindElement(selenium.ByCSSSelector, selector)
			} else {
				element, err = driver.FindElement(selenium.ByID, selector)
			}
			if err == nil {
				isDisplayed, err := element.IsDisplayed()
				if err == nil && isDisplayed {
					time.Sleep(1 * time.Second)
					return
				}
			}
		}
	}
}

func createFile() {
	file, err := os.OpenFile("products.csv", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal("Error:", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"name", "price", "link"}
	if err := writer.Write(headers); err != nil {
		log.Fatal("Error writing headers:", err)
	}

	var records [][]string
	for _, product := range products {
		record := []string{product.name, product.price, product.link}
		records = append(records, record)
	}

	if err := writer.WriteAll(records); err != nil {
		log.Fatal("Error writing records:", err)
	}
}
