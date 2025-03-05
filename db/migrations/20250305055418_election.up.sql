CREATE TABLE "elections"(
    id uuid NOT NULL PRIMARY KEY,
    name_candidate varchar(255) NOT NULL,
    election_no int NOT NULL,
    is_active boolean NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP(6) WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP(6) WITH TIME ZONE NOT NULL DEFAULT now(),
    is_deleted bool NOT NULL DEFAULT FALSE
);
