--
-- PostgreSQL database dump
--

-- Dumped from database version 12.2 (Ubuntu 12.2-4)
-- Dumped by pg_dump version 12.2 (Ubuntu 12.2-4)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: aliases; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.aliases (
    id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    name character varying(255) NOT NULL
);


ALTER TABLE public.aliases OWNER TO postgres;

--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migration (
    version character varying(14) NOT NULL
);


ALTER TABLE public.schema_migration OWNER TO postgres;

--
-- Name: supers; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.supers (
    id uuid NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    original_id integer NOT NULL,
    name character varying(255) NOT NULL,
    full_name character varying(255),
    place_of_birth text,
    first_appearance character varying(255),
    alter_egos text NOT NULL,
    publisher character varying(255) NOT NULL,
    alignment character varying(255) NOT NULL,
    gender character varying(255),
    race character varying(255),
    height_feet character varying(255) NOT NULL,
    height_cm integer NOT NULL,
    weight_lb character varying(255) NOT NULL,
    weight_kg integer NOT NULL,
    eye_color character varying(255),
    hair_color character varying(255),
    occupation character varying(255),
    base character varying(255),
    image character varying(255),
    intelligence integer NOT NULL,
    strength integer NOT NULL,
    speed integer NOT NULL,
    durability integer NOT NULL,
    power integer NOT NULL,
    combat integer NOT NULL
);


ALTER TABLE public.supers OWNER TO postgres;

--
-- Name: aliases aliases_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.aliases
    ADD CONSTRAINT aliases_pkey PRIMARY KEY (id);


--
-- Name: supers supers_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.supers
    ADD CONSTRAINT supers_pkey PRIMARY KEY (id);


--
-- Name: schema_migration_version_idx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE UNIQUE INDEX schema_migration_version_idx ON public.schema_migration USING btree (version);


--
-- PostgreSQL database dump complete
--

