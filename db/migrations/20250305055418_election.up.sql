CREATE TABLE IF NOT EXISTS "parties" (
    id UUID NOT NULL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    logo_path VARCHAR(512),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE
    );

-- Table for election pairs (president and vice president)
CREATE TABLE IF NOT EXISTS "election_pairs" (
    id UUID NOT NULL PRIMARY KEY,
    election_no VARCHAR(20) NOT NULL UNIQUE,
    vote_count INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    pair_photo_path VARCHAR(512),

    -- President candidate info
    president_full_name VARCHAR(255) NOT NULL,
    president_education_history TEXT,
    president_work_experience TEXT,
    president_legal_record_history TEXT,
    president_photo_path VARCHAR(512),

    -- Vice president candidate info
    vice_president_full_name VARCHAR(255) NOT NULL,
    vice_president_education_history TEXT,
    vice_president_work_experience TEXT,
    vice_president_legal_record_history TEXT,
    vice_president_photo_path VARCHAR(512),

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE
    );

-- Table for vision, mission, and work programs
CREATE TABLE IF NOT EXISTS "election_pair_details" (
    id UUID NOT NULL PRIMARY KEY,
    election_pair_id UUID NOT NULL REFERENCES election_pairs(id) ON DELETE CASCADE,
    vision TEXT,
    mission TEXT,
    work_program TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,

    CONSTRAINT uq_pair_details UNIQUE(election_pair_id, is_deleted)
    );


-- Table for supporting parties (many-to-many relationship)
CREATE TABLE IF NOT EXISTS "supporting_parties" (
    id UUID NOT NULL PRIMARY KEY,
    election_pair_id UUID NOT NULL REFERENCES election_pairs(id) ON DELETE CASCADE,
    party_id UUID NOT NULL REFERENCES parties(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,

    CONSTRAINT uq_supporting_party UNIQUE(election_pair_id, party_id, is_deleted)
);

CREATE TABLE IF NOT EXISTS program_documents (
    id UUID NOT NULL PRIMARY KEY,
    election_pair_detail_id UUID NOT NULL,
    document_path VARCHAR(512) NOT NULL,
    document_type VARCHAR(50) NOT NULL,
    original_filename VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,

    CONSTRAINT fk_program_doc_detail
    FOREIGN KEY (election_pair_detail_id)
    REFERENCES election_pair_details(id)
    ON DELETE CASCADE
);


-- Indexes for better query performance
CREATE INDEX idx_election_pairs_election_no ON election_pairs(election_no);
CREATE INDEX idx_election_pairs_is_active ON election_pairs(is_active) WHERE is_deleted = FALSE;
CREATE INDEX idx_supporting_parties_election_pair_id ON supporting_parties(election_pair_id) WHERE is_deleted = FALSE;
CREATE INDEX idx_supporting_parties_party_id ON supporting_parties(party_id) WHERE is_deleted = FALSE;
CREATE INDEX idx_program_docs_detail_id ON program_documents(election_pair_detail_id) WHERE is_deleted = FALSE;