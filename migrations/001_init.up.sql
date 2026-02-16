
CREATE TABLE IF NOT EXISTS public.quiz_types (
    id          bigserial PRIMARY KEY,
    name        text NOT NULL,
    created_at  timestamptz DEFAULT now() NOT NULL,
    updated_at  timestamptz DEFAULT now() NOT NULL,
    deleted_at  timestamptz NULL
);

CREATE TABLE IF NOT EXISTS public.categories (
    id           bigserial PRIMARY KEY,
    name         text NOT NULL,
    quiz_type_id bigint NOT NULL,
    created_at   timestamptz DEFAULT now() NOT NULL,
    updated_at   timestamptz DEFAULT now() NOT NULL,
    deleted_at   timestamptz NULL,
    CONSTRAINT categories_quiz_types_fk
        FOREIGN KEY (quiz_type_id)
        REFERENCES public.quiz_types(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS public.questions (
    id          bigserial PRIMARY KEY,
    question    text NOT NULL,
    category_id bigint NOT NULL,
    created_at  timestamptz DEFAULT now() NOT NULL,
    updated_at  timestamptz DEFAULT now() NOT NULL,
    deleted_at  timestamptz NULL,
    CONSTRAINT questions_categories_fk
        FOREIGN KEY (category_id)
        REFERENCES public.categories(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS public.answers (
    id           bigserial PRIMARY KEY,
    answer       text NOT NULL,
    question_id  bigint NOT NULL,
    is_correct   boolean DEFAULT false,
    created_at   timestamptz DEFAULT now() NOT NULL,
    updated_at   timestamptz DEFAULT now() NOT NULL,
    deleted_at   timestamptz NULL,
    CONSTRAINT answers_questions_fk
        FOREIGN KEY (question_id)
        REFERENCES public.questions(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);


CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,

    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_users_deleted_at
    ON users (deleted_at);


CREATE INDEX IF NOT EXISTS idx_questions_category_id
    ON public.questions(category_id);

CREATE INDEX IF NOT EXISTS idx_answers_question_id
    ON public.answers(question_id);

CREATE INDEX IF NOT EXISTS idx_categories_quiz_type_id
    ON public.categories(quiz_type_id);
