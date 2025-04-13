table "account" {
  schema = schema.beyond-db

  column "user_id" {
    null = false
    type = varchar(36)
  }

  column "firebase_id" {
    null = false
    type = varchar(36)
  }

  column "created_at" {
    type = datetime
    null = false
  }

  primary_key {
    columns = [column.user_id]
  }
}

table "user" {
  schema = schema.beyond-db

  column "id" {
    null = false
    type = varchar(36)
  }

  column "name" {
    null = false
    type = varchar(100)
  }

  column "icon" {
    null = false
    type = text
  }

  column "created_at" {
    type = datetime
    null = false
  }

  primary_key {
    columns = [column.id]
  }

  foreign_key "account_user_fk" {
    columns     = [column.id]
    ref_columns = [table.account.column.user_id]
  }
}

schema "beyond-db" {
  charset = "utf8mb4"
  collate = "utf8mb4_0900_ai_ci"
}
