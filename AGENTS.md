# agent.md

## 🧠 Project Overview

This is a zero-registration, collaborative expense tracking app for ad hoc groups (e.g. trips, events). The app runs as a browser-based PWA with no installation or sign-up requirement. Users are identified via phone numbers and log in through magic links sent by SMS.

The system supports multiple currencies, continuous debt simplification, and editable transactions for all verified users. It is optimized for low-friction, short-term group use.

---

## 🧱 Tech Stack & Architecture

- **Frontend:** Vanilla JavaScript + Web Components + HTML (PWA)
- **Backend:** Go (Golang) using `chi` router
- **Database:** SQLite (or Postgres)
- **Authentication:** SMS magic links, token verification, HTTP-only cookie session
- **Hosting:** Fly.io
- **CI/CD:** GitHub Actions deploys on push to `main`

---

Planned features:

## 👥 User Model

- **Organizer:** creates group, sends invites, manages expenses, triggers "final call" reminders
- **Participant (Verified):** joins via magic link, can add/edit/delete any transaction
- **Placeholder (Unverified):** added by phone number, not yet clicked link, can be replaced or remain inactive

---

## 💸 Expense Model

Each transaction includes:

- `payer` (phone number)
- `amount` + `currency`
- `fx_rate` (fetched at time of entry, stored)
- `converted_amount` (in group base currency)
- `split`: list of phone numbers + individual owed amounts
- `description`
- `is_reimbursement` flag
- `timestamp`
- optional `note`

---

## 💱 Currency Handling

- Each group has a **base currency**
- Expenses can be logged in any currency
- FX rate is fetched at time of transaction using an external API
- Converted amount is stored and used for all calculations

---

## 🔁 Debt Simplification

- Net balances per user are updated after every transaction
- A simplification algorithm reduces debt chains (e.g. A → B → C → A → C)
- There is **no group finalization state**
- Admin can send a "final call" reminder, but group remains active and editable

---

## 🔐 Authentication & Sessions

- Magic link tokens are verified on click
- On success, a session is created using an HTTP-only cookie (or JWT)
- No password or manual login

---

## 🔔 Notifications

- Invite SMS: sent on group creation, includes tokenized magic link
- Final Call SMS/Email: admin-triggered prompt to reconcile expenses
