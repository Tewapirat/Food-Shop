package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/chzyer/readline"

	"github.com/TewApirat/food-shop/pkg/foodShop/model"
	_foodShopService "github.com/TewApirat/food-shop/pkg/foodShop/service"
)

type FoodShopControllerImpl struct {
	in              io.ReadCloser
	out             io.Writer
	foodShopService _foodShopService.FoodShopService
}

func NewFoodShopControllerImpl(in io.ReadCloser, out io.Writer, foodShopService _foodShopService.FoodShopService) FoodShopController {
	return &FoodShopControllerImpl{
		in:              in,
		out:             out,
		foodShopService: foodShopService,
	}
}

func (c *FoodShopControllerImpl) ServeCLI() {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "Select: ",
		Stdin:           c.in,
		Stdout:          c.out,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		fmt.Fprintln(c.out, "Readline init error:", err)
		return
	}
	defer rl.Close()

	for {
		fmt.Fprintln(c.out, "\n==== Food Shop CLI ====")
		fmt.Fprintln(c.out, "1) View all menu items")
		fmt.Fprintln(c.out, "2) View all promotions")
		fmt.Fprintln(c.out, "3) Quote order (JSON input)")
		fmt.Fprintln(c.out, "0) Exit")

		rl.SetPrompt("Select: ")
		choice, err := readLine(rl)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Fprintln(c.out, "\nEOF received. Bye.")
				return
			}
			fmt.Fprintln(c.out, "\nRead error:", err)
			return
		}
		if choice == "" {
			// กรณี Ctrl+C แล้วเราคืน "" -> ให้ loop ต่อ
			continue
		}

		switch choice {
		case "1":
			c.handleViewMenu()
		case "2":
			c.handleViewPromotions()
		case "3":
			if ok := c.handleQuoteOrderJSON(rl); !ok {
				return
			}
		case "0":
			fmt.Fprintln(c.out, "Bye.")
			return
		default:
			fmt.Fprintln(c.out, "Invalid choice. Please select 0-3.")
		}
	}
}

func (c *FoodShopControllerImpl) handleViewMenu() {
	items, err := c.foodShopService.GetMenuCatalog()
	if err != nil {
		fmt.Fprintln(c.out, err)
		return
	}

	fmt.Fprintln(c.out)
	fmt.Fprintln(c.out, "--- Menu Catalog ---")
	fmt.Fprintln(c.out)
	fmt.Fprintln(c.out, "--------+--------------+-----------")
	fmt.Fprintln(c.out, "CODE    | NAME         | PRICE		")
	fmt.Fprintln(c.out, "--------+--------------+-----------")

	for _, it := range items {
		fmt.Fprintf(c.out, "%-7s | %-12s | %s\n", it.Code, it.Name, it.Price.String())
	}
}

func (c *FoodShopControllerImpl) handleViewPromotions() {
	promos, err := c.foodShopService.GetPromotions()
	if err != nil {
		fmt.Fprintln(c.out, err)
		return
	}

	fmt.Fprintln(c.out)
	fmt.Fprintln(c.out, "--- Promotions ---")
	fmt.Fprintln(c.out)

	for _, p := range promos {
		fmt.Fprintf(c.out, "[%s] %s\n", p.Code, p.Title)
		fmt.Fprintf(c.out, " - %s\n", p.Description)
		fmt.Fprintln(c.out)
	}
}

func (c *FoodShopControllerImpl) handleQuoteOrderJSON(rl *readline.Instance) bool {
	fmt.Fprintln(c.out, "\nPaste order JSON in one line, then press Enter.")
	fmt.Fprintln(c.out, `Example: {"items":{"RED":1,"GREEN":2},"member":false}`)

	rl.SetPrompt("Order JSON: ")
	line, err := readLine(rl)
	if err != nil {
		if errors.Is(err, io.EOF) {
			fmt.Fprintln(c.out, "\nEOF received. Bye.")
			return false
		}
		fmt.Fprintln(c.out, "Read error:", err)
		return false
	}

	if strings.TrimSpace(line) == "" {
		fmt.Fprintln(c.out, "Error: empty input")
		return true
	}

	var req model.PurchasingRequest
	if err := json.Unmarshal([]byte(line), &req); err != nil {
		fmt.Fprintln(c.out, "Error: invalid JSON:", err)
		fmt.Fprintln(c.out, `Hint: {"items":{"RED":1,"GREEN":2},"member":false}`)
		return true
	}

	quote, err := c.foodShopService.QuoteOrder(req)
	if err != nil {
		fmt.Fprintln(c.out, err)
		return true
	}

	fmt.Fprintln(c.out, "\n--- Order Quote ---")
	fmt.Fprintln(c.out, "Subtotal      :", quote.Subtotal.String())
	fmt.Fprintln(c.out, "Pair Discount :", quote.PairDiscount.String())
	fmt.Fprintln(c.out, "Member Disc.  :", quote.MemberDiscount.String())
	fmt.Fprintln(c.out, "Total         :", quote.Total.String())

	return true
}

// readline-aware readLine
func readLine(rl *readline.Instance) (string, error) {
	s, err := rl.Readline()
	if err != nil {
		// Ctrl+C
		if errors.Is(err, readline.ErrInterrupt) {
			return "", nil
		}
		// Ctrl+D => io.EOF
		return "", err
	}
	return strings.TrimSpace(s), nil
}
