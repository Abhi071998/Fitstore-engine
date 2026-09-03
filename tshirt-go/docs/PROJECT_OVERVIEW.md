# FitStore Admin Backend — What It Does (Plain-Language Overview)

This explains what this piece of the FitStore system does and why it exists,
without assuming any technical background. For developer-facing details, see
[CONTRIBUTOR_GUIDE.md](CONTRIBUTOR_GUIDE.md).

## The big picture: three separate pieces working together

FitStore isn't one single program — it's three pieces that talk to each other:

1. **The customer-facing store** ("fitstore-core") — the website customers
   actually shop on: browse products, add to cart, place an order.
2. **The admin panel** ("fitstore-ui") — the screens a store admin sees and
   clicks around in: add a product, approve an order, edit the About Us page.
3. **This service** ("fitstore-engine" / admin backend) — the behind-the-
   scenes engine the admin panel talks to. It's the part that actually reads
   from and writes to the store's database whenever an admin does something.

Think of it like a restaurant: the customer-facing store is the dining room,
the admin panel is the manager's tablet, and this service is the kitchen that
actually prepares whatever the manager asks for.

Both the customer store and this admin engine share the **same underlying
database** — so an order a customer places on the storefront is the same
order an admin sees and approves here.

## What an admin can actually do through this system

### Manage the product catalog

- **Category Types** — a fixed, pre-approved list of category names (like
  "T-shirt", "Shoes", "Hoodie"). This exists so that when an admin creates a
  new category, they pick from a dropdown instead of typing a name freely —
  it keeps naming consistent and avoids typos/duplicates like "Tshirt" vs
  "T-Shirt" vs "tshirt" all existing as separate categories.
  - If an admin renames a category type later (say, "T-shirt" → "Shirts"),
    every category already using that type updates its name automatically —
    no need to manually rename each one.
- **Categories** — the actual groupings shown in the store (e.g. a "T-shirt"
  category with its own image). Each one is optionally tied to a Category
  Type from the dropdown above.
  - Deleting a category also removes every product filed under it, and the
    system tells the admin how many products were removed along with it —
    so nothing gets accidentally orphaned or left behind.
- **Products** — the individual items for sale: name, price, description,
  images, and how much stock exists per size (S/M/L/XL/etc).

### Review and approve customer orders

When a customer places an order on the storefront, it doesn't ship
automatically — it waits for an admin to review it here. An admin can:
- See every order still waiting for a decision, grouped by customer.
- Approve or reject an individual order (a reason is required when rejecting).
- Approve a whole batch of orders at once instead of one at a time.

**Current limitation:** this order-review feature depends on the
customer-facing store also being fully set up and connected to the same
database. Until that connection is complete, the "pending orders" screen
will simply show an empty list rather than an error — so admins won't see
broken pages, they'll just see "nothing to review yet" until that other
piece is finished.

### Edit the About Us page

Admins can set (once) and then update the store's About Us content — the
hero image, title, description, a handful of taglines, the founding year,
address, and contact email — all shown on the storefront's About Us page.

## Who admins are, and how they get in

Anyone managing the store first creates an admin account (email, name,
password) and logs in. Every action described above — creating a category,
approving an order, editing About Us — requires being logged in; only
viewing the About Us page and checking that the system is online are open
without logging in.

## Where the "source of truth" for data lives

Even though this service reads and writes almost all of the store's data
(categories, products, orders), it does not decide what the underlying
database structure looks like — that responsibility belongs entirely to the
customer-facing store's side of the system. This service simply assumes the
tables it needs already exist in the right shape, and if the two sides ever
get out of sync (say, a field this service expects hasn't been added on the
other side yet), that specific feature will fail until it's added there
first. This is a deliberate handoff, not an oversight — it keeps one system
in charge of the shape of the data everyone shares.
