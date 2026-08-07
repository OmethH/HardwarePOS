Good choice. For a first real deployment to a hardware shop, **Option 1 (local web app packaged like software)** is a very practical architecture. 

Since you want to use an AI coding agent, the prompt needs to be very specific about: 

- product requirements 

- architecture 

- database design 

- coding standards 

- development phases 

- expected outputs 

Below is a complete project specification + agent prompt. 

You can give this directly to your coding agent. 

# **Hardware POS System — Product Requirements & AI Agent Build Prompt Project Name** 

### **HardwarePOS** 

## **Product Overview** 

Build a local desktop-like POS and inventory management system for a hardware shop. 

The application should run on a shop laptop/PC without requiring internet access. 

The system manages: 

- Products 

- Inventory 

- Purchases 

- Sales 

- Stock movements 

- Customers 

- Suppliers 

- Reports 

- Users and permissions 

The goal is to replace manual stock books and spreadsheets with a simple, reliable business management system. 

# **Technology Stack** 

## **Frontend** 

Use: 

   - Vue 3 

   - TypeScript Vite 

   - Vue Router 

   - Pinia 

   - Tailwind CSS Axios 

- UI: 

- Use a modern dashboard design Responsive layout 

- Sidebar navigation 

- Reusable components 

## **Backend** 

Use: 

- Golang 

- Gin framework 

- REST API architecture 

- JWT authentication 

- Role-based access control 

Structure backend using clean architecture principles. 

Example: 

```
backend/
```

```
cmd/
 └── server/
```

```
internal/
```

```
 ├── auth/
 ├── users/
 ├── products/
```

```
 ├── inventory/
 ├── sales/
 ├── purchases/
```

```
 ├── reports/
```

```
database/
```

```
migrations/
```

```
config/
```

## **Database** 

Use: 

### **SQLite** 

### Reason: 

Runs locally 

- No database installation needed Easy backup 

- Perfect for single-machine POS 

Use: 

GORM ORM 

Database file: 

```
data/hardwarepos.db
```

# **Application Architecture** 

The system should work like: 

```
                Vue Frontend
```

```
                    |
                    |
               REST API
```

```
                    |
```

```
                    |
```

```
              Golang Backend
```

```
                    |
```

```
                    |
```

```
              SQLite Database
```

The final application should be launchable from: 

```
HardwarePOS.exe
```

# **User Roles** 

Implement RBAC. 

## **Admin** 

Full access: 

- Manage users 

- Manage products 

- Manage inventory View reports 

- Manage settings 

## **Manager** 

Can: 

- Manage products Manage suppliers View reports 

- Handle returns 

- Manage inventory 

## **Cashier** 

Can: 

- Create sales 

- Search products Print receipts 

- View own transactions 

Cannot: 

- Delete products Change purchase prices Delete completed sales Modify stock manually 

# **Core Modules** 

# **1. Authentication** 

Features: 

- Login page JWT authentication Password hashing Role permissions 

Database: 

```
users
```

```
id
name
username
password_hash
role_id
created_at
```

# **2. Product Management** 

CRUD operations. 

Product fields: 

```
id
```

```
sku
```

```
barcode
```

```
name
```

```
category_id
```

`brand description unit_type purchase_price selling_price minimum_stock current_stock status created_at updated_at` Unit types: 

```
Piece
Box
Meter
Kg
Liter
Pack
```

### Example: 

```
Electrical Cable
Unit: Meter
PVC Pipe
Unit: Piece
```

# **3. Categories** 

CRUD. 

Example: 

```
Electrical
Plumbing
Tools
Paint
```

```
Construction Materials
Hardware
```

# **4. Supplier Management** 

Fields: 

```
id
```

```
name
```

```
phone
```

```
email
```

```
address
```

```
company
```

# **5. Stock Management** 

IMPORTANT: 

Never directly modify stock. 

Every stock change must create a stock movement. 

Stock movement: 

```
id
```

```
product_id
```

```
quantity
```

```
type
```

```
reference_id
```

```
created_at
```

Types: 

```
PURCHASE
```

```
SALE
```

```
RETURN
```

```
ADJUSTMENT
```

### Example: 

```
PVC Pipe
+100
Purchase
-5
Sale
+2
Return
```

# **6. Purchase Module** 

Used when new inventory arrives. 

Flow: 

```
Create Purchase
```

```
        ↓
Select Supplier
        ↓
Add Products
        ↓
Enter Quantity
        ↓
Enter Purchase Price
        ↓
Confirm
        ↓
```

```
Increase Stock
```

### Database: 

### Purchase: 

```
id
```

```
supplier_id
```

```
total_amount
```

```
created_by
```

```
created_at
```

### Purchase Items: 

```
purchase_id
```

```
product_id
```

```
quantity
```

```
price
```

# **7. POS Sales Module** 

Main cashier screen. 

Features: 

Search products 

- Add products to cart 

- Change quantity 

- Apply discount 

- Calculate totals 

- Complete sale 

Example: 

