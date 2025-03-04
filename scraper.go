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

var siteURL string = "https://www.okeydostavka.ru/msk/ovoshchi-i-frukty/ovoshchi"

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

	for {
		time.Sleep(1 * time.Second)

		productElements, err := driver.FindElements(selenium.ByCSSSelector, ".product")
		if err != nil {
			log.Fatal("Error:", err)
		}

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
			log.Fatal("Error:", err)
		}

		isVisible, _ := element.IsDisplayed()
		if !isVisible {
			break
		}
		element.Click()
	}

	createFile()
	fmt.Println(time.Since(startTime), " - work time")
}

func createFile() {
	file, err := os.Create("products.csv")
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
