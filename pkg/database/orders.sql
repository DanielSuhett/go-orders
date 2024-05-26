CREATE TYPE PMETHODS AS ENUM ('PIX', 'BOLETO');

CREATE TABLE public.orders (
  id UUID NOT NULL DEFAULT gen_random_uuid(),
  product_category VARCHAR(255) NULL,
  product_value INT8 NULL,
  payment_method PMETHODS,
  payment_value INT8 NULL,
  labels text[],
  CONSTRAINT orders_pkey PRIMARY KEY (id ASC)
);
