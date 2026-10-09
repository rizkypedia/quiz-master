ALTER TABLE public.categories
    ADD COLUMN quiz_type_id bigint;

UPDATE public.categories AS category
SET quiz_type_id = (
    SELECT question.quiz_type_id
    FROM public.questions AS question
    WHERE question.category_id = category.id
    ORDER BY question.id
    LIMIT 1
);

ALTER TABLE public.categories
    ADD CONSTRAINT categories_quiz_types_fk
        FOREIGN KEY (quiz_type_id)
        REFERENCES public.quiz_types(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE;

CREATE INDEX idx_categories_quiz_type_id
    ON public.categories(quiz_type_id);

DROP INDEX public.idx_questions_quiz_type_id;

ALTER TABLE public.questions
    DROP CONSTRAINT questions_quiz_types_fk,
    DROP COLUMN quiz_type_id;
