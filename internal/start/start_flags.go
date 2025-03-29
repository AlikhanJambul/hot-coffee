package start

import (
	"encoding/json"
	"flag"
	"fmt"
	"hot-coffee/models"
	"io/ioutil"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
)

func Help() {
	fmt.Print(`Simple Storage Service.

**Usage:**
    triple-s [-port <N>] [-dir <S>]  
    triple-s --help

**Options:**
- --help     Show this screen.
- --port N   Port number
- --dir S    Path to the directory`, "\n")
}

func AllFlags() (int, string) {
	helpFlag := flag.Bool("help", false, "help")
	portFlag := flag.Int("port", 8080, "port")
	dirFlag := flag.String("dir", "data", "dir")
	flag.Usage = Help
	flag.Parse()

	if *helpFlag == true {
		Help()
		os.Exit(0)
	}

	if validPort := isValidPort(*portFlag); !validPort {
		fmt.Fprintf(os.Stderr, "This port isn't a valid\n")
		os.Exit(1)
	}
	if validName := IsValidName(*dirFlag); !validName {
		fmt.Fprintf(os.Stderr, "This name isn't a valid\n")
		os.Exit(1)
	}

	return *portFlag, *dirFlag
}

func isValidPort(portNum int) bool {
	if portNum < 1 || portNum > 65535 {
		return false
	}

	return true
}

func IsValidName(name string) bool {
	re := regexp.MustCompile("^[a-z0-9-\\.]+$")

	if strings.Contains(name, "internal") || strings.Contains(name, "cmd") || strings.Contains(name, "models") || strings.Contains(name, "handler") {
		return false
	}

	if strings.Contains(name, "..") || strings.Contains(name, "--") || strings.Contains(name, "-.") || strings.Contains(name, ".-") {
		return false
	}

	if strings.Contains(name, "service") || strings.Contains(name, "dal") || strings.Contains(name, "start") {
		return false
	}

	if name == "." {
		return false
	}

	if net.ParseIP(name) != nil {
		fmt.Fprintf(os.Stderr, "It's ip address: %s\n", name)
		os.Exit(1)
	}

	return re.MatchString(name)
}

func CreateDir(dirFlag string) {
	_, err := os.Stat(dirFlag)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Папка не существует. Создаем!")
		}
	} else {
		_ = os.RemoveAll(dirFlag)
	}

	_ = os.Mkdir(dirFlag, 0o766)

	ingredients := []map[string]interface{}{
		{
			"ingredient_id": "espresso_shot",
			"name":          "Espresso Shot",
			"quantity":      500,
			"unit":          "shots",
		},
		{
			"ingredient_id": "milk",
			"name":          "Milk",
			"quantity":      5000,
			"unit":          "ml",
		},
		{
			"ingredient_id": "flour",
			"name":          "Flour",
			"quantity":      10000,
			"unit":          "g",
		},
		{
			"ingredient_id": "blueberries",
			"name":          "Blueberries",
			"quantity":      2000,
			"unit":          "g",
		},
		{
			"ingredient_id": "sugar",
			"name":          "Sugar",
			"quantity":      5000,
			"unit":          "g",
		},
	}

	file, err := os.Create(dirFlag + "/inventory.json")
	if err != nil {
		return
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(ingredients, "", "  ")
	if err != nil {
		return
	}

	_, err = file.Write(jsonData)
	if err != nil {
		return
	}

	products := []map[string]interface{}{
		{
			"product_id":  "latte",
			"name":        "Caffe Latte",
			"description": "Espresso with steamed milk",
			"price":       3.50,
			"ingredients": []map[string]interface{}{
				{
					"ingredient_id": "espresso_shot",
					"quantity":      1,
				},
				{
					"ingredient_id": "milk",
					"quantity":      200,
				},
			},
		},
		{
			"product_id":  "muffin",
			"name":        "Blueberry Muffin",
			"description": "Freshly baked muffin with blueberries",
			"price":       2.00,
			"ingredients": []map[string]interface{}{
				{
					"ingredient_id": "flour",
					"quantity":      100,
				},
				{
					"ingredient_id": "blueberries",
					"quantity":      20,
				},
				{
					"ingredient_id": "sugar",
					"quantity":      30,
				},
			},
		},
		{
			"product_id":  "espresso",
			"name":        "Espresso",
			"description": "Strong and bold coffee",
			"price":       2.50,
			"ingredients": []map[string]interface{}{
				{
					"ingredient_id": "espresso_shot",
					"quantity":      1,
				},
			},
		},
	}

	file1, err1 := os.Create(dirFlag + "/menu_items.json")
	if err1 != nil {
		fmt.Println("Error creating file:", err1)
		return
	}
	defer file1.Close()

	jsonData1, err1 := json.MarshalIndent(products, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling data:", err)
		return
	}

	_, err = file1.Write(jsonData1)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	file2, err2 := os.Create(dirFlag + "/orders.json")
	if err2 != nil {
		return
	}
	file2.Close()
}

func ChangeJsonFile() {
	filePath := "data/menu_items.json"
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Ошибка при чтении файла: %v", err)
	}

	var menuItems []models.MenuItem
	err = json.Unmarshal(file, &menuItems)
	if err != nil {
		log.Fatalf("Ошибка при декодировании JSON: %v", err)
	}

	for i, item := range menuItems {
		if item.ID == "latte" {
			menuItems[i].Ingredients[1].Quantity = 100000
			fmt.Println("Обновленная цена Caffe Latte:", menuItems[i].Price)
		}
	}

	updatedData, err := json.MarshalIndent(menuItems, "", "  ")
	if err != nil {
		log.Fatalf("Ошибка при кодировании JSON: %v", err)
	}

	err = ioutil.WriteFile(filePath, updatedData, 0o644)
	if err != nil {
		log.Fatalf("Ошибка при записи в файл: %v", err)
	}

	fmt.Println("Файл успешно обновлен.")
}
