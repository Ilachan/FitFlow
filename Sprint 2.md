# Sprint 2.md

## 1) Entire Team
- Make progress on issues uncompleted in Sprint 1  
- Integrate frontend and backend  

## 2) Frontend
- Write a very simple test using Cypress (can be as simple as clicking a button or filling a form)  
- Write unit tests specific to your chosen frontend framework  
- There is no specific number of unit tests to complete, but aim for a **1:1 unit test to function ratio**  

## 3) Backend
- Detailed documentation of your backend API (**saved in this Sprint 2.md**)  
- Write unit tests specific to your chosen backend framework  
- There is no specific number of unit tests to complete, but aim for a **1:1 unit test to function ratio**  

---

## Detailed Backend API

### Analytics API

#### Main route registration
- **Analytics route is exposed at:** `routes.go:53`

#### Endpoint
- **Method:** `GET`
- **Path:** `/users/:id/analytics`
- **Handler:** `class_api.go:177`

#### Authentication and authorization
- **Requires** `Authorization` header with **Bearer token**
- **Token user id must match** path user id, otherwise **403**
- **Header parsing and token extraction logic:** `class_api.go:211`
- **Forbidden check in analytics handler:** `class_api.go:191`

#### Query parameters
- `range` (optional): accepts `7d`, `1m`, `3m`
- **Default:** `7d`
- **Parsing/default behavior:** `class_api.go:196`
- **Normalization logic:** `Backend/service/class_service.go`

#### Service workflow
- **Entry point:** `class_service.go:152`
- Verifies user exists
- Backfills `UserDailyActivity` from enrollments before computing analytics
- Resolves date range start:
  - `7d` ⇒ now minus 7 days
  - `1m` ⇒ now minus 1 month
  - `3m` ⇒ now minus 3 months
- Computes:
  - total classes
  - active days
  - grouped daily classes
  - grouped category classes
- Category percentage is rounded to **2 decimals**

#### Response model
- **Defined in:** `activity.go:35`
- **Fields:**
  - `user_id`
  - `range`
  - `from_date`
  - `to_date`
  - `total_classes`
  - `total_time`
  - `active_days`
  - `daily`: array of `{ date, classes }`
  - `categories`: array of `{ category, classes, percentage }`

#### Success response
- **HTTP 200**
- **JSON shape:**
```json
{
  "analytics": {
    "user_id": 50,
    "range": "1m",
    "from_date": "2026-02-12",
    "to_date": "2026-03-12",
    "total_classes": 14,
    "total_time": 285,
    "active_days": 12,
    "daily": [
      { "date": "2026-03-10", "classes": 2 }
    ],
    "categories": [
      { "category": "Cardio", "classes": 6, "percentage": 42.86 }
    ]
  }
}
