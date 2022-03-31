CREATE TABLE conn_attemps (
    id bigserial PRIMARY KEY,

    time        timestamp with time zone,
    port        text NOT NULL,
    ip          text NOT NULL,
    country_code text NOT NULL,
    client_port  text NOT NULL,
    bytes       bytea NOT NULL,

    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);
ALTER TABLE conn_attemps OWNER TO honeypot_user;
