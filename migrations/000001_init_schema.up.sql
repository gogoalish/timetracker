CREATE TABLE IF NOT EXISTS "tasks" (
  "id" serial PRIMARY KEY,
  "user_id" int NOT NULL,
  "description" varchar NOT NULL,
  "start_dt" timestamp,
  "end_dt" timestamp,
  "created_at" timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS "people" (
  "id" serial PRIMARY KEY,
  "name" varchar NOT NULL,
  "surname" varchar NOT NULL,
  "patronymic" varchar,
  "passport_number" int NOT NULL,
  "passport_serie" int NOT NULL,
  "address" varchar NOT NULL
);

ALTER TABLE "tasks" ADD FOREIGN KEY ("user_id") REFERENCES "people" ("id");