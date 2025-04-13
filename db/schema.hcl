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
}

schema "beyond-db" {
  charset = "utf8mb4"
  collate = "utf8mb4_0900_ai_ci"
}
