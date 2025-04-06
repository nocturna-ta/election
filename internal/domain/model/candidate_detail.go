package model

import "github.com/google/uuid"

type CandidateDetail struct {
	BaseModel
	ID              uuid.UUID `db:"id"`
	CandidateID     uuid.UUID `db:"candidate_id"`
	Biodata         string    `db:"biodata"`
	Visi            string    `db:"visi"`
	Misi            string    `db:"misi"`
	ProgramKerja    string    `db:"program_kerja"`
	PartaiPendukung Partai    `db:"partai_pendukung"`
}

type Partai struct {
	ID         uuid.UUID `db:"id"`
	NamaPartai string    `db:"nama_partai"`
}
