BEGIN;

CREATE TABLE public.items (
    id bigint NOT NULL,
    invoice_id VARCHAR(10) NOT NULL,
    item_id UUID NOT NULL,
    name character varying(255),
    type character varying(255),
    quantity numeric,
    unit_price numeric,
    amount numeric,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp with time zone
);

CREATE SEQUENCE public.items_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.items_id_seq OWNED BY public.items.id;

ALTER TABLE ONLY public.items ALTER COLUMN id SET DEFAULT nextval('public.items_id_seq'::regclass);

ALTER TABLE ONLY public.items
    ADD CONSTRAINT items_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.items
    ADD CONSTRAINT invoice_id FOREIGN KEY (invoice_id) REFERENCES public.invoices(invoice_id);

COMMIT;