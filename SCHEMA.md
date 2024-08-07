# ER Diagrams

Here are the ER diagrams for the database schema.

## User Auth


```mermaid

erDiagram
    USERS {
        id int
        hash varchar
        username varchar
        email varchar
        validated_at datetime
        created_at datetime
        updated_at datetime
    }

    INTERNAL_LOGINS {
        id int
        user_id int
        email varchar
        password varchar
        password_salt varchar
        algorithm int
        password_last_updated_at datetime
        login_attempts int
        last_login_attempt datetime
        last_login_success datetime
        created_at datetime
        updated_at datetime
    }

    LOGIN_TOKENS {
        id int
        user_id int
        token varchar
        token_type int
        expires_at datetime
        used_at datetime
        created_at datetime
    }

    EXTERNAL_LOGIN_PROVIDERS {
        id int
        name varchar
        type int
        endpoint varchar
        enabled boolean
        created_at datetime
        updated_at datetime
    }

    EXTERNAL_LOGINS {
        id int
        user_id int
        provider_id varchar
        created_at datetime
    }

    EXTERNAL_LOGIN_TOKENS {
        id int
        external_login_id int
        login_ip varchar
        user_agent varchar
        logged_in_at datetime
        token text
        created_at datetime
    }

    USER_SESSIONS {
        id varchar
        user_id int
        login_ip varchar
        user_agent varchar
        logged_in_at datetime
        logged_out_at datetime
        created_at datetime
        updated_at datetime
    }

    EXTERNAL_LOGIN_PROVIDERS ||--|{ EXTERNAL_LOGINS : contains
    EXTERNAL_LOGIN_TOKENS }|--|| EXTERNAL_LOGINS : contains
    EXTERNAL_LOGINS }o--|| USERS : contains
    LOGIN_TOKENS }|--|| USERS : contains
    INTERNAL_LOGINS |o--|| USERS : contains
    USERS ||--|{ USER_SESSIONS : contains
```

## Core Functionalities

```mermaid

erDiagram
    CONTACTLISTS {
        id int
        user_id int
        contact_id int
        created_at datetime
        updated_at datetime
        deleted_at datetime
    }

    GROUPS {
        id int
        name varchar
        created_by int
        created_at datetime
        updated_at datetime
    }

    USER_GROUPS {
        id int
        user_id int
        group_id int
        joined_at datetime
        left_at datetime
        status int
    }

    ACCOUNTS {
        id int
        created_by int
        group_id int
        currency_id int
        balance float
        name varchar
        created_at datetime
        updated_at datetime
    }

    BUDGETS {
        id int
        name varchar
        account_id int
        amount float
        created_by int
        start_date date
        end_date date
        created_at datetime
        updated_at datetime
    }

    TRANSACTIONS {
        id int
        category varchar
        currency varchar
        currency_rate float
        transaction_type int
        budget_id int
        amount float
        date date
        executed_by int
        description text
        receipt_url varchar
        registered_at datetime
        registered_by int
        created_at datetime
        updated_at datetime
    }

    CATEGORIES {
        id int
        name varchar
        icon varchar
        color varchar
        description text
        order int
    }

    USER_PREFERENCES {
        user_id int
        key varchar
        value text
        upserted_at datetime
    }

    USERDATA {
        user_id int
        avatar_url varchar
        display_name varchar
        currency varchar
        language varchar
        timezone varchar
        upserted_at datetime
    }

    ACCOUNTS ||--|{ BUDGETS : contains
    CATEGORIES ||--o{ TRANSACTIONS : contains
    BUDGETS ||--o{ TRANSACTIONS : contains
    GROUPS ||--|{ ACCOUNTS : contains
    USERDATA ||--o{ CONTACTLISTS : contains
    USERDATA ||--|{ USER_GROUPS : contains
    GROUPS ||--|{ USER_GROUPS : contains
    USER_PREFERENCES }|--|| USERDATA : contains
```

## Currency Conversion

```mermaid

erDiagram
    CURRENCIES {
        id int PK
        name varchar UK
        code varchar UK
        symbol varchar
    }

    EXCHANGE_RATES {
        id int PK
        currency_id int FK 
        base_currency_id int FK
        rate float
        date date
    }

    CURRENCY_CONVERSION_PROVIDERS {
        id int PK
        name varchar
        type int
        endpoint varchar UK
        enabled boolean
        params json
        run_at string
    }

    CURRENCIES ||--|{ EXCHANGE_RATES : contains
    CURRENCY_CONVERSION_PROVIDERS }|--|{ CURRENCIES : ""
```

## Notifications

```mermaid

erDiagram

    NOTIFICATION_TEMPLATES {
        id int
        name varchar
        subject varchar
        body text
        notification_type int
        notification_channel int
        created_at datetime
        updated_at datetime
    }

    NOTIFICATION_TYPES {
        id int
        name varchar
        description text
    }

    NOTIFICATION_CHANNELS {
        id int
        channel varchar
        created_at datetime
        updated_at datetime
    }

    NOTIFICATIONS {
        id int
        user_id int
        template_id int
        created_at datetime
        read_at datetime
    }

    NOTIFICATION_CHANNELS ||--o{ NOTIFICATION_TEMPLATES : contains
    NOTIFICATION_TYPES ||--o{ NOTIFICATION_TEMPLATES : contains
    NOTIFICATION_TEMPLATES ||--o{ NOTIFICATIONS : contains    
```
