--
-- PostgreSQL database dump
--

\restrict E4R2YLYmrla6Qk6es0jyyCZ61ekwCxf7SJkuHx2KnCtLBsrM45THugvOawGlaAc

-- Dumped from database version 15.16
-- Dumped by pg_dump version 15.16

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

--
-- Name: pgcrypto; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;


--
-- Name: EXTENSION pgcrypto; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION pgcrypto IS 'cryptographic functions';


--
-- Name: update_daily_usage_summary(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.update_daily_usage_summary() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    INSERT INTO daily_usage_summary (
        user_id,
        date,
        total_requests,
        total_input_tokens,
        total_output_tokens,
        total_tokens,
        total_cost,
        models_used,
        providers_used
    )
    VALUES (
        NEW.user_id,
        DATE(NEW.created_at),
        1,
        NEW.input_tokens,
        NEW.output_tokens,
        NEW.total_tokens,
        NEW.total_cost,
        jsonb_build_object(NEW.model_requested, 1),
        jsonb_build_object(NEW.provider, 1)
    )
    ON CONFLICT (user_id, date)
    DO UPDATE SET
        total_requests = daily_usage_summary.total_requests + 1,
        total_input_tokens = daily_usage_summary.total_input_tokens + NEW.input_tokens,
        total_output_tokens = daily_usage_summary.total_output_tokens + NEW.output_tokens,
        total_tokens = daily_usage_summary.total_tokens + NEW.total_tokens,
        total_cost = daily_usage_summary.total_cost + NEW.total_cost,
        models_used = daily_usage_summary.models_used || jsonb_build_object(NEW.model_requested,
            COALESCE((daily_usage_summary.models_used->NEW.model_requested)::integer, 0) + 1),
        providers_used = daily_usage_summary.providers_used || jsonb_build_object(NEW.provider,
            COALESCE((daily_usage_summary.providers_used->NEW.provider)::integer, 0) + 1),
        updated_at = NOW();

    RETURN NEW;
END;
$$;


ALTER FUNCTION public.update_daily_usage_summary() OWNER TO postgres;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: api_keys; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.api_keys (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid,
    api_key text NOT NULL,
    active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public.api_keys OWNER TO postgres;

--
-- Name: api_usage_analytics; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.api_usage_analytics (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    api_key_id uuid,
    request_id character varying(255) NOT NULL,
    model_requested character varying(255) NOT NULL,
    model_used character varying(255) NOT NULL,
    provider character varying(100) NOT NULL,
    input_tokens integer DEFAULT 0 NOT NULL,
    output_tokens integer DEFAULT 0 NOT NULL,
    total_tokens integer DEFAULT 0 NOT NULL,
    input_cost numeric(15,8) DEFAULT 0 NOT NULL,
    output_cost numeric(15,8) DEFAULT 0 NOT NULL,
    total_cost numeric(15,8) DEFAULT 0 NOT NULL,
    input_price_per_token numeric(15,8) NOT NULL,
    output_price_per_token numeric(15,8) NOT NULL,
    status character varying(50) DEFAULT 'success'::character varying NOT NULL,
    error_message text,
    response_time_ms integer,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    currency character varying(10) DEFAULT 'INR'::character varying
);


ALTER TABLE public.api_usage_analytics OWNER TO postgres;

--
-- Name: api_usage_logs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.api_usage_logs (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid,
    api_key_id uuid,
    model_id uuid,
    input_tokens integer DEFAULT 0,
    output_tokens integer DEFAULT 0,
    cost numeric(12,4) DEFAULT 0,
    created_at timestamp with time zone DEFAULT now(),
    model character varying(255),
    provider character varying(100),
    currency character varying(10) DEFAULT 'INR'::character varying
);


ALTER TABLE public.api_usage_logs OWNER TO postgres;

--
-- Name: chat_messages; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.chat_messages (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    chat_id uuid NOT NULL,
    role character varying(20) NOT NULL,
    content text NOT NULL,
    token_count integer DEFAULT 0,
    cost numeric(12,6) DEFAULT 0,
    created_at timestamp with time zone DEFAULT now(),
    CONSTRAINT chat_messages_role_check CHECK (((role)::text = ANY (ARRAY[('user'::character varying)::text, ('assistant'::character varying)::text, ('system'::character varying)::text])))
);


ALTER TABLE public.chat_messages OWNER TO postgres;

--
-- Name: chats; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.chats (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    title character varying(255) DEFAULT 'New Chat'::character varying NOT NULL,
    model character varying(100) NOT NULL,
    provider character varying(100) NOT NULL,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public.chats OWNER TO postgres;

--
-- Name: credits; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.credits (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid,
    total_credits numeric(12,2) DEFAULT 0,
    used_credits numeric(12,2) DEFAULT 0,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public.credits OWNER TO postgres;

--
-- Name: daily_usage_summary; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.daily_usage_summary (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    date date NOT NULL,
    total_requests integer DEFAULT 0 NOT NULL,
    total_input_tokens bigint DEFAULT 0 NOT NULL,
    total_output_tokens bigint DEFAULT 0 NOT NULL,
    total_tokens bigint DEFAULT 0 NOT NULL,
    total_cost numeric(15,8) DEFAULT 0 NOT NULL,
    models_used jsonb DEFAULT '{}'::jsonb,
    providers_used jsonb DEFAULT '{}'::jsonb,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public.daily_usage_summary OWNER TO postgres;

--
-- Name: models; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.models (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name character varying(255) NOT NULL,
    provider character varying(100) NOT NULL,
    input_token_price numeric(12,4) NOT NULL,
    output_token_price numeric(12,4) NOT NULL,
    active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public.models OWNER TO postgres;

--
-- Name: payments; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.payments (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid,
    razorpay_order_id character varying(255) NOT NULL,
    razorpay_payment_id character varying(255),
    amount numeric(12,2) NOT NULL,
    status character varying(50) DEFAULT 'pending'::character varying,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public.payments OWNER TO postgres;

--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.schema_migrations (
    version character varying NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO postgres;

--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    password_hash text NOT NULL,
    email_verified boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Data for Name: api_keys; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.api_keys (id, user_id, api_key, active, created_at) FROM stdin;
90891216-82f9-43fd-8919-583c0356004c	3e419448-bf28-4683-bed3-f34dbda297c5	rb5tt5LPA4kd4ecQMS_xNVxQ0-yTz-3RcF1Gb1oaWcw=	t	2026-02-26 08:11:56.859964+00
935f5496-ad47-4c68-b2a5-161f4b27d9e6	45da1969-60fd-4248-9efa-d061e34e3a1d	U9pBpLxA1ixkkllNRJlB_cBBeHUZaVmoyaOCP_eR4rE=	t	2026-02-26 09:52:01.983516+00
bc94f290-d8e0-4fcc-87d5-8c0e92a6aed2	45da1969-60fd-4248-9efa-d061e34e3a1d	QL9FqPs_XxpT8X-inYpGZm4mJNxj5_P7Qwui9RYwbkM=	t	2026-02-27 03:24:31.685002+00
49224a0a-a79b-4c1e-9e6c-ccf9d3c36e45	45da1969-60fd-4248-9efa-d061e34e3a1d	r5WJROw4CSpGRjhIdmB-Ed_ZA9-AqaUdV9XUs5WsJsY=	t	2026-02-28 05:39:31.590563+00
d66dc778-6c56-4a1e-a445-32882f46a49a	8dd34eca-8ab9-4228-a9df-b20a4414fb59	F-ypqNSZbMVqPpdGlTDpYKGRCO3fiRfLXBDwYyfs-oU=	t	2026-04-23 02:02:03.963692+00
dc7b302d-eb2a-4516-856b-207ad2f8ee31	8dd34eca-8ab9-4228-a9df-b20a4414fb59	wlVqVTPTCQCdQpHvcmPnmncGRFCarhvFsDin9S_Xygc=	t	2026-05-06 06:27:26.62845+00
\.


--
-- Data for Name: api_usage_analytics; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.api_usage_analytics (id, user_id, api_key_id, request_id, model_requested, model_used, provider, input_tokens, output_tokens, total_tokens, input_cost, output_cost, total_cost, input_price_per_token, output_price_per_token, status, error_message, response_time_ms, created_at, updated_at, currency) FROM stdin;
76a7ac6b-a462-4f8d-b7c9-7398c30c71e9	45da1969-60fd-4248-9efa-d061e34e3a1d	\N	chatcmpl-1772161762	gemini-2.5-flash-lite	gemini-2.5-flash-lite	google	2	9	11	0.00001660	0.00029880	0.00031540	0.00000010	0.00000040	success	\N	0	2026-02-27 03:09:22.822201+00	2026-02-27 03:09:22.822201+00	INR
1569dfb0-625e-42cc-8090-9580845d9fe4	45da1969-60fd-4248-9efa-d061e34e3a1d	\N	chatcmpl-1772161775	gemini-2.5-flash-lite	gemini-2.5-flash-lite	google	6	16	22	0.00004980	0.00053120	0.00058100	0.00000010	0.00000040	success	\N	0	2026-02-27 03:09:35.814848+00	2026-02-27 03:09:35.814848+00	INR
47cd0e40-247f-4f44-b16e-a89200f9a5c6	45da1969-60fd-4248-9efa-d061e34e3a1d	\N	chatcmpl-1772161817	gemini-2.5-flash-lite	gemini-2.5-flash-lite	google	28	246	274	0.00023240	0.00816720	0.00839960	0.00000010	0.00000040	success	\N	0	2026-02-27 03:10:17.247556+00	2026-02-27 03:10:17.247556+00	INR
c76e309a-ac3f-4f88-a12c-68949123fd31	45da1969-60fd-4248-9efa-d061e34e3a1d	bc94f290-d8e0-4fcc-87d5-8c0e92a6aed2	chatcmpl-1772162697	gemini-2.5-flash-lite	gemini-2.5-flash-lite	google	2	10	12	0.00001660	0.00033200	0.00034860	0.00000010	0.00000040	success	\N	0	2026-02-27 03:24:57.550046+00	2026-02-27 03:24:57.550046+00	INR
1864c88e-e609-47aa-b9a2-668ce49013b6	45da1969-60fd-4248-9efa-d061e34e3a1d	bc94f290-d8e0-4fcc-87d5-8c0e92a6aed2	chatcmpl-1772162750	gemini-2.5-flash-lite	gemini-2.5-flash-lite	google	2	10	12	0.00001660	0.00033200	0.00034860	0.00000010	0.00000040	success	\N	0	2026-02-27 03:25:50.369729+00	2026-02-27 03:25:50.369729+00	INR
3f324b31-6e61-419d-8d90-392daf330197	45da1969-60fd-4248-9efa-d061e34e3a1d	\N	chatcmpl-1772243979	gemini-2.5-flash-lite	gemini-2.5-flash-lite	google	4	84	88	0.00003320	0.00278880	0.00282200	0.00000010	0.00000040	success	\N	0	2026-02-28 01:59:39.637343+00	2026-02-28 01:59:39.637343+00	INR
8567c412-0e88-48b5-bf5e-62d1be0de44b	8dd34eca-8ab9-4228-a9df-b20a4414fb59	\N	chatcmpl-1776911647	gemini-3-flash-preview	gemini-3-flash-preview	google	5	408	914	0.00012450	0.08466000	0.08478450	0.00000030	0.00000250	success	\N	0	2026-04-23 02:34:07.250084+00	2026-04-23 02:34:07.250084+00	INR
528b8c77-a1f0-4852-a2c7-bd634726c8f4	8dd34eca-8ab9-4228-a9df-b20a4414fb59	\N	chatcmpl-1777267064	gemini-3-flash-preview	gemini-3-flash-preview	google	2	9	68	0.00004980	0.00186750	0.00191730	0.00000030	0.00000250	success	\N	0	2026-04-27 05:17:44.972507+00	2026-04-27 05:17:44.972507+00	INR
fde9871d-083e-47cb-bc3b-992cad12a002	8dd34eca-8ab9-4228-a9df-b20a4414fb59	d66dc778-6c56-4a1e-a445-32882f46a49a	chatcmpl-1777534770	gemini-2.5-flash-lite	gemini-2.5-flash-lite	google	5	21	26	0.00004150	0.00069720	0.00073870	0.00000010	0.00000040	success	\N	2157	2026-04-30 07:39:30.235362+00	2026-04-30 07:39:30.235362+00	INR
e850c53e-b3f0-4c3f-b082-e18142631d26	8dd34eca-8ab9-4228-a9df-b20a4414fb59	d66dc778-6c56-4a1e-a445-32882f46a49a	chatcmpl-1778039764	gemini-2.5-flash-lite	gemini-2.5-flash-lite	google	4	942	946	0.00003320	0.03127440	0.03130760	0.00000010	0.00000040	success	\N	4566	2026-05-06 03:56:04.428401+00	2026-05-06 03:56:04.428401+00	INR
7f85a1ad-1659-4f31-a662-cb77b5e63a5c	8dd34eca-8ab9-4228-a9df-b20a4414fb59	d66dc778-6c56-4a1e-a445-32882f46a49a	chatcmpl-1778063819	gemini-2.5-flash-lite	gemini-2.5-flash-lite	google	5	12	17	0.00004150	0.00039840	0.00043990	0.00000010	0.00000040	success	\N	1679	2026-05-06 10:36:59.829726+00	2026-05-06 10:36:59.829726+00	INR
bf6ee34a-f7d4-43be-8027-a8f5b5dffd24	8dd34eca-8ab9-4228-a9df-b20a4414fb59	d66dc778-6c56-4a1e-a445-32882f46a49a	chatcmpl-1778063919	gemini-2.5-flash-lite	gemini-2.5-flash-lite	google	22	147	169	0.00018260	0.00488040	0.00506300	0.00000010	0.00000040	success	\N	4415	2026-05-06 10:38:39.695703+00	2026-05-06 10:38:39.695703+00	INR
489bdcd2-101f-4979-bd59-55c297033829	8dd34eca-8ab9-4228-a9df-b20a4414fb59	d66dc778-6c56-4a1e-a445-32882f46a49a	chatcmpl-1778063932	gemini-2.5-flash-lite	gemini-2.5-flash-lite	google	22	199	221	0.00018260	0.00660680	0.00678940	0.00000010	0.00000040	success	\N	2286	2026-05-06 10:38:52.831334+00	2026-05-06 10:38:52.831334+00	INR
7ddd1566-5b7e-4699-87dd-0c8ba0dd394c	8dd34eca-8ab9-4228-a9df-b20a4414fb59	d66dc778-6c56-4a1e-a445-32882f46a49a	chatcmpl-1778136250	gemini-2.5-flash-lite	gemini-2.5-flash-lite	google	5	1007	1012	0.00004150	0.03343240	0.03347390	0.00000010	0.00000040	success	\N	7943	2026-05-07 06:44:10.695609+00	2026-05-07 06:44:10.695609+00	INR
\.


--
-- Data for Name: api_usage_logs; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.api_usage_logs (id, user_id, api_key_id, model_id, input_tokens, output_tokens, cost, created_at, model, provider, currency) FROM stdin;
72559b83-1ec0-43c3-96ca-5e1d61013ea1	8dd34eca-8ab9-4228-a9df-b20a4414fb59	d66dc778-6c56-4a1e-a445-32882f46a49a	\N	5	21	0.0007	2026-04-30 07:39:30.221801+00	gemini-2.5-flash-lite	google	INR
f68fbc61-de4c-4836-b8a7-a7142b7e14f8	8dd34eca-8ab9-4228-a9df-b20a4414fb59	d66dc778-6c56-4a1e-a445-32882f46a49a	\N	4	942	0.0313	2026-05-06 03:56:04.414557+00	gemini-2.5-flash-lite	google	INR
b31aeb0a-a6bc-4266-81ab-07e2a8729c63	8dd34eca-8ab9-4228-a9df-b20a4414fb59	d66dc778-6c56-4a1e-a445-32882f46a49a	\N	5	12	0.0004	2026-05-06 10:36:59.815392+00	gemini-2.5-flash-lite	google	INR
bfa11558-e8ae-412b-871b-a2186c75bd20	8dd34eca-8ab9-4228-a9df-b20a4414fb59	d66dc778-6c56-4a1e-a445-32882f46a49a	\N	22	147	0.0051	2026-05-06 10:38:39.694109+00	gemini-2.5-flash-lite	google	INR
6dd7ce82-f461-47fd-8647-181af4b931c7	8dd34eca-8ab9-4228-a9df-b20a4414fb59	d66dc778-6c56-4a1e-a445-32882f46a49a	\N	22	199	0.0068	2026-05-06 10:38:52.829777+00	gemini-2.5-flash-lite	google	INR
96ccccc2-ecc5-4dab-a51e-2f220fa852bd	8dd34eca-8ab9-4228-a9df-b20a4414fb59	d66dc778-6c56-4a1e-a445-32882f46a49a	\N	5	1007	0.0335	2026-05-07 06:44:10.679287+00	gemini-2.5-flash-lite	google	INR
\.


--
-- Data for Name: chat_messages; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.chat_messages (id, chat_id, role, content, token_count, cost, created_at) FROM stdin;
e83fcf22-e84e-4c10-971c-3a7c72b8e698	2e00a59b-a76d-47fa-bf0b-5384a61f76aa	user	Hi	2	0.000017	2026-02-27 03:09:22.820925+00
06075915-64e2-4d3d-bb36-bc93b82564d5	2e00a59b-a76d-47fa-bf0b-5384a61f76aa	assistant	Hello! How can I help you today?	9	0.000299	2026-02-27 03:09:22.821344+00
8cff9121-ac88-4e50-83ff-0e6fa4437856	8d38af36-5b90-4be6-acee-cf2f39d4c600	user	hi	6	0.000050	2026-02-27 03:09:35.813434+00
1d69dfd3-9b62-4705-bd3e-f9b7ca9818eb	8d38af36-5b90-4be6-acee-cf2f39d4c600	assistant	Hihihi yourself! What's tickling your funny bone today? 😄	16	0.000531	2026-02-27 03:09:35.813966+00
9d8d9587-1adc-4756-a913-40cfad64a3b4	32228c5c-4a66-4250-a666-67716c057057	user	who is tanu	28	0.000232	2026-02-27 03:10:17.246073+00
5fb525dc-953f-4584-87b7-d4d3730e7771	32228c5c-4a66-4250-a666-67716c057057	assistant	"Tanu" is a name that can be used for both **girls and boys**, though it is more commonly a **female given name** in many cultures.\n\nThe meaning and origin of the name can vary:\n\n*   **In Indian cultures (particularly Hindi):** Tanu is a feminine name. It often means:\n    *   **Body**\n    *   **The body of the sun**\n    *   **An infant or child**\n    *   It can also be a diminutive or affectionate form of longer names.\n\n*   **In other cultures:** While less common, it can appear as a given name or surname in other parts of the world.\n\n**To give you a more specific answer, I might need more context.** For example:\n\n*   **Are you asking about a specific person named Tanu?** If so, do you know anything else about them (like their profession, where they are from, etc.)?\n*   **Did you hear the name in a particular movie, book, or TV show?**\n*   **Are you interested in the general meaning of the name?**\n\nLet me know if you have any more information!	246	0.008167	2026-02-27 03:10:17.246514+00
fc358b2a-e208-4f09-9b5c-958099f524f7	b5253ad8-7462-460c-bd36-aa5f6c7a437c	user	Hi	2	0.000017	2026-02-27 03:24:57.546246+00
86c8c345-6919-4fcc-abb3-916fe5ef302c	b5253ad8-7462-460c-bd36-aa5f6c7a437c	assistant	Hi there! How can I help you today?	10	0.000332	2026-02-27 03:24:57.548919+00
51cfef57-9c8f-40bb-b16e-051f39597666	3975b615-a1a9-4318-a577-d50c53eb6bac	user	Hi	2	0.000017	2026-02-27 03:25:50.368531+00
a3813148-5f89-4946-9c5d-0d47b5bd364c	3975b615-a1a9-4318-a577-d50c53eb6bac	assistant	Hi there! How can I help you today?	10	0.000332	2026-02-27 03:25:50.368961+00
a69cde99-0051-420d-8318-d36b14aa89f0	b38438b3-2e02-47b4-a844-93a64c563ee3	user	ji	4	0.000033	2026-02-28 01:59:39.633363+00
ad1d8548-a516-4b1e-bafb-adcf47acd20d	b38438b3-2e02-47b4-a844-93a64c563ee3	assistant	"Hiji" is a word in the **Indonesian language**.\n\nIt means **"one"**.\n\nFor example:\n* **Satu** is also used for "one" in Indonesian, and both are commonly understood.\n* You might hear it in counting or when referring to a single item.\n\nIs there anything specific you'd like to know about "hiji" or its usage?	84	0.002789	2026-02-28 01:59:39.634827+00
d7e0a939-0a78-49d1-bebf-f9e1cb004ace	0e154374-b9d4-41fe-8350-979ece7e17f9	user	who is tanu	5	0.000125	2026-04-23 02:34:07.242181+00
96988150-fad8-44fd-b481-115ffb1f87fa	0e154374-b9d4-41fe-8350-979ece7e17f9	assistant	"Tanu" is a common name and nickname, so who she is depends on the context. Here are the most likely people you might be referring to:\n\n### 1. **Film Character (Tanu Weds Manu)**\nThe most famous "Tanu" in popular culture is **Tanuja "Tanu" Trivedi**, the lead character in the hit Bollywood films *Tanu Weds Manu* (2011) and *Tanu Weds Manu Returns* (2015). She is played by actress **Kangana Ranaut**. In the movies, she is depicted as a rebellious, bold, and free-spirited woman from Kanpur.\n\n### 2. **Television Character (Kumkum Bhagya)**\nIn the long-running Indian soap opera *Kumkum Bhagya*, **Tanushree "Tanu" Mehta** is a major character (usually an antagonist). She has been played by actresses like Madhurima Tuli and Leena Jumani.\n\n### 3. **Famous Real-Life People**\n*   **Dr. Tanu Jain:** A very popular former IAS officer and motivational speaker/educator who is well-known on social media (especially YouTube/Instagram) for her UPSC interview guidance.\n*   **Tanushree Dutta:** A Bollywood actress and former Miss India Universe who is often referred to as Tanu.\n*   **Tanu Rawat:** A popular social media influencer and dancer known for her videos on TikTok and Instagram.\n\n### 4. **Meaning of the Name**\nIn Sanskrit and Hindi, "Tanu" (तनु) generally means **"slim," "delicate," or "body."** It is often used as a shortened version of names like Tanushree, Tanuja, Tanvi, or Tanmayi.\n\n**If you are asking about someone specific—like a person in the news, a book character, or a friend—please provide a little more context!**	408	0.084660	2026-04-23 02:34:07.245068+00
95c67c84-c184-4a67-8c6e-a678c406f8cf	7881a3b4-0e68-434b-a159-969b72984cfc	user	hi	2	0.000050	2026-04-27 05:17:44.965909+00
3e583931-7c24-4770-a52c-ed1b8c9a848e	7881a3b4-0e68-434b-a159-969b72984cfc	assistant	Hello! How can I help you today?	9	0.001868	2026-04-27 05:17:44.969396+00
de659a46-783d-4ff2-8a81-1574f614d1a4	5f83dc4e-3bc1-4bb3-a16b-92c73e0e5a34	user	Hello, world!	5	0.000042	2026-04-30 07:39:30.258822+00
89ea32b9-3355-4ec4-924b-2c9c8257bfd9	5f83dc4e-3bc1-4bb3-a16b-92c73e0e5a34	assistant	Hello to you too! It's nice to hear from you. How can I help you today?	21	0.000697	2026-04-30 07:39:30.262897+00
f040d32e-6094-4be0-b6a3-a048ad726f76	b3a0490b-70a0-4d57-b2e6-24a202fb8a8d	user	What is postgres	4	0.000033	2026-05-06 03:56:04.458821+00
17bc1f32-33b7-4c26-92a9-47418d570b21	b3a0490b-70a0-4d57-b2e6-24a202fb8a8d	assistant	**PostgreSQL** (often shortened to **Postgres**) is a powerful, open-source **relational database management system (RDBMS)**. It's renowned for its robustness, extensibility, and adherence to SQL standards.\n\nHere's a breakdown of what makes Postgres stand out:\n\n**Key Characteristics and Features:**\n\n*   **Relational Database Management System (RDBMS):** Like other RDBMS, Postgres stores data in a structured way using tables, columns, and rows. It enforces relationships between these tables through keys, ensuring data integrity.\n*   **Open-Source and Free:** Postgres is released under a liberal open-source license, meaning it's free to use, distribute, and modify without licensing fees. This makes it a popular choice for both individuals and organizations.\n*   **Highly Extensible:** This is one of Postgres's most defining features. You can extend its functionality by:\n    *   **User-Defined Functions (UDFs):** Write functions in various languages (like SQL, PL/pgSQL, Python, C, etc.) to perform custom operations.\n    *   **Data Types:** Create your own custom data types to better represent your application's data.\n    *   **Operators and Aggregates:** Define new operators and aggregate functions for specialized queries.\n    *   **Extensions:** Install pre-built extensions that add significant new features (e.g., PostGIS for geographic data, timescaledb for time-series data).\n*   **Standards Compliance:** Postgres is known for its strong adherence to SQL standards, making it easier to migrate applications from other SQL databases and write portable SQL queries.\n*   **Advanced Data Types:** Beyond standard SQL types, Postgres supports a wide range of advanced data types, including:\n    *   **JSON/JSONB:** Native support for storing and querying JSON data, making it a great choice for handling semi-structured data.\n    *   **Arrays:** Store collections of values within a single column.\n    *   **Geometric Types:** For storing and manipulating spatial data.\n    *   **Network Address Types:** For IP addresses and network ranges.\n    *   **UUID:** Universally Unique Identifiers.\n*   **ACID Compliance:** Postgres ensures data integrity and reliability through ACID properties:\n    *   **Atomicity:** Transactions are all-or-nothing. Either all operations within a transaction succeed, or none of them do.\n    *   **Consistency:** Transactions bring the database from one valid state to another.\n    *   **Isolation:** Concurrent transactions don't interfere with each other, ensuring that each transaction sees a consistent snapshot of the data.\n    *   **Durability:** Once a transaction is committed, its changes are permanent, even in the event of system failures.\n*   **Concurrency Control:** Postgres uses Multiversion Concurrency Control (MVCC) to handle multiple users accessing the database simultaneously without blocking each other.\n*   **Reliability and Stability:** It's known for its stability and ability to handle large and complex workloads, making it a trusted choice for mission-critical applications.\n*   **Replication:** Postgres offers various replication methods (streaming replication, logical replication) for high availability, disaster recovery, and read scaling.\n*   **Performance:** While performance is highly dependent on configuration and workload, Postgres is generally considered a high-performance database that can be optimized for various use cases.\n*   **Security:** It provides robust security features, including user authentication, role-based access control, SSL encryption, and auditing.\n\n**When is Postgres a good choice?**\n\nPostgres is a versatile database suitable for a wide range of applications, including:\n\n*   **Web Applications:** Powering the backend for dynamic websites and web services.\n*   **Data Warehousing and Analytics:** Handling large datasets for reporting and business intelligence.\n*   **Geographic Information Systems (GIS):** With the PostGIS extension, it's a leading choice for spatial data.\n*   **Scientific Research:** Storing and analyzing complex datasets.\n*   **Financial Applications:** Requiring high data integrity and transactional consistency.\n*   **Enterprise-level applications:** Where scalability, reliability, and advanced features are crucial.\n*   **Applications that need to store and query JSON data efficiently.**\n\n**In summary,** PostgreSQL is a feature-rich, reliable, and highly extensible open-source relational database that is a strong contender for many different types of projects, from small startups to large enterprises.	942	0.031274	2026-05-06 03:56:04.462809+00
13f40c42-2fb3-443f-bf62-cee9932ec39a	da7337d2-19b1-41ab-9c1c-2d1fbe327111	user	Hello, world!	5	0.000042	2026-05-06 10:36:59.851493+00
d7763cf2-6984-471a-8e05-6c6556780f0c	da7337d2-19b1-41ab-9c1c-2d1fbe327111	assistant	Hello to you too! How can I help you today?	12	0.000398	2026-05-06 10:36:59.855135+00
87216391-cef2-4449-a809-8fdfcd2d68a2	f85a2a04-e124-4410-9815-8d4d74ab64d0	user	Say Hariom and Tanu are both nerdy and their vibe match too much with lots of emojis bro!	22	0.000183	2026-05-06 10:38:39.699138+00
c6470c58-fc3c-4872-999b-20c78cc867a2	f85a2a04-e124-4410-9815-8d4d74ab64d0	assistant	OMG, Hariom and Tanu are totally giving me major nerdy vibes! 🤓✨ Their energy is just ✨chef's kiss✨! They're like, the ultimate power couple of intellectual pursuits and adorable chemistry! 💯👫💘 I can just picture them bonding over obscure trivia, deep philosophical debates, and probably a shared love for a good spreadsheet. 📊📚 Their synergy is unreal, bro! It's like they were made for each other, completing each other's sentences and finishing each other's theorems. ⚛️💡 So much mutual nerdy adoration happening, it's truly heartwarming to witness! 🥰🚀 I'm living for their matching energy, it's the cutest! 💖🌟	147	0.004880	2026-05-06 10:38:39.700123+00
9c675f5a-1ac6-4d5e-8fc8-bec89b4b6111	d1e241ca-b96a-4b9c-b2cf-f15e628ff52c	user	Say Hariom and Tanu are both nerdy and their vibe match too much with lots of emojis bro!	22	0.000183	2026-05-06 10:38:52.835626+00
f3816ce8-4b9d-4a08-bcf6-93cb5fdb02a5	d1e241ca-b96a-4b9c-b2cf-f15e628ff52c	assistant	Hariom and Tanu are such nerds! 🤓📚 Their vibe is *chef's kiss* 💯🔥! They're practically twinsies when it comes to their passions. 👯‍♀️✨ I swear, they finish each other's sentences about obscure trivia! 🧠💡 And don't even get me started on their shared love for all things sci-fi and fantasy. 🚀🐉 They're always geeking out together about the latest releases or rewatching their favorite classics. 🎬🍿 Their conversations are a beautiful symphony of memes, inside jokes, and existential ponderings. 😂🌌 It's like they're speaking a secret nerdy language that only they understand. 🗣️🤫 And the way they support each other's intellectual pursuits? Goals! 🥰🏆 Seriously, their connection is so strong, it's like they were coded to be besties. 💻💖 They're the epitome of a perfect nerdy duo! 🌟😎	199	0.006607	2026-05-06 10:38:52.837139+00
5f1667e5-6353-4607-a47a-5fc6e910c15f	995013a7-4ce2-440d-922e-1ad94b0d7e0f	user	what is postgres!	5	0.000042	2026-05-07 06:44:10.728786+00
f4a0d7ac-731c-4c47-9d4d-ff39439c7c6d	995013a7-4ce2-440d-922e-1ad94b0d7e0f	assistant	**PostgreSQL (pronounced POST-grees-que-el, often shortened to Postgres)** is a powerful, open-source **relational database management system (RDBMS)**. It's known for its robustness, extensibility, and compliance with SQL standards.\n\nHere's a breakdown of what that means and why it's so popular:\n\n**Key Concepts:**\n\n* **Relational Database Management System (RDBMS):** This is the core. A relational database stores data in tables, which are structured collections of rows and columns. Relationships between different tables are defined using keys. This structure allows for efficient data organization, retrieval, and manipulation.\n* **Open-Source:** This means the source code for PostgreSQL is freely available. Anyone can download, use, modify, and distribute it without paying licensing fees. This fosters a large and active community that contributes to its development and support.\n* **SQL (Structured Query Language):** This is the standard language used to interact with relational databases. You use SQL commands to create, read, update, and delete (CRUD) data, as well as to define the structure of your database. PostgreSQL is highly compliant with SQL standards.\n\n**Why is PostgreSQL so popular and what are its key features?**\n\n* **Robustness and Reliability:** PostgreSQL is designed for data integrity and stability. It uses a sophisticated transaction system (MVCC - Multi-Version Concurrency Control) that ensures data remains consistent even with many users accessing and modifying it simultaneously. This makes it suitable for mission-critical applications.\n* **Extensibility:** One of PostgreSQL's biggest strengths is its extensibility. You can:\n    * **Define your own data types:** Create custom data types to represent specific kinds of data.\n    * **Define your own functions:** Write functions in various languages (like PL/pgSQL, Python, Perl, Tcl, etc.) that can be executed within the database.\n    * **Create custom operators:** Define new operators that can be used in your queries.\n    * **Develop custom index methods:** Optimize performance for specific data access patterns.\n* **Advanced Features:** PostgreSQL offers a rich set of advanced features that go beyond basic RDBMS capabilities:\n    * **JSON and JSONB support:** Efficiently store and query semi-structured data like JSON documents.\n    * **Full-text search:** Powerful capabilities for searching through large amounts of text.\n    * **Geospatial data support (PostGIS):** A very popular extension that provides advanced spatial data types and functions, making it a leading choice for geographic information systems (GIS).\n    * **Advanced indexing:** Supports various indexing methods (B-tree, Hash, GiST, SP-GiST, GIN, BRIN) for optimizing query performance.\n    * **Foreign Data Wrappers (FDWs):** Allows PostgreSQL to access data from other databases or external data sources as if they were local tables.\n    * **Window functions:** Powerful SQL functions for performing calculations across sets of table rows that are related to the current row.\n    * **Common Table Expressions (CTEs):** Make complex queries more readable and manageable.\n    * **Stored Procedures:** Allow you to group multiple SQL statements into a single callable unit.\n* **ACID Compliance:** PostgreSQL adheres to ACID properties (Atomicity, Consistency, Isolation, Durability), which are essential for ensuring reliable transaction processing.\n* **Scalability:** PostgreSQL can scale to handle large amounts of data and high user loads. It supports replication, partitioning, and other techniques for performance and availability.\n* **Strong Community Support:** Being open-source, PostgreSQL benefits from a large and active community of developers and users. This means readily available documentation, forums, mailing lists, and a constant stream of improvements and bug fixes.\n* **Security:** PostgreSQL offers robust security features, including user authentication, access control, and encryption.\n\n**When is PostgreSQL a good choice?**\n\nPostgreSQL is a versatile database that is suitable for a wide range of applications, including:\n\n* **Web applications:** Its reliability and scalability make it ideal for powering dynamic websites and web services.\n* **E-commerce platforms:** Handles large volumes of transactions and product data with integrity.\n* **Data warehousing and analytics:** Its advanced features and query optimization capabilities make it suitable for complex analytical tasks.\n* **Geographic Information Systems (GIS):** With PostGIS, it's a go-to solution for spatial data management.\n* **Scientific and research applications:** Its extensibility and ability to handle complex data types are valuable for these domains.\n* **Enterprise-level applications:** Its robustness and advanced features meet the demands of large organizations.\n\n**In summary, PostgreSQL is a high-performance, feature-rich, and reliable open-source relational database that is a favorite among developers and organizations looking for a powerful and flexible database solution.**	1007	0.033432	2026-05-07 06:44:10.734875+00
\.


--
-- Data for Name: chats; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.chats (id, user_id, title, model, provider, created_at, updated_at) FROM stdin;
2e00a59b-a76d-47fa-bf0b-5384a61f76aa	45da1969-60fd-4248-9efa-d061e34e3a1d	Hi	gemini-2.5-flash-lite	google	2026-02-27 03:09:22.820476+00	2026-02-27 03:09:22.82176+00
8d38af36-5b90-4be6-acee-cf2f39d4c600	45da1969-60fd-4248-9efa-d061e34e3a1d	hi	gemini-2.5-flash-lite	google	2026-02-27 03:09:35.812955+00	2026-02-27 03:09:35.81429+00
32228c5c-4a66-4250-a666-67716c057057	45da1969-60fd-4248-9efa-d061e34e3a1d	who is tanu	gemini-2.5-flash-lite	google	2026-02-27 03:10:17.245647+00	2026-02-27 03:10:17.247068+00
8491c78c-ffe7-4840-a59c-a465b429e03d	45da1969-60fd-4248-9efa-d061e34e3a1d	New Chat	gemini-2.0-flash	google	2026-02-27 03:13:52.559166+00	2026-02-27 03:13:52.559166+00
b5253ad8-7462-460c-bd36-aa5f6c7a437c	45da1969-60fd-4248-9efa-d061e34e3a1d	Hi	gemini-2.5-flash-lite	google	2026-02-27 03:24:57.543805+00	2026-02-27 03:24:57.549366+00
3975b615-a1a9-4318-a577-d50c53eb6bac	45da1969-60fd-4248-9efa-d061e34e3a1d	Hi	gemini-2.5-flash-lite	google	2026-02-27 03:25:50.368104+00	2026-02-27 03:25:50.369326+00
b38438b3-2e02-47b4-a844-93a64c563ee3	45da1969-60fd-4248-9efa-d061e34e3a1d	ji	gemini-2.5-flash-lite	google	2026-02-28 01:59:39.631932+00	2026-02-28 01:59:39.636576+00
1a07751d-d88a-4390-892b-9920da45b33a	8dd34eca-8ab9-4228-a9df-b20a4414fb59	New Chat	gemini-2.0-flash	google	2026-04-23 02:11:04.291837+00	2026-04-23 02:11:04.291837+00
eb0733e1-640d-4256-9181-1e4ccc042905	8dd34eca-8ab9-4228-a9df-b20a4414fb59	New Chat	gemini-2.0-flash	google	2026-04-23 02:12:16.688927+00	2026-04-23 02:12:16.688927+00
957edb4c-b067-4dbd-8a8a-ae3966eda882	8dd34eca-8ab9-4228-a9df-b20a4414fb59	New Chat	gemini-2.0-flash	google	2026-04-23 02:20:23.53818+00	2026-04-23 02:20:23.53818+00
5e145b78-7a43-4081-b58d-ebfa81bd0298	8dd34eca-8ab9-4228-a9df-b20a4414fb59	New Chat	gemini-2.0-flash	google	2026-04-23 02:27:41.654688+00	2026-04-23 02:27:41.654688+00
0e154374-b9d4-41fe-8350-979ece7e17f9	8dd34eca-8ab9-4228-a9df-b20a4414fb59	who is tanu	gemini-3-flash-preview	google	2026-04-23 02:34:07.238778+00	2026-04-23 02:34:07.247702+00
7afaa3cd-0699-4589-a9aa-5de9186ca212	8dd34eca-8ab9-4228-a9df-b20a4414fb59	New Chat	gemini-3-flash-preview	google	2026-04-23 02:34:14.373711+00	2026-04-23 02:34:14.373711+00
7881a3b4-0e68-434b-a159-969b72984cfc	8dd34eca-8ab9-4228-a9df-b20a4414fb59	hi	gemini-3-flash-preview	google	2026-04-27 05:17:44.956891+00	2026-04-27 05:17:44.970742+00
5f83dc4e-3bc1-4bb3-a16b-92c73e0e5a34	8dd34eca-8ab9-4228-a9df-b20a4414fb59	Hello, world!	gemini-2.5-flash-lite	google	2026-04-30 07:39:30.25577+00	2026-04-30 07:39:30.263799+00
b3a0490b-70a0-4d57-b2e6-24a202fb8a8d	8dd34eca-8ab9-4228-a9df-b20a4414fb59	What is postgres	gemini-2.5-flash-lite	google	2026-05-06 03:56:04.456003+00	2026-05-06 03:56:04.466253+00
da7337d2-19b1-41ab-9c1c-2d1fbe327111	8dd34eca-8ab9-4228-a9df-b20a4414fb59	Hello, world!	gemini-2.5-flash-lite	google	2026-05-06 10:36:59.848949+00	2026-05-06 10:36:59.855797+00
f85a2a04-e124-4410-9815-8d4d74ab64d0	8dd34eca-8ab9-4228-a9df-b20a4414fb59	Say Hariom and Tanu are both nerdy and their vibe match too 	gemini-2.5-flash-lite	google	2026-05-06 10:38:39.698141+00	2026-05-06 10:38:39.701066+00
d1e241ca-b96a-4b9c-b2cf-f15e628ff52c	8dd34eca-8ab9-4228-a9df-b20a4414fb59	Say Hariom and Tanu are both nerdy and their vibe match too 	gemini-2.5-flash-lite	google	2026-05-06 10:38:52.834411+00	2026-05-06 10:38:52.838351+00
995013a7-4ce2-440d-922e-1ad94b0d7e0f	8dd34eca-8ab9-4228-a9df-b20a4414fb59	what is postgres!	gemini-2.5-flash-lite	google	2026-05-07 06:44:10.724766+00	2026-05-07 06:44:10.737538+00
\.


--
-- Data for Name: credits; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.credits (id, user_id, total_credits, used_credits, created_at, updated_at) FROM stdin;
889e6746-2c08-43a2-b515-221bfb107b61	45da1969-60fd-4248-9efa-d061e34e3a1d	100.00	0.01	2026-02-26 09:08:56.201146+00	2026-02-28 01:59:39.62506+00
30e805c8-89c4-4190-8b7c-04d839d9a277	8dd34eca-8ab9-4228-a9df-b20a4414fb59	500.00	0.16	2026-04-23 02:00:51.380001+00	2026-05-07 06:44:10.722888+00
\.


--
-- Data for Name: daily_usage_summary; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.daily_usage_summary (id, user_id, date, total_requests, total_input_tokens, total_output_tokens, total_tokens, total_cost, models_used, providers_used, created_at, updated_at) FROM stdin;
fb6e0671-837c-4e32-86f8-a2a4285a0222	45da1969-60fd-4248-9efa-d061e34e3a1d	2026-02-27	5	40	291	331	0.00999320	{"gemini-2.5-flash-lite": 5}	{"google": 5}	2026-02-27 03:09:22.819934+00	2026-02-27 03:25:50.367402+00
092b729c-babb-4114-a150-3a554bbc62a2	45da1969-60fd-4248-9efa-d061e34e3a1d	2026-02-28	1	4	84	88	0.00282200	{"gemini-2.5-flash-lite": 1}	{"google": 1}	2026-02-28 01:59:39.624827+00	2026-02-28 01:59:39.624827+00
5a3f4a11-8ef0-46a4-9831-c7585559c60a	8dd34eca-8ab9-4228-a9df-b20a4414fb59	2026-04-23	1	5	408	914	0.08478450	{"gemini-3-flash-preview": 1}	{"google": 1}	2026-04-23 02:34:07.233613+00	2026-04-23 02:34:07.233613+00
50cb144c-09d7-4f83-8b44-1492aea61a70	8dd34eca-8ab9-4228-a9df-b20a4414fb59	2026-04-27	1	2	9	68	0.00191730	{"gemini-3-flash-preview": 1}	{"google": 1}	2026-04-27 05:17:44.941967+00	2026-04-27 05:17:44.941967+00
bc8123a0-ee1d-4f82-8e27-4409c4d2b5c7	8dd34eca-8ab9-4228-a9df-b20a4414fb59	2026-04-30	1	5	21	26	0.00073870	{"gemini-2.5-flash-lite": 1}	{"google": 1}	2026-04-30 07:39:30.220959+00	2026-04-30 07:39:30.220959+00
070ad158-0914-4267-a1ef-54221d011576	8dd34eca-8ab9-4228-a9df-b20a4414fb59	2026-05-06	4	53	1300	1353	0.04359990	{"gemini-2.5-flash-lite": 4}	{"google": 4}	2026-05-06 03:56:04.413631+00	2026-05-06 10:38:52.829239+00
f98edac7-28a0-4a70-8f6f-bcfcb897d2f4	8dd34eca-8ab9-4228-a9df-b20a4414fb59	2026-05-07	1	5	1007	1012	0.03347390	{"gemini-2.5-flash-lite": 1}	{"google": 1}	2026-05-07 06:44:10.678302+00	2026-05-07 06:44:10.678302+00
\.


--
-- Data for Name: models; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.models (id, name, provider, input_token_price, output_token_price, active, created_at) FROM stdin;
\.


--
-- Data for Name: payments; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.payments (id, user_id, razorpay_order_id, razorpay_payment_id, amount, status, created_at, updated_at) FROM stdin;
5d3d1c8e-a47e-49c7-bbbf-17f80cfb3267	45da1969-60fd-4248-9efa-d061e34e3a1d	order_SKiYpHbVUlwzBE		100.00	pending	2026-02-26 09:02:34.676283+00	2026-02-26 09:02:34.676283+00
7e693204-3530-4b2e-a5b0-4ffb270b94fe	45da1969-60fd-4248-9efa-d061e34e3a1d	order_SKiacYR8AMArr5		100.00	pending	2026-02-26 09:04:16.688078+00	2026-02-26 09:04:16.688078+00
d22b85ec-5423-4155-a181-65b67932c8a4	45da1969-60fd-4248-9efa-d061e34e3a1d	order_SKidVX5212S1Zn	pay_SKieMhUmU8InPx	100.00	completed	2026-02-26 09:07:00.929149+00	2026-02-26 09:08:56.203365+00
01b9ea03-01c5-42ad-92fb-785f626de50b	45da1969-60fd-4248-9efa-d061e34e3a1d	order_SL1ZA2avE2wB4A		500.00	pending	2026-02-27 03:38:04.457654+00	2026-02-27 03:38:04.457654+00
bfac5a82-ef1a-4373-981c-037cac083f04	45da1969-60fd-4248-9efa-d061e34e3a1d	order_SLSBSW98Bhq5fR		100.00	pending	2026-02-28 05:40:21.66565+00	2026-02-28 05:40:21.66565+00
c1c9e3fe-3247-451a-895c-1dbdff2762a0	8dd34eca-8ab9-4228-a9df-b20a4414fb59	order_SgkvCmmfNTxcDc		100.00	pending	2026-04-23 01:39:11.357426+00	2026-04-23 01:39:11.357426+00
d7718aee-24ff-4b6b-b279-7f205ee3c7f3	8dd34eca-8ab9-4228-a9df-b20a4414fb59	order_Sgl4rb6LW8E29t		100.00	pending	2026-04-23 01:48:19.933577+00	2026-04-23 01:48:19.933577+00
66b02ed1-cf6d-4b58-9ef9-32a983a7f80c	8dd34eca-8ab9-4228-a9df-b20a4414fb59	order_Sgl93Dt58ikK6m		1000.00	pending	2026-04-23 01:52:17.861027+00	2026-04-23 01:52:17.861027+00
73b683be-7bae-4486-a878-646b3b289660	8dd34eca-8ab9-4228-a9df-b20a4414fb59	order_SglHDRM9ud75Du	pay_SglHhoHIlrPoZ5	500.00	completed	2026-04-23 02:00:01.591233+00	2026-04-23 02:00:51.380188+00
65521f33-a938-494b-87f2-88f1aafcf233	8dd34eca-8ab9-4228-a9df-b20a4414fb59	order_SgqXz2iZF40dyb		100.00	pending	2026-04-23 07:09:21.882618+00	2026-04-23 07:09:21.882618+00
e324c1dc-21a8-4e6e-bcaf-aded30e5c4dc	8dd34eca-8ab9-4228-a9df-b20a4414fb59	order_SqKmZzBuj1ZnmF		100.00	pending	2026-05-17 06:35:00.541386+00	2026-05-17 06:35:00.541386+00
fa5cef56-a6e5-468e-9e4e-ab42144d3aeb	8dd34eca-8ab9-4228-a9df-b20a4414fb59	order_SqKoRMeBuShCOY		100.00	pending	2026-05-17 06:36:46.215852+00	2026-05-17 06:36:46.215852+00
5f1d11e4-af4f-4537-a4c5-997a5e89f69e	8dd34eca-8ab9-4228-a9df-b20a4414fb59	order_SqKt83SZKJD5N8		100.00	pending	2026-05-17 06:41:12.528618+00	2026-05-17 06:41:12.528618+00
b3769091-cb3b-43ac-bef3-09c681c4671a	8dd34eca-8ab9-4228-a9df-b20a4414fb59	order_SqKxotydwd2dyF		100.00	pending	2026-05-17 06:45:38.995468+00	2026-05-17 06:45:38.995468+00
\.


--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.schema_migrations (version) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, name, email, password_hash, email_verified, created_at, updated_at) FROM stdin;
6a3f0d2a-73e5-4aa6-95b9-be619f76ef72	Admin	admin@clearrouter.local	$2a$14$kobkg3PQDyE0EDuyMjsDWu6IPCgnRHDUqXFcr5fqIvbsYG5qXoPYa	t	2026-02-26 07:35:54.61762+00	2026-02-26 07:35:54.61762+00
ecb5f1a4-7271-458b-aa9d-92f7117a4609	Hariom OPPP	hariomvirkhare02@gmail.com	$2a$14$CdGAmzcwrcTkx7zzGEl4seSIytcQLk4NnLcqvPXv4bXFUUpcj738i	t	2026-02-26 07:39:03.815432+00	2026-02-26 07:39:03.815432+00
3e419448-bf28-4683-bed3-f34dbda297c5	devmmm	h.ariomvirkhare02@gmail.com	$2a$14$6NJaxwmDpxf1cfANCJ2jw.KMwge1mTFzIR05xT3zdyORW8bpX9YWa	t	2026-02-26 08:05:23.624194+00	2026-02-26 08:05:23.624194+00
ef1f50d8-b1d7-4841-bfdd-13a977ddac87	devmmm12	ddodexpert111@gmail.com	$2a$14$NAOq1mQUYhdLayuGH73KI.nQE/5lf0/WbruSUPu/uD/SQMoJQpez2	t	2026-02-26 08:17:06.674171+00	2026-02-26 08:17:06.674171+00
1424b834-0a06-4afe-ad20-aa9fc85d4a69	John Doe	kobiyo6385@pazuric.com	$2a$14$AvL7BNHYzjMBzclC1dM/XuM6M.TGY9u5ILVZ1wHIcrNbiF2RwklD2	t	2026-02-26 08:23:31.330945+00	2026-02-26 08:23:31.330945+00
3ceaa174-c74a-41d1-9633-a7e465030171	John Doe	dd.odexpert111@gmail.com	$2a$14$4gAlkc6vGfGnym1.J09Fs./n3DMhPSyh1NaTGkKm0np7/GiGX2816	t	2026-02-26 08:25:29.797411+00	2026-02-26 08:25:29.797411+00
45da1969-60fd-4248-9efa-d061e34e3a1d	Test User	ddodexper.t111@gmail.com	$2a$14$kuSpky6C2524dtgnRct3LuoXOTxqYvhk2D.k9BAztPI6O5RbN6BBS	t	2026-02-26 08:33:56.050038+00	2026-02-26 08:34:15.243486+00
32ba5174-7e7e-4845-a31d-85c7df0cf8e5	Hariom	ddod.expert111@gmail.com	$2a$14$d3oKn74O4wkzW2wejyyqB.Z2rx92XVDudnSyG1h/vEDGNYkqczNYS	t	2026-02-28 01:49:28.804658+00	2026-02-28 01:49:57.699343+00
b9c558b1-59ab-4fa1-9a22-dcc3bc2f6802	Temp user 	h.ariomvi.rkhare02@gmail.com	$2a$14$P5rU70eIBtyk4h2KXG4X.OUmRsOZVM9qvVwxImPuxG3OvMTIbJMXS	t	2026-04-16 06:42:49.136086+00	2026-04-16 06:42:49.136086+00
d81dc2b3-651d-4823-86e5-c5ce5d0dfd4e	Tanu 	Tanu@hariom.com	$2a$14$XFSoIo9f8hMwHuJVllgNgese/rMEmbAExdRZiHC6/ZgVz2ceqQhaW	t	2026-04-23 01:36:23.290767+00	2026-04-23 01:36:23.290767+00
8dd34eca-8ab9-4228-a9df-b20a4414fb59	Tanu	hariomv.irkhare02@gmail.com	$2a$14$crxNUg8tXamlvmY2torwH.5MzlHTIZnzfhdLUBnhAowr2CJdfmfN2	t	2026-04-23 01:37:15.172716+00	2026-04-23 01:38:50.691952+00
\.


--
-- Name: api_keys api_keys_api_key_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.api_keys
    ADD CONSTRAINT api_keys_api_key_key UNIQUE (api_key);


--
-- Name: api_keys api_keys_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.api_keys
    ADD CONSTRAINT api_keys_pkey PRIMARY KEY (id);


--
-- Name: api_usage_analytics api_usage_analytics_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.api_usage_analytics
    ADD CONSTRAINT api_usage_analytics_pkey PRIMARY KEY (id);


--
-- Name: api_usage_logs api_usage_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.api_usage_logs
    ADD CONSTRAINT api_usage_logs_pkey PRIMARY KEY (id);


--
-- Name: chat_messages chat_messages_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chat_messages
    ADD CONSTRAINT chat_messages_pkey PRIMARY KEY (id);


--
-- Name: chats chats_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chats
    ADD CONSTRAINT chats_pkey PRIMARY KEY (id);


--
-- Name: credits credits_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.credits
    ADD CONSTRAINT credits_pkey PRIMARY KEY (id);


--
-- Name: daily_usage_summary daily_usage_summary_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.daily_usage_summary
    ADD CONSTRAINT daily_usage_summary_pkey PRIMARY KEY (id);


--
-- Name: daily_usage_summary daily_usage_summary_user_id_date_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.daily_usage_summary
    ADD CONSTRAINT daily_usage_summary_user_id_date_key UNIQUE (user_id, date);


--
-- Name: models models_name_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.models
    ADD CONSTRAINT models_name_key UNIQUE (name);


--
-- Name: models models_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.models
    ADD CONSTRAINT models_pkey PRIMARY KEY (id);


--
-- Name: payments payments_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: credits unique_user_id; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.credits
    ADD CONSTRAINT unique_user_id UNIQUE (user_id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_api_keys_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_api_keys_user_id ON public.api_keys USING btree (user_id);


--
-- Name: idx_api_usage_analytics_created_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_api_usage_analytics_created_at ON public.api_usage_analytics USING btree (created_at);


--
-- Name: idx_api_usage_analytics_model; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_api_usage_analytics_model ON public.api_usage_analytics USING btree (model_requested);


--
-- Name: idx_api_usage_analytics_model_requested; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_api_usage_analytics_model_requested ON public.api_usage_analytics USING btree (model_requested);


--
-- Name: idx_api_usage_analytics_provider; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_api_usage_analytics_provider ON public.api_usage_analytics USING btree (provider);


--
-- Name: idx_api_usage_analytics_status; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_api_usage_analytics_status ON public.api_usage_analytics USING btree (status);


--
-- Name: idx_api_usage_analytics_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_api_usage_analytics_user_id ON public.api_usage_analytics USING btree (user_id);


--
-- Name: idx_api_usage_logs_model_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_api_usage_logs_model_id ON public.api_usage_logs USING btree (model_id);


--
-- Name: idx_api_usage_logs_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_api_usage_logs_user_id ON public.api_usage_logs USING btree (user_id);


--
-- Name: idx_chat_messages_chat_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_chat_messages_chat_id ON public.chat_messages USING btree (chat_id);


--
-- Name: idx_chat_messages_created_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_chat_messages_created_at ON public.chat_messages USING btree (created_at);


--
-- Name: idx_chats_created_at; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_chats_created_at ON public.chats USING btree (created_at DESC);


--
-- Name: idx_chats_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_chats_user_id ON public.chats USING btree (user_id);


--
-- Name: idx_credits_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_credits_user_id ON public.credits USING btree (user_id);


--
-- Name: idx_payments_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_payments_user_id ON public.payments USING btree (user_id);


--
-- Name: api_usage_analytics trigger_update_daily_usage_summary; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER trigger_update_daily_usage_summary AFTER INSERT ON public.api_usage_analytics FOR EACH ROW EXECUTE FUNCTION public.update_daily_usage_summary();


--
-- Name: api_keys api_keys_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.api_keys
    ADD CONSTRAINT api_keys_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: api_usage_analytics api_usage_analytics_api_key_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.api_usage_analytics
    ADD CONSTRAINT api_usage_analytics_api_key_id_fkey FOREIGN KEY (api_key_id) REFERENCES public.api_keys(id) ON DELETE CASCADE;


--
-- Name: api_usage_analytics api_usage_analytics_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.api_usage_analytics
    ADD CONSTRAINT api_usage_analytics_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: api_usage_logs api_usage_logs_api_key_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.api_usage_logs
    ADD CONSTRAINT api_usage_logs_api_key_id_fkey FOREIGN KEY (api_key_id) REFERENCES public.api_keys(id) ON DELETE CASCADE;


--
-- Name: api_usage_logs api_usage_logs_model_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.api_usage_logs
    ADD CONSTRAINT api_usage_logs_model_id_fkey FOREIGN KEY (model_id) REFERENCES public.models(id) ON DELETE CASCADE;


--
-- Name: api_usage_logs api_usage_logs_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.api_usage_logs
    ADD CONSTRAINT api_usage_logs_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: chat_messages chat_messages_chat_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chat_messages
    ADD CONSTRAINT chat_messages_chat_id_fkey FOREIGN KEY (chat_id) REFERENCES public.chats(id) ON DELETE CASCADE;


--
-- Name: chats chats_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.chats
    ADD CONSTRAINT chats_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: credits credits_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.credits
    ADD CONSTRAINT credits_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: daily_usage_summary daily_usage_summary_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.daily_usage_summary
    ADD CONSTRAINT daily_usage_summary_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: payments payments_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

\unrestrict E4R2YLYmrla6Qk6es0jyyCZ61ekwCxf7SJkuHx2KnCtLBsrM45THugvOawGlaAc

