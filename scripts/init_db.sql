-- Table: public.users

-- DROP TABLE public.users;

CREATE TABLE IF NOT EXISTS public.users
(
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    first_name character varying(50) COLLATE pg_catalog."default" NOT NULL,
    last_name character varying(50) COLLATE pg_catalog."default" NOT NULL,
    nickname character varying(30) COLLATE pg_catalog."default" NOT NULL,
    password character varying(50) COLLATE pg_catalog."default" NOT NULL,
    email character varying(100) COLLATE pg_catalog."default" NOT NULL,
    country character varying(10) COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT users_nickname_key UNIQUE (nickname)
    )

    TABLESPACE pg_default;

ALTER TABLE public.users
    OWNER to faceit;

GRANT ALL ON TABLE public.users TO faceit;

-- FUNCTION: public.update_updated_at()

-- DROP FUNCTION public.update_updated_at();

CREATE FUNCTION public.update_updated_at()
    RETURNS trigger
    LANGUAGE 'plpgsql'
    COST 100
    VOLATILE NOT LEAKPROOF
AS $BODY$
BEGIN
    NEW.updated_at := current_timestamp;

    RETURN NEW;
END;
$BODY$;

ALTER FUNCTION public.update_updated_at()
    OWNER TO faceit;


-- Trigger: update_updated_at

-- DROP TRIGGER update_updated_at ON public.users;

CREATE TRIGGER update_updated_at
    BEFORE UPDATE
    ON public.users
    FOR EACH ROW
    EXECUTE FUNCTION public.update_updated_at();