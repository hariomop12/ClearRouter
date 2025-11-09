--
-- PostgreSQL database dump
--

\restrict f1Gtu0u75mFEYcnz3eMhIm0ZNgMiAeCmiwbLa4ThWTrRU7y5dV7tgBXdSIJtHhy

-- Dumped from database version 15.14
-- Dumped by pg_dump version 16.10 (Ubuntu 16.10-0ubuntu0.24.04.1)

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
-- Name: update_daily_usage_summary(); Type: FUNCTION; Schema: public; Owner: clearrouter
--

CREATE FUNCTION public.update_daily_usage_summary() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    -- Insert or update daily summary
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


ALTER FUNCTION public.update_daily_usage_summary() OWNER TO clearrouter;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: api_keys; Type: TABLE; Schema: public; Owner: clearrouter
--

CREATE TABLE public.api_keys (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid,
    api_key text NOT NULL,
    active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public.api_keys OWNER TO clearrouter;

--
-- Name: api_usage_analytics; Type: TABLE; Schema: public; Owner: clearrouter
--

CREATE TABLE public.api_usage_analytics (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    api_key_id uuid,
    request_id character varying(255) NOT NULL,
    model_requested character varying(255) NOT NULL,
    model_used character varying(255) NOT NULL,
    provider character varying(100) NOT NULL,
    input_tokens integer DEFAULT 0,
    output_tokens integer DEFAULT 0,
    total_tokens integer DEFAULT 0,
    input_cost numeric(15,8) DEFAULT 0,
    output_cost numeric(15,8) DEFAULT 0,
    total_cost numeric(15,8) DEFAULT 0,
    currency character varying(10) DEFAULT 'INR'::character varying,
    input_price_per_token numeric(15,8),
    output_price_per_token numeric(15,8),
    status character varying(50) DEFAULT 'success'::character varying,
    response_time_ms integer,
    created_at timestamp with time zone DEFAULT now(),
    error_message text,
    updated_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public.api_usage_analytics OWNER TO clearrouter;

--
-- Name: api_usage_logs; Type: TABLE; Schema: public; Owner: clearrouter
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


ALTER TABLE public.api_usage_logs OWNER TO clearrouter;

--
-- Name: chat_messages; Type: TABLE; Schema: public; Owner: clearrouter
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


ALTER TABLE public.chat_messages OWNER TO clearrouter;

--
-- Name: chats; Type: TABLE; Schema: public; Owner: clearrouter
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


ALTER TABLE public.chats OWNER TO clearrouter;

--
-- Name: credits; Type: TABLE; Schema: public; Owner: clearrouter
--

CREATE TABLE public.credits (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid,
    total_credits numeric(12,2) DEFAULT 0,
    used_credits numeric(12,2) DEFAULT 0,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public.credits OWNER TO clearrouter;

--
-- Name: daily_usage_summary; Type: TABLE; Schema: public; Owner: clearrouter
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


ALTER TABLE public.daily_usage_summary OWNER TO clearrouter;

--
-- Name: model_pricing_history; Type: TABLE; Schema: public; Owner: clearrouter
--

CREATE TABLE public.model_pricing_history (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    model_id character varying(255) NOT NULL,
    provider character varying(100) NOT NULL,
    input_price numeric(15,8) NOT NULL,
    output_price numeric(15,8) NOT NULL,
    effective_from timestamp with time zone DEFAULT now() NOT NULL,
    effective_until timestamp with time zone,
    is_active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT now()
);


ALTER TABLE public.model_pricing_history OWNER TO clearrouter;

--
-- Name: models; Type: TABLE; Schema: public; Owner: clearrouter
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


ALTER TABLE public.models OWNER TO clearrouter;

--
-- Name: payments; Type: TABLE; Schema: public; Owner: clearrouter
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


ALTER TABLE public.payments OWNER TO clearrouter;

--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: clearrouter
--

CREATE TABLE public.schema_migrations (
    version character varying NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO clearrouter;

--
-- Name: users; Type: TABLE; Schema: public; Owner: clearrouter
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


ALTER TABLE public.users OWNER TO clearrouter;

--
-- Data for Name: api_keys; Type: TABLE DATA; Schema: public; Owner: clearrouter
--

COPY public.api_keys (id, user_id, api_key, active, created_at) FROM stdin;
bb6d03ae-97da-4ee9-9973-91708c009059	8ade694a-d9c6-439c-bd5c-e74d9380b4f3	65DvRObWumOSkPxw3OtQ0CtpF2fLGEMNviG_gvI_H2E=	t	2025-10-01 04:25:12.400922+00
7c4846ac-e747-4867-9fa5-de4d3b24f12a	8ade694a-d9c6-439c-bd5c-e74d9380b4f3	10mnt7COLdBzQ9ZFeKZYI0y7uxDAu52PzMv8CJl2zj4=	t	2025-10-01 04:25:27.614629+00
bd3f059d-2048-4839-ad58-0d6cab4a0f04	8ade694a-d9c6-439c-bd5c-e74d9380b4f3	tLhW3i1lyqv2l8Sl6xrVUk-3dy8w6k8oRsMmRafZ8E4=	t	2025-10-01 05:27:10.790245+00
c4c2a363-da77-4278-a78e-0accf98beef0	8af05435-cafa-4fec-b0d9-5086628c3e59	LAaibjXiOGIGR6ljMgP35Rw-xbkFJ3HoAAZ7giJCPpI=	t	2025-10-01 07:19:53.904527+00
7457ec92-19fa-466a-b9cf-2db635347948	8af05435-cafa-4fec-b0d9-5086628c3e59	F4vHFuD2iYrIdnhym-0z-EU9kJFZ0H3SETxWjXP6R9I=	t	2025-10-01 11:04:12.528978+00
124338c6-c795-49d9-9bb1-35444aa5e1a9	8af05435-cafa-4fec-b0d9-5086628c3e59	Z5g8b1Bg8BElUmZYxQhs7r0mFR1biaPfwMpk1WDQqfk=	t	2025-10-01 11:04:14.287384+00
277b89b6-887b-40fc-8587-bbca37fb4496	9990d63c-b429-42b3-bfe2-b339b4658f05	vi_P0xR7wk8Gq-I14s4GtFamfBNV4iGNFXwZnMNUAdI=	t	2025-10-02 07:38:31.016564+00
c787c423-2299-4099-850a-9c4065dec3b4	8af05435-cafa-4fec-b0d9-5086628c3e59	lnFpJDxLHsh8bt2eMs_wSi_oJAxgtMl1z4OPh_hr1pI=	t	2025-10-03 06:58:46.103582+00
f6e54e13-5625-4b45-bc5c-5fe452019378	9990d63c-b429-42b3-bfe2-b339b4658f05	cpISzP090n0VpUuvPSKJXMEEFLeNwQLAfOPsjJ-HJ90=	t	2025-11-09 06:53:30.842064+00
\.


--
-- Data for Name: api_usage_analytics; Type: TABLE DATA; Schema: public; Owner: clearrouter
--

COPY public.api_usage_analytics (id, user_id, api_key_id, request_id, model_requested, model_used, provider, input_tokens, output_tokens, total_tokens, input_cost, output_cost, total_cost, currency, input_price_per_token, output_price_per_token, status, response_time_ms, created_at, error_message, updated_at) FROM stdin;
cf4f5143-690f-45ad-8da3-7c0194d4c4a7	8ade694a-d9c6-439c-bd5c-e74d9380b4f3	7c4846ac-e747-4867-9fa5-de4d3b24f12a	chatcmpl-1759293995	gemini-2.5-flash-lite	gemini-2.5-flash-lite	google	4	780	784	0.00003320	0.02589600	0.02592920	INR	0.00000010	0.00000040	success	0	2025-10-01 04:46:35.620574+00	\N	2025-10-01 04:46:35.620574+00
c40385cb-8e00-4940-8270-8d86275260da	bce02ed2-0883-4f02-97f3-5f4281930a6f	\N	chatcmpl-1759324765	gemini-2.0-flash	gemini-2.0-flash	google	12	6	18	0.00009960	0.00019920	0.00029880	INR	0.00000010	0.00000040	success	0	2025-10-01 13:19:25.447003+00	\N	2025-10-01 13:19:25.447003+00
93661ac0-1c2d-4da8-a395-105b4b01f9c8	8af05435-cafa-4fec-b0d9-5086628c3e59	\N	chatcmpl-1759324801	gemini-2.0-flash	gemini-2.0-flash	google	2	10	12	0.00001660	0.00033200	0.00034860	INR	0.00000010	0.00000040	success	0	2025-10-01 13:20:01.922415+00	\N	2025-10-01 13:20:01.922415+00
d2030380-2f89-4025-b89e-8757afd9f720	8af05435-cafa-4fec-b0d9-5086628c3e59	\N	chatcmpl-1759324876	gemma-3-4b-it	gemma-3-4b-it	google	36	0	36	0.00022410	0.00000000	0.00022410	INR	0.00000008	0.00000030	success	0	2025-10-01 13:21:16.980603+00	\N	2025-10-01 13:21:16.980603+00
db85a1de-97d1-44a8-95e3-691e6881878c	8af05435-cafa-4fec-b0d9-5086628c3e59	\N	chatcmpl-1759325043	gemini-2.0-flash	gemini-2.0-flash	google	8	86	94	0.00006640	0.00285520	0.00292160	INR	0.00000010	0.00000040	success	0	2025-10-01 13:24:03.019084+00	\N	2025-10-01 13:24:03.019084+00
0142f7a9-4cfa-4bc8-8d64-69b26f3081fe	8af05435-cafa-4fec-b0d9-5086628c3e59	\N	chatcmpl-1759325058	gemini-2.0-flash	gemini-2.0-flash	google	1	11	12	0.00000830	0.00036520	0.00037350	INR	0.00000010	0.00000040	success	0	2025-10-01 13:24:18.955512+00	\N	2025-10-01 13:24:18.955512+00
d6d7f400-dd8e-4bf7-bdd6-6e88758b7f64	bce02ed2-0883-4f02-97f3-5f4281930a6f	\N	chatcmpl-1759325368	gemini-2.0-flash	gemini-2.0-flash	google	12	38	50	0.00009960	0.00126160	0.00136120	INR	0.00000010	0.00000040	success	0	2025-10-01 13:29:28.317946+00	\N	2025-10-01 13:29:28.317946+00
14e3dcb4-7bdb-44bb-94dd-ba4cf12c840d	8af05435-cafa-4fec-b0d9-5086628c3e59	\N	chatcmpl-1759325439	gemini-2.0-flash	gemini-2.0-flash	google	3	136	139	0.00002490	0.00451520	0.00454010	INR	0.00000010	0.00000040	success	0	2025-10-01 13:30:39.257688+00	\N	2025-10-01 13:30:39.257688+00
87f9a2cd-87c5-448d-9bb3-43c084a2575a	bce02ed2-0883-4f02-97f3-5f4281930a6f	\N	chatcmpl-1759325658	gemini-2.0-flash	gemini-2.0-flash	google	62	9	71	0.00051460	0.00029880	0.00081340	INR	0.00000010	0.00000040	success	0	2025-10-01 13:34:18.948689+00	\N	2025-10-01 13:34:18.948689+00
d21f4ce3-a326-41aa-8dd8-5b15f9a88d75	bce02ed2-0883-4f02-97f3-5f4281930a6f	\N	chatcmpl-1759325755	gemini-2.0-flash	gemini-2.0-flash	google	13	57	70	0.00010790	0.00189240	0.00200030	INR	0.00000010	0.00000040	success	0	2025-10-01 13:35:55.424821+00	\N	2025-10-01 13:35:55.424821+00
8931ab2b-2149-415f-8e6c-7bf2b5fc2499	8af05435-cafa-4fec-b0d9-5086628c3e59	\N	chatcmpl-1759325782	gemini-2.0-flash	gemini-2.0-flash	google	1	10	11	0.00000830	0.00033200	0.00034030	INR	0.00000010	0.00000040	success	0	2025-10-01 13:36:22.145237+00	\N	2025-10-01 13:36:22.145237+00
c8a197d1-7182-4386-89fd-0f13b6858e21	8af05435-cafa-4fec-b0d9-5086628c3e59	\N	chatcmpl-1759325797	gemini-2.0-flash	gemini-2.0-flash	google	14	861	875	0.00011620	0.02858520	0.02870140	INR	0.00000010	0.00000040	success	0	2025-10-01 13:36:37.181148+00	\N	2025-10-01 13:36:37.181148+00
6ba46fe9-33b4-4ffe-9d7b-6e9f5cbb6187	8af05435-cafa-4fec-b0d9-5086628c3e59	\N	chatcmpl-1759325864	gemini-2.0-flash	gemini-2.0-flash	google	867	837	1704	0.00719610	0.02778840	0.03498450	INR	0.00000010	0.00000040	success	0	2025-10-01 13:37:44.258728+00	\N	2025-10-01 13:37:44.258728+00
b9a3457e-543a-4c4e-9e2a-774f7c208159	8af05435-cafa-4fec-b0d9-5086628c3e59	\N	chatcmpl-1759478932	gemini-2.0-flash	gemini-2.0-flash	google	867	1209	2076	0.00719610	0.04013880	0.04733490	INR	0.00000010	0.00000040	success	0	2025-10-03 08:08:52.866391+00	\N	2025-10-03 08:08:52.866391+00
93c3bb21-acf4-4157-b1a5-2c88b2860f9a	8af05435-cafa-4fec-b0d9-5086628c3e59	\N	chatcmpl-1759478948	gemini-2.0-flash	gemini-2.0-flash	google	2	11	13	0.00001660	0.00036520	0.00038180	INR	0.00000010	0.00000040	success	0	2025-10-03 08:09:08.536262+00	\N	2025-10-03 08:09:08.536262+00
f222bdc8-ec40-4846-8e99-6b9a06a1df89	8af05435-cafa-4fec-b0d9-5086628c3e59	\N	chatcmpl-1759478955	gemini-2.0-flash	gemini-2.0-flash	google	16	12	28	0.00013280	0.00039840	0.00053120	INR	0.00000010	0.00000040	success	0	2025-10-03 08:09:15.845787+00	\N	2025-10-03 08:09:15.845787+00
\.


--
-- Data for Name: api_usage_logs; Type: TABLE DATA; Schema: public; Owner: clearrouter
--

COPY public.api_usage_logs (id, user_id, api_key_id, model_id, input_tokens, output_tokens, cost, created_at, model, provider, currency) FROM stdin;
fb5def42-3d0d-4afa-bffa-ae288ab8801d	8ade694a-d9c6-439c-bd5c-e74d9380b4f3	7c4846ac-e747-4867-9fa5-de4d3b24f12a	\N	4	780	0.0259	2025-10-01 04:46:35.619924+00	gemini-2.5-flash-lite	google	INR
\.


--
-- Data for Name: chat_messages; Type: TABLE DATA; Schema: public; Owner: clearrouter
--

COPY public.chat_messages (id, chat_id, role, content, token_count, cost, created_at) FROM stdin;
d7a7b2d9-d30d-46e1-b17e-b14eebde06d3	2dc3e871-bbd6-4a16-ad07-70a4512b1783	user	What is Friendship?	4	0.000033	2025-10-01 04:46:35.629737+00
91c1a0aa-9577-47eb-bf5d-be9cb54d2ce2	2dc3e871-bbd6-4a16-ad07-70a4512b1783	assistant	Friendship is a multifaceted and deeply human connection characterized by mutual affection, trust, loyalty, and support between individuals. It's a voluntary bond that goes beyond mere acquaintance, often described as a relationship built on shared experiences, understanding, and a genuine care for one another's well-being.\n\nHere are some key elements that define friendship:\n\n**Core Components:**\n\n*   **Affection and Liking:** At its heart, friendship involves a genuine liking and fondness for another person. You enjoy their company, find them pleasant to be around, and have positive feelings towards them.\n*   **Trust:** This is a cornerstone of friendship. You feel safe sharing your thoughts, feelings, vulnerabilities, and secrets with your friends, believing they won't judge you or betray your confidence. They, in turn, extend the same trust to you.\n*   **Loyalty and Support:** Friends are there for each other through thick and thin. They offer emotional support during difficult times, celebrate successes, and stand by you even when things are tough. Loyalty means having each other's back and prioritizing the relationship.\n*   **Reciprocity:** Friendship is a two-way street. It involves give and take, where both individuals contribute to the relationship, offering support, understanding, and affection in return.\n*   **Shared Experiences and Interests:** While not always essential, shared hobbies, passions, values, or life experiences can be powerful bonding agents. These commonalities provide opportunities for connection and deepen understanding.\n*   **Understanding and Empathy:** Friends often have a good grasp of each other's personalities, quirks, and emotional states. They can empathize with what you're going through, offering comfort and understanding without necessarily needing to have experienced the same thing themselves.\n*   **Acceptance:** True friends accept you for who you are, flaws and all. They don't try to change you or expect you to be someone you're not.\n*   **Companionship and Enjoyment:** Friends provide a sense of belonging and companionship. Spending time with them is often enjoyable, fun, and can alleviate feelings of loneliness.\n\n**Why is Friendship Important?**\n\nFriendship plays a vital role in our lives for several reasons:\n\n*   **Emotional Well-being:** Friends provide a crucial support system, helping us cope with stress, anxiety, and sadness. They offer a sense of belonging and reduce feelings of isolation.\n*   **Personal Growth:** Friends can challenge our perspectives, introduce us to new ideas, and encourage us to step outside our comfort zones, fostering personal growth and self-discovery.\n*   **Mental Health:** Strong social connections are linked to better mental health outcomes, including lower rates of depression and anxiety.\n*   **Happiness and Fulfillment:** The joy and laughter shared with friends contribute significantly to our overall happiness and life satisfaction.\n*   **Practical Support:** Friends can offer practical help, whether it's lending a hand with a move, offering advice, or simply being there to listen.\n\n**Types of Friendships:**\n\nFriendships can vary in intensity and depth. You might have:\n\n*   **Close Friends/Best Friends:** These are the people you share your deepest thoughts and feelings with, who are often considered like family.\n*   **Good Friends:** These are people you have a strong, consistent bond with, who you enjoy spending time with and can rely on.\n*   **Casual Friends/Acquaintances:** These are people you enjoy interacting with, perhaps from work, school, or social activities, but the bond may be less intense.\n\nIn essence, friendship is a precious and invaluable relationship that enriches our lives immeasurably. It's a testament to our social nature and our innate need for connection and belonging.	780	0.025896	2025-10-01 04:46:35.632466+00
ea41fdcd-8d00-4aad-a8b7-0bc9ac3c4053	eaa6e440-17c9-446e-aa16-156070771323	user	Hello! Let me test the new chat functionality with delete buttons.	13	0.000108	2025-10-01 13:35:55.423232+00
55fd030f-f220-4487-8b48-859f915fd457	eaa6e440-17c9-446e-aa16-156070771323	assistant	Okay! I understand. Feel free to send me a message or a series of messages that you want to test with the delete buttons. Let me know if there's anything specific you'd like me to say to help you test! I'm ready when you are.\n	57	0.001892	2025-10-01 13:35:55.423801+00
\.


--
-- Data for Name: chats; Type: TABLE DATA; Schema: public; Owner: clearrouter
--

COPY public.chats (id, user_id, title, model, provider, created_at, updated_at) FROM stdin;
2dc3e871-bbd6-4a16-ad07-70a4512b1783	8ade694a-d9c6-439c-bd5c-e74d9380b4f3	What is Friendship?	gemini-2.5-flash-lite	google	2025-10-01 04:46:35.626482+00	2025-10-01 04:46:35.634915+00
eaa6e440-17c9-446e-aa16-156070771323	bce02ed2-0883-4f02-97f3-5f4281930a6f	Hello! Let me test the new chat functionality with delete bu	gemini-2.0-flash	google	2025-10-01 13:35:55.42234+00	2025-10-01 13:35:55.42429+00
09142a3c-3522-45f7-8a92-245e39288c92	8af05435-cafa-4fec-b0d9-5086628c3e59	New Chat	gemini-2.0-flash	google	2025-10-03 08:10:59.714252+00	2025-10-03 08:10:59.714252+00
\.


--
-- Data for Name: credits; Type: TABLE DATA; Schema: public; Owner: clearrouter
--

COPY public.credits (id, user_id, total_credits, used_credits, created_at, updated_at) FROM stdin;
ba881618-697b-45e9-b822-7c7a42a50389	8ade694a-d9c6-439c-bd5c-e74d9380b4f3	100.00	0.03	2025-10-01 04:30:26.693382+00	2025-10-01 04:46:35.624525+00
33be8603-5858-44ea-a7be-b24e19984ae9	bce02ed2-0883-4f02-97f3-5f4281930a6f	100.00	0.00	2025-10-01 13:08:10.915271+00	2025-10-01 13:35:55.421696+00
a5835cff-4b11-4110-8d20-35b47ac6c540	8af05435-cafa-4fec-b0d9-5086628c3e59	10000.00	0.11	2025-10-01 11:44:35.600502+00	2025-10-03 08:09:15.83675+00
cd5acb51-914f-4dc5-ab16-0219deddaada	9990d63c-b429-42b3-bfe2-b339b4658f05	100.00	0.00	2025-11-09 06:56:09.553972+00	2025-11-09 06:56:09.553972+00
\.


--
-- Data for Name: daily_usage_summary; Type: TABLE DATA; Schema: public; Owner: clearrouter
--

COPY public.daily_usage_summary (id, user_id, date, total_requests, total_input_tokens, total_output_tokens, total_tokens, total_cost, models_used, providers_used, created_at, updated_at) FROM stdin;
45fa00d3-de51-46ff-ada4-e809c24a2b99	bce02ed2-0883-4f02-97f3-5f4281930a6f	2025-10-01	4	99	110	209	0.00447370	{"gemini-2.0-flash": 4}	{"google": 4}	2025-10-01 13:19:25.444946+00	2025-10-01 13:35:55.421483+00
4180073e-d4e7-4cab-944e-2844d72bc457	8af05435-cafa-4fec-b0d9-5086628c3e59	2025-10-01	8	932	1951	2883	0.07243410	{"gemma-3-4b-it": 1, "gemini-2.0-flash": 7}	{"google": 8}	2025-10-01 13:20:01.916878+00	2025-10-01 13:37:44.254935+00
4417fc2d-5189-4b6c-8c37-9ae5a901da0f	8af05435-cafa-4fec-b0d9-5086628c3e59	2025-10-03	3	885	1232	2117	0.04824790	{"gemini-2.0-flash": 3}	{"google": 3}	2025-10-03 08:08:52.85847+00	2025-10-03 08:09:15.836626+00
\.


--
-- Data for Name: model_pricing_history; Type: TABLE DATA; Schema: public; Owner: clearrouter
--

COPY public.model_pricing_history (id, model_id, provider, input_price, output_price, effective_from, effective_until, is_active, created_at) FROM stdin;
\.


--
-- Data for Name: models; Type: TABLE DATA; Schema: public; Owner: clearrouter
--

COPY public.models (id, name, provider, input_token_price, output_token_price, active, created_at) FROM stdin;
\.


--
-- Data for Name: payments; Type: TABLE DATA; Schema: public; Owner: clearrouter
--

COPY public.payments (id, user_id, razorpay_order_id, razorpay_payment_id, amount, status, created_at, updated_at) FROM stdin;
14ba428c-da3b-43bd-9c65-8d37914fc43c	8ade694a-d9c6-439c-bd5c-e74d9380b4f3	order_RO4n3FcYoZAKGV	rzp_test_RMdWZKqulEGpGi	100.00	completed	2025-10-01 04:26:12.001107+00	2025-10-01 04:30:26.694101+00
9f3aab74-5b27-4143-bef3-862913b32acd	8af05435-cafa-4fec-b0d9-5086628c3e59	order_RO7nsslxPQct8a		10000.00	pending	2025-10-01 07:23:04.089041+00	2025-10-01 07:23:04.089041+00
9a08fe73-4f6e-4f4c-92b0-4685eb498847	8af05435-cafa-4fec-b0d9-5086628c3e59	order_RO7ssLOEBHJm4H		10000.00	pending	2025-10-01 07:27:47.609862+00	2025-10-01 07:27:47.609862+00
af29683c-ccc9-4485-acb1-a6640fa8f0b6	8af05435-cafa-4fec-b0d9-5086628c3e59	order_ROASQa4Mfn45re		10000.00	pending	2025-10-01 09:58:50.266821+00	2025-10-01 09:58:50.266821+00
02070529-0f91-4a85-a0fa-751ced6d61a8	8af05435-cafa-4fec-b0d9-5086628c3e59	order_ROBeKYPRGkVIkJ		10000.00	pending	2025-10-01 11:08:48.028867+00	2025-10-01 11:08:48.028867+00
6832a8f3-1efd-436d-a41c-491eab53a349	8af05435-cafa-4fec-b0d9-5086628c3e59	order_ROBjPWh51t8PFq		10000.00	pending	2025-10-01 11:13:36.470924+00	2025-10-01 11:13:36.470924+00
06f07331-70f2-44bc-81d8-b0d7f5f20f2f	bce02ed2-0883-4f02-97f3-5f4281930a6f	order_ROBv29nhZ6h3TS		100.00	pending	2025-10-01 11:24:36.758725+00	2025-10-01 11:24:36.758725+00
f1bb09b3-8347-4e66-a621-281e3eafa4d2	8af05435-cafa-4fec-b0d9-5086628c3e59	order_ROBxjQDed4nFuY		100.00	pending	2025-10-01 11:27:09.957139+00	2025-10-01 11:27:09.957139+00
505ca9cc-df02-4b13-9a25-feee7e8ad767	8af05435-cafa-4fec-b0d9-5086628c3e59	order_ROBz5J4jel4llX		10000.00	pending	2025-10-01 11:28:26.767603+00	2025-10-01 11:28:26.767603+00
227fb175-6e83-4b5b-b5f5-cc1c02e4858c	8af05435-cafa-4fec-b0d9-5086628c3e59	order_ROC1Quxc4qTW3b		500000.00	pending	2025-10-01 11:30:40.239075+00	2025-10-01 11:30:40.239075+00
00ccb77a-bae4-4abe-8a08-f147c8c0d463	8af05435-cafa-4fec-b0d9-5086628c3e59	order_ROC2m3btJqZ5e4		10000.00	pending	2025-10-01 11:31:56.395429+00	2025-10-01 11:31:56.395429+00
fdc0900e-549e-434b-a397-c5f25e5a0fab	8af05435-cafa-4fec-b0d9-5086628c3e59	order_ROC6PoGo3pXWdo		10000.00	pending	2025-10-01 11:35:23.193783+00	2025-10-01 11:35:23.193783+00
8fcc72ee-6780-46c0-ad42-8f6c20261328	8af05435-cafa-4fec-b0d9-5086628c3e59	order_ROCAeOMnkYAOXG		10000.00	pending	2025-10-01 11:39:23.786248+00	2025-10-01 11:39:23.786248+00
1760b382-dc8c-4132-b048-9748fb6aed52	8af05435-cafa-4fec-b0d9-5086628c3e59	order_ROCF6zonRsHn5t	pay_ROCFPlq10kXujH	10000.00	completed	2025-10-01 11:43:37.180555+00	2025-10-01 11:44:35.600679+00
f2bd4fe4-fc63-4759-ba0c-300431362f94	8af05435-cafa-4fec-b0d9-5086628c3e59	order_ROCGSloVtfgXa1		10000.00	pending	2025-10-01 11:44:53.968084+00	2025-10-01 11:44:53.968084+00
fb46a9f3-8e58-45c6-be10-37e5b7308334	8af05435-cafa-4fec-b0d9-5086628c3e59	order_ROCHhmOkJ9Xywa		100000.00	pending	2025-10-01 11:46:04.516186+00	2025-10-01 11:46:04.516186+00
285f37c8-457b-4eae-bee0-8cb30cdcc41a	9990d63c-b429-42b3-bfe2-b339b4658f05	order_ROWhJJmQsZZdRN		100.00	pending	2025-10-02 07:44:11.09289+00	2025-10-02 07:44:11.09289+00
cd9e416e-a199-486e-a1b2-f5e49b17d43b	8af05435-cafa-4fec-b0d9-5086628c3e59	order_ROuXEMM73EYFjv		100.00	pending	2025-10-03 07:03:17.324037+00	2025-10-03 07:03:17.324037+00
25084914-1b1e-46ee-ae6f-975e9c3872e0	9990d63c-b429-42b3-bfe2-b339b4658f05	order_RdY8ixtqlwwGZo	pay_RdYAuj6U1ZB1ul	100.00	completed	2025-11-09 06:53:51.47514+00	2025-11-09 06:56:09.554057+00
\.


--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: clearrouter
--

COPY public.schema_migrations (version) FROM stdin;
20250927072843
20250927101618
20250927102719
20250927192900
20250928134500
20250929143000
20250929150700
20251001052126
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: clearrouter
--

COPY public.users (id, name, email, password_hash, email_verified, created_at, updated_at) FROM stdin;
bce02ed2-0883-4f02-97f3-5f4281930a6f	Admin	admin@clearrouter.local	$2a$14$s4T0yRf3J4lgSgOO/uWNaON3wAWKYz4Qnzu9oWJ0HEWVbJPxDeM3e	t	2025-10-01 03:47:15.974009+00	2025-10-01 03:47:15.974009+00
8ade694a-d9c6-439c-bd5c-e74d9380b4f3	Test User	Hariomv.irkhar.e02@gmail.com	$2a$14$z5l2gGYjRVO9U35LtPRgsO/K/GXw/heiL92INCEwvPTTHj/W4.HF6	t	2025-10-01 04:24:07.404266+00	2025-10-01 04:24:33.71688+00
ba32c13e-f36c-4488-adce-bd3ca27d05cc	Test User	test@example.com	$2a$14$UupCK347z6Bwyy0cdRAkTuz4WQTPLuJpWGk2C6oy8QIUeKXRMXMU.	f	2025-10-01 10:41:05.995599+00	2025-10-01 10:41:05.995599+00
17a20bf1-808e-4e8e-89a0-add5298dceff	Test User D	tatejep134@etenx.com	$2a$14$6jBxV5GFDwnJ3x3oyWf9n.PSBXeqV43gv2GFyaIPTDDfYzSIQFmqS	f	2025-10-03 05:33:55.038124+00	2025-10-03 05:33:55.038124+00
76691ea5-e42e-47fa-9e69-b89dcaa32aff	Test User D	awsuse131313@gmail.com	$2a$14$6oMA2Cg4wmp8oJ4OzOhBdeWd8i0tMWb8R9snjIkHkJtNU8HiI.lN.	f	2025-10-03 05:34:36.510537+00	2025-10-03 05:34:36.510537+00
8ec951c5-c3d8-4504-aef1-99901f3572a4	Test User D	awsuse1313.13@gmail.com	$2a$14$9UKf/fJzZV0hRwMdt39vaOmIBs0dOQ.OfaweXX5rqeoVqXvwUuzoy	f	2025-10-03 05:37:08.301787+00	2025-10-03 05:37:08.301787+00
858dbb1d-8bf0-4127-aa82-d701377fd22c	Test User D	awsus.e1313.13@gmail.com	$2a$14$JNrUyZKVik.hq96VgTLncOk0T9pwSn5f0z5M46Ud6rxbk9RIPVOvG	f	2025-10-03 05:38:46.32707+00	2025-10-03 05:38:46.32707+00
c0ef9d61-d4c2-44a5-82dc-29e870d2a0cf	Test User D	8tulxajb7q@zudpck.com	$2a$14$Iom1EV2Unml0RX6gi5uo/ePzlI18xyu8Z2eJrQVMj/BsnwqlKendi	f	2025-10-03 05:42:09.292567+00	2025-10-03 05:42:09.292567+00
0e06d3bb-cec2-4839-b07b-da6dd59ba287	Test User D	snake5798@mailshan.com	$2a$14$Gg4eIguPvqt1uA9P9Jfbn.uZLxU8sf6paVmlly0Z6nP6Cx6HPQl2C	f	2025-10-03 05:44:45.294493+00	2025-10-03 05:44:45.294493+00
f18a1569-15be-40d6-bb82-4f9e5e363dad	Test User E	testuser@example.com	$2a$14$F44eUjBd0KlgEUH/FMxGPuKiPiBD.rOb7anyqVlcIRkggHzItioka	f	2025-10-03 05:46:41.180989+00	2025-10-03 05:46:41.180989+00
22659e0d-8263-464d-955b-e1e209fc4b57	Test User F	testuser2@example.com	$2a$14$ETnlgQ.oXMWyUCRfoLubteo30.Tw9j4F1D2Ud5Ff9mV7Bj4/8d3E6	f	2025-10-03 05:47:29.18906+00	2025-10-03 05:47:29.18906+00
3d1c2df9-4104-42c8-aaf2-756ac78852bb	Test User G	testuser3@example.com	$2a$14$l3qJH80Gp19ACO.ohhATs.jFkzwy.1o2/pqodVP05o7CZ2AaoSYg6	f	2025-10-03 05:48:41.716386+00	2025-10-03 05:48:41.716386+00
d4e48f26-c84e-40a8-8bc8-29ca99fbf2d9	Hariom OP	tat.ejep.134@etenx.com	$2a$14$ET1PKjDArLTLbG/oKv/X9OkbQE5BwjXl6IY.8EkRdhr1G16hD3.yW	f	2025-10-03 06:32:09.662956+00	2025-10-03 06:32:09.662956+00
4e5042c0-12fc-4949-930c-b9af0ea40cae	Updated Test User	testuser4@example.com	$2a$14$uObnu0QBxsAZ4yxiEF7.SuKjW5aA/CwvTSwAfoE38h83n1lITGfKW	t	2025-10-03 05:49:32.338229+00	2025-10-03 05:50:01.086124+00
795042ac-a190-4b2c-8495-f1875c4bf8f6	Test User D	ddodexpert.111@gmail.com	$2a$14$XFmqvXBi4Dzdz6NWbI3NIuAksfBeaRFp0tEdjcjmvOFcV9Pe4cJJi	f	2025-10-03 05:51:35.659347+00	2025-10-03 05:51:35.659347+00
d2556a4b-8710-4d57-8183-91c655f71fa0	Test User D	dd.odex.pert.111@gmail.com	$2a$14$mORZ3oS6Wuk2TufMqie2ROBRgDIf7CRykE5WldU2fqU6zJ6SUbZCm	f	2025-10-03 05:53:17.902902+00	2025-10-03 05:53:17.902902+00
d17623f2-c665-4074-8d49-daf14f6f9f79	Test User D	dd.od.ex.pert.111@gmail.com	$2a$14$AfOaQzXtwdN9mAgy8iVlHO1iGL/UFuTAI8KW9rZRaL17H7HCp6Mo.	f	2025-10-03 05:54:21.696853+00	2025-10-03 05:54:21.696853+00
572c30c6-56a9-4c8e-90b2-94866a34db52	Test User Final	testuser.final@example.com	$2a$14$YoPqatj0IkSB0S2adgy1O.PzSoBsgZVYhVPzaFnOUFoLnmltKWXAa	f	2025-10-03 05:55:55.050333+00	2025-10-03 05:55:55.050333+00
a31680d6-07ec-4c40-af26-d8ba030b27a4	Test User D	dd.odex.per.t.111@gmail.com	$2a$14$WmXVeai1OSqFjplzrKDz6eEm2tdsfTxLxjg6Z8WSG/IVwKvQYbpBi	f	2025-10-03 05:56:02.37219+00	2025-10-03 05:56:02.37219+00
e1b0fad6-5786-4087-904c-428ec205f9ba	Test User With Env	testuser.env@example.com	$2a$14$/kmiEu0T8.N6r054sXTZKOq7/ypHhgDPR1.158FxPTWPMWPiDoDiS	f	2025-10-03 05:57:27.164355+00	2025-10-03 05:57:27.164355+00
57bb1339-1928-45b7-83aa-396e3a8bc6cf	Test User D	dd.o.dex.per.t.111@gmail.com	$2a$14$VuAgImeAjZwQHzEoKHII1eN34.7Ol4rD0096HmhuE6hLskBmF.zSG	f	2025-10-03 05:57:29.732511+00	2025-10-03 05:57:29.732511+00
9b29da1b-0b8c-4773-bce1-f55766f27b08	Test User D	dd.o.dex.per.t.11.1@gmail.com	$2a$14$4ULnTNBeAeIFQlNac3UFAOsUo4.uP9K2vmXwvFc6RQ4QO18ilED9.	f	2025-10-03 06:04:21.178059+00	2025-10-03 06:04:21.178059+00
004bed9a-93b4-4aa2-9179-c47fc299ed16	Test SMTP Fix	test.smtp@example.com	$2a$14$RvMNiQGRnXK.WFkUMEkcmOPNpjkOQFDGsEjqbBDE889fSpe7DvUz2	f	2025-10-03 06:05:19.801826+00	2025-10-03 06:05:19.801826+00
73aba688-99e5-4b71-88e0-163393e54bf7	Test TLS Fix	test.tls@example.com	$2a$14$I.VBjikBsELN2Y50BDLJmOXfyjh5yPbl.mKzBuJrL9.M0aM/txgg6	f	2025-10-03 06:06:41.867341+00	2025-10-03 06:06:41.867341+00
3e76a8e4-280c-4f59-8b5e-84d4197a17b1	Test User D	dd.o.dex.per.t.1.1.1@gmail.com	$2a$14$oGKVesrXPNlUPjh7NtNR3uAXviMpXNKooh/oXTY.4Mi0gk4zA1bHC	f	2025-10-03 06:10:49.143335+00	2025-10-03 06:10:49.143335+00
4cc3f67f-c049-44c5-bfc5-a58f2387a4e6	Test New Gmail	test.newgmail@example.com	$2a$14$LusX62LceFpwnQwbY2KcUOkErD2.HJfEqPqGkLIUq3iWgsa7wOTyW	f	2025-10-03 06:18:50.620822+00	2025-10-03 06:18:50.620822+00
fd5df59a-e013-4f23-8310-73ea5af23814	Test Final Gmail	test.finalgmail@example.com	$2a$14$0Blj7ZHeMLK0WwSU0w3zcOkcpS/Yz/NDULpYoHbD2bbbMFLMSUnVe	f	2025-10-03 06:21:03.238723+00	2025-10-03 06:21:03.238723+00
326cf4be-cd96-497d-a9a5-a42742ddfc10	Test User D	dd.o.dex.p.er.t.1.1.1@gmail.com	$2a$14$3MRvS/mLGl3h1Y1Phs6b2ea09/uUFOxhj.R6MirNK38YAs1sBZxJS	f	2025-10-03 06:22:39.907332+00	2025-10-03 06:22:39.907332+00
1d197287-9469-408a-ab94-77e23ca49398	Test Fixed Password	test.fixed@example.com	$2a$14$nLP8bthVhVTJnBqfjjMg4OMnqqtSjDlK6bZdNKkvphywrdKJYJIrG	f	2025-10-03 06:22:53.727833+00	2025-10-03 06:22:53.727833+00
569d260e-0115-4773-bf22-72d336cbe881	Test Fresh Build	test.fresh@example.com	$2a$14$.sKAOKmP6aEYtdWljqklBeclVZsNYsciYCBgHtwZiF1GIbOM1dmrm	f	2025-10-03 06:24:21.038969+00	2025-10-03 06:24:21.038969+00
d390b370-8c2c-47c7-b455-dcbb46cd2baa	Test User D	dd.o.dex.p.e.r.t.1.1.1@gmail.com	$2a$14$5h4p7gAxDtGzKgcwR/xaZe8p.0Zo/IZGqIl3KdgBQRvNHvFthVcKq	f	2025-10-03 06:24:25.728696+00	2025-10-03 06:24:25.728696+00
0d5b86cb-de08-4021-8f79-3aa34a182113	New Username	d.d.o.d.ex.p.e.r.t.1.1.1@gmail.com	$2a$14$34bGTQKPz4Qx.82nTRHc7e6/YLYznvabwi92OQpFjBpYcJd1JsH8.	t	2025-10-03 06:30:27.201892+00	2025-10-03 06:31:04.341026+00
8af05435-cafa-4fec-b0d9-5086628c3e59	Hariom OP	hariomvirkhare02@gmail.com	$2a$14$RXuZJBNKsbAqtKs94FA0b.mgtBT1kzElsRlVAmEgFZO9RnMQAr9mK	t	2025-10-01 03:58:53.947582+00	2025-10-03 07:10:52.406777+00
9990d63c-b429-42b3-bfe2-b339b4658f05	zxc	ddodexpert111@gmail.com	$2a$14$gh2k47xGP8PXHP5.mYiMV.SmHtH8qjAB70CGBQSMvrp3utDJ1JIaS	t	2025-10-02 07:37:14.518792+00	2025-11-09 06:56:37.516561+00
\.


--
-- Name: api_keys api_keys_api_key_key; Type: CONSTRAINT; Schema: public; Owner: clearrouter
--

ALTER TABLE ONLY public.api_keys
    ADD CONSTRAINT api_keys_api_key_key UNIQUE (api_key);


--
-- Name: api_keys api_keys_pkey; Type: CONSTRAINT; Schema: public; Owner: clearrouter
--

ALTER TABLE ONLY public.api_keys
    ADD CONSTRAINT api_keys_pkey PRIMARY KEY (id);


--
-- Name: api_usage_analytics api_usage_analytics_pkey; Type: CONSTRAINT; Schema: public; Owner: clearrouter
--

ALTER TABLE ONLY public.api_usage_analytics
    ADD CONSTRAINT api_usage_analytics_pkey PRIMARY KEY (id);


--
-- Name: api_usage_logs api_usage_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: clearrouter
--

ALTER TABLE ONLY public.api_usage_logs
    ADD CONSTRAINT api_usage_logs_pkey PRIMARY KEY (id);


--
-- Name: chat_messages chat_messages_pkey; Type: CONSTRAINT; Schema: public; Owner: clearrouter
--

ALTER TABLE ONLY public.chat_messages
    ADD CONSTRAINT chat_messages_pkey PRIMARY KEY (id);


--
-- Name: chats chats_pkey; Type: CONSTRAINT; Schema: public; Owner: clearrouter
--

ALTER TABLE ONLY public.chats
    ADD CONSTRAINT chats_pkey PRIMARY KEY (id);


--
-- Name: credits credits_pkey; Type: CONSTRAINT; Schema: public; Owner: clearrouter
--

ALTER TABLE ONLY public.credits
    ADD CONSTRAINT credits_pkey PRIMARY KEY (id);


--
-- Name: daily_usage_summary daily_usage_summary_pkey; Type: CONSTRAINT; Schema: public; Owner: clearrouter
--

ALTER TABLE ONLY public.daily_usage_summary
    ADD CONSTRAINT daily_usage_summary_pkey PRIMARY KEY (id);


--
-- Name: daily_usage_summary daily_usage_summary_user_id_date_key; Type: CONSTRAINT; Schema: public; Owner: clearrouter
--

ALTER TABLE ONLY public.daily_usage_summary
    ADD CONSTRAINT daily_usage_summary_user_id_date_key UNIQUE (user_id, date);


--
-- Name: model_pricing_history model_pricing_history_pkey; Type: CONSTRAINT; Schema: public; Owner: clearrouter
--

ALTER TABLE ONLY public.model_pricing_history
    ADD CONSTRAINT model_pricing_history_pkey PRIMARY KEY (id);


--
-- Name: models models_name_key; Type: CONSTRAINT; Schema: public; Owner: clearrouter
--

ALTER TABLE ONLY public.models
    ADD CONSTRAINT models_name_key UNIQUE (name);


--
-- Name: models models_pkey; Type: CONSTRAINT; Schema: public; Owner: clearrouter
--

ALTER TABLE ONLY public.models
    ADD CONSTRAINT models_pkey PRIMARY KEY (id);


--
-- Name: payments payments_pkey; Type: CONSTRAINT; Schema: public; Owner: clearrouter
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: clearrouter
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: credits unique_user_id; Type: CONSTRAINT; Schema: public; Owner: clearrouter
--

ALTER TABLE ONLY public.credits
    ADD CONSTRAINT unique_user_id UNIQUE (user_id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: clearrouter
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: clearrouter
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_api_keys_user_id; Type: INDEX; Schema: public; Owner: clearrouter
--

CREATE INDEX idx_api_keys_user_id ON public.api_keys USING btree (user_id);


--
-- Name: idx_api_usage_analytics_created_at; Type: INDEX; Schema: public; Owner: clearrouter
--

CREATE INDEX idx_api_usage_analytics_created_at ON public.api_usage_analytics USING btree (created_at);


--
-- Name: idx_api_usage_analytics_model; Type: INDEX; Schema: public; Owner: clearrouter
--

CREATE INDEX idx_api_usage_analytics_model ON public.api_usage_analytics USING btree (model_requested);


--
-- Name: idx_api_usage_analytics_provider; Type: INDEX; Schema: public; Owner: clearrouter
--

CREATE INDEX idx_api_usage_analytics_provider ON public.api_usage_analytics USING btree (provider);


--
-- Name: idx_api_usage_analytics_status; Type: INDEX; Schema: public; Owner: clearrouter
--

CREATE INDEX idx_api_usage_analytics_status ON public.api_usage_analytics USING btree (status);


--
-- Name: idx_api_usage_analytics_user_id; Type: INDEX; Schema: public; Owner: clearrouter
--

CREATE INDEX idx_api_usage_analytics_user_id ON public.api_usage_analytics USING btree (user_id);


--
-- Name: idx_api_usage_logs_model_id; Type: INDEX; Schema: public; Owner: clearrouter
--

CREATE INDEX idx_api_usage_logs_model_id ON public.api_usage_logs USING btree (model_id);


--
-- Name: idx_api_usage_logs_user_id; Type: INDEX; Schema: public; Owner: clearrouter
--

CREATE INDEX idx_api_usage_logs_user_id ON public.api_usage_logs USING btree (user_id);


--
-- Name: idx_chat_messages_chat_id; Type: INDEX; Schema: public; Owner: clearrouter
--

CREATE INDEX idx_chat_messages_chat_id ON public.chat_messages USING btree (chat_id);


--
-- Name: idx_chat_messages_created_at; Type: INDEX; Schema: public; Owner: clearrouter
--

CREATE INDEX idx_chat_messages_created_at ON public.chat_messages USING btree (created_at);


--
-- Name: idx_chats_created_at; Type: INDEX; Schema: public; Owner: clearrouter
--

CREATE INDEX idx_chats_created_at ON public.chats USING btree (created_at DESC);


--
-- Name: idx_chats_user_id; Type: INDEX; Schema: public; Owner: clearrouter
--

CREATE INDEX idx_chats_user_id ON public.chats USING btree (user_id);


--
-- Name: idx_credits_user_id; Type: INDEX; Schema: public; Owner: clearrouter
--

CREATE INDEX idx_credits_user_id ON public.credits USING btree (user_id);


--
-- Name: idx_daily_usage_summary_date; Type: INDEX; Schema: public; Owner: clearrouter
--

CREATE INDEX idx_daily_usage_summary_date ON public.daily_usage_summary USING btree (date);


--
-- Name: idx_daily_usage_summary_user_date; Type: INDEX; Schema: public; Owner: clearrouter
--

CREATE INDEX idx_daily_usage_summary_user_date ON public.daily_usage_summary USING btree (user_id, date);


--
-- Name: idx_model_pricing_history_active; Type: INDEX; Schema: public; Owner: clearrouter
--

CREATE INDEX idx_model_pricing_history_active ON public.model_pricing_history USING btree (is_active);


--
-- Name: idx_model_pricing_history_model; Type: INDEX; Schema: public; Owner: clearrouter
--

CREATE INDEX idx_model_pricing_history_model ON public.model_pricing_history USING btree (model_id, provider);


--
-- Name: idx_payments_user_id; Type: INDEX; Schema: public; Owner: clearrouter
--

CREATE INDEX idx_payments_user_id ON public.payments USING btree (user_id);


--
-- Name: api_usage_analytics trigger_update_daily_usage_summary; Type: TRIGGER; Schema: public; Owner: clearrouter
--

CREATE TRIGGER trigger_update_daily_usage_summary AFTER INSERT ON public.api_usage_analytics FOR EACH ROW EXECUTE FUNCTION public.update_daily_usage_summary();


--
-- Name: api_keys api_keys_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: clearrouter
--

ALTER TABLE ONLY public.api_keys
    ADD CONSTRAINT api_keys_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: api_usage_logs api_usage_logs_api_key_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: clearrouter
--

ALTER TABLE ONLY public.api_usage_logs
    ADD CONSTRAINT api_usage_logs_api_key_id_fkey FOREIGN KEY (api_key_id) REFERENCES public.api_keys(id) ON DELETE CASCADE;


--
-- Name: api_usage_logs api_usage_logs_model_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: clearrouter
--

ALTER TABLE ONLY public.api_usage_logs
    ADD CONSTRAINT api_usage_logs_model_id_fkey FOREIGN KEY (model_id) REFERENCES public.models(id) ON DELETE CASCADE;


--
-- Name: api_usage_logs api_usage_logs_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: clearrouter
--

ALTER TABLE ONLY public.api_usage_logs
    ADD CONSTRAINT api_usage_logs_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: chat_messages chat_messages_chat_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: clearrouter
--

ALTER TABLE ONLY public.chat_messages
    ADD CONSTRAINT chat_messages_chat_id_fkey FOREIGN KEY (chat_id) REFERENCES public.chats(id) ON DELETE CASCADE;


--
-- Name: chats chats_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: clearrouter
--

ALTER TABLE ONLY public.chats
    ADD CONSTRAINT chats_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: credits credits_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: clearrouter
--

ALTER TABLE ONLY public.credits
    ADD CONSTRAINT credits_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: daily_usage_summary daily_usage_summary_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: clearrouter
--

ALTER TABLE ONLY public.daily_usage_summary
    ADD CONSTRAINT daily_usage_summary_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: api_usage_analytics fk_user_id; Type: FK CONSTRAINT; Schema: public; Owner: clearrouter
--

ALTER TABLE ONLY public.api_usage_analytics
    ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: payments payments_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: clearrouter
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

\unrestrict f1Gtu0u75mFEYcnz3eMhIm0ZNgMiAeCmiwbLa4ThWTrRU7y5dV7tgBXdSIJtHhy

