INSERT INTO public.quiz_types (name)
SELECT seed.name
FROM (VALUES ('MultipleChoice'), ('SingleChoice')) AS seed(name)
WHERE NOT EXISTS (
    SELECT 1
    FROM public.quiz_types AS quiz_type
    WHERE quiz_type.name = seed.name
       OR (seed.name = 'SingleChoice' AND quiz_type.name = 'SingleChoise')
);

ALTER TABLE public.questions
    ADD COLUMN quiz_type_id bigint;

UPDATE public.questions AS question
SET quiz_type_id = category.quiz_type_id
FROM public.categories AS category
WHERE category.id = question.category_id;

ALTER TABLE public.questions
    ALTER COLUMN quiz_type_id SET NOT NULL,
    ADD CONSTRAINT questions_quiz_types_fk
        FOREIGN KEY (quiz_type_id)
        REFERENCES public.quiz_types(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE;

CREATE INDEX idx_questions_quiz_type_id
    ON public.questions(quiz_type_id);

DROP INDEX public.idx_categories_quiz_type_id;

ALTER TABLE public.categories
    DROP CONSTRAINT categories_quiz_types_fk,
    DROP COLUMN quiz_type_id;
