\restrict sd6uJMuYGBcfP3uIuRziE9d6DIJ9dl3FkTeB3VcDF5yYdb2FPkaTdelPNc6tdl1

-- Dumped from database version 15.14
-- Dumped by pg_dump version 17.6

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
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
-- Name: EXTENSION pgcrypto; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION pgcrypto IS 'cryptographic functions';


--
-- Name: update_daily_usage_summary(); Type: FUNCTION; Schema: public; Owner: -
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


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: api_keys; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.api_keys (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid,
    api_key text NOT NULL,
    active boolean DEFAULT true,
    created_at timestamp with time zone DEFAULT now()
);


--
-- Name: api_usage_analytics; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.api_usage_analytics (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid NOT NULL,
    api_key_id uuid NOT NULL,
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
    updated_at timestamp with time zone DEFAULT now()
);


--
-- Name: api_usage_logs; Type: TABLE; Schema: public; Owner: -
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
    provider character varying(100)
);


--
-- Name: chat_messages; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.chat_messages (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    chat_id uuid NOT NULL,
    role character varying(20) NOT NULL,
    content text NOT NULL,
    token_count integer DEFAULT 0,
    cost numeric(12,6) DEFAULT 0,
    created_at timestamp with time zone DEFAULT now(),
    CONSTRAINT chat_messages_role_check CHECK (((role)::text = ANY ((ARRAY['user'::character varying, 'assistant'::character varying, 'system'::character varying])::text[])))
);


--
-- Name: chats; Type: TABLE; Schema: public; Owner: -
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


--
-- Name: credits; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.credits (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid,
    total_credits numeric(12,2) DEFAULT 0,
    used_credits numeric(12,2) DEFAULT 0,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now()
);


--
-- Name: daily_usage_summary; Type: TABLE; Schema: public; Owner: -
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


--
-- Name: model_pricing_history; Type: TABLE; Schema: public; Owner: -
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


--
-- Name: models; Type: TABLE; Schema: public; Owner: -
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


--
-- Name: payments; Type: TABLE; Schema: public; Owner: -
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


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version character varying NOT NULL
);


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
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


--
-- Name: api_keys api_keys_api_key_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.api_keys
    ADD CONSTRAINT api_keys_api_key_key UNIQUE (api_key);


--
-- Name: api_keys api_keys_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.api_keys
    ADD CONSTRAINT api_keys_pkey PRIMARY KEY (id);


--
-- Name: api_usage_analytics api_usage_analytics_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.api_usage_analytics
    ADD CONSTRAINT api_usage_analytics_pkey PRIMARY KEY (id);


--
-- Name: api_usage_logs api_usage_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.api_usage_logs
    ADD CONSTRAINT api_usage_logs_pkey PRIMARY KEY (id);


--
-- Name: chat_messages chat_messages_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.chat_messages
    ADD CONSTRAINT chat_messages_pkey PRIMARY KEY (id);


--
-- Name: chats chats_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.chats
    ADD CONSTRAINT chats_pkey PRIMARY KEY (id);


--
-- Name: credits credits_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.credits
    ADD CONSTRAINT credits_pkey PRIMARY KEY (id);


--
-- Name: daily_usage_summary daily_usage_summary_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.daily_usage_summary
    ADD CONSTRAINT daily_usage_summary_pkey PRIMARY KEY (id);


--
-- Name: daily_usage_summary daily_usage_summary_user_id_date_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.daily_usage_summary
    ADD CONSTRAINT daily_usage_summary_user_id_date_key UNIQUE (user_id, date);


--
-- Name: model_pricing_history model_pricing_history_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.model_pricing_history
    ADD CONSTRAINT model_pricing_history_pkey PRIMARY KEY (id);


--
-- Name: models models_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.models
    ADD CONSTRAINT models_name_key UNIQUE (name);


--
-- Name: models models_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.models
    ADD CONSTRAINT models_pkey PRIMARY KEY (id);


