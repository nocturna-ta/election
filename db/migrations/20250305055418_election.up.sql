CREATE TABLE "candidates"(
    id uuid NOT NULL PRIMARY KEY,
    name_candidate varchar(255) NOT NULL,
    election_no int NOT NULL,
    is_active boolean NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP(6) WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP(6) WITH TIME ZONE NOT NULL DEFAULT now(),
    is_deleted bool NOT NULL DEFAULT FALSE
);

CREATE TABLE "candidate_detail" (
    id uuid NOT NULL PRIMARY KEY,
    candidate_id uuid NOT NULL,
    biodata text NOT NULL,
    visi text NOT NULL,
    misi text NOT NULL,
    program_kerja text NOT NULL,
    created_at TIMESTAMP(6) WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP(6) WITH TIME ZONE NOT NULL DEFAULT now(),
    is_deleted bool NOT NULL DEFAULT FALSE,
    FOREIGN KEY (candidate_id) REFERENCES elections(id) ON DELETE CASCADE ON UPDATE CASCADE
)

CREATE TABLE "partai"(
    id uuid NOT NULL PRIMARY KEY,
    name varchar(255) NOT NULL,
)