```
Product          Qty     Price
```

```
PVC Pipe          5       900
```

```
Wall Plug        10       250
```

```
Subtotal:
1150
```

```
Discount:
50
```

```
Total:
1100
```

### On checkout: 

System must: 

**1.** Create sale 

**2.** Create sale items 

**3.** Reduce stock 

**4.** Create stock movement 

**5.** Calculate profit 

All operations must happen inside a database transaction. 

# **8. Sales History** 

Features: 

View previous sales 

- Search invoice number 

- Filter by date 

- View invoice details 

Sale table: 

```
id
```

```
invoice_number
```

```
customer_id
```

```
subtotal
```

```
discount
```

```
total
```

```
payment_method
```

```
created_by
```

```
created_at
```

Sale items: 

```
sale_id
```

```
product_id
```

```
quantity
```

```
selling_price
```

```
purchase_price
```

```
profit
```

# **9. Returns** 

Never delete sales. 

Create return transactions. Flow: 

```
Select invoice
```

```
       ↓
Select item
       ↓
```

```
Enter return quantity
```

```
       ↓
Increase stock
       ↓
Adjust profit
```

# **10. Customers** 

Optional customer information. 

Fields: 

```
id
```

```
name
```

```
phone
```

```
address
```

Default customer: 

```
Walk-in Customer
```

# **11. Dashboard** 

Create analytics dashboard. 

Show: 

## **Today** 

- Sales amount Profit 

- Number of transactions Items sold 

## **Inventory** 

- Total products 

- Low stock products Out of stock products 

## **Charts** 

Add: 

- Sales over time 

- Top selling products Profit trend 

Use: 

Chart.js 

# **12. Reports** 

Generate: 

## **Sales Report** 

Filters: 

Daily Weekly Monthly Custom range 

Show: 

```
Date
```

```
Invoices
```

```
Revenue
```

```
Discount
```

```
Profit
```

## **Inventory Report** 

Show: 

```
Product
```

```
Current Stock
```

```
Stock Value
```

```
Status
```

**Product Performance** 

Show: 

```
Product
```

```
Quantity Sold
```

```
Revenue
```

```
Profit
```

### Export: 

PDF CSV 

# **13. Receipt Generation** 

Generate printable receipts. 

Support: 

A4 invoice Thermal printer format 

### Receipt: 

#### `HARDWARE SHOP NAME` 

```
Invoice:
INV00001
Date:
```

```
------------------
```

```
PVC Pipe
5 x 180
Wall Plug
10 x 25
```

```
------------------
```

```
TOTAL:
1150
```

```
Thank You
```

# **14. Settings** 

Allow: 

- Shop name Address Phone number Receipt footer Tax percentage 

# **Frontend Pages** 

Create: 

```
/login
```

```
/dashboard
/products
/categories
/inventory
/purchases
/pos
/sales
/customers
/suppliers
/reports
/users
/settings
```

# **Development Rules** 

Follow these rules: 

**1.** Write clean maintainable code. 

**2.** Use TypeScript everywhere in frontend. 

**3.** Use reusable Vue components. 

**4.** Use API services separately. 

**5.** Validate all forms. 

**6.** Handle errors properly. 

**7.** Use database migrations. 

**8.** Add comments where business logic exists. 

**9.** Never hard delete transactions. 

**10.** Use database transactions for sales and returns. 

# **Build Order** 

Implement in this order: 

## **Phase 1** 

Setup: 

- Vue project 

- Go backend 

- SQLite database 

- Authentication 

## **Phase 2** 

Inventory: 

- Products 

- Categories 

- Suppliers 

- Stock movements 

## **Phase 3** 

Purchasing: 

- Purchase creation 

- Stock increases 

## **Phase 4** 

POS: 

- Cart 

- Checkout 

- Invoice generation 

- Stock deduction 

## **Phase 5** 

Reports: 

- Dashboard 

- Analytics PDF exports 

## **Phase 6** 

Desktop packaging: 

Create: 

```
HardwarePOS.exe
```

using: 

Tauri 

# **Final Agent Instruction** 

```
Act as a senior full-stack engineer.
```

```
Build this HardwarePOS system following the complete specification above.
Do not create a simple CRUD application.
```

```
Think like a production business application.
```

```
Prioritize:
```

- `Data integrity` 

- `Transaction safety` 

- `Clean architecture` 

- `Maintainability` 

- `Good UI/UX` 

```
Implement the system module by module.
```

```
Before coding:
```

`1. Create the folder structure.` 

`2. Create the database schema.` 

`3. Create API documentation.` 

`4. Then implement features.` 

```
After each module:
```

- `Test functionality` 

- `Fix errors` 

```
- Explain what was completed
```

```
The final goal is a working local POS application that can be installed on a hardwa
```

This prompt should give an AI agent enough context to start building a **real product** , not just a demo CRUD app. For even better results, I would give the agent this document first, then ask it to create the **database schema + API contract before writing code** . 

Is this conversation helpful so far? 