--
-- Name: payments payments_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: credits unique_user_id; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.credits
    ADD CONSTRAINT unique_user_id UNIQUE (user_id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_api_keys_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_api_keys_user_id ON public.api_keys USING btree (user_id);


--
-- Name: idx_api_usage_analytics_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_api_usage_analytics_created_at ON public.api_usage_analytics USING btree (created_at);


--
-- Name: idx_api_usage_analytics_model; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_api_usage_analytics_model ON public.api_usage_analytics USING btree (model_requested);


--
-- Name: idx_api_usage_analytics_provider; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_api_usage_analytics_provider ON public.api_usage_analytics USING btree (provider);


--
-- Name: idx_api_usage_analytics_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_api_usage_analytics_status ON public.api_usage_analytics USING btree (status);


--
-- Name: idx_api_usage_analytics_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_api_usage_analytics_user_id ON public.api_usage_analytics USING btree (user_id);


--
-- Name: idx_api_usage_logs_model_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_api_usage_logs_model_id ON public.api_usage_logs USING btree (model_id);


--
-- Name: idx_api_usage_logs_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_api_usage_logs_user_id ON public.api_usage_logs USING btree (user_id);


--
-- Name: idx_chat_messages_chat_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_chat_messages_chat_id ON public.chat_messages USING btree (chat_id);


--
-- Name: idx_chat_messages_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_chat_messages_created_at ON public.chat_messages USING btree (created_at);


--
-- Name: idx_chats_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_chats_created_at ON public.chats USING btree (created_at DESC);


--
-- Name: idx_chats_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_chats_user_id ON public.chats USING btree (user_id);


--
-- Name: idx_credits_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_credits_user_id ON public.credits USING btree (user_id);


--
-- Name: idx_daily_usage_summary_date; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_daily_usage_summary_date ON public.daily_usage_summary USING btree (date);


--
-- Name: idx_daily_usage_summary_user_date; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_daily_usage_summary_user_date ON public.daily_usage_summary USING btree (user_id, date);


--
-- Name: idx_model_pricing_history_active; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_model_pricing_history_active ON public.model_pricing_history USING btree (is_active);


--
-- Name: idx_model_pricing_history_model; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_model_pricing_history_model ON public.model_pricing_history USING btree (model_id, provider);


--
-- Name: idx_payments_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_payments_user_id ON public.payments USING btree (user_id);


--
-- Name: api_usage_analytics trigger_update_daily_usage_summary; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trigger_update_daily_usage_summary AFTER INSERT ON public.api_usage_analytics FOR EACH ROW EXECUTE FUNCTION public.update_daily_usage_summary();


--
-- Name: api_keys api_keys_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.api_keys
    ADD CONSTRAINT api_keys_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: api_usage_analytics api_usage_analytics_api_key_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.api_usage_analytics
    ADD CONSTRAINT api_usage_analytics_api_key_id_fkey FOREIGN KEY (api_key_id) REFERENCES public.api_keys(id) ON DELETE CASCADE;


--
-- Name: api_usage_analytics api_usage_analytics_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.api_usage_analytics
    ADD CONSTRAINT api_usage_analytics_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: api_usage_logs api_usage_logs_api_key_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.api_usage_logs
    ADD CONSTRAINT api_usage_logs_api_key_id_fkey FOREIGN KEY (api_key_id) REFERENCES public.api_keys(id) ON DELETE CASCADE;


--
-- Name: api_usage_logs api_usage_logs_model_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.api_usage_logs
    ADD CONSTRAINT api_usage_logs_model_id_fkey FOREIGN KEY (model_id) REFERENCES public.models(id) ON DELETE CASCADE;


--
-- Name: api_usage_logs api_usage_logs_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.api_usage_logs
    ADD CONSTRAINT api_usage_logs_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: chat_messages chat_messages_chat_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.chat_messages
    ADD CONSTRAINT chat_messages_chat_id_fkey FOREIGN KEY (chat_id) REFERENCES public.chats(id) ON DELETE CASCADE;


--
-- Name: chats chats_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.chats
    ADD CONSTRAINT chats_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: credits credits_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.credits
    ADD CONSTRAINT credits_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: daily_usage_summary daily_usage_summary_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.daily_usage_summary
    ADD CONSTRAINT daily_usage_summary_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: payments payments_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

\unrestrict sd6uJMuYGBcfP3uIuRziE9d6DIJ9dl3FkTeB3VcDF5yYdb2FPkaTdelPNc6tdl1


--
-- Dbmate schema migrations
--

INSERT INTO public.schema_migrations (version) VALUES
    ('20250927072843'),
    ('20250927101618'),
    ('20250927102719'),
    ('20250927192900'),
    ('20250928134500'),
    ('20250929143000'),
    ('20250929150700');
