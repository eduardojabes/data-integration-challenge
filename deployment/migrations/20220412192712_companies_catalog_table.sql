-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS companies_catalog_table (
    cc_company_id UUID PRIMARY KEY, 
    cc_name TEXT, 
    cc_zip VARCHAR(5),
    cc_website TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS companies_catalog_table;
-- +goose StatementEnd
