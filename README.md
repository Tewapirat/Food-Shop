# Food-Shop-App
**Published by** Apirat Khiewmeesuan.

## About Food-Shop-App
A simple CLI-based food shop application that calculates an order total with promotion rules (pair discount + member discount), with a clean architecture structure and unit tests.

## Features of Food-Shop-App
- Menu catalog (7 items)
- Order quotation (price calculation)
- Promotions
  - **Member discount:** 10% off the total (applied after pair discount)
  - **Pair discount:** 5% off per bundle of 2 items for eligible sets (Orange/Pink/Green)
- Order history (Optional)
- Unit tests

## Business Rules
Core pricing rules including pair-bundle discount eligibility and member discount stacking logic.
### Food Shop CLI
Command-line interface to browse menus, view promotions, and generate order quotes from a single JSON input.
``` text
==== Food Shop CLI ====
1) View all menu items
2) View all promotions
3) Quote order (JSON input)
4) View order history
0) Exit
Select:  
```
### All menu item
```text
| Code   | Name        | Price (THB) |
|--------|-------------|-------------|
| RED    | Red set     | 50          |
| GREEN  | Green set   | 40          |
| BLUE   | Blue set    | 30          |
| YELLOW | Yellow set  | 50          |
| PINK   | Pink set    | 80          |
| PURPLE | Purple set  | 90          |
| ORANGE | Orange set  | 120         |
```
### All promotions
Promotion catalog showing all active discount policies, eligibility conditions, and discount rates.
```text

--- Promotions ---

[MEMBER] Member card 10% off
 - Get 10% discount on the total bill if customer has a member card.

[PAIR] Pair discount 5% (ORANGE/PINK/GREEN)
 - Every pair (2 items of the same code) for ORANGE/PINK/GREEN gets 5% off that pair value.
```

### Discounts
1. **Pair Discount (5%)**
   - Eligible codes: `ORANGE`, `PINK`, `GREEN`
   - Bundle size: 2 (e.g., qty 2 => 1 bundle, qty 4 => 2 bundles)
   - Discount applies **per bundle subtotal**, not the whole order
2. **Member Discount (10%)**
   - Applied on **total after pair discount**

## Unit Test Cases 

### 1) Discount Policies (Pair / Member)
**Suite:** `TestQuoteOrder_DiscountPolicies`

#### 1.1 Pair Discount — Eligible code & bundle logic
- `Pair: GREEN(4) => 2 bundles`
- `Pair: GREEN(3) => 1 bundle + remainder`
- `Pair: PINK(2) => 1 bundle`
- `Pair: PINK(4) => 2 bundles`

#### 1.2 Pair Discount — Multi-code & non-eligible rules
- `Pair: GREEN(2) + ORANGE(2) => sum per code`
- `Pair: RED(2) not eligible => 0 pair discount`
- `No cross-code pairing: GREEN(1) + ORANGE(1) => 0 pair discount`

#### 1.3 Member Discount — Stacking order
- `Member without pair: RED(2), member=true => 10% off`
- `Member + multi-code: GREEN(2) + RED(1), member=true => member stacks after pair on total`

---

### 2) Success Scenarios (Happy case)
**Suite:** `TestQuoteOrderSuccess`

- `Success: mixed items, pair applies to GREEN(2), no member`
- `Success: member stacks after pair discount (GREEN 2)`

---

### 3) Failure Scenarios (Validation / Error Handling)
**Suite:** `TestQuoteOrderFail`

- `Fail: empty order`
- `Fail: invalid quantity`
- `Fail: invalid item code (blank)`

**Suite:** `TestQuoteOrderFail_UnknownMenuItem`
- `Fail: unknown menu item`

---

### 4) Normalization (Input Cleanup / Canonicalization)
**Suite:** `TestQuoteOrder_Normalization`

- `Normalize code: "_green_" => GREEN`
- `Normalize + merge qty: {"_green_":1,"GREEN":1} => GREEN=2`


## Usage Examples
### Example 1 
Input:
```json
{"items":{"RED":1,"GREEN":2},"member":false}
```
Expected:

```text
--- Order Items ---

--------+--------------+-----+------------+-----------
CODE    | NAME         | QTY | UNIT PRICE |   TOTAL.  
--------+--------------+-----+------------+-----------
RED     | Red set      |   1 | 50.00 THB  | 50.00 THB
GREEN   | Green set    |   2 | 40.00 THB  | 80.00 THB

--- Order Quote ---

Subtotal         : 130.00 THB
Pair Discount    : 4.00 THB
Member Discount  : 0.00 THB
Total            : 126.00 THB
```

### Example 2 (Pair discount)
Input:
```json
{"items":{"GREEN":2},"member":false}
```
Expected:

```text
--- Order Items ---

--------+--------------+-----+------------+-----------
CODE    | NAME         | QTY | UNIT PRICE |   TOTAL.  
--------+--------------+-----+------------+-----------
GREEN   | Green set    |   2 | 40.00 THB  | 80.00 THB

--- Order Quote ---

Subtotal         : 80.00 THB
Pair Discount    : 4.00 THB
Member Discount  : 0.00 THB
Total  
```
### Example 3 (Pair + Member)
Input:
```json
{"items":{"GREEN":2},"member":true}
```
Expected:

```text
--- Order Items ---

--------+--------------+-----+------------+-----------
CODE    | NAME         | QTY | UNIT PRICE |   TOTAL.  
--------+--------------+-----+------------+-----------
GREEN   | Green set    |   2 | 40.00 THB  | 80.00 THB

--- Order Quote ---

Subtotal         : 80.00 THB
Pair Discount    : 4.00 THB
Member Discount  : 7.60 THB
Total            : 68.40 THB
```

## Start Food-Shop-App using Docker

You can run this project in 3 ways:
1) Clone the repo and run locally in your terminal
2) Build the Docker image locally and run 
3) Pull the pre-built image from Docker Hub and run 

## Method 1: Clone and run locally (Terminal)
### Step 1: Clone the repository
```bash
git clone https://github.com/Tewapirat/Food-Shop.git
cd food_shop
```
### Step 2: Download dependencies
```bash
go mod download
```
### Step 3: Run the app
```bash
go run main.go
```
## Method 2: Build locally (Docker)

### Step 1: Clone the repository
```bash
git clone https://github.com/Tewapirat/Food-Shop.git
cd food_shop
```
### Step 2: Build the Docker image
Build the application image from the Dockerfile.
```bash
docker build -t food-shop-app:v1.0.0 .
```
### Step 3: Run in interactive mode
Interactive mode runs the app in your terminal, so you can type input and see the output immediately.
```bash
docker run --rm -it --name food-shop-app food-shop-app:v1.0.0
```
## Method 3: Pull from Docker Hub

### Step 1: Pull the image
```bash
docker pull tewdev/food-shop-app:latest  
```
### Step 2: Run in interactive mode
```bash
docker run --rm -it tewdev/food-shop-app:latest
```